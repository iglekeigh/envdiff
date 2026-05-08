package compare

import (
	"testing"
)

func TestCompare_EmptyInput_ReturnsEmptyReport(t *testing.T) {
	r, err := Compare(map[string]map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(r.Keys))
	}
}

func TestCompare_AllConsistent_NoMissingKeys(t *testing.T) {
	envs := map[string]map[string]string{
		"prod":    {"HOST": "example.com", "PORT": "443"},
		"staging": {"HOST": "example.com", "PORT": "443"},
	}
	r, err := Compare(envs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, kr := range r.Keys {
		if !kr.Consistent() {
			t.Errorf("key %q expected consistent, was not", kr.Key)
		}
	}
}

func TestCompare_ValueDrift_MarkedInconsistent(t *testing.T) {
	envs := map[string]map[string]string{
		"prod":    {"DB_HOST": "prod-db"},
		"staging": {"DB_HOST": "staging-db"},
	}
	r, err := Compare(envs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Keys) != 1 {
		t.Fatalf("expected 1 key report, got %d", len(r.Keys))
	}
	if r.Keys[0].Consistent() {
		t.Error("expected DB_HOST to be inconsistent")
	}
}

func TestCompare_MissingKeyInOneFile(t *testing.T) {
	envs := map[string]map[string]string{
		"prod":    {"SECRET": "abc", "PORT": "80"},
		"staging": {"PORT": "80"},
	}
	r, err := Compare(envs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var secretReport *KeyReport
	for i := range r.Keys {
		if r.Keys[i].Key == "SECRET" {
			secretReport = &r.Keys[i]
		}
	}
	if secretReport == nil {
		t.Fatal("expected SECRET in report")
	}
	if _, ok := secretReport.Values["staging"]; ok {
		t.Error("staging should not have SECRET")
	}
	if _, ok := secretReport.Values["prod"]; !ok {
		t.Error("prod should have SECRET")
	}
}

func TestCompare_LabelsAreSorted(t *testing.T) {
	envs := map[string]map[string]string{
		"z-env": {"K": "1"},
		"a-env": {"K": "1"},
		"m-env": {"K": "1"},
	}
	r, err := Compare(envs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"a-env", "m-env", "z-env"}
	for i, l := range r.Labels {
		if l != expected[i] {
			t.Errorf("label[%d] = %q, want %q", i, l, expected[i])
		}
	}
}

func TestKeyReport_PresentIn_ReturnsSortedLabels(t *testing.T) {
	kr := KeyReport{
		Key:    "FOO",
		Values: map[string]string{"z": "1", "a": "2", "m": "3"},
	}
	got := kr.PresentIn()
	expected := []string{"a", "m", "z"}
	for i, v := range got {
		if v != expected[i] {
			t.Errorf("PresentIn[%d] = %q, want %q", i, v, expected[i])
		}
	}
}
