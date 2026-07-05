package ai

import (
	"fmt"

	"enigma/hwdetect/pkg/hardware"
)

// TorchTarget determines the PyTorch wheel to install
// Per hwdetect/DESIGN.md §TORCH MATCHER
func TorchTarget(hwPath string) (string, error) {
	cfg, err := hardware.ReadFromFile(hwPath)
	if err != nil {
		return "cpu", nil // Default to CPU if no detection
	}

	// If already verified, return that
	if cfg.AI.Verified {
		return cfg.AI.TorchTarget, nil
	}

	// Decision tree: vendor → generation → rocm support
	if len(cfg.GPU) == 0 {
		return "cpu", nil
	}

	gpu := cfg.GPU[0]
	switch gpu.Vendor {
	case "NVIDIA":
		// Turing+ → cu128, Pascal → cu121, older → cpu
		switch gpu.Generation {
		case "Blackwell", "Ada", "Ampere", "Turing", "Maxwell":
			return "cu128", nil
		case "Pascal":
			return "cu121", nil
		default:
			return "cpu", nil
		}

	case "AMD":
		if cfg.AI.ROCmSupported {
			return "rocm6", nil
		}
		return "cpu", nil

	default:
		return "cpu", nil
	}
}

// Setup installs the correct PyTorch wheel for the detected GPU
func Setup(hwPath string) error {
	target, err := TorchTarget(hwPath)
	if err != nil {
		return err
	}

	fmt.Printf("Installing PyTorch for: %s\n", target)

	// In real Phase 5: use uv to install torch wheel matching the target
	// For now: stub with message
	fmt.Printf("TODO Phase 5.5: install torch[%s] via uv\n", target)

	return nil
}
