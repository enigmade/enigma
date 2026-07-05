package portability

import "fmt"

// SgdiskCommands renders the sgdisk invocation that writes the computed GPT
// layout to `device` (e.g. /dev/sda). Returned as an argv slice per command
// so the caller can exec them directly without shell quoting concerns.
func (p *PartitionPlan) SgdiskCommands(device string) [][]string {
	var cmds [][]string
	// Start clean.
	cmds = append(cmds, []string{"sgdisk", "--zap-all", device})

	for _, part := range p.Partitions {
		var sizeSpec string
		if part.SizeMiB < 0 {
			sizeSpec = "0" // sgdisk: 0 means "to end of disk"
		} else {
			sizeSpec = fmt.Sprintf("+%dM", part.SizeMiB)
		}
		startSpec := fmt.Sprintf("%dM", part.StartMiB)

		cmds = append(cmds, []string{
			"sgdisk",
			fmt.Sprintf("--new=%d:%s:%s", part.Number, startSpec, sizeSpec),
			fmt.Sprintf("--typecode=%d:%s", part.Number, part.TypeCode),
			fmt.Sprintf("--change-name=%d:%s", part.Number, part.Name),
			device,
		})
	}
	return cmds
}

// LuksFormatCommand renders the cryptsetup command to LUKS2-format the
// persistence partition. keyFile is read for the passphrase (SPEC §8: TPM
// enrollment happens separately on first boot).
func LuksFormatCommand(partitionDevice, keyFile string) []string {
	return []string{
		"cryptsetup", "luksFormat",
		"--type", "luks2",
		"--key-file", keyFile,
		"--batch-mode",
		partitionDevice,
	}
}

// LuksOpenCommand renders the command to unlock the persistence partition
// to a mapper name (e.g. enigma-persist).
func LuksOpenCommand(partitionDevice, mapperName, keyFile string) []string {
	return []string{
		"cryptsetup", "open",
		"--key-file", keyFile,
		partitionDevice, mapperName,
	}
}
