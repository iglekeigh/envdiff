package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/compare"
)

// WriteCompareReport writes a human-readable multi-file comparison report to w.
func WriteCompareReport(w io.Writer, r compare.Report) {
	if len(r.Keys) == 0 {
		fmt.Fprintln(w, "No keys found across provided files.")
		return
	}

	header := fmt.Sprintf("%-30s", "KEY")
	for _, l := range r.Labels {
		header += fmt.Sprintf("  %-20s", l)
	}
	fmt.Fprintln(w, header)
	fmt.Fprintln(w, strings.Repeat("-", len(header)))

	driftCount := 0
	for _, kr := range r.Keys {
		marker := " "
		if !kr.Consistent() {
			marker = "!"
			driftCount++
		}
		row := fmt.Sprintf("%s %-29s", marker, kr.Key)
		for _, l := range r.Labels {
			val, ok := kr.Values[l]
			if !ok {
				row += fmt.Sprintf("  %-20s", "<missing>")
			} else {
				if len(val) > 18 {
					val = val[:15] + "..."
				}
				row += fmt.Sprintf("  %-20s", val)
			}
		}
		fmt.Fprintln(w, row)
	}

	fmt.Fprintln(w, strings.Repeat("-", len(header)))
	if driftCount == 0 {
		fmt.Fprintln(w, "All keys are consistent across files.")
	} else {
		fmt.Fprintf(w, "%d key(s) have drift or are missing in one or more files.\n", driftCount)
	}
}

// WriteCompareSummary writes a one-line summary of the comparison result.
func WriteCompareSummary(w io.Writer, r compare.Report) {
	total := len(r.Keys)
	drift := 0
	for _, kr := range r.Keys {
		if !kr.Consistent() {
			drift++
		}
	}
	fmt.Fprintf(w, "Compared %d files, %d total keys, %d with drift.\n", len(r.Labels), total, drift)
}
