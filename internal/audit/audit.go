// Package audit provides change tracking and audit log support for env file operations.
package audit

import (
	"fmt"
	"time"
)

// Action represents the type of change recorded in an audit entry.
type Action string

const (
	ActionAdded    Action = "added"
	ActionRemoved  Action = "removed"
	ActionChanged  Action = "changed"
	ActionRedacted Action = "redacted"
)

// Entry represents a single recorded change to an environment variable.
type Entry struct {
	Timestamp time.Time
	Key       string
	Action    Action
	OldValue  string
	NewValue  string
	Source    string
}

// Log holds a sequence of audit entries.
type Log struct {
	entries []Entry
	source  string
}

// New creates a new audit Log tagged with the given source label.
func New(source string) *Log {
	return &Log{source: source}
}

// Record appends an entry to the log.
func (l *Log) Record(key string, action Action, oldVal, newVal string) {
	l.entries = append(l.entries, Entry{
		Timestamp: time.Now().UTC(),
		Key:       key,
		Action:    action,
		OldValue:  oldVal,
		NewValue:  newVal,
		Source:    l.source,
	})
}

// Entries returns a copy of all recorded entries.
func (l *Log) Entries() []Entry {
	out := make([]Entry, len(l.entries))
	copy(out, l.entries)
	return out
}

// Len returns the number of entries in the log.
func (l *Log) Len() int {
	return len(l.entries)
}

// String formats an entry as a human-readable line.
func (e Entry) String() string {
	ts := e.Timestamp.Format(time.RFC3339)
	switch e.Action {
	case ActionAdded:
		return fmt.Sprintf("[%s] %s: ADD %s = %q", ts, e.Source, e.Key, e.NewValue)
	case ActionRemoved:
		return fmt.Sprintf("[%s] %s: REMOVE %s (was %q)", ts, e.Source, e.Key, e.OldValue)
	case ActionChanged:
		return fmt.Sprintf("[%s] %s: CHANGE %s %q -> %q", ts, e.Source, e.Key, e.OldValue, e.NewValue)
	case ActionRedacted:
		return fmt.Sprintf("[%s] %s: REDACT %s", ts, e.Source, e.Key)
	default:
		return fmt.Sprintf("[%s] %s: UNKNOWN %s", ts, e.Source, e.Key)
	}
}
