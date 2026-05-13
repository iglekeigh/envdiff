package dedup_test

import (
	"testing"

	"github.com/user/envdiff/internal/dedup"
)

func TestApply_NoDuplicates_ReturnsAllKeys(t *testing.T) {
	sources := []map[string]string{
		{"A": "1", "B": "2"},
		{"C": "3"},
	}
	res, err := dedup.Apply(sources, dedup.StrategyKeepFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 3 {
		t.Errorf("expected 3 keys, got %d", len(res.Env))
	}
	if len(res.Duplicates) != 0 {
		t.Errorf("expected no duplicates, got %d", len(res.Duplicates))
	}
}

func TestApply_StrategyKeepFirst_RetainsFirstValue(t *testing.T) {
	sources := []map[string]string{
		{"KEY": "first"},
		{"KEY": "second"},
	}
	res, err := dedup.Apply(sources, dedup.StrategyKeepFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "first" {
		t.Errorf("expected 'first', got %q", res.Env["KEY"])
	}
	if len(res.Duplicates) != 1 {
		t.Errorf("expected 1 duplicate, got %d", len(res.Duplicates))
	}
}

func TestApply_StrategyKeepLast_RetainsLastValue(t *testing.T) {
	sources := []map[string]string{
		{"KEY": "first"},
		{"KEY": "second"},
		{"KEY": "third"},
	}
	res, err := dedup.Apply(sources, dedup.StrategyKeepLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "third" {
		t.Errorf("expected 'third', got %q", res.Env["KEY"])
	}
}

func TestApply_StrategyError_ReturnsDuplicateError(t *testing.T) {
	sources := []map[string]string{
		{"KEY": "a"},
		{"KEY": "b"},
	}
	_, err := dedup.Apply(sources, dedup.StrategyError)
	if err == nil {
		t.Fatal("expected error for duplicate key, got nil")
	}
}

func TestApply_EmptySources_ReturnsEmptyResult(t *testing.T) {
	res, err := dedup.Apply([]map[string]string{}, dedup.StrategyKeepFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 0 {
		t.Errorf("expected empty env, got %d keys", len(res.Env))
	}
}

func TestApply_DuplicateValues_RecordsAllValues(t *testing.T) {
	sources := []map[string]string{
		{"X": "alpha"},
		{"X": "beta"},
	}
	res, _ := dedup.Apply(sources, dedup.StrategyKeepFirst)
	if len(res.Duplicates) != 1 {
		t.Fatalf("expected 1 duplicate entry, got %d", len(res.Duplicates))
	}
	d := res.Duplicates[0]
	if d.Key != "X" {
		t.Errorf("expected key 'X', got %q", d.Key)
	}
	if len(d.Values) != 2 {
		t.Errorf("expected 2 values recorded, got %d", len(d.Values))
	}
}
