package validation

import (
	"reflect"
	"testing"
)

func TestParseBenchmarksAllPass(t *testing.T) {
	input := `cold_boot=12.4
service_ready=22.1
gui_ready=38.0
rollback=3.2`

	results, err := ParseBenchmarks(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}
	if !AllPass(results) {
		t.Errorf("expected all to pass, failures: %v", Failures(results))
	}
}

func TestParseBenchmarksWithFailure(t *testing.T) {
	input := `cold_boot=18.9
rollback=3.0`

	results, err := ParseBenchmarks(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if AllPass(results) {
		t.Error("expected cold_boot to fail its 15s threshold")
	}
	if got := Failures(results); !reflect.DeepEqual(got, []string{"cold_boot"}) {
		t.Errorf("expected [cold_boot] to fail, got %v", got)
	}
}

func TestParseBenchmarksBoundaryPasses(t *testing.T) {
	// Exactly at threshold should pass (<=).
	results, err := ParseBenchmarks("cold_boot=15.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Pass {
		t.Error("value exactly at threshold should pass")
	}
}

func TestParseBenchmarksUnknownMetric(t *testing.T) {
	_, err := ParseBenchmarks("teleport_speed=9000")
	if err == nil {
		t.Fatal("expected error for unknown metric")
	}
}

func TestParseBenchmarksMalformed(t *testing.T) {
	_, err := ParseBenchmarks("cold_boot 12.4")
	if err == nil {
		t.Fatal("expected error for missing =")
	}
}

func TestParseBenchmarksSkipsCommentsAndBlanks(t *testing.T) {
	input := "# boot metrics\n\ncold_boot=10.0\n"
	results, err := ParseBenchmarks(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}
