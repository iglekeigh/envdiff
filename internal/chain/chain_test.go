package chain_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/chain"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestApply_NoPaths_ReturnsEmptyResult(t *testing.T) {
	result, err := chain.Apply(nil, chain.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Env) != 0 {
		t.Errorf("expected empty env, got %v", result.Env)
	}
}

func TestApply_SingleFile_LoadsKeys(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	result, err := chain.Apply([]string{path}, chain.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Env["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", result.Env["FOO"])
	}
	if result.Sources["FOO"] != path {
		t.Errorf("expected source %q, got %q", path, result.Sources["FOO"])
	}
}

func TestApply_StrategyFirst_KeepsFirstValue(t *testing.T) {
	a := writeTempEnv(t, "KEY=from_a\n")
	b := writeTempEnv(t, "KEY=from_b\n")
	result, err := chain.Apply([]string{a, b}, chain.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Env["KEY"] != "from_a" {
		t.Errorf("expected from_a, got %q", result.Env["KEY"])
	}
}

func TestApply_StrategyLast_KeepsLastValue(t *testing.T) {
	a := writeTempEnv(t, "KEY=from_a\n")
	b := writeTempEnv(t, "KEY=from_b\n")
	result, err := chain.Apply([]string{a, b}, chain.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Env["KEY"] != "from_b" {
		t.Errorf("expected from_b, got %q", result.Env["KEY"])
	}
}

func TestApply_OverridesTracked(t *testing.T) {
	a := writeTempEnv(t, "KEY=from_a\n")
	b := writeTempEnv(t, "KEY=from_b\n")
	result, err := chain.Apply([]string{a, b}, chain.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Overrides["KEY"]) != 2 {
		t.Errorf("expected 2 override sources, got %d", len(result.Overrides["KEY"]))
	}
}

func TestApply_InvalidPath_ReturnsError(t *testing.T) {
	_, err := chain.Apply([]string{filepath.Join(t.TempDir(), "nonexistent.env")}, chain.StrategyFirst)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
