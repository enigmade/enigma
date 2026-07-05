package gui

// Phase 7.5: Enigma Center GUI (Qt/QML thin client) — STUB
//
// Provides unified dashboard for:
//   - Runtimes tab: installed Python/Node/Go/Rust versions
//   - Services tab: active Ollama/ComfyUI/indexd processes
//   - Projects tab: enigma dev up + port allocations
//   - AI Models tab: installed model list + VRAM usage
//   - Windows Apps tab: Wine Tier 1 shortcuts
//   - System tab: boot time, update status, GPU/CPU load
//
// Architecture: API daemon (Go) + QML frontend (Qt 6)
// IPC: Unix socket /run/user/$UID/enigma-center.sock (JSON-RPC 2.0)
//
// TODO Phase 7.5:
//   - enigma/api: HTTP daemon serving /v1/runtimes, /v1/services, /v1/projects
//   - enigma-center: Qt application with tabbed UI
//   - Integration: call enigma dev/ai/game subcommands via IPC
//   - Per-project service integration
