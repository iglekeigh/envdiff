package defaults

import (
	"testing"
)

func TestApply_AddsOnlyMissingKeys(t *testing.T) {
	target := map[string]string{"A": "1"}
	defs := map[string]string{"A": "overridden", "B": "2", "C": "3"}

	res, err := Apply(target, defs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["A"] != "1" {
		t.Errorf("expected A=1, got %s", res.Env["A"])
	}
	if res.Env["B"] != "2" {
		t.Errorf("expected B=2, got %s", res.Env["B"])
	}
	if res.Env["C"] != "3" {
		t.Errorf("expected C=3, got %s", res.Env["C"])
	}
}

func TestApply_ReportsAppliedAndSkipped(t *testing.T) {
	target := map[string]string{"EXISTING": "yes"}
	defs := map[string]string{"EXISTING": "no", "NEW": "added"}

	res, err := Apply(target, defs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Applied) != 1 || res.Applied[0] != "NEW" {
		t.Errorf("expected Applied=[NEW], got %v", res.Applied)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "EXISTING" {
		t.Errorf("expected Skipped=[EXISTING], got %v", res.Skipped)
	}
}

func TestApply_DoesNotMutateInputs(t *testing.T) {
	target := map[string]string{"A": "1"}
	defs := map[string]string{"B": "2"}

	_, err := Apply(target, defs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := target["B"]; ok {
		t.Error("Apply mutated the target map")
	}
}

func TestApply_EmptyDefaults_ReturnsTargetUnchanged(t *testing.T) {
	target := map[string]string{"X": "10"}
	defs := map[string]string{}

	res, err := Apply(target, defs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 1 || res.Env["X"] != "10" {
		t.Errorf("unexpected env: %v", res.Env)
	}
	if len(res.Applied) != 0 {
		t.Errorf("expected no applied keys, got %v", res.Applied)
	}
}

func TestApply_NilTarget_ReturnsError(t *testing.T) {
	_, err := Apply(nil, map[string]string{})
	if err == nil {
		t.Error("expected error for nil target, got nil")
	}
}

func TestApply_NilDefaults_ReturnsError(t *testing.T) {
	_, err := Apply(map[string]string{}, nil)
	if err == nil {
		t.Error("expected error for nil defaults, got nil")
	}
}

func TestApply_ResultIsSorted(t *testing.T) {
	target := map[string]string{}
	defs := map[string]string{"Z": "z", "A": "a", "M": "m"}

	res, err := Apply(target, defs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"A", "M", "Z"}
	for i, k := range res.Applied {
		if k != expected[i] {
			t.Errorf("Applied[%d] = %q, want %q", i, k, expected[i])
		}
	}
}
