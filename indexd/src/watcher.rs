use notify::{Event, EventKind, RecommendedWatcher, RecursiveMode, Watcher};
use std::path::{Path, PathBuf};
use std::sync::mpsc::{channel, Receiver};

/// A simplified file-change event handed up to the index worker.
#[derive(Debug, Clone, PartialEq)]
pub enum IndexEvent {
    Upsert(PathBuf),
    Remove(PathBuf),
}

/// Wraps `notify`'s cross-platform watcher (FSEvents on macOS, inotify on
/// Linux, fanotify semantics approximated via the same recursive watch) so
/// the same code path can be exercised in local dev and in the Enigma OS
/// target environment.
pub struct FileWatcher {
    _inner: RecommendedWatcher,
    rx: Receiver<IndexEvent>,
}

impl FileWatcher {
    pub fn new(root: &Path) -> notify::Result<Self> {
        let (tx, rx) = channel();

        let mut watcher = notify::recommended_watcher(move |res: notify::Result<Event>| {
            if let Ok(event) = res {
                for path in event.paths {
                    let mapped = match event.kind {
                        EventKind::Remove(_) => Some(IndexEvent::Remove(path)),
                        EventKind::Create(_) | EventKind::Modify(_) => Some(IndexEvent::Upsert(path)),
                        _ => None,
                    };
                    if let Some(ev) = mapped {
                        let _ = tx.send(ev);
                    }
                }
            }
        })?;

        watcher.watch(root, RecursiveMode::Recursive)?;

        Ok(Self { _inner: watcher, rx })
    }

    /// Non-blocking drain of pending events, for use in the daemon's poll loop.
    pub fn try_recv_all(&self) -> Vec<IndexEvent> {
        self.rx.try_iter().collect()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::fs;
    use std::thread;
    use std::time::Duration;

    #[test]
    fn detects_file_create_and_modify() {
        let dir = std::env::temp_dir().join(format!("enigma-watch-test-{}", std::process::id()));
        fs::create_dir_all(&dir).unwrap();

        let watcher = FileWatcher::new(&dir).expect("watcher should start");

        let file_path = dir.join("hello.txt");
        fs::write(&file_path, "hello").unwrap();

        // FSEvents/inotify delivery is async; poll briefly instead of a fixed sleep.
        let mut events = Vec::new();
        for _ in 0..20 {
            events.extend(watcher.try_recv_all());
            if !events.is_empty() {
                break;
            }
            thread::sleep(Duration::from_millis(100));
        }

        assert!(
            events.iter().any(|e| matches!(e, IndexEvent::Upsert(p) if p == &file_path)),
            "expected an Upsert event for {:?}, got {:?}",
            file_path,
            events
        );

        let _ = fs::remove_dir_all(&dir);
    }
}
