#!/bin/bash
# Enigma OS archiso profile definition

iso_name="enigma"
iso_label="ENIGMA_OS"
iso_publisher="Enigma OS"
iso_application="Enigma OS Live ISO"
iso_version="0.1.0"
install_dir="arch"
buildmodes=('iso')
bootmodes=('uefi' 'bios')
arch="x86_64"
pacman_conf="pacman.conf"
airootfs_image_type="erofs"
airootfs_image_tool_options=('-zlz4hc,12' '-Ezcache-raid-stripe-width=32K')
file_permissions=(
  ["/root/.automated_script.sh"]="0755"
  ["/root/.customize_airootfs.sh"]="0755"
  ["/etc/shadow"]="0400"
)
