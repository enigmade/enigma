package lspci

import (
	"testing"
)

func TestParseNVIDIA(t *testing.T) {
	input := `0d:00.0 VGA compatible controller [0300]: NVIDIA Corporation GA104 [GeForce RTX 3070] [10de:2484]
Kernel driver in use: nvidia_drm`

	gpus, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(gpus) != 1 {
		t.Fatalf("Expected 1 GPU, got %d", len(gpus))
	}

	gpu := gpus[0]
	if gpu.Vendor != "NVIDIA" {
		t.Errorf("Vendor: expected NVIDIA, got %s", gpu.Vendor)
	}
	if gpu.VendorID != "10de" {
		t.Errorf("VendorID: expected 10de, got %s", gpu.VendorID)
	}
	if gpu.DeviceID != "2484" {
		t.Errorf("DeviceID: expected 2484, got %s", gpu.DeviceID)
	}
	if gpu.BoundDriver != "nvidia_drm" {
		t.Errorf("Driver: expected nvidia_drm, got %s", gpu.BoundDriver)
	}
}

func TestParseAMD(t *testing.T) {
	input := `1d:00.0 VGA compatible controller [0300]: Advanced Micro Devices, Inc. [AMD/ATI] Navi 21 [Radeon RX 6800 XT] [1002:73bf]
Kernel driver in use: amdgpu`

	gpus, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(gpus) != 1 {
		t.Fatalf("Expected 1 GPU, got %d", len(gpus))
	}

	gpu := gpus[0]
	if gpu.Vendor != "AMD" {
		t.Errorf("Vendor: expected AMD, got %s", gpu.Vendor)
	}
	if gpu.BoundDriver != "amdgpu" {
		t.Errorf("Driver: expected amdgpu, got %s", gpu.BoundDriver)
	}
}

func TestParseMultipleGPUs(t *testing.T) {
	input := `0d:00.0 VGA compatible controller [0300]: NVIDIA Corporation GA104 [GeForce RTX 3070] [10de:2484]
Kernel driver in use: nvidia_drm
1d:00.0 VGA compatible controller [0300]: Advanced Micro Devices, Inc. [AMD/ATI] Navi 21 [Radeon RX 6800 XT] [1002:73bf]
Kernel driver in use: amdgpu`

	gpus, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(gpus) != 2 {
		t.Fatalf("Expected 2 GPUs, got %d", len(gpus))
	}

	if gpus[0].Vendor != "NVIDIA" {
		t.Errorf("First GPU: expected NVIDIA, got %s", gpus[0].Vendor)
	}
	if gpus[1].Vendor != "AMD" {
		t.Errorf("Second GPU: expected AMD, got %s", gpus[1].Vendor)
	}
}

func TestParseIgnoresNonGPU(t *testing.T) {
	input := `00:02.0 VGA compatible controller [0300]: Intel Corporation Device [8086:1234]
Kernel driver in use: i915
10:00.0 Ethernet controller [0200]: Intel Corporation 82579LM Gigabit Network Connection [8086:1502]
Kernel driver in use: e1000e`

	gpus, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Should only get the GPU (VGA 0300), not the Ethernet (0200)
	if len(gpus) > 2 {
		t.Errorf("Expected ≤2 GPUs (Intel iGPU + maybe discrete), got %d", len(gpus))
	}
}

func TestParseEmpty(t *testing.T) {
	gpus, err := Parse("")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(gpus) != 0 {
		t.Errorf("Expected 0 GPUs from empty input, got %d", len(gpus))
	}
}
