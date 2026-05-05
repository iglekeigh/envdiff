// Package output provides formatting utilities for rendering diff results
// to various output formats such as plain text and JSON.
package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the output format type.
type Format string

const (
	// FormatText renders diffs as human-readable plain text.
	FormatText Format = "text"
	// FormatJSON renders diffs as a JSON array.
	FormatJSON Format = "json"
)

// Formatter writes diff results to an io.Writer in a specific format.
type Formatter struct {
	format Format
	w      io.Writer
}

// New creates a new Formatter writing to w using the given format.
func New(w io.Writer, format Format) *Formatter {
	return &Formatter{format: format, w: w}
}

// Write renders the diff result entries to the underlying writer.
func (f *Formatter) Write(result *diff.Result) error {
	switch f.format {
	case FormatJSON:
		return f.writeJSON(result)
	default:
		return f.writeText(result)
	}
}

func (f *Formatter) writeText(result *diff.Result) error {
	for _, entry := range result.Entries {
		var line string
		switch entry.Status {
		case diff.StatusMatch:
			line = fmt.Sprintf("  %s=%s", entry.Key, entry.BaseValue)
		case diff.StatusMissing:
			line = fmt.Sprintf("- %s=%s", entry.Key, entry.BaseValue)
		case diff.StatusExtra:
			line = fmt.Sprintf("+ %s=%s", entry.Key, entry.OtherValue)
		case diff.StatusChanged:
			line = fmt.Sprintf("~ %s: %s -> %s", entry.Key, entry.BaseValue, entry.OtherValue)
		}
		if _, err := fmt.Fprintln(f.w, line); err != nil {
			return err
		}
	}
	return nil
}

func (f *Formatter) writeJSON(result *diff.Result) error {
	enc := json.NewEncoder(f.w)
	enc.SetIndent("", "  ")
	return enc.Encode(result.Entries)
}
