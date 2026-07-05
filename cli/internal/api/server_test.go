package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// fakeProvider is a test StateProvider returning canned data or an error.
type fakeProvider struct {
	state *State
	err   error
}

func (f *fakeProvider) Snapshot() (*State, error) {
	return f.state, f.err
}

func TestHandleStateReturnsJSON(t *testing.T) {
	provider := &fakeProvider{state: &State{
		Runtimes: []Runtime{{Name: "go", Version: "1.26", Active: true}},
		Services: []Service{{Name: "enigma-ollama", Status: "running", Port: 11434}},
		Projects: []Project{{Path: "/home/user/app", Port: 8080, URL: "https://app.test"}},
		Models:   []Model{{Name: "llama3.1:8b", SizeGB: 4.7, Backend: "ollama"}},
		GPU:      &GPU{Vendor: "NVIDIA", Model: "RTX 4090", VRAMGiB: 24},
	}}
	srv := NewServer(provider)

	req := httptest.NewRequest(http.MethodGet, "/v1/state", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var got State
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if len(got.Runtimes) != 1 || got.Runtimes[0].Name != "go" {
		t.Errorf("unexpected runtimes: %+v", got.Runtimes)
	}
	if got.GPU == nil || got.GPU.VRAMGiB != 24 {
		t.Errorf("unexpected GPU: %+v", got.GPU)
	}
}

func TestHandleStateProviderError(t *testing.T) {
	srv := NewServer(&fakeProvider{err: errors.New("boom")})

	req := httptest.NewRequest(http.MethodGet, "/v1/state", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestHandleStateRejectsPost(t *testing.T) {
	srv := NewServer(&fakeProvider{state: &State{}})

	req := httptest.NewRequest(http.MethodPost, "/v1/state", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestHandleHealth(t *testing.T) {
	srv := NewServer(&fakeProvider{state: &State{}})

	req := httptest.NewRequest(http.MethodGet, "/v1/health", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var body map[string]string
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	if body["status"] != "ok" {
		t.Errorf("expected status ok, got %v", body)
	}
}
