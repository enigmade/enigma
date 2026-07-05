# ENIGMA OS — FINAL MASTER SPECIFICATION v3.1
Zero-bloat Arch-based Linux distribution for developers, AI builders, and
gamers. Runs web stacks, LLMs, ComfyUI, Steam games, Windows apps, and
daily use. Installable to disk, installable to USB (portable), OR fully
live from USB (amnesic, zero trace). macOS-level polish, Windows-level
control.

This file is the single source of truth. Claude Code: read it fully
before every task. Section cross-references in this file are verified
and correct — trust them.

---

## 0. NON-NEGOTIABLE PRINCIPLES
- Zero bloat. No telemetry, no accounts, no ads, no Anthropic/vendor
  watermarks or branding anywhere in the OS — Enigma branding only.
- Everything works out of the box: GPU, Wi-Fi, BT, codecs, printers.
- One CLI controls everything: `enigma`. One GUI mirrors it: Enigma
  Center (§18). Same engine underneath — zero duplicated logic.
- Every update reversible (Btrfs snapshots + boot rollback).
- Three modes from the SAME ISO: (a) install to internal disk,
  (b) live-from-USB amnesic session leaving zero trace, (c) full
  portable install to the USB itself ("Enigma To Go", §17).
- Targets: x86_64 v1; ARM64 v2. Apple Silicon out of scope.
- Reuse upstream projects (§12) before writing anything custom.
- Honest docs: every limitation stated plainly (§21).

## 1. BASE SYSTEM
- Arch Linux base + CachyOS repositories (optimized packages,
  x86-64-v3) and linux-cachyos kernel (BORE scheduler) as default;
  linux-lts as fallback boot entry.
- Init: systemd. Bootloader: systemd-boot (UEFI) with Unified Kernel
  Images; GRUB fallback for BIOS.
- Btrfs subvolumes @, @home, @snapshots, @var_log; snapper + snap-pac
  (auto snapshot before every update); snapshots exposed in boot menu.
- zram swap 50% RAM (no disk swap; no hibernation in v1 — see §20.8).
- PipeWire + WirePlumber audio. linux-firmware + fwupd.
- Package sources: pacman (base), paru for AUR (user-invoked only,
  NEVER in base ISO or CI — §19), Flatpak+Flathub for GUI apps.

## 2. HARDWARE / GPU (highest-risk subsystem — build first, test most)
`enigma-hwdetect` runs in installer, first boot, live boot, and EVERY
boot in portable mode (§17). Full decision engine specified in
hwdetect/DESIGN.md — that file is authoritative for Phase 3.
Summary:
- Parses `lspci -nnk` vendor IDs: 10de=NVIDIA, 1002=AMD, 8086=Intel →
  writes /etc/enigma/hardware.toml, consumed by ai/game/win/doctor/
  Center. No other module ever re-detects.
- NVIDIA: nvidia-open-dkms (Turing+) else nvidia-dkms legacy pin;
  nvidia-utils, CUDA, cuDNN, nvidia-container-toolkit; DRM modeset for
  Wayland; hybrid laptops default iGPU-display + NVIDIA render/compute
  offload (prime-run).
- AMD: Mesa + radv default; ROCm only via opt-in `enigma ai enable-rocm`.
- Intel: Mesa + intel-media-driver.
- VM detection FIRST: never install dkms drivers inside a VM without
  passthrough (§20.15 in spirit; full logic in DESIGN.md).
- Acceptance: `enigma doctor` green on vulkaninfo, nvidia-smi + 1-line
  torch CUDA check, VA-API decode, external display, suspend/resume.

## 3. DESKTOP
- KDE Plasma 6 Wayland, stripped (no PIM suite, no KDE games,
  Discover = Flathub only).
- Preinstalled ONLY: Firefox (native, NOT Flatpak — §20.14; uBlock
  preinstalled), Dolphin, Konsole, VLC (all codecs), Gwenview, Kate,
  Bottles, Beekeeper Studio (community), Enigma Center (§18).
- Fonts: Inter (UI), JetBrains Mono Nerd (terminal), Noto fallback
  (full CJK + emoji).
- Dark theme default, sane touchpad defaults, night light on schedule.

## 4. `enigma` CLI (Go, single static binary)
Global commands: `doctor` (full health: GPU, boot time, ports,
snapshots, security), `update` (§11), `rollback`, `ports` (list/kill by
port), `clean`.

