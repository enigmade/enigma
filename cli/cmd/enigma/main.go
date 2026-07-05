package main

import (
	"flag"
	"fmt"
	"os"

	"enigma/cli/internal/dev"
	"enigma/cli/internal/doctor"
	"enigma/cli/internal/ports"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `enigma — zero-bloat dev OS CLI

Usage:
  enigma dev up            Start project environment
  enigma dev down          Stop project services
  enigma dev services      List services + ports
  enigma doctor            Health report
  enigma ports <port>      Kill process on port
  enigma clean             Cleanup

Phases:
  ✓ Phase 4: dev + doctor (this phase)
  ⏳ Phase 5: ai setup/ollama/comfy
  ⏳ Phase 7: game + win + containers

`)
	}

	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	switch args[0] {
	case "dev":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Usage: enigma dev [up|down|services]\n")
			os.Exit(1)
		}
		switch args[1] {
		case "up":
			dev.Up()
		case "down":
			dev.Down()
		case "services":
			dev.Services()
		default:
			fmt.Fprintf(os.Stderr, "Unknown dev command: %s\n", args[1])
			os.Exit(1)
		}

	case "doctor":
		doctor.Report()

	case "ports":
		if len(args) < 2 {
			ports.List()
		} else {
			// `enigma ports <port>` would kill the process (stubbed)
			fmt.Println("Kill port command coming in Phase 4.5")
		}

	case "clean":
		fmt.Println("✓ Clean command stubbed — nothing to clean yet")

	case "ai":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Usage: enigma ai [setup|ollama|comfy|run]\n")
			os.Exit(1)
		}
		switch args[1] {
		case "setup":
			fmt.Println("enigma ai setup: Install matching PyTorch for detected GPU")
			fmt.Println("  TODO Phase 5.5: wire full implementation")
		case "ollama":
			fmt.Println("enigma ai ollama: Manage Ollama LLM service")
		case "comfy":
			fmt.Println("enigma ai comfy: Manage ComfyUI image generation")
		case "run":
			fmt.Println("enigma ai run: Run Pinokio JSON scripts")
		default:
			fmt.Fprintf(os.Stderr, "Unknown ai command: %s\n", args[1])
			os.Exit(1)
		}

	case "game", "win", "containers":
		fmt.Printf("⏳ `enigma %s` coming in Phase 7/8\n", args[0])

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", args[0])
		os.Exit(1)
	}
}
