package ai

import "fmt"

// Ollama manages the Ollama LLM service
type Ollama struct {
	ModelDir string // /var/lib/enigma/models
	Port     int    // default 11434
}

// Start launches Ollama as a systemd user service
func (o *Ollama) Start() error {
	fmt.Println("TODO Phase 5.5: Start Ollama via systemd user service")
	fmt.Printf("  Models: %s\n", o.ModelDir)
	fmt.Printf("  Port: %d\n", o.Port)
	return nil
}

// Stop halts Ollama
func (o *Ollama) Stop() error {
	fmt.Println("TODO Phase 5.5: Stop Ollama systemd user service")
	return nil
}

// Status checks if Ollama is running
func (o *Ollama) Status() (bool, error) {
	fmt.Println("TODO Phase 5.5: Check Ollama systemd user service status")
	return false, nil
}