### 4a. `enigma dev` — universal project environments
`enigma dev up` in any folder:
1. Detect stack: composer.json→PHP/Laravel, package.json→Node,
   pyproject/requirements→Python, go.mod, Cargo.toml, Gemfile,
   docker-compose.yml→delegate to compose.
2. Runtime versions via mise (per-project, zero global conflicts).
3. Services as per-project systemd user units: MariaDB, PostgreSQL,
   Redis, Mailpit, MinIO, Meilisearch.
4. Ports: preferred port or atomic next-free from 8000-8999 (bind(0)
   probe + hold, §20.9), recorded in ~/.config/enigma/ports.db, URL
   printed.
5. Domains: <folder>.test via local dnsmasq (§20.1); TLS via mkcert
   root CA (created first boot, trusted system-wide AND in Firefox NSS
   db, §20.2).
`enigma dev db`: per-project DB+user, writes Beekeeper connection entry.
`enigma dev services`: status. `enigma dev down`: stops all.
Architecture reference: DDEV (Apache-2.0) — copy per-project config +
router concepts; default is native (non-Docker); `enigma dev --docker`
delegates to DDEV.

### 4b. `enigma ai`
- `ai setup`: reads hardware.toml → installs the MATCHING torch (cu12x
  NVIDIA / ROCm wheel AMD / CPU). Kills the #1 Linux AI pain:
  CUDA/torch mismatch. Full matcher logic in hwdetect/DESIGN.md.
- `ai ollama`: systemd service, GPU-wired, models in
  /var/lib/enigma/models.
- `ai comfy`: ComfyUI in isolated uv venv, correct torch, user service
  on auto port, URL printed; `ai comfy nodes add <repo>` via comfy-cli.
- `ai run <script-or-hf-repo>`: adopt Pinokio's (MIT) JSON script
  format; reuse community scripts where licensed.
- `ai status`: VRAM/RAM per model (nvidia-smi/rocm-smi), one-key
  unload; combined-VRAM warnings BEFORE launch (§20.10).
- Unified model cache: HF_HOME + Ollama + ComfyUI symlinked under
  /var/lib/enigma/models — study Stability Matrix shared-model layout
  (AGPL: study only, NO code copying; optionally ship the app as an
  extra GUI).

### 4c. `enigma game`
- `game setup`: multilib Steam + lib32 Vulkan libs (§20.7), Proton-GE
  via protonup-rs auto-update, gamemode, gamescope, MangoHud,
  controller udev rules, esync limits. Based on cachyos-gaming-meta.
- `game extras`: Lutris + Heroic.
- Docs state plainly: kernel-anticheat titles run on NO Linux (§21).

### 4d. `enigma win` — Windows software (two tiers)
Tier 1 (default, no Windows license): Bottles + wine-staging +
winetricks + DXVK/VKD3D. .exe double-click → "Run with Bottles"
auto-bottle with sane defaults (fonts, vcredist). `win run <exe>`,
`win apps` (creates .desktop entries → Windows apps appear in
launcher). Ship compat recipes (yaml in configs/) for top requested
Windows apps.
Tier 2 (opt-in `win setup-vm`; needs user's Windows license + 16GB
RAM): KVM/QEMU/libvirt Win11 VM + WinApps → individual Windows apps as
native-feeling windows via FreeRDP. VM auto-suspends when no Windows
app open (RAM freed), resumes ~2-3s. Shared ~/WinShared, clipboard
sync. GPU passthrough out of scope v1.

### 4e. `enigma containers`
Podman rootless default with docker CLI compat; `containers docker`
swaps in real Docker. NVIDIA container toolkit wired when detected.

## 5. SEARCH + LAUNCHER
- Launcher: ship Vicinae (GPL-3.0 — package/theme as separate process,
  NEVER merge its code into our binaries) as Super+Space launcher.
  Raycast-extension compatible.
- enigma-indexd (Rust, MIT, OUR code): fanotify watcher → Tantivy index
  of filenames, file contents (text/code/pdf), apps, enigma projects.
  Incremental, SCHED_IDLE + ionice, excludes node_modules/.git/venvs.
  Local IPC (unix socket) queried by a Vicinae extension.
