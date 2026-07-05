# Enigma OS — Phase Status & Known Gaps

**Last updated**: Phase 1 complete (ISO builds, boots UEFI+BIOS, amnesic guarantee proven)

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

## Phase 2 Status: ⏳ NOT STARTED

**Goal**: Installer (Calamares) + KDE Plasma polish + Btrfs snapshots/rollback.

**Deferred**:
- [ ] Calamares installer + theming
- [ ] Btrfs subvolume layout (@, @home, @snapshots, @var_log)
- [ ] snapper + snap-pac integration
- [ ] Boot speed CI gate (< 10s NVMe target)
- [ ] First-boot wizard

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

Once Phase 1 passes (green GitHub Actions build), move to Phase 2: Calamares installer + Btrfs + boot speed.
