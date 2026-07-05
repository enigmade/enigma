use tantivy::schema::{Schema, SchemaBuilder, STORED, STRING, TEXT};

/// Builds the Tantivy schema for the Enigma search index (SPEC §5 & §6).
///
/// Fields:
///   - path: exact-match stored field, used as the document's unique key
///   - filename: tokenized + stored, matched against user queries
///   - content: tokenized, not stored (large text bodies), used for full-text search
///   - modified: stored unix timestamp (u64) for recency ranking
pub fn build_schema() -> Schema {
    let mut builder: SchemaBuilder = Schema::builder();
    builder.add_text_field("path", STRING | STORED);
    builder.add_text_field("filename", TEXT | STORED);
    builder.add_text_field("content", TEXT);
    builder.add_u64_field("modified", STORED);
    builder.build()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn schema_has_expected_fields() {
        let schema = build_schema();
        assert!(schema.get_field("path").is_ok());
        assert!(schema.get_field("filename").is_ok());
        assert!(schema.get_field("content").is_ok());
        assert!(schema.get_field("modified").is_ok());
    }
}
