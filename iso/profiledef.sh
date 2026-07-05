#!/bin/bash
# Enigma OS archiso profile definition

iso_name="enigma"
iso_label="ENIGMA_OS"
iso_publisher="Enigma OS"
iso_application="Enigma OS Live ISO"
iso_version="0.1.0"
install_dir="arch"
buildmodes=('iso')
# Modern archiso bootmode identifiers. UEFI via systemd-boot (efiboot/loader);
# BIOS via syslinux is added once the syslinux/ configs are in place.
bootmodes=('uefi-x64.systemd-boot.esp' 'uefi-x64.systemd-boot.eltorito')
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
