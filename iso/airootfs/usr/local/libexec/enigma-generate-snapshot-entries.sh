#!/bin/bash
# Enigma OS: generate systemd-boot entries for snapper snapshots (SPEC §7).
# Phase 2 decision (documented): grub-btrfs handles the BIOS/GRUB path; this
# script handles the UEFI/systemd-boot path so snapshots appear in the boot
# menu. Runs from a pacman hook after snapshots are created.
#
# The entries are derived from the system's own default boot entry rather than
# being written from a template. An earlier template hardcoded the *ISO* paths
# (/arch/boot/x86_64/vmlinuz-linux-cachyos) and root=/dev/mapper/enigma-root,
# so on an installed system every generated entry pointed at a kernel that does
# not exist and a device that usually does not exist either — selecting a
# snapshot from the boot menu could never work. Cloning the real entry keeps
# LUKS/non-LUKS, kernel naming and microcode handling correct automatically.

set -euo pipefail

LOADER_DIR="/boot/loader/entries"
MAX_ENTRIES=10

# Only regenerate if systemd-boot is present.
[ -d "$LOADER_DIR" ] || exit 0
command -v snapper >/dev/null 2>&1 || exit 0

# Find a template: the first entry that is not itself a snapshot entry.
TEMPLATE=""
for f in "$LOADER_DIR"/*.conf; do
    [ -e "$f" ] || continue
    case "${f##*/}" in
        snapshot-*.conf) continue ;;
    esac
    TEMPLATE="$f"
    break
done

# With no real entry to clone there is nothing safe to generate.
[ -n "$TEMPLATE" ] || exit 0

TEMPLATE_LINUX=$(grep -E '^\s*linux\s' "$TEMPLATE" | head -1 || true)
TEMPLATE_INITRD=$(grep -E '^\s*initrd\s' "$TEMPLATE" || true)
TEMPLATE_OPTIONS=$(grep -E '^\s*options\s' "$TEMPLATE" | head -1 | sed -E 's/^\s*options\s+//' || true)

[ -n "$TEMPLATE_LINUX" ] && [ -n "$TEMPLATE_OPTIONS" ] || exit 0

# Remove previously generated entries. A glob cannot be tested with -f, so let
# the loop itself handle the no-match case.
for old in "$LOADER_DIR"/snapshot-*.conf; do
    [ -e "$old" ] && rm -f "$old"
done

SNAPSHOTS=$(snapper -c root list --columns number,date -t single 2>/dev/null \
    | tail -n +3 | sort -rn | head -n "$MAX_ENTRIES" || true)
[ -n "$SNAPSHOTS" ] || exit 0

COUNT=0
while IFS= read -r line; do
    [ -z "$line" ] && continue

    SNAP_NUM=$(printf '%s\n' "$line" | awk '{print $1}')
    SNAP_DATE=$(printf '%s\n' "$line" | awk '{print $2, $3}')

    # Skip header remnants and anything non-numeric.
    case "$SNAP_NUM" in
        ''|*[!0-9]*) continue ;;
    esac
    # Snapshot 0 is snapper's "current" pseudo-entry, not a real snapshot.
    [ "$SNAP_NUM" = "0" ] && continue

    # Point the existing rootflags at this snapshot's subvolume, keeping every
    # other kernel parameter (root=UUID=..., cryptdevice=, etc.) intact.
    SNAP_SUBVOL="@/.snapshots/${SNAP_NUM}/snapshot"
    if printf '%s\n' "$TEMPLATE_OPTIONS" | grep -q 'rootflags=subvol='; then
        SNAP_OPTIONS=$(printf '%s\n' "$TEMPLATE_OPTIONS" \
            | sed -E "s#(rootflags=[^ ]*subvol=)[^ ,]*#\1${SNAP_SUBVOL}#")
    else
        SNAP_OPTIONS="${TEMPLATE_OPTIONS} rootflags=subvol=${SNAP_SUBVOL}"
    fi
    # A snapshot is read-only; booting it read-write would corrupt it.
    SNAP_OPTIONS=$(printf '%s\n' "$SNAP_OPTIONS" | sed -E 's/(^| )rw($| )/\1ro\2/')

    ENTRY_FILE="$LOADER_DIR/snapshot-$SNAP_NUM.conf"
    {
        printf 'title   Enigma OS — Snapshot %s (%s)\n' "$SNAP_NUM" "$SNAP_DATE"
        printf 'sort-key snapshot-%05d\n' "$((MAX_ENTRIES - COUNT))"
        printf '%s\n' "$TEMPLATE_LINUX"
        [ -n "$TEMPLATE_INITRD" ] && printf '%s\n' "$TEMPLATE_INITRD"
        printf 'options %s\n' "$SNAP_OPTIONS"
    } > "$ENTRY_FILE"

    COUNT=$((COUNT + 1))
done <<< "$SNAPSHOTS"

exit 0
