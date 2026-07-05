package api

import (
	"strconv"
	"strings"
)

// ParseSystemctlServices parses the output of
// `systemctl --user list-units --type=service --no-legend --plain`
// into Service records. Only enigma-managed units are kept.
//
// Sample line:
//   enigma-ollama.service loaded active running Enigma Ollama LLM Service
func ParseSystemctlServices(output string) []Service {
	var services []Service
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		unit := fields[0]
		if !strings.HasPrefix(unit, "enigma-") || !strings.HasSuffix(unit, ".service") {
			continue
		}
		name := strings.TrimSuffix(unit, ".service")
		// fields: [UNIT LOAD ACTIVE SUB ...]; SUB is the fine-grained state.
		status := fields[3]
		services = append(services, Service{Name: name, Status: status})
	}
	return services
}

// ParseMiseRuntimes parses `mise ls --json`-free plain output of the form
// produced by `mise ls`:
//   node   22.3.0   ~/.config/mise/config.toml  (active)
//   python 3.12.4
func ParseMiseRuntimes(output string) []Runtime {
	var runtimes []Runtime
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		rt := Runtime{Name: fields[0], Version: fields[1]}
		if strings.Contains(line, "(active)") {
			rt.Active = true
		}
		runtimes = append(runtimes, rt)
	}
	return runtimes
}

// ParseOllamaModels parses the tabular output of `ollama list`:
//   NAME            ID           SIZE     MODIFIED
//   llama3.1:8b     abc123       4.7 GB   2 days ago
func ParseOllamaModels(output string) []Model {
	var models []Model
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		// Skip the header row.
		if i == 0 && strings.EqualFold(fields[0], "NAME") {
			continue
		}
		name := fields[0]
		// SIZE is two tokens: value + unit, e.g. "4.7 GB".
		size := parseSizeGB(fields[2], fields[3])
		models = append(models, Model{Name: name, SizeGB: size, Backend: "ollama"})
	}
	return models
}

// parseSizeGB converts a value+unit pair (e.g. "4.7", "GB") to gibibytes.
func parseSizeGB(value, unit string) float64 {
	n, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	switch strings.ToUpper(unit) {
	case "TB":
		return n * 1024
	case "GB":
		return n
	case "MB":
		return n / 1024
	case "KB":
		return n / (1024 * 1024)
	default:
		return n
	}
}
