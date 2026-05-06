package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a saved state of an env file at a point in time.
type Snapshot struct {
	Label     string            `json:"label"`
	CreatedAt time.Time         `json:"created_at"`
	Env       map[string]string `json:"env"`
}

// New creates a new Snapshot from the given env map and label.
func New(label string, env map[string]string) *Snapshot {
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return &Snapshot{
		Label:     label,
		CreatedAt: time.Now().UTC(),
		Env:       copy,
	}
}

// Save writes the snapshot to the given file path as JSON.
func Save(s *Snapshot, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create file %q: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open file %q: %w", path, err)
	}
	defer f.Close()

	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &s, nil
}
