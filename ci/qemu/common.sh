#!/bin/bash
# Enigma OS QEMU CI helper functions
# shellcheck disable=SC2034

set -e

# OVMF (UEFI) firmware paths. The edk2-ovmf package layout has changed
# names across Arch releases (OVMF_CODE.fd vs OVMF_CODE.4m.fd etc), so
# discover the actual files on disk rather than hardcoding one.
find_firmware() {
    local pattern="$1"
    find /usr/share/edk2 /usr/share/OVMF /usr/share/qemu -iname "$pattern" 2>/dev/null | sort | head -1
}
OVMF_CODE="$(find_firmware 'OVMF_CODE*.fd')"
OVMF_VARS="$(find_firmware 'OVMF_VARS*.fd')"
SEABIOS="$(find_firmware 'bios-256k.bin')"
[ -z "$SEABIOS" ] && SEABIOS="$(find_firmware 'bios.bin')"

# QEMU binary
QEMU_BIN="qemu-system-x86_64"

# Timeout for boot tests (seconds)
BOOT_TIMEOUT=120

# Fail fast with a clear message if /dev/kvm isn't usable, rather than
# letting every boot test silently burn its full timeout under TCG.
require_kvm() {
    if [ ! -e /dev/kvm ]; then
        echo "ERROR: /dev/kvm not present — hardware virtualization is not"
        echo "available to this runner/container. QEMU boot tests need KVM"
        echo "acceleration (unaccelerated TCG boot of a full desktop session"
        echo "would take far longer than any reasonable CI timeout)."
        return 1
    fi
    if [ ! -r /dev/kvm ] || [ ! -w /dev/kvm ]; then
        echo "ERROR: /dev/kvm exists but is not read/write accessible to this user."
        ls -la /dev/kvm
        return 1
    fi
    return 0
}

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
    if [ -z "$file" ] || [ ! -f "$file" ]; then
        echo "ERROR: File not found: '$file'"
        echo "  Searched: /usr/share/edk2, /usr/share/OVMF, /usr/share/qemu"
        find /usr/share/edk2 /usr/share/OVMF /usr/share/qemu -type f 2>/dev/null | head -20
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
