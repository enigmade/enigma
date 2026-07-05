package api

import (
	"encoding/json"
	"net/http"
)

// Server is the HTTP layer of the Enigma Center daemon. It is transport
// agnostic: Handler() returns a mux that can be served over a Unix socket
// (production) or a TCP listener (tests via httptest).
type Server struct {
	provider StateProvider
}

// NewServer builds a Server backed by the given StateProvider.
func NewServer(provider StateProvider) *Server {
	return &Server{provider: provider}
}

// Handler returns the routed HTTP handler for the daemon API.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/state", s.handleState)
	mux.HandleFunc("/v1/health", s.handleHealth)
	return mux
}

func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	state, err := s.provider.Snapshot()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, state)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
