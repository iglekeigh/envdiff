package output

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/export"
)

// WriteExportResult writes a summary of an export operation to w.
func WriteExportResult(w io.Writer, result export.Result) {
	fmt.Fprintf(w, "Exported %d key(s) to %s [format: %s]\n", result.Count, result.Path, result.Format)
	if len(result.Skipped) > 0 {
		sort.Strings(result.Skipped)
		fmt.Fprintf(w, "Skipped %d key(s) incompatible with format:\n", len(result.Skipped))
		for _, k := range result.Skipped {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}
}

// WriteExportPreview writes the env map as it would appear in the given format,
// without writing to disk. Useful for dry-run display.
func WriteExportPreview(w io.Writer, env map[string]string, format export.Format) {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintf(w, "Preview [format: %s]:\n", format)
	switch format {
	case export.FormatJSON:
		fmt.Fprintln(w, "{")
		for i, k := range keys {
			comma := ","
			if i == len(keys)-1 {
				comma = ""
			}
			fmt.Fprintf(w, "  %q: %q%s\n", k, env[k], comma)
		}
		fmt.Fprintln(w, "}")
	case export.FormatShell:
		for _, k := range keys {
			fmt.Fprintf(w, "export %s=%q\n", k, env[k])
		}
	default:
		for _, k := range keys {
			fmt.Fprintf(w, "%s=%s\n", k, env[k])
		}
	}
}
