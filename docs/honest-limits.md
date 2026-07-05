# Enigma OS — Honest Limits

We state our limitations plainly. This document is written in the user's voice and kept up-to-date as features ship.

## Amnesic ≠ Anonymous

The live amnesic mode leaves zero trace on your USB and internal disks, but **does not hide your network traffic or identity**.

- We do NOT route through Tor.
- Your ISP/network operator can see what you access.
- Website operators can log your IP, User-Agent, etc. normally.

**If you need anonymity**, use [Tails](https://tails.boum.org/), which routes through Tor by design.

Enigma's amnesic mode is for **privacy from local disk inspection**, not from network-level adversaries. It's the right tool if you:
- Use a public computer and want to leave no trace.
- Boot into live mode to run portable tools (password managers, offline docs).
- Need to prevent accidental writes to internal storage.

## No Cold-Boot-Attack Protection

After a live session ends, we do NOT:
- Physically wipe RAM (impossible; power-off achieves this naturally).
- Claim protection against cold-boot attacks (freezing RAM to read contents).

**Why**: Consumer hardware does not expose RAM-wipe primitives. Professional threat models requiring this protection should use hardware designed for it (e.g., systems with self-wiping memory on power loss).

## Kernel-Anticheat Games Don't Run

Games that require kernel-level anticheat (Valorant, some EA/Activision titles) will **not run on Linux in any form**:
- Not in native Wine
- Not in Proton
- Not in the optional Windows VM

This is not an Enigma limitation — it's a game publisher choice. The anticheat is fundamentally incompatible with Linux kernel architecture.

We document this plainly because users sometimes expect "just run it on Linux" to be possible. It isn't.

## Windows Tier 2 Requires Your License

The optional VM tier (running a full Windows 11 inside QEMU) requires:
- **Your own Windows license** (you must provide it; we do not sell or bundle one)
- At least **16GB of RAM** (8GB for the host, 8GB for the VM)
- A CPU with nested virtualization support

If you don't have these, Tier 1 (Bottles + Wine) is the only path forward.

## No Hibernation in v1

We do NOT support hibernation (suspend-to-disk). We ship:
- **Suspend** (suspend-to-RAM): works, ~2-3s resume
- **zram swap** (50% of RAM): keeps RAM pressure low without touching disk

**Why hibernation is deferred**:
- Btrfs + swapfile hibernation is fragile and often fails after updates.
- The extra complexity is not worth the upside on modern SSDs with fast suspend.
- We'd rather ship a reliable suspend/resume than a flaky hibernation.

If you need hibernation, use a distro that has solved it more thoroughly (Fedora, Debian) or file an issue once we hit Phase 2.

## Apple Silicon Not Supported

Enigma OS targets:
- **x86_64-v1** (baseline, runs on all x86_64)
- **ARM64-v2** (coming later, for ARM servers/SBCs)

We do **not** support:
- Apple Silicon (M1/M2/M3/etc.)
- 32-bit systems

If you need Enigma on Apple Silicon, run it in UTM or VMware Fusion as a guest VM (ARM64 build, when ready).

## Cannot Install from Inside Another OS (Without Rebooting)

Contrary to some marketing claims, **no operating system can be installed from inside another OS without rebooting**. The installer must boot.

Our USB Creator app + one reboot is the minimum possible flow. We document this because some users expect a Windows-to-Linux one-click installer that doesn't require reboot (it's not possible).

## What Works Well

- ✓ GPUs (NVIDIA, AMD, Intel auto-detected, drivers pre-installed)
- ✓ Wi-Fi, Bluetooth, codecs out of the box
- ✓ Steam games on Proton (non-anticheat)
- ✓ Web stacks (PHP, Node, Python) with per-project isolation
- ✓ LLMs (Ollama, ComfyUI) with GPU wiring
- ✓ Portable/amnesic mode (zero local disk trace)
- ✓ Snapshots + boot rollback (Btrfs + snapper)

---

See [SPEC.md](../SPEC.md) §21 for the same points in technical language.