- `ai ` prefix in launcher → streams from local Ollama.
- CI gate: top-10 results <50ms over a 1M-file synthetic tree on the
  NVMe reference machine.

## 6. INSTALLER
Calamares, Enigma-branded. Simple path: pick disk, username, done
(auto Btrfs layout). Advanced: dual-boot aware (detects Windows,
resizes), LUKS2 checkbox. Third target: "Install to this USB drive" =
Enigma To Go (§17).
First-boot wizard: GPU confirm + profile toggles (Web dev / AI /
Gaming / Windows apps / Minimal) — each just runs the matching
`enigma ... setup`. Study Omakub (MIT) for wizard UX.

## 7. UPDATES & RELIABILITY
`enigma update`: snapper pre-snapshot → apply → dkms/initramfs/
boot-entry verify (§20.4) → post-snapshot → summary. Boot failure →
previous snapshot selectable in boot menu; `enigma rollback` from TTY.
Stable channel: packages held 7 days behind upstream unless security.

## 8. LIVE USB / AMNESIC MODE (same ISO as installer)
Boot menu entries:
1. "Install Enigma OS"
2. "Live session (amnesic)" — runs from USB; ALL writes go to a tmpfs
   RAM overlay (archiso cowspace). Nothing ever written to USB or any
   internal disk. Power off = everything gone.
3. "Live session — load to RAM" (copytoram=y): copies full image to
   RAM (~30-60s), then the USB CAN BE PHYSICALLY REMOVED mid-session.
   Entry greys out with a message if RAM < image size + 6GB headroom
   (§20.5).
Amnesic guarantees in live mode:
- Internal disks NOT auto-mounted; udisks polkit rule requires explicit
  user action with a warning (§20.6).
- No swap on any disk (zram only). machine-id randomized per boot.
- MAC address randomization toggle in live greeter (NetworkManager
  cloned-mac=random), default ON in live mode.
- On shutdown tmpfs is freed; RAM clears on power-off. We do NOT claim
  cold-boot-attack protection (documented, §21).
4. Optional persistence (opt-in, Tails-style): `enigma live persist
   create` makes a LUKS2-encrypted second partition on the USB; boot
   entry "Live with encrypted persistence" overlays /home + chosen
   configs. Wrong/no passphrase → falls back to pure amnesic.
HONEST SCOPE LINE (goes in docs verbatim): amnesic ≠ anonymous. We do
not route through Tor and do not claim Tails-level anonymity. Users
needing anonymity should use Tails.
USB image requirements: hybrid ISO (dd + Rufus + Etcher + Ventoy all
work — §17), Secure Boot bootable via signed shim (§20.11),
persistence survives ISO rebuilds.

## 9. SECURITY BASELINE (installed, portable, and live)
- Secure Boot: sbctl-signed UKIs; shim path documented for locked
  firmware; "disable SB" fallback documented (§20.11).
- Firewall: nftables default-deny inbound (dev/ai services bind
  localhost only unless `--expose`).
- AppArmor enabled with default profiles.
- LUKS2 full-disk encryption = first-class installer checkbox.
- dnsmasq for .test binds 127.0.0.1:53 ONLY; systemd-resolved stub
  listener disabled (§20.1).
- fwupd firmware updates surfaced in `enigma update`.
- No listening network services in default install. `enigma doctor`
  includes a security section (open ports, SB state, encryption state).

## 10. BOOT SPEED — HARD BUDGET
Power button → usable desktop <10s on NVMe reference machine (LUKS
passphrase typing excluded). UKI + systemd-boot timeout 0 (hold key
for menu). dracut hostonly + zstd for disk installs; GENERIC initramfs
for To Go (§17, §20.16). NetworkManager/bluetooth/fwupd/indexd start
after graphical.target. Plymouth flicker-free splash.
CI: systemd-analyze archived every build; regression >1.5s fails the
build. Live-USB target: <25s to desktop from a USB 3.0 stick
(copytoram variant excluded from this number).

## 11. UPDATE SPEED — "FEELS INSTANT"
Replicates the macOS feel: download invisibly, apply in seconds.
- enigma-updated (systemd timer; AC power + unmetered network only):
  pre-downloads pacman packages (pacman -Syuw) + Flatpak OSTree deltas
  to cache; NEVER auto-applies.
