package game

import (
	"reflect"
	"testing"
)

func TestParseWinetricksLog(t *testing.T) {
	log := "d3dcompiler_47\nvcrun2019\n\ncorefonts\n"
	got := ParseWinetricksLog(log)
	want := []string{"d3dcompiler_47", "vcrun2019", "corefonts"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestMissingRedistributablesAllPresent(t *testing.T) {
	installed := []string{"d3dcompiler_47", "d3dx9", "d3dx11_43", "vcrun2019", "corefonts"}
	missing := MissingRedistributables(installed, RequiredWindowsRedistributables())
	if len(missing) != 0 {
		t.Errorf("expected no missing redistributables, got %v", missing)
	}
}

func TestMissingRedistributablesNonePresent(t *testing.T) {
	missing := MissingRedistributables([]string{}, RequiredWindowsRedistributables())
	if !reflect.DeepEqual(missing, RequiredWindowsRedistributables()) {
		t.Errorf("expected all required to be missing, got %v", missing)
	}
}

func TestMissingRedistributablesPartial(t *testing.T) {
	installed := []string{"d3dcompiler_47", "vcrun2019"}
	missing := MissingRedistributables(installed, RequiredWindowsRedistributables())
	want := []string{"d3dx9", "d3dx11_43"}
	if !reflect.DeepEqual(missing, want) {
		t.Errorf("got %v, want %v", missing, want)
	}
}

func TestSetupWineTier1ReturnsMissing(t *testing.T) {
	missing := SetupWineTier1("d3dcompiler_47\nd3dx9\n")
	want := []string{"d3dx11_43", "vcrun2019"}
	if !reflect.DeepEqual(missing, want) {
		t.Errorf("got %v, want %v", missing, want)
	}
}
