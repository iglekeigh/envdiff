package output

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/user/envdiff/internal/audit"
)

// WriteAuditLog writes a formatted audit log to w.
// Entries are printed in chronological order with tabular alignment.
func WriteAuditLog(w io.Writer, log *audit.Log) {
	entries := log.Entries()
	if len(entries) == 0 {
		fmt.Fprintln(w, "audit log: no changes recorded")
		return
	}

	fmt.Fprintf(w, "audit log: %d change(s)\n", len(entries))

	tw := tabwriter.NewWriter(w, 2, 2, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tSOURCE\tACTION\tKEY")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02T15:04:05Z"),
			e.Source,
			string(e.Action),
			e.Key,
		)
	}
	tw.Flush()
}

// WriteAuditSummary writes a grouped summary of audit actions to w.
func WriteAuditSummary(w io.Writer, log *audit.Log) {
	entries := log.Entries()
	counts := map[audit.Action]int{}
	for _, e := range entries {
		counts[e.Action]++
	}

	if len(counts) == 0 {
		fmt.Fprintln(w, "audit summary: nothing to report")
		return
	}

	actions := make([]string, 0, len(counts))
	for a := range counts {
		actions = append(actions, string(a))
	}
	sort.Strings(actions)

	fmt.Fprintln(w, "audit summary:")
	for _, a := range actions {
		fmt.Fprintf(w, "  %-10s %d\n", a, counts[audit.Action(a)])
	}
}
