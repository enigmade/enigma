package stack

import (
	"os"
	"path/filepath"
)

type Stack string

const (
	PHP         Stack = "php"
	Node        Stack = "node"
	Python      Stack = "python"
	Go          Stack = "go"
	Rust        Stack = "rust"
	Ruby        Stack = "ruby"
	DockerCompose Stack = "docker-compose"
	Unknown     Stack = "unknown"
)

// Detect identifies the project stack from files in the given directory
func Detect(projectDir string) Stack {
	// Check for stack indicator files in priority order
	checks := []struct {
		files []string
		stack Stack
	}{
		{[]string{"go.mod"}, Go},
		{[]string{"Cargo.toml"}, Rust},
		{[]string{"composer.json", "composer.lock"}, PHP},
		{[]string{"package.json"}, Node},
		{[]string{"requirements.txt", "pyproject.toml", "setup.py", "Pipfile"}, Python},
		{[]string{"Gemfile", "Gemfile.lock"}, Ruby},
		{[]string{"docker-compose.yml", "docker-compose.yaml"}, DockerCompose},
	}

	for _, check := range checks {
		for _, file := range check.files {
			if _, err := os.Stat(filepath.Join(projectDir, file)); err == nil {
				return check.stack
			}
		}
	}

	return Unknown
}

// Description returns a human-readable description of the stack
func (s Stack) Description() string {
	switch s {
	case PHP:
		return "PHP (Composer)"
	case Node:
		return "Node.js (npm/yarn)"
	case Python:
		return "Python (pip/poetry)"
	case Go:
		return "Go"
	case Rust:
		return "Rust (Cargo)"
	case Ruby:
		return "Ruby (Bundler)"
	case DockerCompose:
		return "Docker Compose"
	default:
		return "Unknown"
	}
}
