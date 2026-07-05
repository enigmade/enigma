package detect

// AMDArchs lists AMD GPU architectures that support ROCm
// Per hwdetect/DESIGN.md §4
var AMDArchs = struct {
	RDNA2Official []string
	CDNAOfficial  []string
	RDNA1Comment  []string
	Unsupported   []string
}{
	RDNA2Official: []string{"RDNA2", "Navi 21", "Navi 22", "Navi 23"},
	CDNAOfficial:  []string{"CDNA", "CDNA2", "MI100", "MI200"},
	RDNA1Comment:  []string{"RDNA", "RDNA1", "Navi 10", "Navi 14"}, // Community support
	Unsupported:   []string{"Vega", "Polaris", "Fury"},
}

// AMDTorchTarget returns the torch target for an AMD GPU
// Based on ROCm support per DESIGN.md §4
func AMDTorchTarget(supported bool) string {
	if supported {
		return "rocm6"
	}
	return "cpu"
}

// AMDSupportsROCm checks if an AMD GPU supports ROCm
// based on its architecture name
func AMDSupportsROCm(arch string) bool {
	for _, name := range AMDArchs.RDNA2Official {
		if arch == name {
			return true
		}
	}
	for _, name := range AMDArchs.CDNAOfficial {
		if arch == name {
			return true
		}
	}
	// RDNA1 is community-supported but listed as such in notes
	// Returning false here; caller can check and add note
	return false
}
