package lspci

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

// Device represents a GPU device from lspci output
type Device struct {
	Slot        string // e.g., "0d:00.0"
	VendorID    string // e.g., "10de" (NVIDIA)
	DeviceID    string // e.g., "2786"
	Vendor      string // Human name: NVIDIA, AMD, Intel
	Name        string // Device name from lspci
	BoundDriver string // Currently bound kernel driver (e.g., "nouveau", "nvidia_drm")
}

// GPUVendorID maps PCI vendor IDs to human names
var GPUVendorID = map[string]string{
	"10de": "NVIDIA",
	"1002": "AMD",
	"8086": "Intel",
}

// IsGPU checks if a device class is a GPU (0300=VGA, 0302=3D controller)
func IsGPU(class string) bool {
	return strings.HasPrefix(class, "0300") || strings.HasPrefix(class, "0302")
}

// Parse parses lspci -nnk output and returns a list of GPU devices
// Expected format (per hwdetect/DESIGN.md §1):
//
//	0d:00.0 VGA compatible controller [0300]: NVIDIA Corporation GA104 [GeForce RTX 3070] [10de:2484]
//	Kernel driver in use: nvidia_drm
//
func Parse(input string) ([]Device, error) {
	var gpus []Device
	scanner := bufio.NewScanner(strings.NewReader(input))

	// Regex patterns
	driverRe := regexp.MustCompile(`Kernel driver in use: (\S+)`)
	// Matches: slot type-and-class vendor:device at end
	vendorDeviceRe := regexp.MustCompile(`\[([0-9a-f]{4}):([0-9a-f]{4})\]\s*$`)

	for scanner.Scan() {
		line := scanner.Text()

		// Extract slot (first token)
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		slot := parts[0]

		// Find class code in brackets (second bracket in the line)
		classIdx := strings.Index(line, "[")
		if classIdx == -1 {
			continue
		}
		classEndIdx := strings.Index(line[classIdx:], "]")
		if classEndIdx == -1 {
			continue
		}
		classCode := line[classIdx+1 : classIdx+classEndIdx]

		// Only process GPU devices
		if !IsGPU(classCode) {
			continue
		}

		// Extract vendor:device ID (last [...] pair)
		vendorDeviceMatch := vendorDeviceRe.FindStringSubmatch(line)
		if len(vendorDeviceMatch) < 3 {
			continue
		}
		vendorID := vendorDeviceMatch[1]
		deviceID := vendorDeviceMatch[2]

		// Check if vendor is known
		vendor, ok := GPUVendorID[vendorID]
		if !ok {
			continue
		}

		// Extract device name (text between class bracket and vendor:device bracket)
		colonIdx := strings.Index(line, "]:")
		if colonIdx == -1 {
			continue
		}
		nameStart := colonIdx + 2
		nameEnd := strings.LastIndex(line, "[")
		name := ""
		if nameEnd > nameStart {
			name = strings.TrimSpace(line[nameStart:nameEnd])
		}

		device := Device{
			Slot:     slot,
			VendorID: vendorID,
			DeviceID: deviceID,
			Vendor:   vendor,
			Name:     name,
		}

		// Read next line for driver info
		if scanner.Scan() {
			nextLine := scanner.Text()
			if driverMatches := driverRe.FindStringSubmatch(nextLine); len(driverMatches) > 0 {
				device.BoundDriver = driverMatches[1]
			}
		}

		gpus = append(gpus, device)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("parse lspci: %w", err)
	}

	return gpus, nil
}
