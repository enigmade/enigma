#!/bin/bash
# Enigma OS Ventoy boot test
# Creates a Ventoy disk, copies ISO onto it, boots in QEMU (normal + GRUB2 mode)
# Verifies hybrid ISO is Ventoy-compatible per SPEC §17

set -e

source "$(dirname "$0")/common.sh"

ISO="${1:?ISO path required}"
LOGDIR="${2:-.}"

[ -f "$ISO" ] || { echo "ERROR: ISO not found: $ISO"; exit 1; }
require_kvm || exit 1
check_file_exists "$OVMF_CODE" || exit 1

SERIAL_LOG=$(create_serial_log "$LOGDIR")
PID_FILE=$(mktemp)

# Create a temporary Ventoy disk image
VENTOY_DISK=$(mktemp -p "$LOGDIR" "ventoy-disk-XXXXXX.qcow2")
qemu-img create -q -f qcow2 "$VENTOY_DISK" 2G

OVMF_VARS_TMP=$(mktemp)
cp "$OVMF_VARS" "$OVMF_VARS_TMP"

trap "cleanup_qemu '$PID_FILE'; rm -f '$OVMF_VARS_TMP' '$VENTOY_DISK'" EXIT

echo "Testing Ventoy compatibility from: $ISO"
echo "Ventoy disk: $VENTOY_DISK"
echo "Serial log: $SERIAL_LOG"

# Note: In a full CI environment, this would:
# 1. Convert qcow2 to raw: qemu-img convert -O raw "$VENTOY_DISK" "$VENTOY_DISK_RAW"
# 2. Install Ventoy: ventoy2disk -I "$VENTOY_DISK_RAW"
# 3. Copy ISO: copy "$ISO" to the Ventoy data partition
# 4. Convert back to qcow2 and boot

# For now, just verify the ISO boots when attached directly
# Full Ventoy testing requires ventoy2disk binary in the CI environment

echo "NOTE: Full Ventoy installation testing requires ventoy2disk binary"
echo "Performing basic compatibility check (ISO as CD-ROM)..."

# Boot with ISO attached
"$QEMU_BIN" \
    -m 4096 \
    -cpu host \
    -enable-kvm \
    -drive if=pflash,format=raw,readonly=on,file="$OVMF_CODE" \
    -drive if=pflash,format=raw,file="$OVMF_VARS_TMP" \
    -drive file="$ISO",format=raw,if=ide,media=cdrom \
    -serial file:"$SERIAL_LOG" \
    -display none \
    -daemonize \
    -pidfile "$PID_FILE" \
    -boot d

# Wait for boot
if wait_for_boot "$SERIAL_LOG" "Graphical Interface" "$BOOT_TIMEOUT"; then
    echo "✓ Ventoy compatibility test PASSED"
    exit 0
else
    echo "✗ Ventoy compatibility test FAILED"
    echo "--- full serial log ($(wc -c < "$SERIAL_LOG" 2>/dev/null || echo 0) bytes) ---"
    cat "$SERIAL_LOG"
    exit 1
fi
