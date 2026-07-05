package game

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const protonGELatestReleaseURL = "https://api.github.com/repos/GloriousEggroll/proton-ge-custom/releases/latest"

// SetupSteam fetches the latest Proton-GE release metadata and reports
// where it would be installed under the user's Steam directory
// (SPEC §7). Actual package install (multilib repo + steam) is left to
// the pacman-driven install hook; this handles the Proton-GE side.
func SetupSteam(steamRoot string) (*ProtonGERelease, string, error) {
	resp, err := http.Get(protonGELatestReleaseURL)
	if err != nil {
		return nil, "", fmt.Errorf("fetch Proton-GE release info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("read Proton-GE response: %w", err)
	}

	release, err := ParseLatestProtonGERelease(body)
	if err != nil {
		return nil, "", err
	}

	return release, InstallDir(steamRoot, release.Version), nil
}

// SetupContainers writes the rootless Podman config and verifies the
// current user has a sufficient subuid/subgid range for user namespaces
// (SPEC §7).
func SetupContainers(configDir, subuidPath, subgidPath, username string) (string, error) {
	path, err := WriteContainersConf(configDir)
	if err != nil {
		return "", err
	}

	subuid, err := os.ReadFile(subuidPath)
	if err != nil {
		return path, fmt.Errorf("read %s: %w", subuidPath, err)
	}
	subuidRanges, err := ParseSubIDFile(string(subuid))
	if err != nil {
		return path, fmt.Errorf("parse %s: %w", subuidPath, err)
	}
	if !HasSufficientSubIDRange(subuidRanges, username) {
		return path, fmt.Errorf("user %s lacks a subuid range >= %d in %s", username, MinRootlessSubIDCount, subuidPath)
	}

	subgid, err := os.ReadFile(subgidPath)
	if err != nil {
		return path, fmt.Errorf("read %s: %w", subgidPath, err)
	}
	subgidRanges, err := ParseSubIDFile(string(subgid))
	if err != nil {
		return path, fmt.Errorf("parse %s: %w", subgidPath, err)
	}
	if !HasSufficientSubIDRange(subgidRanges, username) {
		return path, fmt.Errorf("user %s lacks a subgid range >= %d in %s", username, MinRootlessSubIDCount, subgidPath)
	}

	return path, nil
}
