/* Enigma OS Calamares slideshow (minimal placeholder) */
import QtQuick 2.15
import "branding"

Presentation {
    id: presentation

    TextSlide {
        title: "Welcome to Enigma OS"
        description: "Zero-bloat, zero-telemetry Linux for developers, AI builders, and gamers."
    }

    TextSlide {
        title: "GPU Auto-Detection"
        description: "NVIDIA, AMD, and Intel GPUs are detected and configured automatically."
    }

    TextSlide {
        title: "Web & AI Ready"
        description: "PHP, Node, Python per-project. Ollama and ComfyUI pre-wired for GPU acceleration."
    }

    TextSlide {
        title: "Snapshots & Rollback"
        description: "Btrfs snapshots are created before every update. Boot a previous snapshot to roll back."
    }

    TextSlide {
        title: "Installation Complete"
        description: "Your Enigma OS system is ready. Reboot to begin."
    }
}
