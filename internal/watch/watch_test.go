package watch_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/watch"
)

func writeTempEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestSnapshot_RecordsState(t *testing.T) {
	dir := t.TempDir()
	p := writeTempEnv(t, dir, ".env", "KEY=value\n")

	w := watch.New(p)
	if err := w.Snapshot(); err != nil {
		t.Fatalf("Snapshot: %v", err)
	}

	state, ok := w.State(p)
	if !ok {
		t.Fatal("expected state to be recorded")
	}
	if state.Checksum == "" {
		t.Error("expected non-empty checksum")
	}
}

func TestChanged_UnmodifiedFile_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	p := writeTempEnv(t, dir, ".env", "KEY=value\n")

	w := watch.New(p)
	_ = w.Snapshot()

	changed, err := w.Changed()
	if err != nil {
		t.Fatalf("Changed: %v", err)
	}
	if len(changed) != 0 {
		t.Errorf("expected no changes, got %v", changed)
	}
}

func TestChanged_ModifiedFile_ReturnsPath(t *testing.T) {
	dir := t.TempDir()
	p := writeTempEnv(t, dir, ".env", "KEY=value\n")

	w := watch.New(p)
	_ = w.Snapshot()

	// Modify the file
	if err := os.WriteFile(p, []byte("KEY=changed\n"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	changed, err := w.Changed()
	if err != nil {
		t.Fatalf("Changed: %v", err)
	}
	if len(changed) != 1 || changed[0] != p {
		t.Errorf("expected [%s], got %v", p, changed)
	}
}

func TestChanged_NoSnapshot_ReportsAllAsChanged(t *testing.T) {
	dir := t.TempDir()
	p := writeTempEnv(t, dir, ".env", "KEY=value\n")

	w := watch.New(p)
	// No Snapshot call
	changed, err := w.Changed()
	if err != nil {
		t.Fatalf("Changed: %v", err)
	}
	if len(changed) != 1 {
		t.Errorf("expected 1 changed file, got %d", len(changed))
	}
}

func TestSnapshot_MissingFile_ReturnsError(t *testing.T) {
	w := watch.New("/nonexistent/.env")
	if err := w.Snapshot(); err == nil {
		t.Error("expected error for missing file")
	}
}
