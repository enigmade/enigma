package api

// This package implements the Enigma Center daemon API (Phase 7.5 - SPEC §7.5).
//
// The Qt/QML frontend (enigma-center) is a thin client: it renders whatever
// this daemon reports and issues action requests back over the same HTTP
// interface bound to a Unix socket at /run/user/$UID/enigma-center.sock.
//
// Everything user-facing in the GUI maps to one of the State fields below,
// so the GUI holds no business logic of its own.

// Runtime is an installed language runtime version (Runtimes tab).
type Runtime struct {
	Name    string `json:"name"`    // python, node, go, rust
	Version string `json:"version"` // 3.12.4, 22.3.0, ...
	Active  bool   `json:"active"`  // is this the mise-selected default
}

// Service is a managed background service (Services tab).
type Service struct {
	Name   string `json:"name"`   // enigma-ollama, enigma-indexd, ...
	Status string `json:"status"` // running, stopped, failed
	Port   int    `json:"port,omitempty"`
}

// Project is a dev project with an allocated port (Projects tab).
type Project struct {
	Path string `json:"path"`
	Port int    `json:"port"`
	URL  string `json:"url,omitempty"`
}

// Model is an installed AI model (AI Models tab).
type Model struct {
	Name    string `json:"name"`
	SizeGB  float64 `json:"size_gb"`
	Backend string `json:"backend"` // ollama, comfyui
}

// GPU mirrors the detected GPU from hardware.toml (System tab).
type GPU struct {
	Vendor  string `json:"vendor"`
	Model   string `json:"model"`
	VRAMGiB int    `json:"vram_gib"`
}

// State is the complete snapshot the GUI renders. A single GET /v1/state
// returns this so the frontend can paint every tab from one request.
type State struct {
	Runtimes []Runtime `json:"runtimes"`
	Services []Service `json:"services"`
	Projects []Project `json:"projects"`
	Models   []Model   `json:"models"`
	GPU      *GPU      `json:"gpu,omitempty"`
}

// StateProvider is the seam between the HTTP layer and the live system.
// The real implementation reads ports.db, hardware.toml, systemctl, and
// mise; tests supply a fake so the HTTP layer can be exercised in isolation.
type StateProvider interface {
	Snapshot() (*State, error)
}
