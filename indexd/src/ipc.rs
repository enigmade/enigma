use crate::indexer::SearchIndex;
use serde::{Deserialize, Serialize};
use std::io::{BufRead, BufReader, Write};
use std::os::unix::net::{UnixListener, UnixStream};
use std::path::Path;
use std::sync::{Arc, Mutex};

/// Request/response shapes for the Vicinae launcher extension.
/// Wire format: newline-delimited JSON over a Unix socket
/// at /run/user/$UID/enigma-search.sock (SPEC §5).
#[derive(Debug, Deserialize)]
pub struct SearchRequest {
    pub query: String,
    #[serde(default = "default_limit")]
    pub limit: usize,
}

fn default_limit() -> usize {
    10
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SearchResponse {
    pub results: Vec<String>,
}

/// Handles a single client connection: read one JSON line, write one JSON
/// line back, then close. Kept simple since the launcher issues one query
/// per socket connection.
pub fn handle_client(stream: UnixStream, index: Arc<Mutex<SearchIndex>>) -> std::io::Result<()> {
    let mut reader = BufReader::new(stream.try_clone()?);
    let mut line = String::new();
    reader.read_line(&mut line)?;

    let response = match serde_json::from_str::<SearchRequest>(line.trim()) {
        Ok(req) => {
            let idx = index.lock().unwrap();
            match idx.search(&req.query, req.limit) {
                Ok(results) => SearchResponse { results },
                Err(_) => SearchResponse { results: vec![] },
            }
        }
        Err(_) => SearchResponse { results: vec![] },
    };

    let mut writer = stream;
    let payload = serde_json::to_string(&response).unwrap_or_else(|_| "{}".to_string());
    writer.write_all(payload.as_bytes())?;
    writer.write_all(b"\n")?;
    Ok(())
}

/// Binds the Unix socket, removing any stale socket file left by a prior
/// crashed run (systemd user services can leave these behind).
pub fn bind_socket(socket_path: &Path) -> std::io::Result<UnixListener> {
    if socket_path.exists() {
        std::fs::remove_file(socket_path)?;
    }
    UnixListener::bind(socket_path)
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::io::Read;
    use std::thread;

    #[test]
    fn search_request_parses_default_limit() {
        let req: SearchRequest = serde_json::from_str(r#"{"query":"enigma"}"#).unwrap();
        assert_eq!(req.query, "enigma");
        assert_eq!(req.limit, 10);
    }

    #[test]
    fn search_request_parses_explicit_limit() {
        let req: SearchRequest = serde_json::from_str(r#"{"query":"enigma","limit":5}"#).unwrap();
        assert_eq!(req.limit, 5);
    }

    #[test]
    fn end_to_end_socket_query() -> std::io::Result<()> {
        let socket_path =
            std::env::temp_dir().join(format!("enigma-ipc-test-{}.sock", std::process::id()));
        let _ = std::fs::remove_file(&socket_path);

        let mut index = SearchIndex::create_in_ram().expect("index");
        index
            .index_file(Path::new("/home/user/enigma-notes.md"), "enigma os plan")
            .unwrap();
        index.commit().unwrap();
        let shared = Arc::new(Mutex::new(index));

        let listener = bind_socket(&socket_path)?;
        let server_shared = Arc::clone(&shared);
        let server = thread::spawn(move || {
            if let Ok((stream, _)) = listener.accept() {
                let _ = handle_client(stream, server_shared);
            }
        });

        // Give the acceptor a moment to be ready; connect will retry briefly.
        let mut stream = None;
        for _ in 0..20 {
            if let Ok(s) = UnixStream::connect(&socket_path) {
                stream = Some(s);
                break;
            }
            thread::sleep(std::time::Duration::from_millis(50));
        }
        let mut stream = stream.expect("should connect to socket");

        stream.write_all(br#"{"query":"enigma","limit":10}"#)?;
        stream.write_all(b"\n")?;
        stream.shutdown(std::net::Shutdown::Write)?;

        let mut response = String::new();
        stream.read_to_string(&mut response)?;
        server.join().unwrap();

        let parsed: SearchResponse = serde_json::from_str(response.trim()).unwrap();
        assert_eq!(parsed.results.len(), 1);

        let _ = std::fs::remove_file(&socket_path);
        Ok(())
    }
}
