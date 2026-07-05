#!/bin/bash
# Enigma OS first-boot wizard
# Runs after installation; collects GPU choice + profile selection
# Per Phase 2 scope: UI stub + config writer only (no CLI calls yet — Phase 4+)
# shellcheck disable=SC2034

set -e

PROFILE_DIR="/etc/enigma"
PROFILE_FILE="$PROFILE_DIR/profile.toml"

mkdir -p "$PROFILE_DIR"

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║           Welcome to Enigma OS First-Boot Setup                ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

# GPU Detection Confirmation
echo "GPU Detection:"
echo "  Enigma will auto-detect your GPU (NVIDIA/AMD/Intel) on first boot."
echo "  GPU drivers will be installed automatically in Phase 3."
echo "  ✓ Auto-detection enabled"
echo ""

# Profile Selection (UI stub for Phase 2)
echo "Workload Profiles:"
echo "  Select which tools you plan to use. You can change these later."
echo "  (In Phase 4, selecting a profile will run the matching enigma commands.)"
echo ""
echo "  1) Minimal (no extra tools)"
echo "  2) Web Development (PHP, Node, Python)"
echo "  3) AI & LLMs (Ollama, ComfyUI, torch)"
echo "  4) Gaming (Steam, Proton-GE, gamemode)"
echo "  5) Windows Apps (Bottles + Wine)"
echo "  6) All of the above"
echo ""

# Simple menu-driven selection (stub for now)
# In a real Phase 2, this would use zenity/kdialog for GUI
read -rp "Enter your choice (1-6) [default: 1]: " profile_choice
profile_choice="${profile_choice:-1}"

case "$profile_choice" in
    1) profiles=("minimal") ;;
    2) profiles=("webdev") ;;
    3) profiles=("ai") ;;
    4) profiles=("gaming") ;;
    5) profiles=("windows") ;;
    6) profiles=("webdev" "ai" "gaming" "windows") ;;
    *) profiles=("minimal"); echo "Invalid choice; defaulting to Minimal." ;;
esac

echo ""
echo "Selected profiles: ${profiles[*]}"
echo ""

# Write profile config (TOML format, consumable by Phase 4+ CLI)
{
    echo "[system]"
    echo "# GPU auto-detection runs on every boot (Phase 3)"
    echo "gpu_detected = true"
    echo ""
    echo "[profiles]"
    for profile in "${profiles[@]}"; do
        echo "enabled = [\"$profile\"]"
    done
    echo ""
    echo "# First-boot completed"
    echo "completed = true"
} > "$PROFILE_FILE"

echo "Configuration saved to $PROFILE_FILE"
echo ""
echo "✓ First-boot setup complete!"
echo ""
echo "Next steps:"
echo "  1. Reboot your system (GPU drivers will be installed automatically)"
echo "  2. In Phase 4+, run 'enigma dev up' in a project folder to start coding"
echo "  3. See README.md for more details"
echo ""
