package game

import (
	"fmt"
	"strings"
)

// RequiredWindowsRedistributables lists the winetricks verbs Wine Tier 1
// needs so legacy Windows apps and games have their DirectX + Visual C++
// runtime dependencies satisfied out of the box (SPEC §7).
//
//   - d3dcompiler_47: DirectX shader compiler, required by most DX9-DX11 titles
//   - d3dx9 / d3dx11_43: legacy DirectX helper libraries many older games link against
//   - vcrun2019: Visual C++ 2015-2022 redistributable (covers the vast majority of apps)
//   - dxvk: DirectX-to-Vulkan translation layer (installed separately via package, not winetricks)
func RequiredWindowsRedistributables() []string {
	return []string{
		"d3dcompiler_47",
		"d3dx9",
		"d3dx11_43",
		"vcrun2019",
	}
}

// ParseWinetricksLog parses a prefix's winetricks.log, which records one
// installed verb per line, so we can tell what's already present without
// re-running installs on every `enigma game setup`.
func ParseWinetricksLog(content string) []string {
	var verbs []string
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		verbs = append(verbs, line)
	}
	return verbs
}

// MissingRedistributables returns the subset of `required` not present in
// `installed`, preserving the order of `required`.
func MissingRedistributables(installed, required []string) []string {
	installedSet := make(map[string]bool, len(installed))
	for _, v := range installed {
		installedSet[v] = true
	}

	var missing []string
	for _, v := range required {
		if !installedSet[v] {
			missing = append(missing, v)
		}
	}
	return missing
}

// SetupWineTier1 reports which DirectX/VC++ redistributable verbs still
// need installing into the given Wine prefix, based on its winetricks.log.
// Actual invocation of `winetricks <verb>` is left to the caller so this
// stays testable without a Wine installation present.
func SetupWineTier1(winetricksLogContent string) []string {
	installed := ParseWinetricksLog(winetricksLogContent)
	required := RequiredWindowsRedistributables()
	missing := MissingRedistributables(installed, required)

	if len(missing) == 0 {
		fmt.Println("✓ All required DirectX/VC++ redistributables already installed")
	} else {
		fmt.Printf("Missing redistributables to install via winetricks: %s\n", strings.Join(missing, ", "))
	}

	return missing
}
