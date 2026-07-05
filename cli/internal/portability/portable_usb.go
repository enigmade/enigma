package portability

import "fmt"

// PreparePortableUSB configures USB persistence for Enigma OS (Phase 8 - SPEC §8)
//
// Layout: GPT partition table
//   - /dev/sda1 (512 MB, EFI FAT32) — systemd-boot + kernel
//   - /dev/sda2 (20 GB, Btrfs) — root filesystem (generic initramfs, auto-detect hardware)
//   - /dev/sda3 (remainder, LUKS2) — encrypted persistence (~/.config, /var/lib/enigma)
//
// Behavior:
//   - Boot with generic initramfs (auto-detects GPU via hwdetect)
//   - First boot: auto-populate LUKS2 key in TPM (with fallback to USB passphrase)
//   - Reboot into persistent system (mounts /dev/sda3)
//   - Updates: btrfs snapshot + overlay mount (rollback safe)
// PreparePortableUSB computes the GPT layout for a USB of diskSizeMiB and
// prints the exact sgdisk + cryptsetup command sequence that would create
// it. It does not execute anything destructive — the caller decides whether
// to run the returned plan — so this is safe to invoke for a dry run.
func PreparePortableUSB(device string, diskSizeMiB int) (*PartitionPlan, error) {
	plan, err := PlanPartitions(diskSizeMiB)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Partition plan for %s (%d MiB):\n", device, diskSizeMiB)
	for _, part := range plan.Partitions {
		size := "rest of disk"
		if part.SizeMiB >= 0 {
			size = fmt.Sprintf("%d MiB", part.SizeMiB)
		}
		fmt.Printf("  %d. %-14s %-6s start=%d MiB size=%s\n",
			part.Number, part.Name, part.Filesystem, part.StartMiB, size)
	}

	fmt.Println("\nCommands (not executed):")
	for _, cmd := range plan.SgdiskCommands(device) {
		fmt.Printf("  %v\n", cmd)
	}
	fmt.Printf("  %v\n", LuksFormatCommand(device+"3", "/run/enigma/persist.key"))

	return plan, nil
}
