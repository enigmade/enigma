package game

import "testing"

const sampleReleaseJSON = `{
  "tag_name": "GE-Proton9-20",
  "assets": [
    {
      "name": "GE-Proton9-20.sha512sum",
      "browser_download_url": "https://github.com/GloriousEggroll/proton-ge-custom/releases/download/GE-Proton9-20/GE-Proton9-20.sha512sum"
    },
    {
      "name": "GE-Proton9-20.tar.gz",
      "browser_download_url": "https://github.com/GloriousEggroll/proton-ge-custom/releases/download/GE-Proton9-20/GE-Proton9-20.tar.gz"
    }
  ]
}`

func TestParseLatestProtonGERelease(t *testing.T) {
	rel, err := ParseLatestProtonGERelease([]byte(sampleReleaseJSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rel.Version != "GE-Proton9-20" {
		t.Errorf("expected version GE-Proton9-20, got %s", rel.Version)
	}
	if rel.DownloadURL != "https://github.com/GloriousEggroll/proton-ge-custom/releases/download/GE-Proton9-20/GE-Proton9-20.tar.gz" {
		t.Errorf("unexpected download URL: %s", rel.DownloadURL)
	}
}

func TestParseLatestProtonGEReleaseMissingTarball(t *testing.T) {
	body := `{"tag_name": "GE-Proton9-20", "assets": [{"name": "notes.txt", "browser_download_url": "https://example.com/notes.txt"}]}`
	_, err := ParseLatestProtonGERelease([]byte(body))
	if err == nil {
		t.Fatal("expected error when no .tar.gz asset present")
	}
}

func TestParseLatestProtonGEReleaseInvalidJSON(t *testing.T) {
	_, err := ParseLatestProtonGERelease([]byte("not json"))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParseLatestProtonGEReleaseMissingTag(t *testing.T) {
	body := `{"assets": [{"name": "x.tar.gz", "browser_download_url": "https://example.com/x.tar.gz"}]}`
	_, err := ParseLatestProtonGERelease([]byte(body))
	if err == nil {
		t.Fatal("expected error when tag_name is missing")
	}
}

func TestInstallDir(t *testing.T) {
	got := InstallDir("/home/user/.steam/steam", "GE-Proton9-20")
	want := "/home/user/.steam/steam/compatibilitytools.d/GE-Proton9-20"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
