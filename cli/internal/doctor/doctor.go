package doctor

import (
	"fmt"
	"os"
	"path/filepath"

	"enigma/hwdetect/pkg/hardware"
)

// Report generates a full health report (SPEC §4)
func Report() {
	fmt.Println("Enigma OS Health Report")
	fmt.Println("═════════════════════════════════════════")
	fmt.Println()

	// GPU section (reads hardware.toml if present)
	reportGPU()

	// Ports section
	reportPorts()

	// Stubs for later phases
	fmt.Println("Boot Time:")
	fmt.Println("  ◐ Not yet implemented (Phase 9)")
	fmt.Println()

	fmt.Println("Snapshots:")
	fmt.Println("  ◐ Not yet implemented (Phase 2 local-only)")
	fmt.Println()

	fmt.Println("Security:")
	fmt.Println("  ◐ Not yet implemented (Phase 9)")
}

func reportGPU() {
	fmt.Println("GPU:")

	// Try to read hardware.toml (if hwdetect has run)
	hwPath := filepath.Join(os.Getenv("HOME"), ".config", "enigma", "hardware.toml")
	if hwPath == filepath.Join("", ".config", "enigma", "hardware.toml") {
		hwPath = "/etc/enigma/hardware.toml"
	}

	cfg, err := hardware.ReadFromFile(hwPath)
	if err != nil {
		fmt.Println("  ◐ Hardware not yet detected (run 'enigma-hwdetect detect' first)")
		return
	}

	if cfg.System.IsVM {
		fmt.Println("  ⚠ Running in VM (drivers disabled)")
	}

	if len(cfg.GPU) == 0 {
		fmt.Println("  ◐ No GPUs detected")
		return
	}

	for _, gpu := range cfg.GPU {
		fmt.Printf("  • %s (%s): %s\n", gpu.Vendor, gpu.Generation, gpu.Name)
		fmt.Printf("    VRAM: %d MB, Driver: %s\n", gpu.VRAMMb, gpu.DriverState)
	}

	fmt.Printf("  AI Torch target: %s (%s)\n", cfg.AI.TorchTarget, cfg.AI.Tier)
	fmt.Println()
}

func reportPorts() {
	fmt.Println("Dev Ports:")
	fmt.Println("  (Run 'enigma ports' to see allocated ports)")
	fmt.Println()
}
