package game

import "fmt"

// SetupSteam configures Steam + Proton-GE for gaming (Phase 7 - SPEC §7)
func SetupSteam() error {
	fmt.Println("TODO Phase 7: Setup Steam + Proton-GE")
	fmt.Println("  - Install multilib repo + steam package")
	fmt.Println("  - Configure ~/Games/proton-ge/ symlink")
	fmt.Println("  - Set PROTON_USE_WINED3D env for AMD GPU support")
	fmt.Println("  - Fallback: bottles + wine-staging for non-Steam games")
	return nil
}

// SetupContainers configures Podman rootless containers (Phase 7 - SPEC §7)
func SetupContainers() error {
	fmt.Println("TODO Phase 7: Setup Podman rootless containers")
	fmt.Println("  - Enable user namespaces in /etc/subuid & /etc/subgid")
	fmt.Println("  - Configure ~/.config/containers/containers.conf")
	fmt.Println("  - Integrate with enigma dev (port allocation)")
	return nil
}
