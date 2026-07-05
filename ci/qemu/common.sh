#!/bin/bash
# Enigma OS QEMU CI helper functions
# shellcheck disable=SC2034

set -e

# OVMF (UEFI) firmware paths
OVMF_CODE="/usr/share/edk2/x64/OVMF_CODE.fd"
OVMF_VARS="/usr/share/edk2/x64/OVMF_VARS.fd"
SEABIOS="/usr/share/qemu/bios-256k.bin"

# QEMU binary
QEMU_BIN="qemu-system-x86_64"

# Timeout for boot tests (seconds)
BOOT_TIMEOUT=120

# Wait for a condition with timeout
# Usage: wait_for_condition "command" "timeout_seconds" "description"
wait_for_condition() {
    local cmd="$1"
    local timeout="$2"
    local desc="$3"
    local start
    start=$(date +%s)

    echo "Waiting for: $desc"
    while ! eval "$cmd" 2>/dev/null; do
        local now
        now=$(date +%s)
        local elapsed=$((now - start))
        if [ $elapsed -gt "$timeout" ]; then
            echo "TIMEOUT: $desc after $timeout seconds"
            return 1
        fi
        sleep 1
    done
    echo "SUCCESS: $desc"
    return 0
}

# Wait for QEMU to reach a certain systemd target via serial output
# Usage: wait_for_boot "logfile" "target" "timeout"
wait_for_boot() {
    local logfile="$1"
    local target="$2"
    local timeout="${3:-$BOOT_TIMEOUT}"

    wait_for_condition "grep -q 'Reached target $target' '$logfile'" "$timeout" \
        "boot to $target"
}

# Check if a file exists (for disk write verification)
check_file_exists() {
    local file="$1"
    if [ ! -f "$file" ]; then
        echo "ERROR: File not found: $file"
        return 1
    fi
    return 0
}

# Create a temporary QEMU serial log file
create_serial_log() {
    local logdir="${1:-.}"
    local logfile
    logfile=$(mktemp -p "$logdir" "qemu-serial-XXXXXX.log")
    echo "$logfile"
}

# Cleanup QEMU temporary files
cleanup_qemu() {
    local pidfile="$1"
    if [ -f "$pidfile" ]; then
        local pid
        pid=$(cat "$pidfile")
        if kill -0 "$pid" 2>/dev/null; then
            kill "$pid" || true
            sleep 2
            kill -9 "$pid" 2>/dev/null || true
        fi
    fi
}
