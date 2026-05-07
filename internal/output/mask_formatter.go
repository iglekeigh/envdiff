package output

import (
	"fmt"
	"io"
	"sort"
)

// WriteMaskedEnv writes the masked environment map to w in KEY=VALUE format,
// sorted by key. sensitiveKeys lists the keys whose values were masked.
func WriteMaskedEnv(w io.Writer, env map[string]string, sensitiveKeys []string) {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sensitiveSet := make(map[string]bool, len(sensitiveKeys))
	for _, k := range sensitiveKeys {
		sensitiveSet[k] = true
	}

	for _, k := range keys {
		marker := ""
		if sensitiveSet[k] {
			marker = " # masked"
		}
		fmt.Fprintf(w, "%s=%s%s\n", k, env[k], marker)
	}
}

// WriteMaskSummary writes a summary of how many keys were masked to w.
func WriteMaskSummary(w io.Writer, total int, maskedKeys []string) {
	sort.Strings(maskedKeys)
	fmt.Fprintf(w, "Mask summary: %d/%d keys masked\n", len(maskedKeys), total)
	if len(maskedKeys) == 0 {
		fmt.Fprintln(w, "  No sensitive keys detected.")
		return
	}
	for _, k := range maskedKeys {
		fmt.Fprintf(w, "  - %s\n", k)
	}
}
