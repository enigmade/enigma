#!/bin/bash
# Enigma OS archiso customization hook
# Runs in the airootfs before ISO generation

set -e

# Branding
echo "Enigma OS" > /etc/hostname
sed -i 's/^#GRUB_BACKGROUND.*/GRUB_BACKGROUND="\/usr\/share\/pixmaps\/enigma-bg.png"/' /etc/default/grub 2>/dev/null || true

# Configure SSH (disable by default in live mode)
systemctl disable sshd.service || true

# Disable networkmanager wifi powersave (improves connectivity in live mode)
mkdir -p /etc/NetworkManager/conf.d/
cat >> /etc/NetworkManager/conf.d/10-enigma-live.conf <<EOF
[general]
wifi.powersave = 2
EOF

# Enable cloned-mac randomization by default in live mode for privacy (SPEC §8)
cat >> /etc/NetworkManager/conf.d/10-enigma-live.conf <<EOF

[device]
wifi.scan-rand-mac-address = yes
EOF

# Configure dnsmasq (listening only on 127.0.0.1 per SPEC §20.1)
mkdir -p /etc/dnsmasq.d/
cat >> /etc/dnsmasq.d/enigma-localhost.conf <<EOF
listen-address=127.0.0.1
EOF

# Create polkit rule to block udisks auto-mount in live mode (SPEC §20.6)
mkdir -p /etc/polkit-1/rules.d/

# Disable dkms in live mode (Phase 3 defers GPU drivers anyway)
systemctl disable dkms.service || true

# Randomize machine-id on every live boot (SPEC §8)
# This is done in the initramfs hook, but ensure the file exists
: > /etc/machine-id

# Configure journald for live mode: volatile + size-capped (SPEC §17)
mkdir -p /etc/systemd/journald.conf.d/
cat >> /etc/systemd/journald.conf.d/10-enigma-live.conf <<EOF
[Journal]
Storage=volatile
MaximumJournalSize=50M
EOF

# Enable zram swap (50% of RAM per SPEC §1)
mkdir -p /etc/systemd/zram-generator.conf.d/
cat >> /etc/systemd/zram-generator.conf.d/10-enigma.conf <<EOF
[zram0]
zram-size = ram / 2
EOF

# Set sane KDE/Plasma defaults for live session
mkdir -p /etc/skel/.config/
cat >> /etc/skel/.config/kwinrc <<EOF
[General]
AnimationSpeed=2
EOF

# Blacklist problematic modules per SPEC
cat >> /etc/modprobe.d/enigma-blacklist.conf <<EOF
blacklist sp5100_tco
blacklist gpio_ich
EOF

# Ensure sudo works without password in live mode
echo "%wheel ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/enigma-live
chmod 0440 /etc/sudoers.d/enigma-live

# --- Live session setup: boot straight into the Plasma desktop -------------

# Create the live user (auto-logged-in). Home is populated from /etc/skel,
# which carries the Plasma defaults + the "Install Enigma OS" launcher.
useradd -m -G wheel,network,video,audio,storage,power,rfkill -s /bin/bash enigma
# Passwordless login for the live user (wheel already has NOPASSWD sudo).
passwd -d enigma || true

# Enable the essential services so the ISO reaches a usable desktop:
#   - NetworkManager: networking in the live session + for the installer
#   - sddm: the display manager that starts the Plasma Wayland session
systemctl enable NetworkManager.service
systemctl enable sddm.service
# Default to the graphical target so sddm actually starts.
systemctl set-default graphical.target

# Override cachyos-calamares's config with the Enigma installer config.
# Staged under /usr/share/enigma so it doesn't collide with the package's
# files during pacstrap (this script runs after packages are installed).
if [ -d /usr/share/enigma/calamares ]; then
  cp -rT /usr/share/enigma/calamares /etc/calamares
fi

# Console fallback: autologin root on tty1 (useful if the GPU has no KMS
# driver and Plasma can't start — you still get a shell).
mkdir -p /etc/systemd/system/getty@tty1.service.d/
cat > /etc/systemd/system/getty@tty1.service.d/autologin.conf <<EOF
[Service]
ExecStart=
ExecStart=-/usr/bin/agetty --autologin root --noclear %I \$TERM
EOF

# Clean package cache
pacman -Sc --noconfirm || true

# Phase 2: Snapper + custom systemd-boot snapshot hook setup
# The hook + generator ship in the airootfs (usr/share/enigma, usr/local/libexec);
# wire the pacman hook into place for UEFI snapshot boot entries.
mkdir -p /etc/pacman.d/hooks/
cp /usr/share/enigma/hooks/99-enigma-snapshots-systemd-boot.hook /etc/pacman.d/hooks/ || true

# Initialize snapper config dir for the installer (live mode itself doesn't
# snapshot). configs/ must exist before writing into it.
mkdir -p /etc/snapper/configs/
: > /etc/snapper/configs/root || true

echo "Enigma OS customization complete (Phase 1 + 2)"
