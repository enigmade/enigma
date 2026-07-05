# Installing Enigma OS on an x86_64 machine

## 1. Get the ISO

Every push to `iso/**` builds a fresh ISO in CI. Download the latest:

1. Open the **Actions** tab → **ISO Build (Arch container)** workflow →
   newest green run.
2. Download the **`enigma-iso`** artifact (a zip containing
   `enigma-0.1.0-x86_64.iso`, ~1 GB uncompressed).
3. Unzip it.

You can also trigger a build on demand: Actions → *ISO Build (Arch container)*
→ **Run workflow**.

> Verify it's genuinely bootable before writing: the CI run's *Verify boot
> structure* step confirms both a BIOS and a UEFI El Torito image are present.

## 2. Write the ISO to a USB stick (≥ 4 GB)

The ISO is **isohybrid** — write it raw; do not "burn as files".

**Linux:**
```sh
sudo dd if=enigma-0.1.0-x86_64.iso of=/dev/sdX bs=4M status=progress oflag=sync
#                                        ^^^^ your USB device (NOT a partition)
```

**macOS:**
```sh
diskutil list                      # find /dev/diskN for the USB
diskutil unmountDisk /dev/diskN
sudo dd if=enigma-0.1.0-x86_64.iso of=/dev/rdiskN bs=4m
```

**Windows:** use [Rufus](https://rufus.ie) in **DD image** mode, or
[balenaEtcher](https://etcher.balena.io).

Ventoy also works — copy the `.iso` onto a Ventoy drive.

## 3. Boot it

1. Insert the USB, power on, open the firmware boot menu (usually `F12`,
   `F10`, `Esc`, or `Del`).
2. Pick the USB. **UEFI** (systemd-boot) and **legacy BIOS** (syslinux) are
   both supported.
3. At the menu choose **Enigma OS (Install)** or a **Live** entry.

The live session **auto-logs into the KDE Plasma (Wayland) desktop** as the
`enigma` user (no password). Networking (NetworkManager) is enabled.

If the GPU has no working KMS driver and Plasma can't start, you drop to an
autologin **root console** on tty1 — the system is still usable from there.

## 4. Install to disk

Double-click **Install Enigma OS** on the desktop (or run `sudo calamares`).
Calamares walks you through:

- language / keyboard / timezone
- **partitioning** — "Erase disk" creates a Btrfs layout with subvolumes
  `@`, `@home`, `@var_log`, `@snapshots`; optional **LUKS2** full-disk
  encryption
- **user** — username + password (added to `wheel` for sudo)
- confirm → install (≈ 5–10 min)

The installer copies the live system to disk (squashfs → Btrfs), removes the
live-only bits (installer, passwordless sudo, live autologin), installs
**systemd-boot** on UEFI (GRUB on BIOS), and enables NetworkManager + SDDM.

Reboot, remove the USB, and log in.

## 5. After first boot

The `enigma` CLI and `enigma-hwdetect` ship in `/usr/local/bin`:

```sh
enigma-hwdetect detect     # write /etc/enigma/hardware.toml (GPU/AI tiering)
enigma doctor              # health report: GPU, ports, services
enigma dev up              # spin up a per-project dev environment
enigma ai setup            # install the matching PyTorch wheel for your GPU
```

See [README.md](README.md) and [SPEC.md](SPEC.md) for the full feature set.

## Troubleshooting

- **Black screen after selecting an entry** — try the **Live** entry, or add
  `nomodeset` at the boot menu (press `e` on systemd-boot / `Tab` on syslinux).
- **Installer won't launch** — open Konsole and run `sudo calamares -d`; the
  `-d` flag prints a debug log so you can report what failed.
- **Wi-Fi not connecting** — click the network applet in the system tray, or
  from a console: `nmtui`.
