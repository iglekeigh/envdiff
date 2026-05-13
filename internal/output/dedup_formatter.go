package output

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/dedup"
)

// WriteDedupResult writes a human-readable summary of a dedup.Result to w.
func WriteDedupResult(w io.Writer, res *dedup.Result) {
	if len(res.Duplicates) == 0 {
		fmt.Fprintln(w, "✔ No duplicate keys found.")
	} else {
		fmt.Fprintf(w, "⚠ %d duplicate key(s) resolved:\n", len(res.Duplicates))
		dupes := make([]dedup.Duplicate, len(res.Duplicates))
		copy(dupes, res.Duplicates)
		sort.Slice(dupes, func(i, j int) bool {
			return dupes[i].Key < dupes[j].Key
		})
		for _, d := range dupes {
			fmt.Fprintf(w, "  - %s (%d occurrences)\n", d.Key, len(d.Values))
			for i, v := range d.Values {
				fmt.Fprintf(w, "      [%d] %s\n", i+1, v)
			}
		}
	}

	fmt.Fprintf(w, "\nResulting env (%d keys):\n", len(res.Env))
	keys := make([]string, 0, len(res.Env))
	for k := range res.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(w, "  %s=%s\n", k, res.Env[k])
	}
}
