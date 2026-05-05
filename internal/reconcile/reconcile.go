// Package reconcile provides functionality to merge and reconcile
// differences between two .env file maps, producing a merged result
// with optional conflict resolution strategies.
package reconcile

import (
	"fmt"
	"io"

	"github.com/user/envdiff/internal/diff"
)

// Strategy defines how conflicts are resolved during reconciliation.
type Strategy int

const (
	// PreferBase keeps the base value on conflict.
	PreferBase Strategy = iota
	// PreferOther uses the other (incoming) value on conflict.
	PreferOther
	// ErrorOnConflict returns an error if any conflicting values exist.
	ErrorOnConflict
)

// Result holds the reconciled environment map and any warnings.
type Result struct {
	Merged   map[string]string
	Warnings []string
}

// Reconcile merges base and other env maps using the provided strategy.
// Keys missing from base are added from other. Keys missing from other
// are retained from base. Conflicting keys are resolved by strategy.
func Reconcile(base, other map[string]string, strategy Strategy) (*Result, error) {
	result := &Result{
		Merged:   make(map[string]string),
		Warnings: []string{},
	}

	// Start with a copy of base.
	for k, v := range base {
		result.Merged[k] = v
	}

	cmp := diff.Compare(base, other)

	for _, entry := range cmp.Entries {
		switch entry.Status {
		case diff.Extra:
			// Key exists in other but not base — add it.
			result.Merged[entry.Key] = entry.OtherValue
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("added missing key %q from other", entry.Key))
		case diff.Changed:
			switch strategy {
			case PreferOther:
				result.Merged[entry.Key] = entry.OtherValue
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("conflict on key %q: used other value", entry.Key))
			case ErrorOnConflict:
				return nil, fmt.Errorf("conflict on key %q: base=%q other=%q",
					entry.Key, entry.BaseValue, entry.OtherValue)
			default: // PreferBase
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("conflict on key %q: kept base value", entry.Key))
			}
		}
	}

	return result, nil
}

// WriteWarnings writes any reconciliation warnings to w.
func (r *Result) WriteWarnings(w io.Writer) {
	for _, warn := range r.Warnings {
		fmt.Fprintf(w, "warning: %s\n", warn)
	}
}
