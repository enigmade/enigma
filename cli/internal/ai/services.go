package ai

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ServiceManager handles Ollama and ComfyUI systemd user services
type ServiceManager struct {
	ServiceDir string // ~/.config/systemd/user/
}

// NewServiceManager creates a service manager for the current user
func NewServiceManager() (*ServiceManager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	serviceDir := filepath.Join(home, ".config", "systemd", "user")
	if err := os.MkdirAll(serviceDir, 0o700); err != nil {
		return nil, err
	}

	return &ServiceManager{ServiceDir: serviceDir}, nil
}

// CreateOllamaService creates a systemd user service for Ollama
func (sm *ServiceManager) CreateOllamaService(modelDir string, port int) error {
	serviceFile := filepath.Join(sm.ServiceDir, "enigma-ollama.service")

	content := fmt.Sprintf(`[Unit]
Description=Enigma Ollama LLM Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/bin/ollama serve
Environment="OLLAMA_MODELS=%s"
Environment="OLLAMA_HOST=127.0.0.1:%d"
Restart=on-failure
RestartSec=10

[Install]
WantedBy=default.target
`, modelDir, port)

	if err := os.WriteFile(serviceFile, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write ollama.service: %w", err)
	}

	// Reload systemd user daemon
	cmd := exec.Command("systemctl", "--user", "daemon-reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("daemon-reload: %w", err)
	}

	fmt.Printf("✓ Ollama service created: %s\n", serviceFile)
	return nil
}

// StartService starts a systemd user service
func (sm *ServiceManager) StartService(serviceName string) error {
	cmd := exec.Command("systemctl", "--user", "start", serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("start %s: %w", serviceName, err)
	}
	fmt.Printf("✓ Started %s\n", serviceName)
	return nil
}

// StopService stops a systemd user service
func (sm *ServiceManager) StopService(serviceName string) error {
	cmd := exec.Command("systemctl", "--user", "stop", serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("stop %s: %w", serviceName, err)
	}
	fmt.Printf("✓ Stopped %s\n", serviceName)
	return nil
}

// StatusService checks the status of a service
func (sm *ServiceManager) StatusService(serviceName string) (string, error) {
	cmd := exec.Command("systemctl", "--user", "status", serviceName)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
