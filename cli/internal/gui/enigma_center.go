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
// IPC: HTTP over Unix socket /run/user/$UID/enigma-center.sock
//
// The daemon API is implemented in cli/internal/api:
//   - GET /v1/state   -> full snapshot the GUI renders (all tabs)
//   - GET /v1/health  -> liveness probe
//
// The QML frontend (enigma-center, packaged separately) is a thin client:
// on launch and on a timer it GETs /v1/state and repaints; action buttons
// (start/stop service, open project) POST back to the same socket. It holds
// no business logic — everything comes from api.State.
//
// Remaining GUI work (not Go, lives in the enigma-center Qt package):
//   - main.qml tabbed shell (Runtimes/Services/Projects/AI/Windows/System)
//   - socket client that unmarshals api.State
//   - action POST handlers wiring to enigma dev/ai/game subcommands
