package output

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/reconcile"
)

// WriteReconcileResult writes the reconciled environment map and a summary
// of changes to the provided writer in a human-readable format.
func WriteReconcileResult(w io.Writer, result reconcile.Result) error {
	// Print summary header
	fmt.Fprintf(w, "Reconcile summary: %d added, %d kept, %d overwritten\n",
		result.Added, result.Kept, result.Overwritten)

	if len(result.Conflicts) > 0 {
		fmt.Fprintf(w, "\nConflicts resolved (%d):\n", len(result.Conflicts))
		keys := make([]string, 0, len(result.Conflicts))
		for k := range result.Conflicts {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			c := result.Conflicts[k]
			fmt.Fprintf(w, "  %-30s base=%q other=%q => kept=%q\n", k, c.Base, c.Other, c.Resolved)
		}
	}

	fmt.Fprintln(w, "\nReconciled env:")
	keys := make([]string, 0, len(result.Env))
	for k := range result.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(w, "  %s=%s\n", k, result.Env[k])
	}
	return nil
}
