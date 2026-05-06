package merge_test

import (
	"testing"

	"github.com/user/envdiff/internal/merge"
)

func TestMerge_NoConflicts_CombinesAllKeys(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"C": "3", "D": "4"}

	res, err := merge.Merge([]map[string]string{a, b}, nil, merge.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 4 {
		t.Errorf("expected 4 keys, got %d", len(res.Env))
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(res.Conflicts))
	}
}

func TestMerge_StrategyFirst_KeepsFirstValue(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}

	res, err := merge.Merge([]map[string]string{a, b}, nil, merge.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "first" {
		t.Errorf("expected 'first', got %q", res.Env["KEY"])
	}
	if len(res.Conflicts["KEY"]) != 2 {
		t.Errorf("expected conflict entry with 2 sources")
	}
}

func TestMerge_StrategyLast_KeepsLastValue(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}

	res, err := merge.Merge([]map[string]string{a, b}, nil, merge.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "second" {
		t.Errorf("expected 'second', got %q", res.Env["KEY"])
	}
}

func TestMerge_StrategyError_ReturnsError(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}

	_, err := merge.Merge([]map[string]string{a, b}, nil, merge.StrategyError)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMerge_Labels_UsedInConflicts(t *testing.T) {
	a := map[string]string{"X": "1"}
	b := map[string]string{"X": "2"}
	labels := []string{".env.dev", ".env.prod"}

	res, err := merge.Merge([]map[string]string{a, b}, labels, merge.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	srcs := res.Conflicts["X"]
	if len(srcs) != 2 || srcs[0] != ".env.dev" || srcs[1] != ".env.prod" {
		t.Errorf("unexpected conflict sources: %v", srcs)
	}
}

func TestMerge_EmptySources_ReturnsEmptyEnv(t *testing.T) {
	res, err := merge.Merge([]map[string]string{}, nil, merge.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 0 {
		t.Errorf("expected empty env, got %d keys", len(res.Env))
	}
}

func TestMerge_StrategyLast_ConflictsStillRecorded(t *testing.T) {
	// Even when using StrategyLast, conflicts should still be recorded
	// so callers can inspect which keys had differing values across sources.
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}

	res, err := merge.Merge([]map[string]string{a, b}, nil, merge.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts["KEY"]) != 2 {
		t.Errorf("expected conflict entry with 2 sources, got %d", len(res.Conflicts["KEY"]))
	}
}
