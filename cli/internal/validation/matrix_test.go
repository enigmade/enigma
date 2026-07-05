package validation

import (
	"reflect"
	"testing"
)

func TestDetectionOutcomeCorrect(t *testing.T) {
	cfg := HardwareConfig{"NVIDIA RTX 3080 (Ampere)", "NVIDIA", "Ampere", "cu128"}
	o := DetectionOutcome{
		Config:         cfg,
		DetectedVendor: "NVIDIA",
		DetectedGen:    "Ampere",
		DetectedTorch:  "cu128",
	}
	if !o.Correct() {
		t.Errorf("expected correct, mismatches: %v", o.Mismatches())
	}
	if len(o.Mismatches()) != 0 {
		t.Errorf("expected no mismatches, got %v", o.Mismatches())
	}
}

func TestDetectionOutcomeMismatch(t *testing.T) {
	cfg := HardwareConfig{"AMD RX 6800 (RDNA2)", "AMD", "RDNA2", "rocm6"}
	o := DetectionOutcome{
		Config:         cfg,
		DetectedVendor: "AMD",
		DetectedGen:    "RDNA2",
		DetectedTorch:  "cpu", // wrong: ROCm not detected
	}
	if o.Correct() {
		t.Error("expected outcome to be incorrect")
	}
	got := o.Mismatches()
	want := []string{`torch: got "cpu" want "rocm6"`}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAccuracy(t *testing.T) {
	outcomes := []DetectionOutcome{
		{Config: HardwareConfig{ExpectedVendor: "NVIDIA", ExpectedGeneration: "Ada", ExpectedTorch: "cu128"},
			DetectedVendor: "NVIDIA", DetectedGen: "Ada", DetectedTorch: "cu128"},
		{Config: HardwareConfig{ExpectedVendor: "AMD", ExpectedGeneration: "RDNA2", ExpectedTorch: "rocm6"},
			DetectedVendor: "AMD", DetectedGen: "RDNA2", DetectedTorch: "rocm6"},
		{Config: HardwareConfig{ExpectedVendor: "Intel", ExpectedGeneration: "Alchemist", ExpectedTorch: "cpu"},
			DetectedVendor: "Intel", DetectedGen: "Xe", DetectedTorch: "cpu"}, // wrong gen
	}
	got := Accuracy(outcomes)
	want := 2.0 / 3.0
	if got != want {
		t.Errorf("got accuracy %v, want %v", got, want)
	}
}

func TestAccuracyEmpty(t *testing.T) {
	if Accuracy(nil) != 0 {
		t.Error("empty matrix should have 0 accuracy")
	}
}

func TestStandardMatrixCoversAllVendors(t *testing.T) {
	matrix := StandardMatrix()
	vendors := map[string]bool{}
	for _, c := range matrix {
		vendors[c.ExpectedVendor] = true
	}
	for _, v := range []string{"NVIDIA", "AMD", "Intel"} {
		if !vendors[v] {
			t.Errorf("standard matrix missing vendor %s", v)
		}
	}
}
