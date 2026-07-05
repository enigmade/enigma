package stack

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectStack(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		expected Stack
	}{
		{"PHP project", []string{"composer.json"}, PHP},
		{"Node project", []string{"package.json"}, Node},
		{"Python project", []string{"requirements.txt"}, Python},
		{"Go project", []string{"go.mod"}, Go},
		{"Rust project", []string{"Cargo.toml"}, Rust},
		{"Ruby project", []string{"Gemfile"}, Ruby},
		{"Docker Compose", []string{"docker-compose.yml"}, DockerCompose},
		{"No recognizable files", []string{}, Unknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpdir := t.TempDir()

			// Create the test files
			for _, file := range tt.files {
				fpath := filepath.Join(tmpdir, file)
				if err := os.WriteFile(fpath, []byte(""), 0o644); err != nil {
					t.Fatalf("Failed to create %s: %v", file, err)
				}
			}

			got := Detect(tmpdir)
			if got != tt.expected {
				t.Errorf("Detect() = %s, want %s", got, tt.expected)
			}
		})
	}
}
