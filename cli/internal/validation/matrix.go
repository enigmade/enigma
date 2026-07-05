package validation

import "fmt"

// HardwareConfig is one entry in the validation matrix (SPEC §9): a known
// hardware setup with the GPU vendor/generation hwdetect is expected to
// report and the torch target it should select.
type HardwareConfig struct {
	Label              string // "NVIDIA RTX 3080 (Ampere)"
	ExpectedVendor     string // NVIDIA, AMD, Intel
	ExpectedGeneration string // Ampere, RDNA2, ...
	ExpectedTorch      string // cu128, cu121, rocm6, cpu
}

// DetectionOutcome pairs a config with what hwdetect actually reported on
// that machine during a validation run.
type DetectionOutcome struct {
	Config          HardwareConfig
	DetectedVendor  string
	DetectedGen     string
	DetectedTorch   string
}

// Correct reports whether the outcome matched all expected fields.
func (o DetectionOutcome) Correct() bool {
	return o.DetectedVendor == o.Config.ExpectedVendor &&
		o.DetectedGen == o.Config.ExpectedGeneration &&
		o.DetectedTorch == o.Config.ExpectedTorch
}

// Mismatches returns a human-readable list of the fields that differed,
// empty if the outcome was fully correct.
func (o DetectionOutcome) Mismatches() []string {
	var m []string
	if o.DetectedVendor != o.Config.ExpectedVendor {
		m = append(m, fmt.Sprintf("vendor: got %q want %q", o.DetectedVendor, o.Config.ExpectedVendor))
	}
	if o.DetectedGen != o.Config.ExpectedGeneration {
		m = append(m, fmt.Sprintf("generation: got %q want %q", o.DetectedGen, o.Config.ExpectedGeneration))
	}
	if o.DetectedTorch != o.Config.ExpectedTorch {
		m = append(m, fmt.Sprintf("torch: got %q want %q", o.DetectedTorch, o.Config.ExpectedTorch))
	}
	return m
}

// Accuracy returns the fraction of outcomes that were fully correct, in
// [0,1]. An empty slice returns 0.
func Accuracy(outcomes []DetectionOutcome) float64 {
	if len(outcomes) == 0 {
		return 0
	}
	correct := 0
	for _, o := range outcomes {
		if o.Correct() {
			correct++
		}
	}
	return float64(correct) / float64(len(outcomes))
}

// StandardMatrix is the baseline set of configs Phase 9 validates against.
func StandardMatrix() []HardwareConfig {
	return []HardwareConfig{
		{"NVIDIA RTX 4090 (Ada)", "NVIDIA", "Ada", "cu128"},
		{"NVIDIA RTX 3080 (Ampere)", "NVIDIA", "Ampere", "cu128"},
		{"NVIDIA GTX 1080 (Pascal)", "NVIDIA", "Pascal", "cu121"},
		{"AMD RX 6800 (RDNA2)", "AMD", "RDNA2", "rocm6"},
		{"AMD RX 5700 (RDNA1)", "AMD", "RDNA1", "cpu"},
		{"Intel Arc A770", "Intel", "Alchemist", "cpu"},
	}
}
