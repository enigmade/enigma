package hardware

import (
	"fmt"
	"os"

	toml "github.com/pelletier/go-toml/v2"
)

// Config represents the complete hardware detection output (/etc/enigma/hardware.toml)
// Per SPEC §2 and hwdetect/DESIGN.md §OUTPUT CONTRACT
type Config struct {
	System  SystemInfo      `toml:"system"`
	GPU     []GPUInfo       `toml:"gpu"`
	AI      AIConfig        `toml:"ai"`
	Display DisplayConfig   `toml:"display"`
}

// SystemInfo: CPU, RAM, VM detection, battery, laptop detection
type SystemInfo struct {
	RAMGb       int      `toml:"ram_gb"`
	CPUVendor   string   `toml:"cpu_vendor"`
	CPUFlags    []string `toml:"cpu_flags"`
	IsVM        bool     `toml:"is_vm"`
	IsLaptop    bool     `toml:"is_laptop"`
	BatteryPath string   `toml:"battery_path,omitempty"`
}

// GPUInfo: vendor, device ID, driver state, generation, role
type GPUInfo struct {
	Index         int    `toml:"-"` // Not written to TOML, used internally
	Vendor        string `toml:"vendor"`
	PCIId         string `toml:"pci_id"`
	Name          string `toml:"name"`
	VRAMMb        int    `toml:"vram_mb"`
	DriverState   string `toml:"driver_state"` // none|installed|mismatch
	Generation    string `toml:"generation"`   // Turing, Ada, Blackwell, RDNA2, etc.
	Role          string `toml:"role"`         // primary|render|compute
	CurrentDriver string `toml:"current_driver,omitempty"`
}

// AIConfig: torch target, tier, ROCm support, verification
type AIConfig struct {
	TorchTarget    string   `toml:"torch_target"`    // cu128|cu121|rocm6|cpu
	Tier           string   `toml:"tier"`            // entry|standard|creator|studio
	ROCmSupported  bool     `toml:"rocm_supported"`
	Verified       bool     `toml:"verified"`
	VerifyDate     string   `toml:"verify_date,omitempty"`
	Notes          []string `toml:"notes,omitempty"`
}

// DisplayConfig: session type, hybrid mode (if applicable)
type DisplayConfig struct {
	Session    string `toml:"session"`    // wayland|x11
	HybridMode string `toml:"hybrid_mode"` // none|prime|switcheroo
}

// WriteToFile marshals the config to TOML and writes to path
func (c *Config) WriteToFile(path string) error {
	// Ensure directory exists
	dir := os.Getenv("ENIGMA_CONFIG_DIR")
	if dir == "" {
		dir = "/etc/enigma"
	}

	data, err := toml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshal hardware config: %w", err)
	}

	// Write to file (mode 0644, world-readable)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write hardware config: %w", err)
	}
	return nil
}

// ReadFromFile reads and parses a hardware.toml file
func ReadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read hardware config: %w", err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal hardware config: %w", err)
	}
	return &cfg, nil
}

// String returns a human-readable summary
func (c *Config) String() string {
	summary := fmt.Sprintf("System: %d GB RAM, CPU: %s, VM: %v, Laptop: %v\n",
		c.System.RAMGb, c.System.CPUVendor, c.System.IsVM, c.System.IsLaptop)

	for i, gpu := range c.GPU {
		summary += fmt.Sprintf("GPU %d: %s (%s) - %s - Driver: %s - Role: %s\n",
			i, gpu.Vendor, gpu.PCIId, gpu.Name, gpu.DriverState, gpu.Role)
	}

	summary += fmt.Sprintf("AI: Torch target=%s, Tier=%s, ROCm supported=%v, Verified=%v\n",
		c.AI.TorchTarget, c.AI.Tier, c.AI.ROCmSupported, c.AI.Verified)

	return summary
}
