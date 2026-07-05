package dev

import (
	"fmt"
	"os"

	"enigma/cli/internal/ports"
	"enigma/cli/internal/stack"
)

// Up detects stack and starts a dev environment (SPEC §4a)
func Up() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// Detect the project stack
	detected := stack.Detect(cwd)
	fmt.Printf("Detected stack: %s\n", detected.Description())

	// Allocate a port (SPEC §20.9: bind(0) probe, no race)
	port, err := ports.Allocate(8000)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error allocating port: %v\n", err)
		return
	}

	// Record in ports.db
	db, _ := ports.LoadDB()
	db = append(db, ports.ProjectPort{
		ProjectPath: cwd,
		Port:        port,
		URL:         fmt.Sprintf("http://localhost:%d", port),
	})
	if err := ports.SaveDB(db); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving port: %v\n", err)
		return
	}

	fmt.Printf("✓ Project environment ready\n")
	fmt.Printf("  Port: %d\n", port)
	fmt.Printf("  URL: http://localhost:%d\n", port)
	fmt.Printf("  (TODO Phase 4.5: run mise/systemd services/mkcert .test domain)\n")
}

// Down stops dev services for the current project
func Down() {
	fmt.Println("TODO Phase 4.5: stop systemd user services for this project")
}

// Services lists active services
func Services() {
	fmt.Println("TODO Phase 4.5: list systemd user services with status/ports")
}
