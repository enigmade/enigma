package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"enigma/cli/internal/ai"
	"enigma/cli/internal/dev"
	"enigma/cli/internal/doctor"
	"enigma/cli/internal/game"
	"enigma/cli/internal/portability"
	"enigma/cli/internal/ports"
	"enigma/cli/internal/validation"
)

const hardwarePath = "/etc/enigma/hardware.toml"

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	switch args[0] {
	case "dev":
		cmdDev(args[1:])
	case "doctor":
		doctor.Report()
	case "ports":
		if len(args) < 2 {
			ports.List()
		} else {
			fmt.Println("Kill port command coming in Phase 4.5")
		}
	case "clean":
		fmt.Println("✓ Clean command stubbed — nothing to clean yet")
	case "ai":
		cmdAI(args[1:])
	case "game":
		cmdGame(args[1:])
	case "win":
		cmdWin(args[1:])
	case "containers":
		cmdContainers(args[1:])
	case "usb":
		cmdUSB(args[1:])
	case "validate":
		cmdValidate(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", args[0])
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), `enigma — zero-bloat dev OS CLI

Usage:
  enigma dev up|down|services   Project environment
  enigma doctor                 Health report
  enigma ports [<port>]         List / kill ports
  enigma ai setup|ollama|comfy  AI subsystem
  enigma game setup             Steam + Proton-GE
  enigma win setup              Wine Tier 1 (DirectX/VC++ redists)
  enigma containers setup       Rootless Podman
  enigma usb plan <dev> <MiB>   Portable USB partition plan (dry run)
  enigma validate <file>        Check benchmark results vs SPEC §9

`)
}

func cmdDev(args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: enigma dev [up|down|services]\n")
		os.Exit(1)
	}
	switch args[0] {
	case "up":
		dev.Up()
	case "down":
		dev.Down()
	case "services":
		dev.Services()
	default:
		fmt.Fprintf(os.Stderr, "Unknown dev command: %s\n", args[0])
		os.Exit(1)
	}
}

func cmdAI(args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: enigma ai [setup|ollama|comfy]\n")
		os.Exit(1)
	}
	switch args[0] {
	case "setup":
		if err := ai.Setup(hardwarePath); err != nil {
			fatal(err)
		}
	case "ollama":
		sm, err := ai.NewServiceManager()
		if err != nil {
			fatal(err)
		}
		out, _ := sm.StatusService("enigma-ollama")
		fmt.Print(out)
	case "comfy":
		fmt.Println("enigma ai comfy: manage ComfyUI service (see internal/ai/comfy.go)")
	default:
		fmt.Fprintf(os.Stderr, "Unknown ai command: %s\n", args[0])
		os.Exit(1)
	}
}

func cmdGame(args []string) {
	if len(args) < 1 || args[0] != "setup" {
		fmt.Println("Usage: enigma game setup")
		return
	}
	home, _ := os.UserHomeDir()
	steamRoot := filepath.Join(home, ".steam", "steam")
	release, dir, err := game.SetupSteam(steamRoot)
	if err != nil {
		fatal(err)
	}
	fmt.Printf("✓ Latest Proton-GE: %s\n  install dir: %s\n", release.Version, dir)
}

func cmdWin(args []string) {
	if len(args) < 1 || args[0] != "setup" {
		fmt.Println("Usage: enigma win setup")
		return
	}
	home, _ := os.UserHomeDir()
	logPath := filepath.Join(home, ".local", "share", "enigma", "wine", "pfx", "winetricks.log")
	content, _ := os.ReadFile(logPath) // absent log => all redists reported missing
	missing := game.SetupWineTier1(string(content))
	if len(missing) > 0 {
		fmt.Printf("Run: winetricks -q %s\n", joinSpace(missing))
	}
}

func cmdContainers(args []string) {
	if len(args) < 1 || args[0] != "setup" {
		fmt.Println("Usage: enigma containers setup")
		return
	}
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".config")
	username := currentUsername()
	path, err := game.SetupContainers(configDir, "/etc/subuid", "/etc/subgid", username)
	if err != nil {
		fmt.Fprintf(os.Stderr, "containers setup incomplete: %v\n", err)
		if path != "" {
			fmt.Printf("(wrote %s)\n", path)
		}
		os.Exit(1)
	}
	fmt.Printf("✓ Rootless Podman configured: %s\n", path)
}

func cmdUSB(args []string) {
	if len(args) < 3 || args[0] != "plan" {
		fmt.Println("Usage: enigma usb plan <device> <sizeMiB>")
		return
	}
	device := args[1]
	var sizeMiB int
	if _, err := fmt.Sscanf(args[2], "%d", &sizeMiB); err != nil {
		fatal(fmt.Errorf("invalid size %q: %w", args[2], err))
	}
	if _, err := portability.PreparePortableUSB(device, sizeMiB); err != nil {
		fatal(err)
	}
}

func cmdValidate(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: enigma validate <benchmark-file>")
		return
	}
	content, err := os.ReadFile(args[0])
	if err != nil {
		fatal(err)
	}
	if err := validation.ValidateHardwareMatrix(string(content)); err != nil {
		fmt.Fprintf(os.Stderr, "✗ %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ All benchmarks within SPEC §9 thresholds")
}

func currentUsername() string {
	if u, err := user.Current(); err == nil {
		return u.Username
	}
	return os.Getenv("USER")
}

func joinSpace(items []string) string {
	out := ""
	for i, it := range items {
		if i > 0 {
			out += " "
		}
		out += it
	}
	return out
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
