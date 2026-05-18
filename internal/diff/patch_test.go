package diff

import (
	"testing"
)

func buildResult(entries []Entry) Result {
	return Result{Entries: entries}
}

func TestPatch_AddOnly_AddsMissingKeys(t *testing.T) {
	base := map[string]string{"A": "1"}
	result := buildResult([]Entry{
		{Key: "B", Status: StatusMissing, OtherValue: "2"},
		{Key: "A", Status: StatusChanged, BaseValue: "1", OtherValue: "99"},
	})

	pr, err := Patch(base, result, PatchAddOnly)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if base["B"] != "2" {
		t.Errorf("expected B=2, got %q", base["B"])
	}
	if base["A"] != "1" {
		t.Errorf("expected A unchanged, got %q", base["A"])
	}
	if len(pr.Added) != 1 || pr.Added[0] != "B" {
		t.Errorf("expected Added=[B], got %v", pr.Added)
	}
	if len(pr.Skipped) != 1 || pr.Skipped[0] != "A" {
		t.Errorf("expected Skipped=[A], got %v", pr.Skipped)
	}
}

func TestPatch_AddAndUpdate_UpdatesChanged(t *testing.T) {
	base := map[string]string{"A": "old"}
	result := buildResult([]Entry{
		{Key: "A", Status: StatusChanged, BaseValue: "old", OtherValue: "new"},
	})

	pr, err := Patch(base, result, PatchAddAndUpdate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if base["A"] != "new" {
		t.Errorf("expected A=new, got %q", base["A"])
	}
	if len(pr.Updated) != 1 {
		t.Errorf("expected 1 updated key, got %d", len(pr.Updated))
	}
}

func TestPatch_Full_RemovesExtraKeys(t *testing.T) {
	base := map[string]string{"A": "1", "EXTRA": "x"}
	result := buildResult([]Entry{
		{Key: "EXTRA", Status: StatusExtra, BaseValue: "x"},
	})

	pr, err := Patch(base, result, PatchFull)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := base["EXTRA"]; ok {
		t.Error("expected EXTRA to be removed")
	}
	if len(pr.Removed) != 1 || pr.Removed[0] != "EXTRA" {
		t.Errorf("expected Removed=[EXTRA], got %v", pr.Removed)
	}
}

func TestPatch_NilBase_ReturnsError(t *testing.T) {
	_, err := Patch(nil, Result{}, PatchAddOnly)
	if err == nil {
		t.Error("expected error for nil base map")
	}
}

func TestPatchResult_HasChanges(t *testing.T) {
	empty := PatchResult{}
	if empty.HasChanges() {
		t.Error("expected HasChanges=false for empty result")
	}
	withAdd := PatchResult{Added: []string{"X"}}
	if !withAdd.HasChanges() {
		t.Error("expected HasChanges=true when keys added")
	}
}

func TestPatchResult_Summary_Format(t *testing.T) {
	pr := PatchResult{
		Added:   []string{"A", "B"},
		Updated: []string{"C"},
		Removed: []string{},
		Skipped: []string{"D"},
	}
	got := pr.Summary()
	want := "patch: +2 ~1 -0 (skipped 1)"
	if got != want {
		t.Errorf("Summary() = %q, want %q", got, want)
	}
}
