package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
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

func buildBinary(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "envdiff")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func TestMain_NoDifferences_ExitsZero(t *testing.T) {
	bin := buildBinary(t)
	base := writeTempFile(t, "FOO=bar\nBAZ=qux\n")
	other := writeTempFile(t, "FOO=bar\nBAZ=qux\n")
	cmd := exec.Command(bin, base, other)
	if err := cmd.Run(); err != nil {
		t.Errorf("expected exit 0, got: %v", err)
	}
}

func TestMain_Differences_ExitsTwo(t *testing.T) {
	bin := buildBinary(t)
	base := writeTempFile(t, "FOO=bar\n")
	other := writeTempFile(t, "FOO=changed\n")
	cmd := exec.Command(bin, base, other)
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit")
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 2 {
			t.Errorf("expected exit code 2, got %d", exitErr.ExitCode())
		}
	}
}

func TestMain_MissingArgs_ExitsOne(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit")
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 1 {
			t.Errorf("expected exit code 1, got %d", exitErr.ExitCode())
		}
	}
}
