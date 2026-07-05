package detect

// IntelGPUInfo represents Intel discrete GPU info
// Per hwdetect/DESIGN.md §5
type IntelGPUInfo struct {
	Name          string
	SupportsMesa  bool
	SupportsOneAPI bool // Arc GPUs
}

// IntelGPUSupport checks if an Intel GPU supports various stacks
func IntelGPUSupport(name string) (bool, bool) {
	// All Intel iGPUs and dGPUs support Mesa + i965/iris/iHD drivers
	// Arc discrete GPUs (DG1, DG2) also support oneAPI (ipex) — opt-in for v1

	// For v1: Mesa + intel-media-driver is sufficient
	// oneAPI support deferred to v1.5
	return true, false // Mesa: true, oneAPI: false (deferred)
}

// IntelTorchTarget returns torch target for Intel GPUs
// Most Intel GPUs will use CPU torch in v1, oneAPI in v1.5
func IntelTorchTarget(gpu string) string {
	// Intel Arc could use ipex (Intel PyTorch Extension) but that's v1.5
	// For now, CPU is the safe default
	return "cpu"
}
