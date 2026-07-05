#!/bin/bash
# Enigma OS QEMU BIOS boot test
# Verifies ISO boots to graphical/multi-user target via SeaBIOS

set -e

source "$(dirname "$0")/common.sh"

ISO="${1:?ISO path required}"
LOGDIR="${2:-.}"

[ -f "$ISO" ] || { echo "ERROR: ISO not found: $ISO"; exit 1; }
check_file_exists "$SEABIOS" || exit 1

SERIAL_LOG=$(create_serial_log "$LOGDIR")
PID_FILE=$(mktemp)

echo "Testing BIOS boot from: $ISO"
echo "Serial log: $SERIAL_LOG"

trap "cleanup_qemu '$PID_FILE'" EXIT

# Boot with SeaBIOS firmware
"$QEMU_BIN" \
    -m 4096 \
    -cpu host \
    -enable-kvm \
    -bios "$SEABIOS" \
    -drive file="$ISO",format=raw,if=ide,media=cdrom \
    -serial file:"$SERIAL_LOG" \
    -display none \
    -nographic \
    -daemonize \
    -pidfile "$PID_FILE" \
    -boot d

# Wait for boot to complete
if wait_for_boot "$SERIAL_LOG" "graphical" "$BOOT_TIMEOUT"; then
    echo "✓ BIOS boot test PASSED"
    exit 0
else
    echo "✗ BIOS boot test FAILED"
    echo "=== SERIAL LOG ===="
    tail -50 "$SERIAL_LOG"
    exit 1
fi
