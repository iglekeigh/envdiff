// Package profile manages named environment profiles, allowing users to
// store and switch between multiple named sets of env file paths.
package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Profile represents a named collection of env file paths.
type Profile struct {
	Name  string   `json:"name"`
	Files []string `json:"files"`
}

// Registry holds a set of named profiles persisted to a JSON file.
type Registry struct {
	path     string
	profiles map[string]Profile
}

// Load reads the registry from the given file path.
// If the file does not exist, an empty registry is returned.
func Load(path string) (*Registry, error) {
	r := &Registry{path: path, profiles: make(map[string]Profile)}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return r, nil
	}
	if err != nil {
		return nil, fmt.Errorf("profile: read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, &r.profiles); err != nil {
		return nil, fmt.Errorf("profile: parse %s: %w", path, err)
	}
	return r, nil
}

// Save persists the registry to disk.
func (r *Registry) Save() error {
	if err := os.MkdirAll(filepath.Dir(r.path), 0o755); err != nil {
		return fmt.Errorf("profile: mkdir: %w", err)
	}
	data, err := json.MarshalIndent(r.profiles, "", "  ")
	if err != nil {
		return fmt.Errorf("profile: marshal: %w", err)
	}
	return os.WriteFile(r.path, data, 0o644)
}

// Set adds or replaces a profile.
func (r *Registry) Set(p Profile) {
	r.profiles[p.Name] = p
}

// Get retrieves a profile by name.
func (r *Registry) Get(name string) (Profile, bool) {
	p, ok := r.profiles[name]
	return p, ok
}

// Delete removes a profile by name.
func (r *Registry) Delete(name string) {
	delete(r.profiles, name)
}

// List returns all profiles sorted by name.
func (r *Registry) List() []Profile {
	out := make([]Profile, 0, len(r.profiles))
	for _, p := range r.profiles {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
