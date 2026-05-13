package filter_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/filter"
)

var sampleEnv = map[string]string{
	"APP_HOST":     "localhost",
	"APP_PORT":     "8080",
	"DB_HOST":      "db.local",
	"DB_PASSWORD":  "secret",
	"SECRET_TOKEN": "tok123",
	"LOG_LEVEL":    "info",
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	result, err := filter.Apply(sampleEnv, filter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != len(sampleEnv) {
		t.Errorf("expected %d keys, got %d", len(sampleEnv), len(result))
	}
}

func TestApply_PrefixFilter_RetainsMatchingKeys(t *testing.T) {
	result, err := filter.Apply(sampleEnv, filter.Options{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in result")
	}
}

func TestApply_SuffixFilter_RetainsMatchingKeys(t *testing.T) {
	result, err := filter.Apply(sampleEnv, filter.Options{Suffix: "_HOST"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}

func TestApply_PatternFilter_RetainsMatchingKeys(t *testing.T) {
	result, err := filter.Apply(sampleEnv, filter.Options{Pattern: "^(DB|SECRET)_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 keys, got %d", len(result))
	}
}

func TestApply_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := filter.Apply(sampleEnv, filter.Options{Pattern: "["})
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestApply_Invert_ExcludesMatched(t *testing.T) {
	result, err := filter.Apply(sampleEnv, filter.Options{Prefix: "APP_", Invert: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != len(sampleEnv)-2 {
		t.Errorf("expected %d keys, got %d", len(sampleEnv)-2, len(result))
	}
	if _, ok := result["APP_HOST"]; ok {
		t.Error("expected APP_HOST to be excluded")
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"APP_X": "1", "DB_Y": "2"}
	result, _ := filter.Apply(env, filter.Options{Prefix: "APP_"})
	result["INJECTED"] = "yes"
	if _, ok := env["INJECTED"]; ok {
		t.Error("original map was mutated")
	}
}
