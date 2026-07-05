# CLAUDE.md
Project: Enigma OS — Arch/CachyOS-based zero-bloat distro for dev + AI
+ gaming + Windows apps. Installable to disk, installable to USB
("Enigma To Go"), or live-amnesic from USB. One CLI (`enigma`) + one
GUI (Enigma Center) over the same engine.

## Rules — always
- ALWAYS read SPEC.md fully before any task. It is the single source
  of truth; its section references are verified.
- Follow the build order in SPEC §14. A phase is done only when its CI
  gate is green (SPEC §19).
- Check the reuse map (SPEC §12) BEFORE implementing anything.
- Check the known failure points (SPEC §20) BEFORE designing any
  subsystem they mention — they are pre-solved; do not rediscover them.
- For Phase 3, hwdetect/DESIGN.md is authoritative — implement it
  exactly; TDD from its fixture matrix.

## Languages & code
- Go for cli/ and the daemon. Rust for indexd/. Qt/QML for center/.
  bash ONLY inside archiso hooks, shellcheck-clean.
- No new dependencies without a one-line justification.
- Licensing: MIT/Apache/BSD code may be vendored. GPL/AGPL projects
  are packaged and called as separate processes — NEVER copied into
  our Go/Rust sources. AGPL (Stability Matrix, WinApps) = study
  architecture only.

## Process
- Every change keeps ci/build-iso and ALL QEMU boot variants green
  (UEFI, BIOS, SecureBoot, live amnesic zero-write, copytoram detach,
  Ventoy).
- No AUR anywhere in the base ISO or CI. Ever.
- When unsure about Arch specifics, consult the Arch Wiki before
  guessing.
- Enter plan mode first for every new phase. Ask before architectural
  decisions not covered by SPEC.md.
- Branding: "Enigma OS" / `enigma` everywhere. No third-party or AI
  vendor branding/watermarks in any user-facing artifact.
