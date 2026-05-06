package output

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/schema"
)

// WriteSchemaResult writes schema validation results to w.
// It prints a summary line and, if violations exist, lists each one.
func WriteSchemaResult(w io.Writer, violations []schema.Violation) {
	if len(violations) == 0 {
		fmt.Fprintln(w, "schema: ✓ all fields valid")
		return
	}

	fmt.Fprintf(w, "schema: ✗ %d violation(s) found\n", len(violations))

	// Sort for deterministic output
	sorted := make([]schema.Violation, len(violations))
	copy(sorted, violations)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, v := range sorted {
		fmt.Fprintf(w, "  [%s] %s\n", v.Key, v.Message)
	}
}

// WriteSchemaFields writes a summary of all fields defined in a schema.
func WriteSchemaFields(w io.Writer, s schema.Schema) {
	if len(s.Fields) == 0 {
		fmt.Fprintln(w, "schema: no fields defined")
		return
	}

	fmt.Fprintf(w, "schema: %d field(s) defined\n", len(s.Fields))

	sorted := make([]schema.Field, len(s.Fields))
	copy(sorted, s.Fields)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, f := range sorted {
		required := "optional"
		if f.Required {
			required = "required"
		}
		patternStr := ""
		if f.Pattern != nil {
			patternStr = fmt.Sprintf(" pattern=%s", f.Pattern)
		}
		fmt.Fprintf(w, "  %-30s type=%-8s %s%s\n", f.Key, f.Type, required, patternStr)
	}
}
