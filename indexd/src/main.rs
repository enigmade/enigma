// Enigma OS search index daemon
// Per SPEC §5 & §6: Tantivy-based indexer (fanotify watcher + Vicinae integration)

fn main() {
    println!("enigma-indexd v0.1.0");
    println!("TODO Phase 6: Implement Tantivy index + fanotify watcher");
    println!("  - Indexer: filenames, code contents, apps");
    println!("  - Watcher: fanotify on ~/.config/enigma monitoring");
    println!("  - IPC: unix socket for Vicinae launcher extension");
    println!("  - Excludes: node_modules, .git, venvs, .cache");
}
