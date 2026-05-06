package watch

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

// FileState holds the last known state of a watched file.
type FileState struct {
	Path    string
	Checksum string
	ModTime time.Time
}

// Watcher monitors one or more .env files for changes.
type Watcher struct {
	files  []string
	states map[string]FileState
}

// New creates a Watcher for the given file paths.
func New(paths ...string) *Watcher {
	return &Watcher{
		files:  paths,
		states: make(map[string]FileState),
	}
}

// Snapshot records the current checksum and modification time for all files.
func (w *Watcher) Snapshot() error {
	for _, path := range w.files {
		state, err := fileState(path)
		if err != nil {
			return fmt.Errorf("watch: snapshot %s: %w", path, err)
		}
		w.states[path] = state
	}
	return nil
}

// Changed returns the paths of files that have changed since the last Snapshot.
func (w *Watcher) Changed() ([]string, error) {
	var changed []string
	for _, path := range w.files {
		current, err := fileState(path)
		if err != nil {
			return nil, fmt.Errorf("watch: check %s: %w", path, err)
		}
		prev, seen := w.states[path]
		if !seen || current.Checksum != prev.Checksum {
			changed = append(changed, path)
		}
	}
	return changed, nil
}

// State returns the last recorded FileState for the given path.
func (w *Watcher) State(path string) (FileState, bool) {
	s, ok := w.states[path]
	return s, ok
}

func fileState(path string) (FileState, error) {
	f, err := os.Open(path)
	if err != nil {
		return FileState{}, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return FileState{}, err
	}

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return FileState{}, err
	}

	return FileState{
		Path:    path,
		Checksum: fmt.Sprintf("%x", h.Sum(nil)),
		ModTime: info.ModTime(),
	}, nil
}
