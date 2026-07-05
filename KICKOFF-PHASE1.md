# PHASE 1 KICKOFF PROMPT
Paste this as your FIRST message in Claude Code, in a fresh git repo
containing SPEC.md, CLAUDE.md, and hwdetect/DESIGN.md.

---

Read SPEC.md and CLAUDE.md completely. We are starting Phase 1
(SPEC §14).

Goal of this phase ONLY: a branded, hybrid Enigma OS ISO based on
archiso that:
(a) boots to a live KDE Plasma session in QEMU under UEFI and BIOS,
(b) has the three boot menu entries from SPEC §8 (Install / Live
    amnesic / Live copytoram),
(c) proves the amnesic guarantee in CI by verifying an attached
    virtual disk is byte-identical after a live session (qemu-img
    compare),
(d) boots under Ventoy in QEMU (normal + GRUB2 mode) per SPEC §17,
(e) builds automatically in GitHub Actions with the ISO as an
    artifact.

Constraints: pull base packages from Arch + CachyOS repos per SPEC §1;
minimal package list per SPEC §3 (no enigma CLI yet — later phases);
no AUR anywhere per SPEC §19; all hooks shellcheck-clean; Enigma
branding only, per SPEC §0.

Enter plan mode first: propose the repo scaffold (SPEC §13), the
archiso profile structure, the CI workflow, and the QEMU test scripts.
List any decisions SPEC.md leaves open and ask me before coding. Do
not write code until I approve the plan.
