package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"enigma/hwdetect/internal/detect"
	"enigma/hwdetect/internal/hardware"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `enigma-hwdetect — GPU detection + driver decision engine

Usage:
  enigma-hwdetect detect          Write /etc/enigma/hardware.toml
  enigma-hwdetect doctor          Show health report
  enigma-hwdetect -h             Show this help

`)
	}

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	switch args[0] {
	case "detect":
		handleDetect()
	case "doctor":
		handleDoctor()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", args[0])
		os.Exit(1)
	}
}

func handleDetect() {
	detector := detect.NewDetector()
	cfg, err := detector.Detect()
	if err != nil {
		log.Fatalf("Detection failed: %v", err)
	}

	// Write to /etc/enigma/hardware.toml (or ENIGMA_CONFIG_DIR)
	configPath := "/etc/enigma/hardware.toml"
	if env := os.Getenv("ENIGMA_CONFIG_DIR"); env != "" {
		configPath = env + "/hardware.toml"
	}

	if err := cfg.WriteToFile(configPath); err != nil {
		log.Fatalf("Write config: %v", err)
	}

	fmt.Printf("✓ Hardware detection complete: %s\n", configPath)
	fmt.Print(cfg.String())
}

func handleDoctor() {
	// Read existing hardware.toml
	configPath := "/etc/enigma/hardware.toml"
	if env := os.Getenv("ENIGMA_CONFIG_DIR"); env != "" {
		configPath = env + "/hardware.toml"
	}

	cfg, err := hardware.ReadFromFile(configPath)
	if err != nil {
		fmt.Printf("⚠ Hardware config not found: %v\n", err)
		fmt.Println("Run 'enigma-hwdetect detect' first")
		os.Exit(1)
	}

	// Print health report
	fmt.Println("Enigma OS Hardware Health Report")
	fmt.Println("═════════════════════════════════")
	fmt.Print(cfg.String())

	if cfg.System.IsVM {
		fmt.Println("⚠ Running in virtual machine (driver install skipped)")
	}

	if len(cfg.GPU) == 0 {
		fmt.Println("⚠ No GPUs detected")
	}

	if cfg.AI.Verified {
		fmt.Printf("✓ PyTorch verified on %s\n", cfg.AI.TorchTarget)
	} else {
		fmt.Printf("◐ PyTorch not yet verified (target: %s)\n", cfg.AI.TorchTarget)
	}
}
