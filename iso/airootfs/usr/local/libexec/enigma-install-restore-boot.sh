#!/bin/bash
# Restore the kernel images into the installation target's /boot.
#
# WHY THIS EXISTS
# mkarchiso's _cleanup_pacstrap_dir() runs `find "${pacstrap_dir}/boot"
# -mindepth 1 -delete` before building airootfs.sfs: the kernel and initramfs
# are deliberately NOT duplicated inside the squashfs, they live in the ISO
# tree at <install_dir>/boot/<arch>/ instead. Calamares' unpackfs copies that
# squashfs onto the target verbatim, so straight after unpackfs the target's
# /boot is completely EMPTY — no vmlinuz, no initramfs.
#
# Without this step the bootloader module happily writes systemd-boot entries
# pointing at /vmlinuz-linux-cachyos and /initramfs-linux-cachyos.img, the
# install reports success, and the machine then fails to boot because neither
# file exists. Reinstalling the kernel with pacman is not an option here: the
# pacman database copied from the live system already lists linux-cachyos as
# installed, and the installer must work with no network.
#
# Runs UN-chrooted (dontChroot: true) so it can read the live boot medium at
# /run/archiso/bootmnt — the same path unpackfs.conf reads airootfs.sfs from.
#
# Usage: enigma-install-restore-boot.sh <target-root>

set -euo pipefail

TARGET="${1:?target root required}"
MEDIUM="/run/archiso/bootmnt"
BOOTSRC="${MEDIUM}/arch/boot"

if [ ! -d "$TARGET" ]; then
    echo "ERROR: target root '$TARGET' does not exist" >&2
    exit 1
fi

if [ ! -d "$BOOTSRC" ]; then
    echo "ERROR: live boot medium not found at ${BOOTSRC}." >&2
    echo "       Expected the ISO to be mounted at ${MEDIUM}." >&2
    exit 1
fi

mkdir -p "${TARGET}/boot"

# The default kernel (SPEC §1: linux-cachyos with the BORE scheduler). This one
# is mandatory — the bootloader entries reference it by name.
if [ ! -f "${BOOTSRC}/x86_64/vmlinuz-linux-cachyos" ]; then
    echo "ERROR: vmlinuz-linux-cachyos missing from the boot medium." >&2
    exit 1
fi
install -Dm644 "${BOOTSRC}/x86_64/vmlinuz-linux-cachyos" \
    "${TARGET}/boot/vmlinuz-linux-cachyos"
echo "restored: /boot/vmlinuz-linux-cachyos"

# linux-lts is the fallback boot entry (SPEC §1); optional, so never fatal.
if [ -f "${BOOTSRC}/x86_64/vmlinuz-linux-lts" ]; then
    install -Dm644 "${BOOTSRC}/x86_64/vmlinuz-linux-lts" \
        "${TARGET}/boot/vmlinuz-linux-lts"
    echo "restored: /boot/vmlinuz-linux-lts"
fi

# Early microcode images, when the profile ships them. The `microcode` hook in
# the target's mkinitcpio config picks these up.
for ucode in intel-ucode.img amd-ucode.img; do
    if [ -f "${BOOTSRC}/${ucode}" ]; then
        install -Dm644 "${BOOTSRC}/${ucode}" "${TARGET}/boot/${ucode}"
        echo "restored: /boot/${ucode}"
    fi
done

# The initramfs is NOT copied from the medium: the live one is built with the
# archiso hooks and would look for a squashfs that isn't there on a disk boot.
# enigma-install-initramfs.sh generates a proper one in the target chroot.

echo "Kernel images restored into the target. Initramfs is generated later."
