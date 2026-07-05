package dev

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// DatabaseManager handles per-project database creation (SPEC §4a)
type DatabaseManager struct {
	ServiceDir string // systemd user service dir
}

// NewDatabaseManager creates a database manager
func NewDatabaseManager() (*DatabaseManager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	serviceDir := filepath.Join(home, ".config", "systemd", "user")
	if err := os.MkdirAll(serviceDir, 0o700); err != nil {
		return nil, err
	}

	return &DatabaseManager{ServiceDir: serviceDir}, nil
}

// CreateDatabase creates a per-project database and Beekeeper connection
func (d *DatabaseManager) CreateDatabase(projectName, dbType string) error {
	// dbType: mysql, postgres
	fmt.Printf("TODO Phase 4.5: Create %s database '%s'\n", dbType, projectName)
	fmt.Printf("  - Start MariaDB/PostgreSQL user service\n")
	fmt.Printf("  - Create project-specific DB user\n")
	fmt.Printf("  - Write Beekeeper connection file to ~/.config/beekeeper-studio/connections/\n")
	return nil
}

// ListDatabases lists active databases for the project
func (d *DatabaseManager) ListDatabases() error {
	cmd := exec.Command("systemctl", "--user", "list-units", "--type=service", "--state=running")
	output, _ := cmd.CombinedOutput()
	fmt.Printf("Active services:\n%s\n", string(output))
	return nil
}
