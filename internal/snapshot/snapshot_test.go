package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/snapshot"
)

func TestNew_CopiesEnvMap(t *testing.T) {
	env := map[string]string{"KEY": "value", "DB_HOST": "localhost"}
	s := snapshot.New("test", env)

	env["KEY"] = "mutated"
	if s.Env["KEY"] != "value" {
		t.Errorf("expected snapshot to be independent of original map")
	}
}

func TestNew_SetsLabelAndTime(t *testing.T) {
	before := time.Now().UTC()
	s := snapshot.New("prod", map[string]string{})
	after := time.Now().UTC()

	if s.Label != "prod" {
		t.Errorf("expected label %q, got %q", "prod", s.Label)
	}
	if s.CreatedAt.Before(before) || s.CreatedAt.After(after) {
		t.Errorf("created_at %v out of expected range", s.CreatedAt)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s := snapshot.New("staging", env)

	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	if err := snapshot.Save(s, path); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if loaded.Label != s.Label {
		t.Errorf("label mismatch: got %q, want %q", loaded.Label, s.Label)
	}
	for k, v := range env {
		if loaded.Env[k] != v {
			t.Errorf("env[%q] = %q, want %q", k, loaded.Env[k], v)
		}
	}
}

func TestLoad_InvalidPath_ReturnsError(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestSave_InvalidPath_ReturnsError(t *testing.T) {
	s := snapshot.New("x", map[string]string{})
	err := snapshot.Save(s, "/nonexistent/dir/snap.json")
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}

func TestLoad_CorruptJSON_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not json{"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error for corrupt JSON, got nil")
	}
}
