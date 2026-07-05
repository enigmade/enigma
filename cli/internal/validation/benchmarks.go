package validation

import (
	"fmt"
	"strconv"
	"strings"
)

// Benchmark thresholds from SPEC §9 (seconds).
const (
	ThresholdColdBootSec    = 15.0
	ThresholdServiceReadySec = 30.0
	ThresholdGUIReadySec    = 45.0
	ThresholdRollbackSec    = 5.0
)

// BenchmarkResult is a single measured metric plus its pass/fail verdict.
type BenchmarkResult struct {
	Name      string
	ValueSec  float64
	Threshold float64
	Pass      bool
}

// thresholdFor maps a metric name to its SPEC §9 threshold.
func thresholdFor(metric string) (float64, bool) {
	switch metric {
	case "cold_boot":
		return ThresholdColdBootSec, true
	case "service_ready":
		return ThresholdServiceReadySec, true
	case "gui_ready":
		return ThresholdGUIReadySec, true
	case "rollback":
		return ThresholdRollbackSec, true
	default:
		return 0, false
	}
}

// ParseBenchmarks parses "metric=value" lines (value in seconds) and
// evaluates each against its SPEC §9 threshold. Unknown metrics are an
// error so a typo can't silently pass validation.
//
// Sample input:
//   cold_boot=12.4
//   service_ready=22.1
//   gui_ready=38.0
//   rollback=3.2
func ParseBenchmarks(input string) ([]BenchmarkResult, error) {
	var results []BenchmarkResult
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed benchmark line: %q", line)
		}

		metric := strings.TrimSpace(parts[0])
		value, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value for %s: %w", metric, err)
		}

		threshold, ok := thresholdFor(metric)
		if !ok {
			return nil, fmt.Errorf("unknown benchmark metric: %q", metric)
		}

		results = append(results, BenchmarkResult{
			Name:      metric,
			ValueSec:  value,
			Threshold: threshold,
			Pass:      value <= threshold,
		})
	}
	return results, nil
}

// AllPass reports whether every result met its threshold.
func AllPass(results []BenchmarkResult) bool {
	for _, r := range results {
		if !r.Pass {
			return false
		}
	}
	return true
}

// Failures returns the metric names that exceeded their threshold.
func Failures(results []BenchmarkResult) []string {
	var failed []string
	for _, r := range results {
		if !r.Pass {
			failed = append(failed, r.Name)
		}
	}
	return failed
}