- Desktop notification when staged: "Updates ready — install in ~20s"
  with size already downloaded.
- `enigma update` then: pre-snapshot → apply from cache → verify
  dkms/kernel (§20.4) → post-snapshot. No download wait.
- pacman.conf: ParallelDownloads=8, zstd; keep 2 versions in cache.
- Kernel/driver-only updates → offer `systemctl soft-reboot`
  (userspace restart, seconds); state clearly when a REAL reboot is
  required (kernel/nvidia).
- Flatpak apps update via OSTree deltas in background; GUI apps never
  block on updates.
- CI gate: simulated 2GB staged update applies in <60s in the test VM
  (excluding download).

## 12. UPSTREAM REUSE MAP (check BEFORE building anything)
Rule: wrap or package the upstream project; write original code ONLY
for glue and gaps. Copy code into our repos ONLY from MIT/Apache/BSD
projects. GPL/AGPL = ship as separate packages/processes.

| Need | Project | License | Strategy |
|---|---|---|---|
| Optimized base+kernel | CachyOS repos/kernel | mixed/GPL | consume as packages; study archiso profile |
| Gaming stack | cachyos-gaming-meta, protonup-rs | GPL/MIT | package + thin wrapper |
| GPU image tricks | Bazzite (ublue-os) | Apache-2.0 | study NVIDIA build scripts + quirks lists only |
| Launcher | Vicinae | GPL-3.0 | ship + theme as separate process; contribute upstream |
| Search index lib | Tantivy | MIT | build enigma-indexd on it (our custom moat) |
| Web env architecture | DDEV | Apache-2.0 | copy concepts; native reimplementation |
| Runtime versions | mise (jdx/mise) | MIT | vendor as-is; dev shells out to it |
| Python venvs | uv (astral-sh) | MIT/Apache | vendor as-is for all ai venvs |
| AI install scripts | Pinokio format | MIT | adopt JSON format + reuse licensed scripts |
| Shared model store | Stability Matrix | AGPL-3.0 | STUDY ONLY, clean reimplementation; optional GUI ship |
| LLM serving | Ollama | MIT | package + systemd unit + GPU wiring |
| Image/video gen | ComfyUI + Manager + comfy-cli | GPL-3.0 | package, uv-isolate, shell out to comfy-cli |
| Windows apps T1 | Bottles / Wine / DXVK | GPL/LGPL | package + recipes |
| Windows apps T2 | WinApps + KVM/QEMU | AGPL/GPL | package + wizard |
| DB GUI | Beekeeper Studio community | GPL/MIT | ship; dev db writes its connection config |
| Snapshots | snapper / snap-pac / grub-btrfs | GPL | package + preconfigure |
| Local TLS/DNS | mkcert, dnsmasq | BSD/GPL | package; we orchestrate |
| Setup UX ideas | Omakub (basecamp) | MIT | study first-boot wizard + app choices |

Custom code = ONLY: enigma CLI (Go), enigma daemon + Center GUI (§18),
enigma-indexd (Rust), hwdetect, first-boot wizard, live-persistence +
To Go tooling, USB Creator (§17), and glue. ~5 components instead of
~12. Everything else is packaging.

## 13. REPO LAYOUT
enigma/
  iso/        archiso profile (installer + live), boot entries, branding
  installer/  calamares config/theme + first-boot wizard
  cli/        Go `enigma` binary + daemon (dev,ai,game,win,containers,
              live,togo modules; unix-socket API for Center)
  center/     Qt/QML Enigma Center GUI (thin client of the daemon API)
  indexd/     Rust Tantivy daemon + Vicinae extension
  hwdetect/   DESIGN.md + detection/driver logic (lspci fixtures)
  creator/    cross-platform Enigma USB Creator (§17)
  configs/    /etc skeletons, systemd units + slices, snapper, dnsmasq,
              nftables, compat recipes
  ci/         ISO build, QEMU boot tests (UEFI+BIOS+SecureBoot+live+
              copytoram+Ventoy), CLI tests, search bench, update-speed
              test, QoS latency test
  docs/       user docs incl. honest-limits pages (§21)

## 14. BUILD ORDER (a phase is done ONLY when its CI gate is green)
P1  Branded hybrid ISO; boots QEMU UEFI+BIOS+Ventoy; live amnesic
    verified zero-write                                        2-3wk
