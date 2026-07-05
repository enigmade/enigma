package ports

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ProjectPort tracks a port allocation for a project
type ProjectPort struct {
	ProjectPath string `json:"project_path"`
	Port        int    `json:"port"`
	PID         int    `json:"pid,omitempty"`
	URL         string `json:"url,omitempty"`
}

// DB manages the ports database at ~/.config/enigma/ports.db
func dbPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dbDir := filepath.Join(home, ".config", "enigma")
	if err := os.MkdirAll(dbDir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(dbDir, "ports.db"), nil
}

// LoadDB reads the ports database
func LoadDB() ([]ProjectPort, error) {
	path, err := dbPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []ProjectPort{}, nil // Empty DB is OK
		}
		return nil, err
	}

	var ports []ProjectPort
	if err := json.Unmarshal(data, &ports); err != nil {
		return nil, fmt.Errorf("parse ports.db: %w", err)
	}
	return ports, nil
}

// SaveDB writes the ports database
func SaveDB(ports []ProjectPort) error {
	path, err := dbPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(ports, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o600)
}
