#!/bin/bash
# Enigma OS custom pacman hook: generate systemd-boot entries for snapshots
# Phase 2 decision (documented): grub-btrfs handles BIOS/GRUB path;
# this script handles UEFI/systemd-boot path to show snapshots in boot menu
# Runs after every snapper snapshot creation via pacman hook
# shellcheck disable=SC2086

set -e

LOADER_DIR="/boot/loader/entries"
SNAPSHOT_DIR="/.snapshots"
MAX_ENTRIES=10

# Only regenerate if systemd-boot is present
if [ ! -d "$LOADER_DIR" ]; then
    exit 0
fi

# Back up existing snapshot entries
mkdir -p "$LOADER_DIR/.backup"
if [ -f "$LOADER_DIR"/snapshot-*.conf ]; then
    cp "$LOADER_DIR"/snapshot-*.conf "$LOADER_DIR/.backup/" 2>/dev/null || true
fi

# Remove old snapshot entries
rm -f "$LOADER_DIR"/snapshot-*.conf || true

# Get the last N snapshots from snapper
if ! command -v snapper &> /dev/null; then
    exit 0
fi

SNAPSHOTS=$(snapper -c root list --columns number,date,description -t single 2>/dev/null | tail -n +2 | sort -rn | head -n $MAX_ENTRIES)

if [ -z "$SNAPSHOTS" ]; then
    # No snapshots found; don't create entries
    exit 0
fi

COUNT=0
while IFS= read -r line; do
    [ -z "$line" ] && continue

    SNAP_NUM=$(echo "$line" | awk '{print $1}')
    SNAP_DATE=$(echo "$line" | awk '{print $2, $3}')

    # Skip if snap number is invalid
    [ -z "$SNAP_NUM" ] || [ "$SNAP_NUM" = "number" ] && continue

    # Create a boot entry for this snapshot
    ENTRY_FILE="$LOADER_DIR/snapshot-$SNAP_NUM.conf"

    cat > "$ENTRY_FILE" << EOF
title   Enigma OS — Snapshot $SNAP_NUM ($SNAP_DATE)
sort-key snapshot-$(printf "%05d" $((MAX_ENTRIES - COUNT)))
linux   /arch/boot/x86_64/vmlinuz-linux-cachyos
initrd  /arch/boot/x86_64/initramfs-linux-cachyos.img
options root=/dev/mapper/enigma-root rootflags=subvol=@/.snapshots/$SNAP_NUM/snapshot rw
EOF

    COUNT=$((COUNT + 1))
done <<< "$SNAPSHOTS"

exit 0
