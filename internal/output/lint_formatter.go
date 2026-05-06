package output

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/user/envdiff/internal/lint"
)

// WriteLintFindings writes a formatted lint report to w.
// If there are no findings it prints a single pass line.
func WriteLintFindings(w io.Writer, findings []lint.Finding) {
	if len(findings) == 0 {
		fmt.Fprintln(w, "lint: OK — no issues found")
		return
	}

	// Sort for deterministic output: severity desc (error first), then key.
	sorted := make([]lint.Finding, len(findings))
	copy(sorted, findings)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Severity != sorted[j].Severity {
			// errors before warnings
			return sorted[i].Severity == lint.SeverityError
		}
		return sorted[i].Key < sorted[j].Key
	})

	errCount, warnCount := 0, 0
	for _, f := range sorted {
		if f.Severity == lint.SeverityError {
			errCount++
		} else {
			warnCount++
		}
	}

	fmt.Fprintf(w, "lint: %d error(s), %d warning(s)\n\n", errCount, warnCount)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SEVERITY\tKEY\tRULE\tMESSAGE")
	for _, f := range sorted {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", f.Severity, f.Key, f.Rule, f.Message)
	}
	_ = tw.Flush()
}
