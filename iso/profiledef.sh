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
airootfs_image_type="erofs"
# lz4hc level 12 compression with tail-packing + dedupe (all valid mkfs.erofs
# extended options). The previous -Ezcache-raid-stripe-width flag does not
# exist and made mkfs.erofs fail with "Invalid argument".
airootfs_image_tool_options=('-zlz4hc,12' '-Eztailpacking,dedupe')
file_permissions=(
  ["/root/.automated_script.sh"]="0755"
  ["/root/.customize_airootfs.sh"]="0755"
  ["/etc/shadow"]="0400"
)
