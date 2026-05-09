package output

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/chain"
)

// WriteChainResult writes the resolved environment from a chain operation to w.
// Keys are printed in sorted order with their resolved value and source file.
func WriteChainResult(w io.Writer, result *chain.Result) {
	keys := make([]string, 0, len(result.Env))
	for k := range result.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintln(w, "# Chained environment (resolved)")
	for _, k := range keys {
		fmt.Fprintf(w, "%s=%s\n", k, result.Env[k])
	}
}

// WriteChainSummary writes a human-readable summary of the chain operation,
// highlighting keys that were overridden by multiple files.
func WriteChainSummary(w io.Writer, result *chain.Result, strategy chain.Strategy) {
	strategyName := "first"
	if strategy == chain.StrategyLast {
		strategyName = "last"
	}

	fmt.Fprintf(w, "Strategy : %s\n", strategyName)
	fmt.Fprintf(w, "Total keys: %d\n", len(result.Env))

	var conflicted []string
	for k, sources := range result.Overrides {
		if len(sources) > 1 {
			conflicted = append(conflicted, k)
		}
	}
	sort.Strings(conflicted)

	if len(conflicted) == 0 {
		fmt.Fprintln(w, "No conflicts detected.")
		return
	}

	fmt.Fprintf(w, "Conflicts : %d\n", len(conflicted))
	for _, k := range conflicted {
		sources := result.Overrides[k]
		fmt.Fprintf(w, "  %s defined in %d files (resolved from: %s)\n",
			k, len(sources), result.Sources[k])
	}
}
