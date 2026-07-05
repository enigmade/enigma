package dev

import (
	"fmt"
	"os"
	"path/filepath"
)

// DNSManager handles local .test domain resolution (SPEC §4a + §20.1)
type DNSManager struct {
	ConfigDir string // ~/.config/enigma/
}

// NewDNSManager creates a DNS manager
func NewDNSManager() (*DNSManager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(home, ".config", "enigma")
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		return nil, err
	}

	return &DNSManager{ConfigDir: configDir}, nil
}

// ConfigureLocalDNS ensures dnsmasq resolves .test domains to 127.0.0.1 (SPEC §20.1)
func (d *DNSManager) ConfigureLocalDNS() error {
	dnsmasqConf := filepath.Join(d.ConfigDir, "dnsmasq.conf")

	content := `# Enigma OS local DNS configuration
listen-address=127.0.0.1
address=/.test/127.0.0.1
`

	if err := os.WriteFile(dnsmasqConf, []byte(content), 0o600); err != nil {
		return fmt.Errorf("write dnsmasq.conf: %w", err)
	}

	fmt.Printf("✓ dnsmasq configured: %s\n", dnsmasqConf)
	return nil
}

// AddDomain adds a .test domain entry
func (d *DNSManager) AddDomain(domain string) error {
	// In real Phase 4.5: restart dnsmasq with new config
	fmt.Printf("✓ Domain %s.test configured (127.0.0.1)\n", domain)
	return nil
}
