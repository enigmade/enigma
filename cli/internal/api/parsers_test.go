package api

import "testing"

func TestParseSystemctlServices(t *testing.T) {
	output := `enigma-ollama.service loaded active running Enigma Ollama LLM Service
enigma-indexd.service loaded active running Enigma Search Indexer
dbus.service loaded active running D-Bus System Message Bus
enigma-comfy.service loaded inactive dead Enigma ComfyUI Service`

	services := ParseSystemctlServices(output)
	if len(services) != 3 {
		t.Fatalf("expected 3 enigma services, got %d: %+v", len(services), services)
	}
	if services[0].Name != "enigma-ollama" || services[0].Status != "running" {
		t.Errorf("unexpected first service: %+v", services[0])
	}
	if services[2].Name != "enigma-comfy" || services[2].Status != "dead" {
		t.Errorf("unexpected third service: %+v", services[2])
	}
}

func TestParseSystemctlServicesSkipsNonEnigma(t *testing.T) {
	output := "sshd.service loaded active running OpenSSH\nNetworkManager.service loaded active running Network Manager"
	services := ParseSystemctlServices(output)
	if len(services) != 0 {
		t.Errorf("expected no enigma services, got %+v", services)
	}
}

func TestParseMiseRuntimes(t *testing.T) {
	output := `node   22.3.0   ~/.config/mise/config.toml  (active)
python 3.12.4
go     1.26.0   (active)`

	runtimes := ParseMiseRuntimes(output)
	if len(runtimes) != 3 {
		t.Fatalf("expected 3 runtimes, got %d", len(runtimes))
	}
	if runtimes[0].Name != "node" || runtimes[0].Version != "22.3.0" || !runtimes[0].Active {
		t.Errorf("unexpected node runtime: %+v", runtimes[0])
	}
	if runtimes[1].Active {
		t.Errorf("python should not be active: %+v", runtimes[1])
	}
}

func TestParseOllamaModels(t *testing.T) {
	output := `NAME            ID           SIZE     MODIFIED
llama3.1:8b     abc123       4.7 GB   2 days ago
qwen2.5:14b     def456       9.0 GB   1 week ago`

	models := ParseOllamaModels(output)
	if len(models) != 2 {
		t.Fatalf("expected 2 models, got %d: %+v", len(models), models)
	}
	if models[0].Name != "llama3.1:8b" || models[0].SizeGB != 4.7 || models[0].Backend != "ollama" {
		t.Errorf("unexpected first model: %+v", models[0])
	}
}

func TestParseSizeGBUnits(t *testing.T) {
	cases := []struct {
		value, unit string
		want        float64
	}{
		{"4.7", "GB", 4.7},
		{"1", "TB", 1024},
		{"512", "MB", 0.5},
		{"garbage", "GB", 0},
	}
	for _, c := range cases {
		got := parseSizeGB(c.value, c.unit)
		if got != c.want {
			t.Errorf("parseSizeGB(%q,%q) = %v, want %v", c.value, c.unit, got, c.want)
		}
	}
}