P2  Installer + Plasma + Btrfs/snapper/rollback + boot-speed CI 3-4wk
P3  hwdetect + GPU stacks + doctor + QoS slices (§19a)         3-4wk
P4  enigma dev + Beekeeper wiring                              4-6wk
P5  enigma ai (torch matcher, Ollama, ComfyUI, Pinokio format) 3-4wk
P6  indexd + Vicinae integration                               ~2wk
P7  enigma game + containers + win Tier 1                      ~2wk
P7.5 Enigma Center GUI (thin skin over the daemon)             2-3wk
P8  Enigma To Go + live persistence + hardening + win Tier 2 +
    USB Creator                                                3-4wk
P9  Update pipeline polish (§11) + real-hardware matrix        ongoing

## 15–16. (reserved — do not renumber later sections)

## 17. USB WORKFLOWS — FLASH, COPY-ONE-FILE, AND "TO GO"
Three ways to get Enigma onto a USB stick; ALL must work from the same
hybrid ISO (CI-gated):

A) DIRECT FLASH (primary, no Ventoy needed): the ISO is a hybrid image
   — Rufus, balenaEtcher, Fedora Media Writer, GNOME Disks, and plain
   `dd` all produce a bootable stick. This is the documented default
   path.
B) ENIGMA USB CREATOR (creator/): small cross-platform app (Windows,
   macOS, Linux builds) that downloads/verifies the latest ISO
   (sha256 + signature) and flashes it in one click, with an explicit
   disk picker that HIDES internal disks by default (§20.15). This is
   the "just click" experience for non-technical users.
C) VENTOY (optional convenience, not required): copy enigma.iso onto a
   Ventoy stick as a plain file; new release = replace the file. CI
   gate: ISO boots under Ventoy in QEMU (normal + GRUB2 mode).

HONEST NOTE (docs verbatim): no operating system can be "installed by
clicking inside Windows without rebooting" — an OS must boot to
install. The closest possible flow, which we ship, is: run USB Creator
→ reboot → boot menu → done.

ENIGMA TO GO (full portable install): Calamares "Install to this USB
drive" path — a REAL install targeting the stick:
- Btrfs zstd compression + commit=120 (flash-friendly), fstrim timer,
  journald volatile + size-capped, no disk swap (zram only), relatime.
- GENERIC UKI — NOT hostonly initramfs; must boot on any machine
  (§20.16). hwdetect runs EVERY boot so GPU drivers adapt per machine
  (NVIDIA desktop today, AMD laptop tomorrow).
- Minimum stick: USB 3.0, 64GB; warn below that.
All three USB experiences in one docs table: (a) amnesic live (§8),
(b) live + encrypted persistence (§8), (c) To Go full install (here).

## 18. ENIGMA CENTER (GUI — the visual face of the CLI)
Single Qt/QML app, launchable from tray + launcher.
STRICT RULE: thin GUI over the SAME Go core the CLI uses (local
unix-socket API from the enigma daemon) — zero duplicated logic.
Anything the GUI does, the CLI can do, and vice versa.
Tabs:
- Runtimes: card per language (PHP, Python, Node, Go, Rust, Ruby,
  Java). Dropdown of installed versions + one-click install of others
  (mise behind it). Global default AND per-project override, shown
  visually.
- Services: every enigma-managed service (MariaDB, Postgres, Redis,
  Mailpit, Ollama, ComfyUI, per-project servers) with start/stop
  toggle, port, RAM/VRAM usage, "open in browser", log viewer.
- AI: installed models with sizes, download/remove, combined VRAM
  meter, per-model load/unload, ComfyUI custom-node manager, "won't
  fit" warnings (§20.10) shown BEFORE launching.
- Projects: detected dev projects, .test links, DBs, one-click
  up/down.
