package output

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/clone"
)

// WriteCloneResult writes a human-readable summary of a clone operation.
func WriteCloneResult(w io.Writer, res *clone.Result) {
	keys := make([]string, 0, len(res.Mapped))
	for k := range res.Mapped {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintf(w, "Cloned %d key(s)\n", len(res.Env))
	for _, orig := range keys {
		newKey := res.Mapped[orig]
		if orig == newKey {
			fmt.Fprintf(w, "  %s (unchanged)\n", orig)
		} else {
			fmt.Fprintf(w, "  %s -> %s\n", orig, newKey)
		}
	}

	if len(res.Skipped) > 0 {
		sorted := make([]string, len(res.Skipped))
		copy(sorted, res.Skipped)
		sort.Strings(sorted)
		fmt.Fprintf(w, "\nSkipped %d key(s):\n", len(sorted))
		for _, k := range sorted {
			fmt.Fprintf(w, "  %s\n", k)
		}
	}
}

// WriteClonedEnv writes the resulting env map in KEY=VALUE format.
func WriteClonedEnv(w io.Writer, env map[string]string) {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(w, "%s=%s\n", k, env[k])
	}
}
