#!/bin/bash
# Enigma OS archiso profile definition

iso_name="enigma"
iso_label="ENIGMA_OS"
iso_publisher="Enigma OS"
iso_application="Enigma OS Live ISO"
iso_version="0.1.0"
install_dir="arch"
buildmodes=('iso')
# Consolidated archiso bootmode identifiers (non-deprecated in current
# archiso). bios.syslinux covers MBR + El Torito; uefi.systemd-boot covers
# ESP + El Torito. Together they satisfy the SPEC §1 UEFI+BIOS matrix.
bootmodes=('bios.syslinux' 'uefi.systemd-boot')
arch="x86_64"
pacman_conf="pacman.conf"
# squashfs (zstd) for the root image: Calamares' unpackfs module is built
# around squashfs (unsquashfs), so this maximizes installer reliability on
# real hardware vs. erofs.
airootfs_image_type="squashfs"
airootfs_image_tool_options=('-comp' 'zstd' '-Xcompression-level' '15' '-b' '1M')
# EVERY executable shipped in airootfs/ must be listed here. mkarchiso copies
# the tree with `cp -af --no-preserve=ownership,mode`, so it deliberately
# discards the modes the files have in git — this table is the only thing that
# makes anything executable on the ISO.
#
# The value format is "owner:group:mode" (mkarchiso splits it on ':' into
# chown owner:group + chmod mode). The previous values were a bare mode
# ("0755"), which made archiso run `chown -fh 0755: <file>` and
# `chmod -f "" <file>` — both invalid, and both silenced by their -f flag, so
# every entry below was quietly a no-op. Most visibly that left the installer
# launcher on the live desktop non-executable, and Plasma refuses to run a
# .desktop file that is not executable: clicking "Install Enigma OS" did
# nothing at all.
file_permissions=(
  ["/etc/shadow"]="0:0:400"
  ["/etc/gshadow"]="0:0:400"
  ["/root"]="0:0:750"
  ["/root/customize_airootfs.sh"]="0:0:755"
  # Launched by the live user from the desktop icon / autostart.
  ["/usr/local/bin/enigma"]="0:0:755"
  ["/usr/local/bin/enigma-hwdetect"]="0:0:755"
  ["/usr/local/bin/enigma-install"]="0:0:755"
  ["/usr/local/bin/enigma-install-autostart"]="0:0:755"
  # Installer-side repair scripts invoked by Calamares shellprocess modules.
  ["/usr/local/libexec/enigma-install-cleanup.sh"]="0:0:755"
  ["/usr/local/libexec/enigma-install-restore-boot.sh"]="0:0:755"
  ["/usr/local/libexec/enigma-install-initramfs.sh"]="0:0:755"
  ["/usr/local/libexec/enigma-generate-snapshot-entries.sh"]="0:0:755"
  # Plasma will not launch a desktop entry that lacks the executable bit.
  ["/etc/skel/Desktop/install-enigma-os.desktop"]="0:0:755"
  ["/etc/skel/.config/autostart/enigma-install.desktop"]="0:0:755"
)
