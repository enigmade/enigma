package detect

import (
	"fmt"
	"os/exec"

	"enigma/hwdetect/pkg/hardware"
	"enigma/hwdetect/internal/lspci"
)

// Detector orchestrates the full detection pipeline per hwdetect/DESIGN.md §DETECTION PIPELINE
type Detector struct {
	VirtDetector VirtDetector
}

// NewDetector creates a new Detector with real systemd-detect-virt
func NewDetector() *Detector {
	return &Detector{
		VirtDetector: &RealVirtDetector{},
	}
}

// Detect runs the full hardware detection pipeline
func (d *Detector) Detect() (*hardware.Config, error) {
	cfg := &hardware.Config{}

	// Step 0: Get system info (placeholder for v1; real parsing deferred to Phase 3.5)
	cfg.System = hardware.SystemInfo{
		RAMGb:     8, // Placeholder
		CPUVendor: "Intel",
		CPUFlags:  []string{"v3", "v4"},
		IsLaptop:  false,
	}

	// Step 1: VM Check FIRST (per DESIGN.md step 2)
	isVM, err := d.VirtDetector.IsVM()
	if err != nil {
		return nil, fmt.Errorf("VM detection: %w", err)
	}
	cfg.System.IsVM = isVM

	// Step 2: Parse lspci
	cmd := exec.Command("lspci", "-nnk")
	lspciOutput, err := cmd.CombinedOutput()
	if err != nil {
		// If lspci not available (common on non-Linux), return empty GPU list
		// This is expected on macOS, handled gracefully
		cfg.GPU = []hardware.GPUInfo{}
		cfg.AI.TorchTarget = "cpu"
		// Still respect VM session choice
		if isVM {
			cfg.Display.Session = "x11"
		} else {
			cfg.Display.Session = "wayland"
		}
		cfg.Display.HybridMode = "none"
		return cfg, nil
	}

	gpus, err := lspci.Parse(string(lspciOutput))
	if err != nil {
		return nil, fmt.Errorf("lspci parse: %w", err)
	}

	// Step 3: Process each GPU
	var detectedGPUs []hardware.GPUInfo

	for i, gpu := range gpus {
		hwGPU := hardware.GPUInfo{
			Index:         i,
			Vendor:        gpu.Vendor,
			PCIId:         gpu.DeviceID,
			Name:          gpu.Name,
			CurrentDriver: gpu.BoundDriver,
			Role:          "primary",
		}

		if isVM {
			// In VMs, skip dkms driver installation
			hwGPU.DriverState = "none"
			hwGPU.Role = "compute"
		} else {
			switch gpu.Vendor {
			case "NVIDIA":
				gen, _, _, err := DetectNVIDIAGeneration(gpu.DeviceID)
				if err != nil {
					gen = "unknown"
				}
				hwGPU.Generation = gen
				hwGPU.DriverState = "none" // Will be set to "installed" after Phase 3's driver install
				// (placeholder for Phase 3.5 integration with installer)

			case "AMD":
				hwGPU.Generation = "RDNA2" // Placeholder
				hwGPU.DriverState = "none"
				// ROCm support determined during driver install phase

			case "Intel":
				hwGPU.Generation = "Iris"
				hwGPU.DriverState = "none"
			}
		}

		detectedGPUs = append(detectedGPUs, hwGPU)
	}

	cfg.GPU = detectedGPUs

	// Step 4: Set AI defaults
	if len(detectedGPUs) == 0 {
		cfg.AI.TorchTarget = "cpu"
	} else {
		// Use first GPU's target as default (can be overridden later)
		gpuVendor := detectedGPUs[0].Vendor
		switch gpuVendor {
		case "NVIDIA":
			cfg.AI.TorchTarget = "cu128" // Default to newest
		case "AMD":
			cfg.AI.TorchTarget = "rocm6"
		default:
			cfg.AI.TorchTarget = "cpu"
		}
	}

	cfg.AI.Tier = "standard"
	cfg.AI.Verified = false

	// Step 5: Display config
	if isVM {
		cfg.Display.Session = "x11" // VMs often default to x11
	} else {
		cfg.Display.Session = "wayland" // Enigma primary (SPEC §3)
	}
	cfg.Display.HybridMode = "none"

	return cfg, nil
}

// DetectWithVirtFake allows tests to inject a fake VirtDetector
func DetectWithVirt(virt VirtDetector) (*hardware.Config, error) {
	d := &Detector{VirtDetector: virt}
	return d.Detect()
}
