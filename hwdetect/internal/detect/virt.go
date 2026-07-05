package detect

import (
	"os/exec"
	"strings"
)

// VirtDetector detects if running in a VM
// Interface allows testing with fake implementations
type VirtDetector interface {
	IsVM() (bool, error)
	VirtType() (string, error)
}

// RealVirtDetector calls systemd-detect-virt
type RealVirtDetector struct{}

func (r *RealVirtDetector) IsVM() (bool, error) {
	cmd := exec.Command("systemd-detect-virt")
	output, err := cmd.CombinedOutput()
	result := strings.TrimSpace(string(output))

	if err != nil {
		// If systemd-detect-virt returns exit code 1, we're not in a VM
		if strings.Contains(result, "none") {
			return false, nil
		}
		return false, err
	}

	// If output is not "none", we're in a VM
	return result != "none", nil
}

func (r *RealVirtDetector) VirtType() (string, error) {
	cmd := exec.Command("systemd-detect-virt")
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

// FakeVirtDetector for testing
type FakeVirtDetector struct {
	VM   bool
	Type string
}

func (f *FakeVirtDetector) IsVM() (bool, error) {
	return f.VM, nil
}

func (f *FakeVirtDetector) VirtType() (string, error) {
	return f.Type, nil
}
