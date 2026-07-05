package game

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateContainersConfContainsRootlessSettings(t *testing.T) {
	conf := GenerateContainersConf()
	for _, want := range []string{`userns = "auto"`, `netns = "private"`, `runtime = "crun"`} {
		if !strings.Contains(conf, want) {
			t.Errorf("expected containers.conf to contain %q, got:\n%s", want, conf)
		}
	}
}

func TestWriteContainersConf(t *testing.T) {
	tmp := t.TempDir()

	path, err := WriteContainersConf(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantPath := filepath.Join(tmp, "containers", "containers.conf")
	if path != wantPath {
		t.Errorf("got path %s, want %s", path, wantPath)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}
	if string(data) != GenerateContainersConf() {
		t.Error("written file content does not match GenerateContainersConf()")
	}
}
