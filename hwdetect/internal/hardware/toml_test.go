package hardware

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMarshalUnmarshal(t *testing.T) {
	cfg := &Config{
		System: SystemInfo{
			RAMGb:     16,
			CPUVendor: "Intel",
			CPUFlags:  []string{"v3", "v4"},
			IsVM:      false,
			IsLaptop:  false,
		},
		GPU: []GPUInfo{
			{
				Vendor:      "NVIDIA",
				PCIId:       "2484",
				Name:        "GA104 [GeForce RTX 3070]",
				VRAMMb:      8192,
				Generation: "Ampere",
				DriverState: "installed",
				Role:        "primary",
			},
		},
		AI: AIConfig{
			TorchTarget:   "cu128",
			Tier:          "standard",
			ROCmSupported: false,
			Verified:      true,
		},
		Display: DisplayConfig{
			Session:    "wayland",
			HybridMode: "none",
		},
	}

	// Write to temp file
	tmpFile := filepath.Join(t.TempDir(), "hardware.toml")
	if err := cfg.WriteToFile(tmpFile); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Read it back
	loaded, err := ReadFromFile(tmpFile)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	// Verify
	if loaded.System.RAMGb != 16 {
		t.Errorf("RAM mismatch: expected 16, got %d", loaded.System.RAMGb)
	}
	if len(loaded.GPU) != 1 {
		t.Errorf("GPU count mismatch: expected 1, got %d", len(loaded.GPU))
	}
	if loaded.GPU[0].Vendor != "NVIDIA" {
		t.Errorf("GPU vendor mismatch: expected NVIDIA, got %s", loaded.GPU[0].Vendor)
	}
	if loaded.AI.TorchTarget != "cu128" {
		t.Errorf("Torch target mismatch: expected cu128, got %s", loaded.AI.TorchTarget)
	}

	// Cleanup
	os.Remove(tmpFile)
}

func TestString(t *testing.T) {
	cfg := &Config{
		System: SystemInfo{
			RAMGb:     8,
			CPUVendor: "Intel",
			IsVM:      false,
			IsLaptop:  true,
		},
		GPU: []GPUInfo{},
		AI: AIConfig{
			TorchTarget: "cpu",
			Tier:        "entry",
		},
	}

	str := cfg.String()
	if len(str) == 0 {
		t.Error("String() returned empty")
	}
	if !contains(str, "8 GB RAM") {
		t.Errorf("String() missing RAM info: %s", str)
	}
}

func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
