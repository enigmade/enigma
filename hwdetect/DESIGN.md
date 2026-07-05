# hwdetect/DESIGN.md — GPU DETECTION & DRIVER DECISION ENGINE (Phase 3)
Authoritative design for SPEC §2. Place this file at hwdetect/DESIGN.md.

## OUTPUT CONTRACT
Single source of truth: /etc/enigma/hardware.toml — every other module
(ai, game, win, doctor, Center) READS this, never re-detects.

[system]  ram_gb, cpu_vendor, cpu_flags(v3/v4), is_vm, is_laptop, battery
[gpu.N]   vendor(nvidia|amd|intel), pci_id, name, vram_mb, driver_state
          (none|installed|mismatch), generation, role(primary|render|compute)
[ai]      torch_target(cu128|cu121|rocm6|cpu), tier(entry|standard|
          creator|studio), rocm_supported(bool), verified(bool), notes[]
[display] session(wayland|x11), hybrid_mode(none|prime|switcheroo)

## DETECTION PIPELINE (order matters)
1. Enumerate: lspci -nnk → class 0300/0302 devices. Parse vendor
   (10de/1002/8086), device ID, and the CURRENTLY BOUND kernel driver.
2. VM check FIRST: systemd-detect-virt. If VM → virtio/vmware GPU
   path, torch=cpu, SKIP all dkms logic (installing nvidia in a VM
   without passthrough is the classic self-brick). GPU passthrough
   detected = treat that device as bare metal.
3. NVIDIA generation from PCI device ID ranges (maintained table,
   fixtures in tests):
   - Blackwell(RTX50)/Ada(40)/Ampere(30)/Turing(20/16) →
     nvidia-open-dkms + CUDA 12.8 stack, torch cu128.
   - Pascal(10)/Maxwell(9) → nvidia-dkms LEGACY branch pin, torch
     cu121, WARN "last supported driver branch; AI works but frozen".
   - Older → nouveau only, torch=cpu, honest note.
   - UNKNOWN ID newer than table → assume newest path (open driver) +
     flag "unverified GPU, report ID" (future-proofs next RTX gen).
4. AMD: Mesa/radv always (gaming/desktop solved). ROCm decision is
   SEPARATE and lazy: rocm_supported computed from a gfx-arch
   allowlist (RDNA2+/CDNA official; RDNA1 = HSA_OVERRIDE hack, mark
   "community"). Never install ROCm in Phase 3 — only
   `enigma ai enable-rocm` does, using this precomputed answer.
5. Intel: Mesa + intel-media-driver. Arc discrete → note oneAPI
   possible; torch stays cpu until user opts into ipex (v1.5).
6. Hybrid laptop (2 GPUs, one Intel/AMD + one NVIDIA): DEFAULT = iGPU
   drives the display, NVIDIA as render-offload (prime-run) AND
   compute (CUDA works without driving a display — key fact).
   Full-NVIDIA mode is a Center toggle, not the default (battery).
   Record hybrid_mode.
7. Multi-dGPU / eGPU: role=compute for secondaries; the ai module
   exposes CUDA_VISIBLE_DEVICES per service in Center. eGPU
   (thunderbolt) → note hotplug; doctor rechecks on every boot.
8. Portable mode (Enigma To Go, SPEC §17): this whole pipeline runs on
   EVERY boot; on hardware change, reconcile drivers before the
   session starts (generic initramfs guarantees we reach userspace).

## DRIVER INSTALL RULES (the part that bricks systems — be paranoid)
- NEVER switch driver stacks in one shot on a running system without:
  pre-snapshot → install → dkms build VERIFY for every installed
  kernel → initramfs rebuild → boot-entry check → THEN reboot prompt.
- nouveau blacklist + nvidia-drm.modeset=1 fbdev=1 appended to the UKI
  cmdline ONLY after dkms verify passes (else the black-screen trap,
  SPEC §20.3).
- Live-USB mode: load nvidia open modules if present on the ISO, else
  nouveau/simpledrm — live must ALWAYS reach a desktop; never
  dkms-build in live sessions.
- Rollback story: driver-change failure → `enigma rollback` restores
  the pre-driver snapshot; doctor detects "boot after failed driver
  change" (flag file) and offers rollback proactively on next boot.

## TORCH MATCHER (kills the #1 AI pain)
Map (vendor, generation, rocm_supported) → exact wheel index URL:
cu128 (NVIDIA Turing+/newer) | cu121 (Pascal legacy) |
rocm6.x (AMD allowlisted) | cpu (everything else / VMs).
`enigma ai setup` reads hardware.toml, writes a uv constraint file
used by EVERY venv (ollama native, comfy, ai run) — one decision,
applied everywhere. Verification = `import torch` + device smoke test
in the venv; result cached into hardware.toml [ai].verified=true.

## `enigma doctor` GPU SECTION (green/yellow/red + plain language)
vulkaninfo summary | driver-vs-kernel dkms status per kernel |
nvidia-smi / rocm-smi presence + versions | torch device smoke test |
VA-API decode | external display count | suspend/resume flag |
hybrid: prime-run glxinfo renders on the dGPU | VRAM free/total.
Every red comes with ONE suggested command, never a wall of text.

## TEST MATRIX (fixtures, no hardware needed in CI)
lspci fixture files: blackwell.txt ada.txt ampere.txt turing.txt
pascal.txt rdna3.txt rdna2.txt rdna1.txt vega.txt arc.txt iris.txt
hybrid-intel-nvidia.txt dual-nvidia.txt vm-virtio.txt
vm-passthrough.txt egpu.txt unknown-future-nvidia.txt
Each fixture → assert the full expected hardware.toml.
Plus: mocked dkms-verify failure → assert install ABORTS before any
cmdline change; VM fixture → assert zero dkms calls.
Real-hardware smoke checklist (manual, release-gated): 1x NVIDIA
desktop, 1x hybrid laptop, 1x AMD, 1x 8GB iGPU-only machine.
