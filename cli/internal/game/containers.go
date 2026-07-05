package game

import (
	"fmt"
	"os"
	"path/filepath"
)

// GenerateContainersConf renders the Podman rootless containers.conf
// content (SPEC §7). Kept as a pure string function so it can be
// unit-tested without touching the filesystem.
func GenerateContainersConf() string {
	return `[containers]
netns = "private"
userns = "auto"

[engine]
runtime = "crun"

[network]
network_backend = "netavark"
`
}

// WriteContainersConf writes the rendered config to
// ~/.config/containers/containers.conf, creating the directory if needed.
func WriteContainersConf(configDir string) (string, error) {
	dir := filepath.Join(configDir, "containers")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create containers config dir: %w", err)
	}

	path := filepath.Join(dir, "containers.conf")
	if err := os.WriteFile(path, []byte(GenerateContainersConf()), 0o644); err != nil {
		return "", fmt.Errorf("write containers.conf: %w", err)
	}

	return path, nil
}
