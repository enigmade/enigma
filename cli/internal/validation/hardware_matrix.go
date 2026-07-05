package validation

import "fmt"

// ValidateHardwareMatrix runs real-hardware smoke tests (Phase 9 - SPEC §9)
//
// Test matrix:
//   GPUs: NVIDIA (Ampere, Turing), AMD (RDNA2, RDNA1), Intel (Iris, UHD)
//   Boot: UEFI + BIOS, SSD/HDD, 8GB/16GB/32GB RAM
//   Features: KDE Plasma 6 animations (60fps), AI services, Steam/Proton, Wine apps
//
// Benchmarks:
//   - Cold boot: <15s (SPEC §9 target)
//   - AI service ready: <30s (Ollama + torch imports)
//   - GUI ready: <45s (KDE Plasma 6 login, dock visible)
//   - Rollback: <5s (snapshot switch + reboot)
// ValidateHardwareMatrix evaluates a run's benchmark output against the
// SPEC §9 thresholds and reports pass/fail. benchmarkOutput is the
// "metric=seconds" text produced by the on-device measurement harness.
// Returns an error if any metric exceeded its threshold, so this can gate
// a release in CI or a smoke-test script.
func ValidateHardwareMatrix(benchmarkOutput string) error {
	results, err := ParseBenchmarks(benchmarkOutput)
	if err != nil {
		return err
	}

	for _, r := range results {
		status := "PASS"
		if !r.Pass {
			status = "FAIL"
		}
		fmt.Printf("  [%s] %-14s %.1fs (threshold %.1fs)\n", status, r.Name, r.ValueSec, r.Threshold)
	}

	if !AllPass(results) {
		return fmt.Errorf("benchmark failures: %v", Failures(results))
	}
	return nil
}
