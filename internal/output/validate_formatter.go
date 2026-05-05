package output

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/validate"
)

// WriteViolations writes a human-readable summary of validation violations to w.
// Returns the number of violations written and any write error.
func WriteViolations(w io.Writer, violations []validate.Violation) (int, error) {
	if len(violations) == 0 {
		_, err := fmt.Fprintln(w, "validation passed: no violations found")
		return 0, err
	}

	// Sort for deterministic output
	sorted := make([]validate.Violation, len(violations))
	copy(sorted, violations)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Key != sorted[j].Key {
			return sorted[i].Key < sorted[j].Key
		}
		return sorted[i].Rule < sorted[j].Rule
	})

	_, err := fmt.Fprintf(w, "validation failed: %d violation(s)\n", len(sorted))
	if err != nil {
		return len(sorted), err
	}

	for _, v := range sorted {
		_, err = fmt.Fprintf(w, "  [%s] %s: %s\n", v.Rule, v.Key, v.Message)
		if err != nil {
			return len(sorted), err
		}
	}
	return len(sorted), nil
}
