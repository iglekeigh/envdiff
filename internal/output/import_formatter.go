package output

import (
	"fmt"
	"io"
	"sort"

	importer "github.com/user/envdiff/internal/import"
)

// WriteImportResult writes a human-readable summary of an import operation.
func WriteImportResult(w io.Writer, res *importer.Result) {
	fmt.Fprintf(w, "Imported %d key(s) from %s [format: %s]\n",
		len(res.Env), res.Source, res.Format)

	if len(res.Skipped) > 0 {
		sort.Strings(res.Skipped)
		fmt.Fprintf(w, "Skipped %d key(s) (unsupported type):\n", len(res.Skipped))
		for _, k := range res.Skipped {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}
}

// WriteImportedEnv writes the imported env map in dotenv format.
func WriteImportedEnv(w io.Writer, res *importer.Result) {
	keys := make([]string, 0, len(res.Env))
	for k := range res.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := res.Env[k]
		if needsQuoting(v) {
			fmt.Fprintf(w, "%s=%q\n", k, v)
		} else {
			fmt.Fprintf(w, "%s=%s\n", k, v)
		}
	}
}

// needsQuoting returns true if the value contains characters that require quoting.
func needsQuoting(v string) bool {
	for _, c := range v {
		if c == ' ' || c == '\t' || c == '\n' || c == '"' || c == '\'' {
			return true
		}
	}
	return false
}
