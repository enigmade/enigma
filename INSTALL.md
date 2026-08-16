# Installing Enigma OS on an x86_64 machine

> **Read first — two things that will cost you data or an evening:**
>
> - **Secure Boot must be disabled.** We do not sign the bootloader yet, so
>   firmware with Secure Boot on will refuse the USB outright.
> - **The installer defaults to erasing the whole disk.** Dual-boot detection
>   is packaged but untested. Do not run it on a machine with data you want.
>
> To verify a build without risking anything, use a VM — see
> [TESTING.md](TESTING.md).

## 1. Get the ISO

Every push to `iso/**` builds a fresh ISO in CI. Download the latest:

1. Open the **Actions** tab → **ISO Build (Arch container)** workflow →
   newest green run.
2. Download the **`enigma-iso`** artifact (a zip containing
   `enigma-0.1.0-x86_64.iso`, ~3.2 GB uncompressed).
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

1. **Disable Secure Boot** in firmware settings (see the note at the top).
2. Insert the USB, power on, open the firmware boot menu (usually `F12`,
   `F10`, `Esc`, or `Del`).
3. Pick the USB. **UEFI** (systemd-boot) and **legacy BIOS** (syslinux) are
   both supported.
4. At the menu choose one of the three entries. The menu waits 10 seconds and
   defaults to **Install**.

The three entries are:

| Entry | What it does |
|---|---|
| **Enigma OS (Install)** | Live desktop, and **auto-starts the installer** |
| **Enigma OS (Live — Amnesic)** | Live only; all writes go to RAM, internal disks are never auto-mounted |
| **Enigma OS (Live — Load to RAM)** | Copies to RAM first, so the USB can be removed mid-session |

To look around without any chance of touching a disk, pick **Live — Amnesic**.

The live session **auto-logs into the KDE Plasma (Wayland) desktop** as the
`enigma` user (no password). Networking (NetworkManager) is enabled.

If the GPU has no working KMS driver and Plasma can't start, you drop to an
autologin **root console** on tty1 — the system is still usable from there.

## 4. Install to disk

On the **Install** entry the installer opens by itself. Otherwise double-click
**Install Enigma OS** on the desktop, or run `enigma-install` from Konsole.

Calamares walks you through:

- language / keyboard / timezone
- **partitioning** — "Erase disk" creates a Btrfs layout with subvolumes
  `@`, `@home`, `@var_log`, `@snapshots`; optional **LUKS2** full-disk
  encryption
- **user** — username + password (added to `wheel` for sudo)
- confirm → install (≈ 5–10 min)

The installer copies the live system to disk (squashfs → Btrfs), strips the
live-only bits (installer, passwordless sudo, live autologin, amnesic disk
policy), then repairs what the copy cannot carry: `mkarchiso` deletes `/boot`
from the live filesystem, so the kernel is restored from the boot medium and a
fresh initramfs is generated before **systemd-boot** (GRUB on BIOS) writes its
entries.

Reboot, remove the USB, and log in with the username you created — the live
`enigma` account does not exist on the installed system.

## 5. After first boot

The `enigma` CLI and `enigma-hwdetect` ship in `/usr/local/bin`:

```sh
enigma-hwdetect detect     # write /etc/enigma/hardware.toml (GPU/AI tiering)
enigma doctor              # health report: GPU, ports, snapshots, security
enigma dev up              # detect the project stack and allocate a port
enigma ai setup            # choose the matching PyTorch build for your GPU
```

**These are not all finished.** `hwdetect` and `doctor` do what they say.
`dev up` detects the stack and allocates a port, then stops — mise, per-project
services, mkcert and `.test` domains are not implemented. `ai setup` picks the
right torch build but does not install it. `enigma update` and `enigma
rollback` do not exist at all yet.

[SPEC.md](SPEC.md) describes the intended feature set; [STATUS.md](STATUS.md)
records what is actually built.

## Troubleshooting

- **Firmware won't boot the USB at all** — Secure Boot. Disable it in firmware
  settings; we don't sign the bootloader yet.
- **Black screen after selecting an entry** — try the **Live** entry, or add
  `nomodeset` at the boot menu (press `e` on systemd-boot / `Tab` on syslinux).
  A root shell auto-logs in on tty1 (**Ctrl+Alt+F1**); `journalctl -b -u sddm`
  says why Plasma didn't start.
- **Installer won't launch** — open Konsole and run `enigma-install`, or
  `sudo calamares -d` for a debug log. Note that plain `pkexec calamares` will
  *not* work in a Wayland session; that's what the `enigma-install` wrapper is
  for.
- **Installed system won't boot** — check `ls -la /boot` from a live session
  with the target mounted: a kernel and an initramfs must both be present.
- **Wi-Fi not connecting** — click the network applet in the system tray, or
  from a console: `nmtui`.
- **Keyboard and trackpad dead on a MacBook** — T2 chip; see
  [TESTING.md](TESTING.md). Use a USB keyboard and mouse.
