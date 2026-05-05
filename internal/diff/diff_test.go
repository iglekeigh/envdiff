package diff_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/diff"
)

func TestCompare_MatchingKeys(t *testing.T) {
	base := map[string]string{"FOO": "bar", "BAZ": "qux"}
	target := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := diff.Compare(base, target)

	if result.HasDifferences() {
		t.Errorf("expected no differences, got some")
	}
	if len(result.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result.Entries))
	}
}

func TestCompare_MissingKey(t *testing.T) {
	base := map[string]string{"FOO": "bar", "SECRET": "abc"}
	target := map[string]string{"FOO": "bar"}

	result := diff.Compare(base, target)

	if !result.HasDifferences() {
		t.Errorf("expected differences")
	}
	found := findEntry(result, "SECRET")
	if found == nil {
		t.Fatalf("expected entry for SECRET")
	}
	if found.Status != diff.StatusMissing {
		t.Errorf("expected StatusMissing, got %s", found.Status)
	}
}

func TestCompare_ExtraKey(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	target := map[string]string{"FOO": "bar", "EXTRA": "val"}

	result := diff.Compare(base, target)

	if !result.HasDifferences() {
		t.Errorf("expected differences")
	}
	found := findEntry(result, "EXTRA")
	if found == nil {
		t.Fatalf("expected entry for EXTRA")
	}
	if found.Status != diff.StatusExtra {
		t.Errorf("expected StatusExtra, got %s", found.Status)
	}
}

func TestCompare_ChangedValue(t *testing.T) {
	base := map[string]string{"DB_HOST": "localhost"}
	target := map[string]string{"DB_HOST": "prod.db.example.com"}

	result := diff.Compare(base, target)

	if !result.HasDifferences() {
		t.Errorf("expected differences")
	}
	found := findEntry(result, "DB_HOST")
	if found == nil {
		t.Fatalf("expected entry for DB_HOST")
	}
	if found.Status != diff.StatusChanged {
		t.Errorf("expected StatusChanged, got %s", found.Status)
	}
	if found.BaseValue != "localhost" || found.TargetValue != "prod.db.example.com" {
		t.Errorf("unexpected values: base=%q target=%q", found.BaseValue, found.TargetValue)
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	result := diff.Compare(map[string]string{}, map[string]string{})
	if result.HasDifferences() {
		t.Errorf("expected no differences for empty maps")
	}
}

func findEntry(r *diff.Result, key string) *diff.Entry {
	for i := range r.Entries {
		if r.Entries[i].Key == key {
			return &r.Entries[i]
		}
	}
	return nil
}
