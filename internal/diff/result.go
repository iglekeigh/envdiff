package diff

import (
	"fmt"
	"io"
	"sort"
)

// Status represents the kind of difference found for a key.
type Status string

const (
	StatusMatch   Status = "match"
	StatusMissing Status = "missing"
	StatusExtra   Status = "extra"
	StatusChanged Status = "changed"
)

// Entry represents a single comparison result for one key.
type Entry struct {
	Key      string
	Status   Status
	BaseVal  string
	OtherVal string
}

// Result holds all comparison entries from a diff operation.
type Result struct {
	Entries []Entry
}

// HasDifferences returns true if any entry is not a match.
func (r *Result) HasDifferences() bool {
	for _, e := range r.Entries {
		if e.Status != StatusMatch {
			return true
		}
	}
	return false
}

// Filter returns a new Result containing only entries with the given status.
func (r *Result) Filter(s Status) *Result {
	var filtered []Entry
	for _, e := range r.Entries {
		if e.Status == s {
			filtered = append(filtered, e)
		}
	}
	return &Result{Entries: filtered}
}

// WriteTo writes a human-readable diff summary to the provided writer.
func (r *Result) WriteTo(w io.Writer) {
	keys := make([]string, 0, len(r.Entries))
	index := make(map[string]Entry, len(r.Entries))
	for _, e := range r.Entries {
		keys = append(keys, e.Key)
		index[e.Key] = e
	}
	sort.Strings(keys)

	for _, k := range keys {
		e := index[k]
		switch e.Status {
		case StatusMatch:
			fmt.Fprintf(w, "  %s=%s\n", e.Key, e.BaseVal)
		case StatusMissing:
			fmt.Fprintf(w, "- %s=%s\n", e.Key, e.BaseVal)
		case StatusExtra:
			fmt.Fprintf(w, "+ %s=%s\n", e.Key, e.OtherVal)
		case StatusChanged:
			fmt.Fprintf(w, "~ %s: %s -> %s\n", e.Key, e.BaseVal, e.OtherVal)
		}
	}
}
