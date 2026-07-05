package ai

import (
	"os"
	"testing"
)

func TestTorchTargetCPU(t *testing.T) {
	// Non-existent hardware.toml should default to CPU
	target, err := TorchTarget("/nonexistent/hardware.toml")
	if err != nil {
		t.Fatalf("Expected error to be nil for missing file, got %v", err)
	}
	if target != "cpu" {
		t.Errorf("Expected cpu for missing hardware.toml, got %s", target)
	}
}

func TestTorchTargetParsing(t *testing.T) {
	// Create a minimal valid hardware.toml
	tmpfile, err := os.CreateTemp("", "hardware-*.toml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	// Write minimal config
	if _, err := tmpfile.WriteString(`
[system]
ram_gb = 16
cpu_vendor = "Intel"
cpu_flags = ["v3"]
is_vm = false
is_laptop = false

[ai]
torch_target = "cpu"
tier = "entry"
rocm_supported = false
verified = false
`); err != nil {
		t.Fatalf("Failed to write hardware.toml: %v", err)
	}
	tmpfile.Close()

	// Should parse without error
	target, err := TorchTarget(tmpfile.Name())
	if err != nil {
		t.Fatalf("Parsing failed: %v", err)
	}
	if target != "cpu" {
		t.Errorf("Expected cpu, got %s", target)
	}
}
