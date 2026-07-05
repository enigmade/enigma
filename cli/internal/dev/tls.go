package dev

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// TLSManager handles local TLS certificate generation (SPEC §4a + §20.2)
type TLSManager struct {
	CertDir string // ~/.local/share/enigma/certs/
}

// NewTLSManager creates a TLS certificate manager
func NewTLSManager() (*TLSManager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	certDir := filepath.Join(home, ".local", "share", "enigma", "certs")
	if err := os.MkdirAll(certDir, 0o700); err != nil {
		return nil, err
	}

	return &TLSManager{CertDir: certDir}, nil
}

// InstallCARoot ensures mkcert CA root is installed (runs once)
func (t *TLSManager) InstallCARoot() error {
	cmd := exec.Command("mkcert", "-install")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mkcert install: %w", err)
	}
	fmt.Println("✓ mkcert CA root installed")
	return nil
}

// CreateCert creates a certificate for a .test domain
func (t *TLSManager) CreateCert(domain string) (string, string, error) {
	certPath := filepath.Join(t.CertDir, domain+".pem")
	keyPath := filepath.Join(t.CertDir, domain+"-key.pem")

	cmd := exec.Command("mkcert", "-cert-file", certPath, "-key-file", keyPath, domain)
	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("create cert for %s: %w", domain, err)
	}

	return certPath, keyPath, nil
}
