package portability

import (
	"strings"
	"testing"
)

func TestPlanPartitionsLayout(t *testing.T) {
	// 64 GiB stick.
	plan, err := PlanPartitions(64 * GiB)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(plan.Partitions) != 3 {
		t.Fatalf("expected 3 partitions, got %d", len(plan.Partitions))
	}

	esp := plan.Partitions[0]
	if esp.StartMiB != 1 || esp.SizeMiB != 512 || esp.Filesystem != "vfat" {
		t.Errorf("unexpected ESP: %+v", esp)
	}

	root := plan.Partitions[1]
	if root.StartMiB != 513 || root.Filesystem != "btrfs" {
		t.Errorf("unexpected root: %+v", root)
	}

	persist := plan.Partitions[2]
	if persist.SizeMiB != -1 || persist.Filesystem != "luks2" {
		t.Errorf("unexpected persist: %+v", persist)
	}
}

func TestPlanPartitionsTooSmall(t *testing.T) {
	_, err := PlanPartitions(1 * GiB)
	if err == nil {
		t.Fatal("expected error for undersized disk")
	}
}

func TestPlanPartitionsExactMinimum(t *testing.T) {
	_, err := PlanPartitions(MinDiskMiB)
	if err != nil {
		t.Fatalf("minimum-size disk should succeed, got: %v", err)
	}
}

func TestPersistenceSizeMiB(t *testing.T) {
	plan, _ := PlanPartitions(64 * GiB)
	// Persistence = disk - (1 MiB align + 512 ESP + 20 GiB root)
	want := 64*GiB - (1 + 512 + 20*GiB)
	if got := plan.PersistenceSizeMiB(); got != want {
		t.Errorf("got %d MiB, want %d MiB", got, want)
	}
}

func TestSgdiskCommands(t *testing.T) {
	plan, _ := PlanPartitions(64 * GiB)
	cmds := plan.SgdiskCommands("/dev/sda")

	// zap-all + one per partition.
	if len(cmds) != 4 {
		t.Fatalf("expected 4 commands, got %d", len(cmds))
	}
	if cmds[0][1] != "--zap-all" {
		t.Errorf("first command should zap: %v", cmds[0])
	}

	// The persistence partition uses size 0 (to end of disk).
	last := cmds[3]
	joined := strings.Join(last, " ")
	if !strings.Contains(joined, "--new=3:") || !strings.Contains(joined, ":0") {
		t.Errorf("persistence partition should extend to end of disk: %v", last)
	}
	if !strings.Contains(joined, "8309") {
		t.Errorf("persistence partition should be LUKS type 8309: %v", last)
	}
}

func TestLuksCommands(t *testing.T) {
	format := LuksFormatCommand("/dev/sda3", "/run/enigma/key")
	joined := strings.Join(format, " ")
	if !strings.Contains(joined, "luksFormat") || !strings.Contains(joined, "luks2") {
		t.Errorf("unexpected luksFormat command: %v", format)
	}

	open := LuksOpenCommand("/dev/sda3", "enigma-persist", "/run/enigma/key")
	if open[len(open)-1] != "enigma-persist" {
		t.Errorf("luksOpen should target mapper name: %v", open)
	}
}