- Windows apps: bottles + VM state (Tier1/Tier2), launch/add.
- System: snapshots + one-click rollback, staged updates ("install in
  ~20s"), boot-time breakdown, doctor results with plain-language
  fixes.
CI: Center must pass headless tests against a mocked enigma daemon.

## 19. TESTING RULES
- Never mark a phase done without a green CI gate. No AUR in base ISO
  or CI, ever (§20.13).
- Every ISO change boots in QEMU: UEFI, BIOS, SecureBoot (OVMF), live
  amnesic (verify ZERO writes to an attached virtual disk via qemu-img
  compare), copytoram (verify session survives USB device detach),
  Ventoy (normal + GRUB2 mode).
- CLI: unit + integration tests in an Arch container. hwdetect:
  fixture lspci outputs per DESIGN.md test matrix.
- Shellcheck-clean on all hooks; no bash beyond archiso hooks.

### 19a. MULTITASKING QoS ("like butter" under AI load)
- systemd resource-control slices: AI inference (Ollama, ComfyUI) in
  enigma-ai.slice with CPUWeight=40 + lowered IO weight; the
  interactive session stays responsive while models generate.
- enigma-indexd stays SCHED_IDLE (§5). Compilations via `enigma dev`
  get nice=10 default (override flag).
- systemd-oomd tuned so a runaway model is killed before the desktop
  freezes; Ollama keep-alive releases VRAM after idle timeout
  (configurable in Center → AI).
- CI gate: while synthetic GPU+CPU load runs in enigma-ai.slice,
  terminal input latency and window-switch time in the test VM stay
  under thresholds; regression fails the build.

## 20. KNOWN FAILURE POINTS — SOLVED IN ADVANCE (do not rediscover)
1.  dnsmasq vs systemd-resolved both want port 53 → ship resolved with
    DNSStubListener=no + NetworkManager dns=dnsmasq from day one.
2.  mkcert + Firefox: system trust store is NOT enough; install into
    Firefox NSS db (certutil) or cert errors appear only in Firefox.
3.  NVIDIA + Wayland: needs nvidia-drm.modeset=1 AND fbdev=1 on newer
    drivers in the cmdline; without it → black screen after install.
    Append to UKI cmdline ONLY after dkms verify passes. Keep nouveau
    blacklisted on NVIDIA driver installs.
4.  DKMS vs kernel updates: update flow MUST verify dkms status for
    EVERY installed kernel before declaring success, else next boot
    fails. #1 rolling-distro breakage.
5.  copytoram: refuse (grey out) when RAM < image size + 6GB, checked
    in the boot hook — not after OOM.
6.  Live amnesic disk-mount leak: udisks auto-mount MUST be
    policy-blocked in live mode or the "zero trace" claim is false the
    first time a user clicks an internal drive in Dolphin.
7.  Steam needs multilib + 32-bit Vulkan libs (lib32-mesa /
    lib32-nvidia-utils); forgetting them = games silently fall back to
    CPU rendering.
8.  Btrfs + swapfile hibernation is fragile: do NOT promise
    hibernation in v1; zram + suspend only.
9.  Port allocator race: two `enigma dev up` at once must not grab the
    same port — allocate via bind(0)+hold, never scan-then-bind.
10. Ollama + ComfyUI VRAM collision: `ai status` shows combined VRAM;
    `ai run` warns when the requested model won't fit with what's
    loaded and offers unload. Never let both OOM the GPU silently.
11. Secure Boot: unsigned live USBs fail on many laptops with SB on —
    ship a signed shim path AND document the "disable SB" fallback.
12. Calamares + systemd-boot + UKI is a less-trodden path than GRUB;
    keep a tested GRUB fallback flag in the installer until P8.
13. AUR in CI: never — builds become unreproducible. User-space only.
14. Flatpak Firefox breaks mkcert/NSS and .test resolution inside the
    sandbox → ship NATIVE Firefox.
15. USB Creator wrong-disk catastrophe: the disk picker must hide
    internal/system disks by default, require an explicit "show all"
    toggle + typed confirmation of the device name before flashing.
16. Enigma To Go with a hostonly initramfs boots ONLY on the machine
    it was created on — portable installs MUST use a generic initramfs
    (all storage/USB/GPU modules included) or "boots anywhere" is
    false.

## 21. HONEST LIMITS (docs pages, written plainly)
- Amnesic ≠ anonymous (no Tor; use Tails for anonymity).
- No cold-boot-attack protection claim.
- Kernel-anticheat games (Valorant, some EA/Activision titles) run on
  no Linux, neither Wine nor VM tier.
- Windows Tier 2 requires the user's own Windows license + ~16GB RAM.
- No hibernation in v1. Apple Silicon not supported.
- An OS cannot be installed from inside another OS without rebooting;
  USB Creator + one reboot is the minimum possible flow.
