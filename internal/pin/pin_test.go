package pin

import (
	"testing"
)

func TestApply_PinsAllKeysWhenNoneSpecified(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	res, err := Apply(env, nil, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Pinned["A"] != "1" || res.Pinned["B"] != "2" {
		t.Errorf("expected all keys pinned, got %v", res.Pinned)
	}
}

func TestApply_PinsSpecificKeys(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2", "C": "3"}
	res, err := Apply(env, nil, Options{Keys: []string{"A", "C"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Pinned["B"]; ok {
		t.Error("key B should not be pinned")
	}
	if res.Pinned["A"] != "1" || res.Pinned["C"] != "3" {
		t.Errorf("unexpected pinned values: %v", res.Pinned)
	}
}

func TestApply_SkipsMissingKeys(t *testing.T) {
	env := map[string]string{"A": "1"}
	res, err := Apply(env, nil, Options{Keys: []string{"A", "MISSING"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "MISSING" {
		t.Errorf("expected MISSING in skipped, got %v", res.Skipped)
	}
}

func TestApply_Release_RemovesPinnedKey(t *testing.T) {
	existing := map[string]string{"A": "1", "B": "2"}
	env := map[string]string{"A": "1", "B": "2"}
	res, err := Apply(env, existing, Options{Keys: []string{"A"}, Release: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Pinned["A"]; ok {
		t.Error("key A should have been released")
	}
	if res.Pinned["B"] != "2" {
		t.Error("key B should remain pinned")
	}
	if len(res.Released) != 1 || res.Released[0] != "A" {
		t.Errorf("expected A in released, got %v", res.Released)
	}
}

func TestApply_NilEnv_ReturnsError(t *testing.T) {
	_, err := Apply(nil, nil, Options{})
	if err == nil {
		t.Error("expected error for nil env")
	}
}

func TestApply_DoesNotMutateExisting(t *testing.T) {
	existing := map[string]string{"A": "old"}
	env := map[string]string{"A": "new"}
	_, err := Apply(env, existing, Options{Keys: []string{"A"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if existing["A"] != "old" {
		t.Error("Apply must not mutate existing map")
	}
}

func TestCheck_DetectsDrift(t *testing.T) {
	env := map[string]string{"A": "changed", "B": "same"}
	pinned := map[string]string{"A": "original", "B": "same"}
	drifted := Check(env, pinned)
	if len(drifted) != 1 || drifted[0] != "A" {
		t.Errorf("expected [A] drifted, got %v", drifted)
	}
}

func TestCheck_MissingKeyCountsAsDrift(t *testing.T) {
	env := map[string]string{}
	pinned := map[string]string{"A": "1"}
	drifted := Check(env, pinned)
	if len(drifted) != 1 || drifted[0] != "A" {
		t.Errorf("expected [A] as drifted, got %v", drifted)
	}
}
