mod indexer;
mod ipc;
mod schema;
mod watcher;

use indexer::SearchIndex;
use ipc::{bind_socket, handle_client};
use std::env;
use std::path::PathBuf;
use std::sync::{Arc, Mutex};
use watcher::{FileWatcher, IndexEvent};

fn main() -> anyhow::Result<()> {
    println!("enigma-indexd — local search index daemon");
    println!("Per SPEC §5 & §6: Tantivy-based file + code indexer");

    let home = env::var("HOME").unwrap_or_else(|_| "/root".to_string());
    let watch_root = PathBuf::from(&home);

    let socket_path = env::var("XDG_RUNTIME_DIR")
        .map(|d| PathBuf::from(d).join("enigma-search.sock"))
        .unwrap_or_else(|_| PathBuf::from("/tmp/enigma-search.sock"));

    let index = Arc::new(Mutex::new(SearchIndex::create_in_ram()?));
    let watcher = FileWatcher::new(&watch_root)?;
    let listener = bind_socket(&socket_path)?;

    println!("✓ Watching {} (excludes: node_modules, .git, venv, __pycache__, ...)", watch_root.display());
    println!("✓ Listening on {}", socket_path.display());

    listener.set_nonblocking(true)?;

    loop {
        for event in watcher.try_recv_all() {
            let mut idx = index.lock().unwrap();
            match event {
                IndexEvent::Upsert(path) => {
                    if let Ok(content) = std::fs::read_to_string(&path) {
                        let _ = idx.index_file(&path, &content);
                    }
                }
                IndexEvent::Remove(path) => {
                    let _ = idx.remove_file(&path);
                }
            }
            let _ = idx.commit();
        }

        match listener.accept() {
            Ok((stream, _)) => {
                let idx_clone = Arc::clone(&index);
                std::thread::spawn(move || {
                    let _ = handle_client(stream, idx_clone);
                });
            }
            Err(ref e) if e.kind() == std::io::ErrorKind::WouldBlock => {
                std::thread::sleep(std::time::Duration::from_millis(50));
            }
            Err(e) => eprintln!("accept error: {e}"),
        }
    }
}
