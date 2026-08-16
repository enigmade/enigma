# Testing Enigma OS

How to actually verify a build, in the order that wastes the least time.

**Test in a VM first.** Booting the live USB and installing to disk are two
very different risks, and only one of them can destroy data. The VM answers
the important question — *does the install complete and reboot into a working
desktop?* — in about twenty minutes, with nothing at stake.

---

## 1. Get the ISO

Every push to `iso/**` builds one in CI.

1. **Actions** → **ISO Build (Arch container)** → newest green run.
2. Download the **`enigma-iso`** artifact (a zip; the ISO inside is ~3.2 GB).
3. Unzip.

Artifacts expire after 90 days. You can also build on demand: Actions →
*ISO Build (Arch container)* → **Run workflow**.

---

## 2. Test in a virtual machine

This is the recommended loop. It costs nothing to retry and cannot harm the
host.

### VirtualBox (Windows, macOS, Linux)

New VM → Type **Linux**, Version **Arch Linux (64-bit)**, then:

| Setting | Value | Why |
|---|---|---|
| Memory | 4096 MB | Below ~3 GB the live session struggles |
| Disk | 40 GB | Enough for a full install |
| **System → Enable EFI** | **on** | Otherwise you test the BIOS path only |
| Storage | attach the `.iso` to the optical drive | |
| Display → Video Memory | 128 MB | Plasma needs headroom |

### UTM (macOS)

New → **Virtualize** → **Other**. Set the boot ISO, 4 GB RAM, 40 GB drive, and
make sure **UEFI Boot** is ticked.

### QEMU (command line)

```sh
qemu-img create -f qcow2 enigma-test.qcow2 40G

qemu-system-x86_64 \
  -m 4096 -smp 2 \
  -bios /usr/share/edk2/x64/OVMF.4m.fd \
  -drive file=enigma-test.qcow2,format=qcow2,if=virtio \
  -cdrom enigma-0.1.0-x86_64.iso \
  -boot d -vga virtio
```

Drop `-bios` to test the legacy BIOS path instead. Add `-enable-kvm` on Linux
(or `-accel hvf` on macOS) or it will be painfully slow.

---

## 3. What to check, in order

Work down this list. Each stage failing points at a different subsystem, so
note *where* it stops.

| # | Stage | Expected |
|---|---|---|
| 1 | Boot menu appears | Three entries, 10s timeout, "Install" preselected |
| 2 | Kernel boots | Scrolling messages, no hang |
| 3 | Desktop appears | KDE Plasma, auto-logged in as `enigma` |
| 4 | Installer opens | Only on the **Install** entry — it now auto-starts |
| 5 | Install completes | No error dialog |
| 6 | **Reboot into the installed system** | **The one that matters** |
| 7 | Login screen | SDDM, *your* username — not `enigma` |

Stage 6 is the point of the whole exercise. Everything before it was already
passing when the installed system could not boot at all.

### If it fails, where it fails tells you what broke

| Symptom | Likely cause |
|---|---|
| Firmware refuses the USB entirely | Secure Boot (see below) |
| Boot menu, then nothing | Kernel or initramfs on the live medium |
| Desktop, but no installer on the Install entry | `enigma-install` launcher / Wayland root |
| Error *during* the install | One of the `enigma-install-*` repair scripts |
| Installs, then "no bootable device" | Kernel restore into the target |
| Installs, then black screen or emergency shell | Initramfs generation |
| Desktop, but keyboard and trackpad dead | T2 Mac (see below) |

Useful once you're in a session:

```sh
journalctl -b -u sddm          # why the desktop did not start
sudo calamares -d              # run the installer with a debug log
ls -la /boot                   # after install: kernel + initramfs present?
```

If Plasma cannot start, a root shell auto-logs in on **tty1** —
**Ctrl+Alt+F1**.

---

## 4. Testing on real hardware

### Disable Secure Boot — required

**We do not sign the bootloader yet.** There is no `sbctl` and no signed shim
in the image, so any machine with Secure Boot enabled will refuse to boot the
USB. This is usually a vague "security violation", or the firmware silently
skipping the stick and booting your normal OS.

