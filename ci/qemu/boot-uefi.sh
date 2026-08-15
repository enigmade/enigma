#!/bin/bash
# Enigma OS QEMU UEFI boot test
# Verifies ISO boots to graphical/multi-user target via OVMF UEFI firmware

set -e

source "$(dirname "$0")/common.sh"

ISO="${1:?ISO path required}"
LOGDIR="${2:-.}"

[ -f "$ISO" ] || { echo "ERROR: ISO not found: $ISO"; exit 1; }
require_kvm || exit 1
check_file_exists "$OVMF_CODE" || exit 1

SERIAL_LOG=$(create_serial_log "$LOGDIR")
PID_FILE=$(mktemp)

echo "Testing UEFI boot from: $ISO"
echo "Serial log: $SERIAL_LOG"

# Create temporary OVMF vars (copy to avoid modifying system copy)
OVMF_VARS_TMP=$(mktemp)
cp "$OVMF_VARS" "$OVMF_VARS_TMP"

trap 'cleanup_qemu "$PID_FILE"; rm -f "$OVMF_VARS_TMP"' EXIT

# Boot with UEFI firmware
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

# Wait for boot to complete
if wait_for_boot "$SERIAL_LOG" "graphical" "$BOOT_TIMEOUT"; then
    echo "✓ UEFI boot test PASSED"
    exit 0
else
    echo "✗ UEFI boot test FAILED"
    echo "--- full serial log ($(wc -c < "$SERIAL_LOG" 2>/dev/null || echo 0) bytes) ---"
    cat "$SERIAL_LOG"
    exit 1
fi
