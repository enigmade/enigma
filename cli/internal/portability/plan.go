package portability

import "fmt"

// Sizes in mebibytes. GPT layout for a portable Enigma USB (SPEC §8).
const (
	MiB = 1
	GiB = 1024 * MiB

	espSizeMiB  = 512      // EFI System Partition (systemd-boot + kernel)
	rootSizeMiB = 20 * GiB // Btrfs root with generic initramfs

	// A usable stick needs ESP + root + at least a little persistence.
	minPersistMiB = 2 * GiB
	MinDiskMiB    = espSizeMiB + rootSizeMiB + minPersistMiB
)

// Partition is one entry in the computed GPT layout.
type Partition struct {
	Number     int
	Name       string
	StartMiB   int
	SizeMiB    int // -1 means "rest of disk"
	Filesystem string // vfat, btrfs, luks2
	TypeCode   string // sgdisk type code
}

// PartitionPlan is the full computed layout for a disk of a given size.
type PartitionPlan struct {
	DiskSizeMiB int
	Partitions  []Partition
}

// PlanPartitions computes the GPT layout for a USB of diskSizeMiB.
// Returns an error if the disk is too small to hold ESP + root + a
// minimum persistence area.
func PlanPartitions(diskSizeMiB int) (*PartitionPlan, error) {
	if diskSizeMiB < MinDiskMiB {
		return nil, fmt.Errorf("disk too small: %d MiB, need at least %d MiB", diskSizeMiB, MinDiskMiB)
	}

	esp := Partition{
		Number:     1,
		Name:       "EFI",
		StartMiB:   1, // leave 1 MiB for GPT/alignment
		SizeMiB:    espSizeMiB,
		Filesystem: "vfat",
		TypeCode:   "ef00",
	}

	root := Partition{
		Number:     2,
		Name:       "enigma-root",
		StartMiB:   esp.StartMiB + esp.SizeMiB,
		SizeMiB:    rootSizeMiB,
		Filesystem: "btrfs",
		TypeCode:   "8300",
	}

	persist := Partition{
		Number:     3,
		Name:       "enigma-persist",
		StartMiB:   root.StartMiB + root.SizeMiB,
		SizeMiB:    -1, // remainder of disk
		Filesystem: "luks2",
		TypeCode:   "8309", // Linux LUKS
	}

	return &PartitionPlan{
		DiskSizeMiB: diskSizeMiB,
		Partitions:  []Partition{esp, root, persist},
	}, nil
}

// PersistenceSizeMiB returns the size the LUKS2 persistence partition will
// occupy given the plan's disk size (disk minus ESP, root, and alignment).
func (p *PartitionPlan) PersistenceSizeMiB() int {
	last := p.Partitions[len(p.Partitions)-1]
	return p.DiskSizeMiB - last.StartMiB
}
