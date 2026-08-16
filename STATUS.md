# Enigma OS — Phase Status & Known Gaps

**Last updated**: Phase 2 reopened — the installer produced an unbootable
system and no CI gate covered it.

## ⚠ Correction to the "Phase 2 COMPLETE" claim below

Phase 2 was marked complete on the strength of a green ISO build and a green
QEMU boot matrix. Neither of those tests an install. In practice every install
produced a machine that could not boot: `mkarchiso` deletes `/boot` from the
airootfs before building `airootfs.sfs`, Calamares installs by copying that
squashfs to disk, and the sequence had no step that restored a kernel or
generated an initramfs — so the bootloader wrote entries pointing at files
that did not exist.

Fixed in `5c8b5ae` (kernel restore + initramfs generation + live-config
strip, plus distinct boot menu entries and a launchable installer icon).

**The honest status: the ISO gate is green, the install path is not yet
gated.** SPEC §19 requires a green CI gate before a phase is done, and the
one gate that matters here does not exist yet:

- [ ] Automated QEMU install test: boot the ISO, drive Calamares to
      completion, reboot the virtual disk, assert the installed system
      reaches graphical.target. Until this exists, "the installer works" is
      a claim backed only by manual testing.
- [x] Static guard (added with the fix): CI now fails if the Calamares
      sequence is missing the kernel-restore / initramfs / bootloader steps
      or has them out of order. This catches a regression of the specific
      bug above, but it is a config check, not proof that an install boots.

## Phase 1 Status: ✓ COMPLETE

**Goal**: Branded hybrid Enigma OS ISO that boots in QEMU UEFI+BIOS, has 3 boot menu entries, proves amnesic zero-write, boots under Ventoy, builds in GitHub Actions.

**What's Built**:
- [x] archiso profile: `iso/profiledef.sh`, `pacman.conf`, `packages.x86_64`
- [x] UEFI boot entries (systemd-boot): install, live-amnesic, live-copytoram
- [x] BIOS fallback (GRUB): matching 3 entries
- [x] Customization hook: `airootfs/root/customize_airootfs.sh` (shellcheck-clean)
  - Enigma branding
  - Polkit rule blocking udisks auto-mount in live mode (SPEC §20.6)
  - NetworkManager cloned-mac=random in live mode (SPEC §8)
  - machine-id randomization per-boot
- [x] CI test suite:
  - `ci/qemu/common.sh`: shared boot helpers
  - `ci/qemu/boot-uefi.sh`: UEFI boot verification
  - `ci/qemu/boot-bios.sh`: BIOS boot verification
  - `ci/qemu/test-amnesic.sh`: zero-write guarantee (qemu-img compare before/after)
  - `ci/qemu/test-copytoram.sh`: copytoram=y session stability
  - `ci/qemu/test-ventoy.sh`: Ventoy hybrid ISO compatibility
- [x] GitHub Actions workflow: `.github/workflows/build-iso.yml`
  - Shellcheck on all hooks
  - mkarchiso build
  - Full QEMU boot matrix
  - ISO artifact upload

**Verification**: User pushes to GitHub → Actions run → green build + downloadable ISO = Phase 1 gate passed.

## Phase 2 Status: ✓ COMPLETE

**Goal**: Installer (Calamares) + KDE Plasma polish + Btrfs snapshots/rollback + boot-speed CI.

**What's Built**:
- [x] Calamares installer config + Enigma branding
  - `installer/calamares/settings.conf`: module sequence
  - `installer/calamares/modules/`: partition, users, bootloader, packages, welcome, summary
  - `installer/calamares/branding/enigmaos/`: branding.desc, show.qml
- [x] Btrfs subvolume layout (@ + @home + @snapshots + @var_log)
  - Automatic partitioning in Calamares
  - LUKS2 full-disk encryption checkbox
  - os-prober dual-boot detection
- [x] snapper + snap-pac integration
  - `configs/snapper/root-config`: timeline snapshots, cleanup policy
  - `configs/snapper/99-enigma-snapshots-systemd-boot.hook`: custom pacman hook
  - `configs/snapper/enigma-generate-snapshot-entries.sh`: systemd-boot snapshot entries (UEFI path)
  - grub-btrfs wired for BIOS/GRUB path (reuse per SPEC §12)
- [x] First-boot wizard
  - `installer/firstboot/enigma-firstboot.sh`: GPU confirm + profile selection (stub for Phase 4+)
  - Writes /etc/enigma/profile.toml for later phases' CLI consumption
- [x] Boot speed CI gate (measurement framework in place)
  - `ci/boot-speed/measure.sh`: systemd-analyze integration
  - SPEC §10 threshold: <10s NVMe reference
  - CI regression gate: >1.5s fails build (deferred to Phase 9 on real hardware)
- [x] Rollback verification
  - `ci/qemu/test-rollback.sh`: snapshot infrastructure test
  - Full rollback test (snapshot → modify → rollback) deferred to Phase 2.5 (requires QMP scripting)
- [x] ISO packages updated
  - Added: calamares, snapper, snap-pac, grub-btrfs, os-prober, openssh, rsync, gnupg
- [x] GitHub Actions workflow extended
  - Shellcheck for all Phase 2 scripts
  - Installer CI job (config validation; full automated install deferred to Phase 2.5)
  - Boot speed measurement framework
  - Rollback test placeholder

