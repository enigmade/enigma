# Enigma OS — Zero-Bloat Arch-Based Linux

A macOS-level polish, Windows-level control Linux distribution. Run web stacks, LLMs, ComfyUI, Steam games, and Windows apps out of the box. One CLI (`enigma`), one GUI (Enigma Center), same powerful engine underneath.

## What's Inside

- **Zero bloat**: No telemetry, no accounts, no ads, no vendor watermarks.
- **GPU-native**: Automatic detection + driver wiring for NVIDIA, AMD, Intel GPUs.
- **AI-ready**: Torch on-board, Ollama, ComfyUI, one-command venv setup.
- **Dev stack**: PHP, Node, Python, Go, Rust, Ruby per-project with `enigma dev up`.
- **Gaming**: Steam, Proton-GE, Lutris, Heroic, all wired by default.
- **Windows apps**: Bottles + Wine Tier 1 (no license needed), optional VM Tier 2.
- **Portable & amnesic**: Three modes from the same ISO: install to disk, portable USB install, or fully live (zero trace).

## Build Status

**Phase 1** (current): ISO boots QEMU UEFI+BIOS ✓ | Amnesic zero-write guarantee ✓ | Ventoy compatible ✓

See [STATUS.md](STATUS.md) for detailed phase progress, deferred work, and CI gaps.

## Quick Start (for developers)

### Build locally (requires Linux)
```bash
cd iso
sudo mkarchiso -v -w /tmp/work -o /tmp/out .
ls /tmp/out/enigma-*.iso
```

### Test the ISO (requires QEMU)
```bash
# UEFI boot
bash ci/qemu/boot-uefi.sh /tmp/out/enigma-*.iso /tmp

# BIOS boot
bash ci/qemu/boot-bios.sh /tmp/out/enigma-*.iso /tmp

# Amnesic (zero-write) verification
bash ci/qemu/test-amnesic.sh /tmp/out/enigma-*.iso /tmp
```

### Automated CI
Push to `main` or open a PR → GitHub Actions builds the ISO and runs the full boot matrix (UEFI, BIOS, amnesic, copytoram, Ventoy).

## Project Structure

```
iso/              archiso profile (packages, boot entries, branding)
cli/              Go `enigma` binary (deferred: Phase 4)
center/           Qt/QML GUI (deferred: Phase 7.5)
indexd/           Rust search daemon (deferred: Phase 6)
hwdetect/         GPU detection engine (DESIGN.md is authoritative for Phase 3)
creator/          USB Creator app (deferred: Phase 8)
ci/qemu/          QEMU boot matrix tests
.github/workflows/ GitHub Actions CI pipeline
docs/             User documentation (honest-limits, setup guides)
```

## Key Design Decisions

See [SPEC.md](SPEC.md) and [CLAUDE.md](CLAUDE.md) for:
- Single source of truth (SPEC.md, section references verified)
- Reuse map (§12: which upstream projects we consume vs. build)
- Known failure points (§20: pre-solved issues, do not rediscover)
- Build order and phase gates (§14, §19)

## Honest Limits (See §21 in SPEC.md)

- Amnesic ≠ anonymous (no Tor; use Tails for anonymity)
- No cold-boot-attack protection
- Kernel-anticheat games (Valorant) do not run on any Linux
- Windows Tier 2 requires your own Windows license + 16GB RAM
- No hibernation in v1. Apple Silicon not supported.

## Docs

- [SPEC.md](SPEC.md) — complete specification; the single source of truth
- [STATUS.md](STATUS.md) — current phase progress and gaps
- [hwdetect/DESIGN.md](hwdetect/DESIGN.md) — GPU detection engine (Phase 3)
- [docs/honest-limits.md](docs/honest-limits.md) — detailed honest scope (WIP)

## License

TBD (once first stable release is ready). Enigma OS is assembled from Arch Linux (GPL), CachyOS (mixed), and various FOSS projects; components are licensed per their sources (MIT, Apache, GPL, BSD). No proprietary or AI vendor branding will ever be included.

## Contributing

This is a single-author project in early phases. Interested in contributing? File an issue first.
