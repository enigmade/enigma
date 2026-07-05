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
file_permissions=(
  ["/root/.automated_script.sh"]="0755"
  ["/root/.customize_airootfs.sh"]="0755"
  ["/etc/shadow"]="0400"
)
