package output

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/user/envdiff/internal/diff"
)

// WriteDiffSummary writes a human-readable summary of diff results to w.
// It groups entries by status and prints a final count line.
func WriteDiffSummary(w io.Writer, result diff.Result) {
	entries := result.Entries()
	if len(entries) == 0 {
		fmt.Fprintln(w, "✔ No differences found.")
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	var missing, extra, changed int
	for _, e := range entries {
		switch e.Status {
		case diff.StatusMissing:
			fmt.Fprintf(tw, "  MISSING\t%s\t(not in second file)\n", e.Key)
			missing++
		case diff.StatusExtra:
			fmt.Fprintf(tw, "  EXTRA  \t%s\t(not in first file)\n", e.Key)
			extra++
		case diff.StatusChanged:
			fmt.Fprintf(tw, "  CHANGED\t%s\t%q -> %q\n", e.Key, e.BaseValue, e.OtherValue)
			changed++
		}
	}
	tw.Flush()

	fmt.Fprintf(w, "\nSummary: %d missing, %d extra, %d changed\n", missing, extra, changed)
}
