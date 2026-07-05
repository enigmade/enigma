package detect

import (
	"fmt"
	"strconv"
)

// NVIDIAGeneration maps device ID ranges to GPU generations
// Per hwdetect/DESIGN.md §3
type NVIDIAGeneration struct {
	Name          string
	DeviceIDMin   int
	DeviceIDMax   int
	Driver        string // nvidia-open-dkms or nvidia-dkms (legacy)
	CUDAVersion   string // cu128, cu121, cu121, etc.
	SupportStatus string // supported, legacy, unsupported
}

// NVIDIAGenerations lists all known NVIDIA GPU generations
var NVIDIAGenerations = []NVIDIAGeneration{
	// Blackwell (RTX 50 series)
	{Name: "Blackwell", DeviceIDMin: 0x2800, DeviceIDMax: 0x28FF, Driver: "nvidia-open-dkms", CUDAVersion: "cu128", SupportStatus: "supported"},
	// Ada (RTX 40 series)
	{Name: "Ada", DeviceIDMin: 0x2600, DeviceIDMax: 0x26FF, Driver: "nvidia-open-dkms", CUDAVersion: "cu128", SupportStatus: "supported"},
	// Ampere (RTX 30 series)
	{Name: "Ampere", DeviceIDMin: 0x2200, DeviceIDMax: 0x25FF, Driver: "nvidia-open-dkms", CUDAVersion: "cu128", SupportStatus: "supported"},
	// Turing (RTX 20 series)
	{Name: "Turing", DeviceIDMin: 0x1C00, DeviceIDMax: 0x1FFF, Driver: "nvidia-open-dkms", CUDAVersion: "cu128", SupportStatus: "supported"},
	// Maxwell (GTX 750, 950, 960, 970, 980, 980 Ti, Titan X (Maxwell))
	{Name: "Maxwell", DeviceIDMin: 0x1180, DeviceIDMax: 0x12FF, Driver: "nvidia-open-dkms", CUDAVersion: "cu128", SupportStatus: "supported"},
	// Pascal (GTX 10 series, Titan X (Pascal))
	{Name: "Pascal", DeviceIDMin: 0x1000, DeviceIDMax: 0x117F, Driver: "nvidia-dkms", CUDAVersion: "cu121", SupportStatus: "legacy"},
	// Older cards fall back to nouveau or CPU-only
}

// DetectNVIDIAGeneration determines the generation and driver for an NVIDIA GPU
// by device ID
func DetectNVIDIAGeneration(deviceIDHex string) (string, string, string, error) {
	// Parse hex device ID
	deviceID, err := strconv.ParseInt(deviceIDHex, 16, 32)
	if err != nil {
		return "", "", "", fmt.Errorf("invalid NVIDIA device ID: %s", deviceIDHex)
	}

	intDeviceID := int(deviceID)

	// Check against known generations (in order, newest first)
	for _, gen := range NVIDIAGenerations {
		if intDeviceID >= gen.DeviceIDMin && intDeviceID <= gen.DeviceIDMax {
			return gen.Name, gen.Driver, gen.CUDAVersion, nil
		}
	}

	// Unknown device: assume newest driver for forward compatibility
	// (DESIGN.md: "UNKNOWN ID newer than table → assume newest path")
	return "unknown", "nvidia-open-dkms", "cu128", nil
}

// NVIDIATorchTarget returns the torch target for an NVIDIA GPU
// based on its generation and CUDA version
func NVIDIATorchTarget(generation string) string {
	switch generation {
	case "Turing", "Ampere", "Ada", "Blackwell", "Maxwell":
		return "cu128"
	case "Pascal":
		return "cu121"
	default:
		return "cpu"
	}
}
