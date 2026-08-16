#!/bin/bash
# Turn the copied live root into a proper installed system. Runs INSIDE the
# target chroot (dontChroot: false) after unpackfs.
#
# unpackfs copies the live filesystem verbatim, so everything that makes the
# ISO a *live* medium lands on the installed disk too: passwordless sudo, SDDM
# autologin as the throwaway `enigma` user, the root autologin getty, the
# installer launcher, the amnesic udisks policy, and the archiso initramfs
# config. All of it has to go.

set -euo pipefail

# --- live session credentials and autologin --------------------------------
# Without these removals the installed system logs straight into a passwordless
# throwaway account with NOPASSWD sudo — i.e. no security boundary at all.
rm -f /etc/sddm.conf.d/10-enigma-autologin.conf
rm -f /etc/sudoers.d/enigma-live
rm -rf /etc/systemd/system/getty@tty1.service.d
rm -f /root/.automated_script.sh /root/.zlogin

# --- installer launchers ---------------------------------------------------
rm -f /etc/polkit-1/rules.d/49-enigma-calamares.rules
rm -f /etc/skel/Desktop/install-enigma-os.desktop
rm -f /etc/skel/.config/autostart/enigma-install.desktop
rm -f /home/enigma/Desktop/install-enigma-os.desktop
rm -f /usr/local/bin/enigma-install /usr/local/bin/enigma-install-autostart

# --- live-only policy ------------------------------------------------------
# The amnesic guarantees (SPEC §8) are live-mode promises. On an installed
# system, blocking udisks mounts and randomising the MAC just breaks normal
# desktop use.
rm -f /etc/polkit-1/rules.d/10-enigma-udisks-live.rules
rm -f /etc/NetworkManager/conf.d/10-enigma-live.conf
rm -f /etc/systemd/journald.conf.d/10-enigma-live.conf

# systemd-firstboot is masked on the ISO because it blocks an unattended live
# boot waiting on stdin. Calamares sets locale/timezone/hostname on the target,
# so the mask is simply unnecessary here — drop it rather than ship a masked
# unit forever.
rm -f /etc/systemd/system/systemd-firstboot.service

# A fresh machine-id must be generated on first boot of the installed system;
# an empty file is systemd's "generate me one" marker. (Calamares' machineid
# module normally handles this — belt and braces, since a machine-id copied
# from the live image would be identical on every install.)
: > /etc/machine-id

# --- initramfs configuration ----------------------------------------------
# Replace the archiso live config with one that can actually boot a disk.
# This MUST happen before the packages module removes mkinitcpio-archiso,
# otherwise any pacman-triggered mkinitcpio run references hooks that are
# already gone. See enigma-install-initramfs.sh for the full rationale.
rm -f /etc/mkinitcpio.conf.d/archiso.conf
mkdir -p /etc/mkinitcpio.conf.d
cat > /etc/mkinitcpio.conf.d/enigma.conf <<'EOF'
# Enigma OS installed-system initramfs (written by the installer).
#
# btrfs is the root filesystem for every Enigma install (SPEC §1), so the
# module belongs in the image rather than relying on autodetect finding it.
MODULES=(btrfs)
# `autodetect` is correct here and NOT on the ISO: a disk install only ever
# boots the machine it was made on, so a trimmed hostonly image is what we
# want for the SPEC §10 boot budget. (Enigma To Go is the opposite case and
# must stay generic — SPEC §17/§20.16.)
# `encrypt` is what makes the LUKS2 installer checkbox actually unlock at boot.
HOOKS=(base udev autodetect microcode modconf kms keyboard keymap consolefont block encrypt filesystems fsck)
COMPRESSION="zstd"
EOF

echo "Target cleaned: live credentials, live policy and archiso initramfs config removed."
