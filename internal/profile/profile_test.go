package profile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/profile"
)

func tempRegistry(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "profiles.json")
}

func TestLoad_NonExistentFile_ReturnsEmptyRegistry(t *testing.T) {
	r, err := profile.Load(tempRegistry(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := r.List(); len(got) != 0 {
		t.Errorf("expected empty list, got %d profiles", len(got))
	}
}

func TestSet_And_Get(t *testing.T) {
	r, _ := profile.Load(tempRegistry(t))
	p := profile.Profile{Name: "staging", Files: []string{".env.staging"}}
	r.Set(p)
	got, ok := r.Get("staging")
	if !ok {
		t.Fatal("expected profile to exist")
	}
	if got.Name != "staging" || len(got.Files) != 1 {
		t.Errorf("unexpected profile: %+v", got)
	}
}

func TestDelete_RemovesProfile(t *testing.T) {
	r, _ := profile.Load(tempRegistry(t))
	r.Set(profile.Profile{Name: "dev", Files: []string{".env.dev"}})
	r.Delete("dev")
	if _, ok := r.Get("dev"); ok {
		t.Error("expected profile to be deleted")
	}
}

func TestList_ReturnsSorted(t *testing.T) {
	r, _ := profile.Load(tempRegistry(t))
	r.Set(profile.Profile{Name: "prod", Files: []string{".env.prod"}})
	r.Set(profile.Profile{Name: "dev", Files: []string{".env.dev"}})
	r.Set(profile.Profile{Name: "staging", Files: []string{".env.staging"}})
	list := r.List()
	names := []string{list[0].Name, list[1].Name, list[2].Name}
	want := []string{"dev", "prod", "staging"}
	for i, n := range want {
		if names[i] != n {
			t.Errorf("position %d: got %s, want %s", i, names[i], n)
		}
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	path := tempRegistry(t)
	r, _ := profile.Load(path)
	r.Set(profile.Profile{Name: "ci", Files: []string{".env.ci", ".env.secrets"}})
	if err := r.Save(); err != nil {
		t.Fatalf("save error: %v", err)
	}
	r2, err := profile.Load(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	p, ok := r2.Get("ci")
	if !ok {
		t.Fatal("profile 'ci' not found after reload")
	}
	if len(p.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(p.Files))
	}
}

func TestLoad_InvalidJSON_ReturnsError(t *testing.T) {
	path := tempRegistry(t)
	os.WriteFile(path, []byte("not json{"), 0o644)
	_, err := profile.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
