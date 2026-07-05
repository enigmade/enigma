package game

import "testing"

const sampleSubuid = `# comment line
alice:100000:65536
bob:165536:1000
`

func TestParseSubIDFile(t *testing.T) {
	ranges, err := ParseSubIDFile(sampleSubuid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ranges) != 2 {
		t.Fatalf("expected 2 ranges, got %d", len(ranges))
	}
	if ranges[0].Username != "alice" || ranges[0].Start != 100000 || ranges[0].Count != 65536 {
		t.Errorf("unexpected first range: %+v", ranges[0])
	}
}

func TestParseSubIDFileInvalidLine(t *testing.T) {
	_, err := ParseSubIDFile("not:a:valid:line:here")
	if err == nil {
		t.Fatal("expected error for malformed line")
	}
}

func TestParseSubIDFileInvalidNumber(t *testing.T) {
	_, err := ParseSubIDFile("alice:notanumber:65536")
	if err == nil {
		t.Fatal("expected error for non-numeric start")
	}
}

func TestHasSufficientSubIDRange(t *testing.T) {
	ranges, _ := ParseSubIDFile(sampleSubuid)

	if !HasSufficientSubIDRange(ranges, "alice") {
		t.Error("expected alice to have sufficient range")
	}
	if HasSufficientSubIDRange(ranges, "bob") {
		t.Error("expected bob to lack sufficient range (only 1000)")
	}
	if HasSufficientSubIDRange(ranges, "nobody") {
		t.Error("expected unknown user to lack sufficient range")
	}
}
