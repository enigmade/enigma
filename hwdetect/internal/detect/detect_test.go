package detect

import (
	"testing"
)

func TestDetectVM(t *testing.T) {
	// Test with fake VM detector
	fakeVirt := &FakeVirtDetector{VM: true, Type: "qemu"}
	cfg, err := DetectWithVirt(fakeVirt)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !cfg.System.IsVM {
		t.Errorf("Expected IsVM=true, got %v", cfg.System.IsVM)
	}

	if cfg.AI.TorchTarget != "cpu" {
		t.Errorf("Expected torch=cpu in VM, got %s", cfg.AI.TorchTarget)
	}

	if cfg.Display.Session != "x11" {
		t.Errorf("Expected x11 session in VM, got %s", cfg.Display.Session)
	}
}

func TestDetectNonVM(t *testing.T) {
	fakeVirt := &FakeVirtDetector{VM: false, Type: "none"}
	cfg, err := DetectWithVirt(fakeVirt)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if cfg.System.IsVM {
		t.Errorf("Expected IsVM=false, got %v", cfg.System.IsVM)
	}

	if cfg.Display.Session != "wayland" {
		t.Errorf("Expected wayland session (primary), got %s", cfg.Display.Session)
	}
}

func TestNVIDIAGenerationBlackwell(t *testing.T) {
	gen, driver, cuda, err := DetectNVIDIAGeneration("2800")
	if err != nil {
		t.Fatalf("Detection failed: %v", err)
	}
	if gen != "Blackwell" {
		t.Errorf("Expected Blackwell, got %s", gen)
	}
	if driver != "nvidia-open-dkms" {
		t.Errorf("Expected nvidia-open-dkms, got %s", driver)
	}
	if cuda != "cu128" {
		t.Errorf("Expected cu128, got %s", cuda)
	}
}

func TestNVIDIAGenerationPascal(t *testing.T) {
	gen, driver, cuda, err := DetectNVIDIAGeneration("1000")
	if err != nil {
		t.Fatalf("Detection failed: %v", err)
	}
	if gen != "Pascal" {
		t.Errorf("Expected Pascal, got %s", gen)
	}
	if driver != "nvidia-dkms" {
		t.Errorf("Expected nvidia-dkms legacy, got %s", driver)
	}
	if cuda != "cu121" {
		t.Errorf("Expected cu121 for Pascal, got %s", cuda)
	}
}

func TestAMDROCmSupport(t *testing.T) {
	tests := []struct {
		arch      string
		supported bool
	}{
		{"RDNA2", true},
		{"CDNA", true},
		{"RDNA1", false}, // Community, marked false but noted
		{"Vega", false},
	}

	for _, tt := range tests {
		if got := AMDSupportsROCm(tt.arch); got != tt.supported {
			t.Errorf("AMDSupportsROCm(%s) = %v, want %v", tt.arch, got, tt.supported)
		}
	}
}
