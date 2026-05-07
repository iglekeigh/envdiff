package env_test

import (
	"os"
	"testing"

	"github.com/user/envdiff/internal/env"
)

func TestResolve_FirstSourceWins(t *testing.T) {
	sources := []env.Source{
		{Label: "prod", Values: map[string]string{"DB_HOST": "prod-db"}},
		{Label: "dev", Values: map[string]string{"DB_HOST": "dev-db"}},
	}
	res := env.Resolve("DB_HOST", sources, false)
	if !res.Found {
		t.Fatal("expected key to be found")
	}
	if res.Value != "prod-db" {
		t.Errorf("expected prod-db, got %s", res.Value)
	}
	if res.SourceLabel != "prod" {
		t.Errorf("expected source prod, got %s", res.SourceLabel)
	}
}

func TestResolve_FallsBackToOS(t *testing.T) {
	os.Setenv("_TEST_ENVDIFF_KEY", "from-os")
	defer os.Unsetenv("_TEST_ENVDIFF_KEY")

	res := env.Resolve("_TEST_ENVDIFF_KEY", nil, true)
	if !res.Found {
		t.Fatal("expected key to be found in OS")
	}
	if res.Value != "from-os" {
		t.Errorf("expected from-os, got %s", res.Value)
	}
	if res.SourceLabel != "os" {
		t.Errorf("expected source os, got %s", res.SourceLabel)
	}
}

func TestResolve_NotFound(t *testing.T) {
	res := env.Resolve("NONEXISTENT_KEY_XYZ", nil, false)
	if res.Found {
		t.Error("expected key not to be found")
	}
}

func TestResolveAll_CollectsAllKeys(t *testing.T) {
	sources := []env.Source{
		{Label: "a", Values: map[string]string{"FOO": "1", "BAR": "2"}},
		{Label: "b", Values: map[string]string{"BAZ": "3", "BAR": "99"}},
	}
	results := env.ResolveAll(sources, false)
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
	if results["BAR"].SourceLabel != "a" {
		t.Errorf("expected BAR from source a, got %s", results["BAR"].SourceLabel)
	}
	if results["BAZ"].Value != "3" {
		t.Errorf("expected BAZ=3, got %s", results["BAZ"].Value)
	}
}

func TestExpand_ResolvesVariables(t *testing.T) {
	env_ := map[string]string{"HOST": "localhost", "PORT": "5432"}
	out, err := env.Expand("postgres://${HOST}:${PORT}/db", env_)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "postgres://localhost:5432/db" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestExpand_ReturnsErrorOnMissing(t *testing.T) {
	_, err := env.Expand("${MISSING_VAR}", map[string]string{})
	if err == nil {
		t.Error("expected error for missing variable")
	}
}
