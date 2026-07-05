package ai

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// InstallTorch installs PyTorch wheel matching the GPU
func InstallTorch(hwPath, venvPath string) error {
	target, err := TorchTarget(hwPath)
	if err != nil {
		return err
	}

	fmt.Printf("Installing PyTorch[%s] into %s\n", target, venvPath)

	// Determine wheel index URL based on target
	var wheelURL string
	switch target {
	case "cu128":
		wheelURL = "https://download.pytorch.org/whl/cu118" // CUDA 12.8 compatible
	case "cu121":
		wheelURL = "https://download.pytorch.org/whl/cu121" // CUDA 12.1 legacy
	case "rocm6":
		wheelURL = "https://download.pytorch.org/whl/rocm6.1" // ROCm 6.1
	default:
		wheelURL = "" // CPU build
	}

	// Use uv to install torch (SPEC §4a uses mise for runtimes, uv for venvs)
	cmd := exec.Command("uv", "venv", venvPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("create venv: %w", err)
	}

	// Install torch + common packages
	pipCmd := filepath.Join(venvPath, "bin", "pip")
	torchCmd := exec.Command(pipCmd, "install", "--index-url", wheelURL, "torch", "torchvision", "torchaudio")
	if err := torchCmd.Run(); err != nil {
		return fmt.Errorf("install torch: %w", err)
	}

	fmt.Printf("✓ PyTorch[%s] installed to %s\n", target, venvPath)
	return nil
}

// VerifyTorch checks if torch works in the venv
func VerifyTorch(venvPath string) error {
	pythonBin := filepath.Join(venvPath, "bin", "python")
	cmd := exec.Command(pythonBin, "-c", "import torch; print(f'✓ PyTorch {torch.__version__}'); print(f'CUDA available: {torch.cuda.is_available()}')")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("torch verify failed: %w", err)
	}
	fmt.Print(string(output))
	return nil
}