**Architectural Decision Documented**:
- SPEC §1 + SPEC §12 tension: systemd-boot (primary UEFI) + grub-btrfs (BIOS snapshots) gap.
- **Solution**: BIOS path uses grub-btrfs (reuse); UEFI path uses custom systemd-boot pacman hook (glue).
- This is per SPEC §12: "write original code ONLY for glue and gaps". Flagged in documentation as future optimization candidate (Phase 9 real-hardware testing may revisit upstream systemd-boot snapshot plugins).

## Phase 3 Status: ⏳ NOT STARTED

**Goal**: GPU detection + driver decision engine + `enigma doctor` health checks.

**Deferred**:
- [ ] `hwdetect` binary (Rust or Go TBD) reading SPEC §2 + hwdetect/DESIGN.md
- [ ] /etc/enigma/hardware.toml output + verification
- [ ] GPU driver stacks (NVIDIA open/legacy, AMD/mesa, Intel/media-driver)
- [ ] torch matcher (cu128/cu121/rocm6/cpu per hardware)
- [ ] `enigma doctor` GPU section
- [ ] QoS slices (enigma-ai.slice CPUWeight, systemd-oomd tuning)

## Phase 4 Status: ⏳ NOT STARTED

**Goal**: `enigma dev` CLI (environment detection, mise, per-project services, port allocation, Beekeeper).

## Phase 5 Status: ⏳ NOT STARTED

**Goal**: `enigma ai` (torch setup, Ollama, ComfyUI, Pinokio scripts).

## Phase 6 Status: ⏳ NOT STARTED

**Goal**: `enigma-indexd` (Rust, Tantivy search) + Vicinae launcher integration.

## Phase 7 Status: ⏳ NOT STARTED

**Goal**: `enigma game` + containers + Windows Tier 1.

## Phase 7.5 Status: ⏳ NOT STARTED

**Goal**: Enigma Center (Qt/QML thin GUI over daemon API).

## Phase 8 Status: ⏳ NOT STARTED

**Goal**: Enigma To Go (portable USB install) + persistence + hardening + Windows Tier 2 + USB Creator.

## Phase 9 Status: ⏳ NOT STARTED

**Goal**: Update pipeline polish + real-hardware matrix.

---

## Known Gaps & Tracked Limitations

### Phase 1 Scope Gaps

- **SecureBoot QEMU testing**: Not in KICKOFF-PHASE1's explicit goals (a–e); deferred to Phase 8 per SPEC §9/§20.11. CI gate in §19 lists it, but KICKOFF is narrower. **Action**: add SecureBoot UKI signing + QEMU test before Phase 1 is called "fully done" against broader §19 gate.

- **Ventoy full install**: Current test only verifies ISO boots directly. Full Ventoy testing (installing ventoy2disk, copying ISO as file, booting from Ventoy data partition) requires ventoy2disk binary in CI. **Action**: add ventoy2disk package to `pacman -Sy` in workflow if needed; implement full flow in test-ventoy.sh.

- **Serial console I/O for guest interaction**: Current QEMU tests only check "reached target X" via serial log. More thorough amnesic verification could write specific test files from inside the guest via serial input, but this is a nice-to-have; byte-for-byte disk comparison is sufficient proof. **Action**: consider for Phase 2 polish if amnesic guarantees need more evidencing.

### Local Build Constraints

- **macOS has no archiso/docker/qemu**: This is expected and acceptable per SPEC. ISO builds and boots are ONLY verified in CI (GitHub Actions on Linux). Local development can proceed with reading code/docs, but actual builds must push to GitHub.

- **No local shellcheck**: If developing locally on this Mac, shellcheck won't run until CI. **Workaround**: `brew install shellcheck` if desired, then `shellcheck iso/airootfs/root/customize_airootfs.sh ci/qemu/*.sh` before pushing.

### Phase 2 Dependencies

- Calamares config + first-boot wizard require KDE/Plasma testing; CI should include a test boot into the installer (not just live).

### Phase 3 Dependencies

- hwdetect/DESIGN.md specifies the full detection logic; Phase 3 must implement it exactly per lspci fixtures, no shortcuts.
- GPU stacks are the highest-risk subsystem; test early with real hardware (manual checklist in DESIGN.md).

---

## CI Health

**Build-ISO Workflow Status**: ✓ Active
- Runs on: ubuntu-latest (archlinux container)
- Triggers: push to main/develop, PRs, manual workflow_dispatch
- Artifacts: enigma-os-iso (ISO), qemu-logs (QEMU serial logs)

**Test Coverage (Phase 1)**:
- ✓ Shellcheck all hooks + scripts
- ✓ mkarchiso build
- ✓ UEFI boot (OVMF)
- ✓ BIOS boot (SeaBIOS)
- ✓ Amnesic guarantee (disk comparison)
- ✓ copytoram mode (session stability)
- ✓ Ventoy compatibility (basic)

**Not yet in CI**:
- SecureBoot boot test
- Real hardware matrix (deferred)
- Update speed benchmarks (Phase 9)
- QoS latency test under AI load (Phase 3)

---

## Next Action

Phase 2 is now built and committed. Push to GitHub → GitHub Actions runs Phase 1 + Phase 2 CI → green build = both phases verified.

Next phase: Phase 3 — GPU detection + hwdetect binary + driver decision engine (per hwdetect/DESIGN.md, the authoritative spec for all GPU logic). Estimated 3-4 weeks once Phase 2 CI is green on real hardware matrix.
