#!/bin/bash
# Generate the installed system's initramfs. Runs INSIDE the target chroot
# (dontChroot: false), after the kernel images have been restored to /boot and
# after the live-only packages have been removed.
#
# WHY THIS EXISTS
# unpackfs copies the live root verbatim, which includes
# /etc/mkinitcpio.conf.d/archiso.conf:
#
#   HOOKS=(base udev microcode modconf kms memdisk archiso archiso_loop_mnt
#          block filesystems keyboard)
#
# Those hooks are correct for a live ISO and wrong for every disk install:
#   - `archiso`/`archiso_loop_mnt`/`memdisk` search for a squashfs on a boot
#     medium that no longer exists once you boot from disk, and they are
#     deleted outright when the packages module removes mkinitcpio-archiso;
#   - there is no `autodetect`, so the image carries every module ever built;
#   - there is no `encrypt`, so a LUKS2 install (SPEC §9) can never unlock;
#   - there is no `fsck`.
#
# enigma-install-cleanup.sh replaces that drop-in with a disk-install config
# before the packages module runs; this script does the actual generation once
# /boot is populated.

set -euo pipefail

if [ ! -f /boot/vmlinuz-linux-cachyos ]; then
    echo "ERROR: /boot/vmlinuz-linux-cachyos missing — the kernel restore step" >&2
    echo "       must run before the initramfs is generated." >&2
    exit 1
fi

if [ -f /etc/mkinitcpio.conf.d/archiso.conf ]; then
    echo "ERROR: the archiso mkinitcpio drop-in is still present; the cleanup" >&2
    echo "       step must run before this one." >&2
    exit 1
fi

# -P builds every preset in /etc/mkinitcpio.d (linux-cachyos, and linux-lts
# when installed), producing both the default and fallback images that
# bootloader.conf references.
echo "Generating initramfs for all installed kernels..."
mkinitcpio -P

# Fail loudly rather than leaving a system that boots to an emergency shell.
if [ ! -f /boot/initramfs-linux-cachyos.img ]; then
    echo "ERROR: mkinitcpio did not produce /boot/initramfs-linux-cachyos.img" >&2
    exit 1
fi

echo "Initramfs generated:"
ls -la /boot/initramfs-*.img
