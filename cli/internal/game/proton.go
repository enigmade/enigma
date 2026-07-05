package game

import (
	"encoding/json"
	"fmt"
	"strings"
)

// GitHubRelease models the subset of the GitHub Releases API response
// needed to locate the Proton-GE tarball asset.
type GitHubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []GitHubAsset `json:"assets"`
}

type GitHubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// ProtonGERelease is the resolved release info used to install Proton-GE.
type ProtonGERelease struct {
	Version     string
	DownloadURL string
}

// ParseLatestProtonGERelease parses a GitHub Releases API JSON response
// (from GET /repos/GloriousEggroll/proton-ge-custom/releases/latest) and
// picks out the .tar.gz asset (SPEC §7: Proton-GE for Steam Play).
func ParseLatestProtonGERelease(body []byte) (*ProtonGERelease, error) {
	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, fmt.Errorf("parse release JSON: %w", err)
	}

	if release.TagName == "" {
		return nil, fmt.Errorf("release has no tag_name")
	}

	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.Name, ".tar.gz") {
			return &ProtonGERelease{
				Version:     release.TagName,
				DownloadURL: asset.BrowserDownloadURL,
			}, nil
		}
	}

	return nil, fmt.Errorf("no .tar.gz asset found in release %s", release.TagName)
}

// InstallDir returns the destination directory for a given Proton-GE
// version under the Steam compatibilitytools.d layout (SPEC §7).
func InstallDir(steamRoot, version string) string {
	return steamRoot + "/compatibilitytools.d/" + version
}
