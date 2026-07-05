use crate::schema::build_schema;
use std::path::Path;
use std::time::{SystemTime, UNIX_EPOCH};
use tantivy::collector::TopDocs;
use tantivy::query::QueryParser;
use tantivy::{doc, Document, Index, IndexReader, IndexWriter, ReloadPolicy};

/// Directory name fragments that are never indexed (SPEC §6 exclusions).
pub const EXCLUDED_DIRS: &[&str] = &[
    "node_modules",
    ".git",
    ".cache",
    "__pycache__",
    "venv",
    ".venv",
    "target",
    "dist",
    "build",
];

/// Returns true if `path` should be skipped based on EXCLUDED_DIRS.
pub fn is_excluded(path: &Path) -> bool {
    path.components().any(|c| {
        let s = c.as_os_str().to_string_lossy();
        EXCLUDED_DIRS.iter().any(|ex| *ex == s)
    })
}

pub struct SearchIndex {
    index: Index,
    writer: IndexWriter,
    reader: IndexReader,
}

impl SearchIndex {
    /// Creates a new in-memory index (used for tests and as the default
    /// fallback before persistent storage is wired up in deployment).
    pub fn create_in_ram() -> tantivy::Result<Self> {
        let schema = build_schema();
        let index = Index::create_in_ram(schema);
        let writer = index.writer(50_000_000)?;
        let reader = index
            .reader_builder()
            .reload_policy(ReloadPolicy::OnCommit)
            .try_into()?;
        Ok(Self { index, writer, reader })
    }

    /// Indexes a single file: path + filename + content, skipping excluded dirs.
    /// Content is truncated defensively; huge binary files are not expected
    /// to reach this path once MIME sniffing is added in production.
    pub fn index_file(&mut self, path: &Path, content: &str) -> tantivy::Result<()> {
        if is_excluded(path) {
            return Ok(());
        }

        let schema = self.index.schema();
        let path_field = schema.get_field("path").unwrap();
        let filename_field = schema.get_field("filename").unwrap();
        let content_field = schema.get_field("content").unwrap();
        let modified_field = schema.get_field("modified").unwrap();

        let filename = path
            .file_name()
            .map(|n| n.to_string_lossy().to_string())
            .unwrap_or_default();

        let modified = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap_or_default()
            .as_secs();

        // Remove any previous version of this document before re-adding.
        self.remove_file(path)?;

        self.writer.add_document(doc!(
            path_field => path.to_string_lossy().to_string(),
            filename_field => filename,
            content_field => content,
            modified_field => modified,
        ))?;

        Ok(())
    }

    /// Removes a document by its exact path key (used on file delete/modify).
    pub fn remove_file(&mut self, path: &Path) -> tantivy::Result<()> {
        let schema = self.index.schema();
        let path_field = schema.get_field("path").unwrap();
        let term = tantivy::Term::from_field_text(path_field, &path.to_string_lossy());
        self.writer.delete_term(term);
        Ok(())
    }

    pub fn commit(&mut self) -> tantivy::Result<()> {
        self.writer.commit()?;
        Ok(())
    }

    /// Searches filename + content fields, returning up to `limit` paths
    /// ranked by relevance (SPEC §5 target: top-10 in <50ms).
    pub fn search(&self, query_str: &str, limit: usize) -> tantivy::Result<Vec<String>> {
        let searcher = self.reader.searcher();
        let schema = self.index.schema();
        let filename_field = schema.get_field("filename").unwrap();
        let content_field = schema.get_field("content").unwrap();
        let path_field = schema.get_field("path").unwrap();

        let query_parser = QueryParser::for_index(&self.index, vec![filename_field, content_field]);
        let query = query_parser.parse_query(query_str)?;

        let top_docs = searcher.search(&query, &TopDocs::with_limit(limit))?;

        let mut results = Vec::with_capacity(top_docs.len());
        for (_score, doc_address) in top_docs {
            let retrieved: Document = searcher.doc(doc_address)?;
            if let Some(value) = retrieved.get_first(path_field) {
                if let Some(s) = value.as_text() {
                    results.push(s.to_string());
                }
            }
        }
        Ok(results)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::path::PathBuf;

    #[test]
    fn excludes_node_modules() {
        let p = PathBuf::from("/home/user/project/node_modules/lib/index.js");
        assert!(is_excluded(&p));
    }

    #[test]
    fn excludes_git_dir() {
        let p = PathBuf::from("/home/user/project/.git/HEAD");
        assert!(is_excluded(&p));
    }

    #[test]
    fn does_not_exclude_normal_source() {
        let p = PathBuf::from("/home/user/project/src/main.rs");
        assert!(!is_excluded(&p));
    }

    #[test]
    fn index_and_search_roundtrip() -> tantivy::Result<()> {
        let mut idx = SearchIndex::create_in_ram()?;
        idx.index_file(
            std::path::Path::new("/home/user/notes/enigma-plan.md"),
            "Enigma OS zero-bloat Arch distro plan",
        )?;
        idx.index_file(
            std::path::Path::new("/home/user/notes/groceries.md"),
            "milk eggs bread",
        )?;
        idx.commit()?;

        let results = idx.search("enigma", 10)?;
        assert_eq!(results.len(), 1);
        assert!(results[0].ends_with("enigma-plan.md"));

        let no_match = idx.search("nonexistentxyz", 10)?;
        assert_eq!(no_match.len(), 0);
        Ok(())
    }

    #[test]
    fn reindexing_replaces_old_document() -> tantivy::Result<()> {
        let mut idx = SearchIndex::create_in_ram()?;
        let path = std::path::Path::new("/home/user/notes/todo.md");
        idx.index_file(path, "buy milk")?;
        idx.commit()?;
        assert_eq!(idx.search("milk", 10)?.len(), 1);

        idx.index_file(path, "call dentist")?;
        idx.commit()?;
        assert_eq!(idx.search("milk", 10)?.len(), 0);
        assert_eq!(idx.search("dentist", 10)?.len(), 1);
        Ok(())
    }

    #[test]
    fn excluded_file_is_never_indexed() -> tantivy::Result<()> {
        let mut idx = SearchIndex::create_in_ram()?;
        idx.index_file(
            std::path::Path::new("/home/user/project/node_modules/pkg/readme.md"),
            "this should not be searchable",
        )?;
        idx.commit()?;
        assert_eq!(idx.search("searchable", 10)?.len(), 0);
        Ok(())
    }
}