Turn Secure Boot off in firmware settings before blaming the ISO.

### On a Windows machine

Also **disable Fast Startup** first: Control Panel → Power Options → *Choose
what the power buttons do* → uncheck *Turn on fast startup*. Fast Startup
leaves the disk in a half-hibernated state; booting another OS against it can
corrupt the Windows filesystem.

Then boot the stick and choose **Live — Amnesic**. That mode writes nothing to
any internal disk and does not auto-mount them, so it is safe on a machine you
care about.

### ⚠ Do not run the installer on a machine with data on it

The installer currently defaults to **erase the entire disk**
(`initialPartitioningChoice: erase` in
`iso/airootfs/usr/share/enigma/calamares/modules/partition.conf`). Dual-boot
detection is packaged but **untested**. Install only on a spare machine, a
spare disk, or a VM.

### Known-bad hardware: Macs with the T2 chip

Intel MacBooks from ~2018–2020 (MacBookAir8,x, MacBookPro15,x and similar)
carry Apple's T2 security chip. The architecture is right — they are x86_64 —
but two things get in the way:

1. **The T2 blocks external boot by default.** Boot into Recovery (**⌘R**) →
   *Utilities* → *Startup Security Utility* → set **No Security** and **Allow
   booting from external media**. Without this you get *"unable to verify
   startup disk"*.
2. **The built-in keyboard and trackpad will not work.** T2 Macs route them
   through the T2 over a proprietary interface that needs the `apple-bce`
   driver, which is not in our kernel and exists only in the AUR — and
   CLAUDE.md forbids AUR in the base ISO. Wi-Fi (BCM4377) likewise needs
   firmware extracted from macOS.

A USB keyboard and mouse work around the input half. **These machines are a
poor testing target — use a VM or an ordinary PC.**

---

## 5. Working on the code from another machine

Everything is on GitHub; do not copy the working tree around.

```sh
git clone https://github.com/enigmade/enigma.git
cd enigma
```

The tracked source is ~9 MB across ~156 files. The working directory on disk
is much larger because `indexd/target` holds several hundred MB of Rust build
artifacts that are not tracked and are not worth moving.

**Do not move the tree through a USB stick.** FAT32 and exFAT cannot store
Unix executable bits, so every script would silently lose its `+x` — which is
exactly the class of bug that left the installer launcher unlaunchable. Git
records the mode in its index, so cloning preserves it.

### What runs where

| Task | Windows | macOS | Linux |
|---|---|---|---|
| Edit code | yes | yes | yes |
| `go test ./cli/... ./hwdetect/...` | yes | yes | yes |
| `cargo test` in `indexd/` | yes | yes | yes |
| **Build the ISO** | **no** | **no** | Arch only |
| Run the QEMU boot matrix | no | no | needs KVM |

`mkarchiso` only runs on Arch Linux, so **you cannot build the ISO locally on
Windows or macOS** — and you do not need to. Push, and CI builds it in about
fifteen minutes.

---

## 6. What CI proves, and what it does not

The `ISO Build (Arch container)` workflow gates every change on:

- shellcheck across `iso/` and `ci/` (pinned version)
- static validation of the Calamares install sequence
- every package in `packages.x86_64` resolving
- `mkarchiso` completing
- hybrid BIOS + UEFI El Torito structure
- the live initramfs carrying the archiso hooks and USB input drivers
- the QEMU boot matrix: UEFI, BIOS, amnesic zero-write, copytoram, Ventoy

**It does not install anything.** The boot matrix starts the ISO and watches
the serial console; no step runs the installer. That is how an installer which
produced an unbootable system stayed green for 27 consecutive runs — the
target had no kernel and no initramfs, and nothing checked.

A static guard now fails the build if the repair steps go missing or fall out
of order, but that is a configuration check, not proof that an install boots.
**Until an automated QEMU install test exists, stage 6 above is only ever
verified by a human.** See [STATUS.md](STATUS.md).
