package output

import (
	"fmt"
	"io"

	"github.com/user/envdiff/internal/template"
)

// WriteTemplateResult writes a human-readable summary of a template check result.
func WriteTemplateResult(w io.Writer, result template.Result) {
	if !result.HasIssues() {
		fmt.Fprintln(w, "✔ env matches template — all required keys present, no extra keys")
		return
	}

	if len(result.Missing) > 0 {
		fmt.Fprintf(w, "✘ Missing keys (%d):\n", len(result.Missing))
		for _, mk := range result.Missing {
			if mk.Comment != "" {
				fmt.Fprintf(w, "  - %s  # %s\n", mk.Key, mk.Comment)
			} else {
				fmt.Fprintf(w, "  - %s\n", mk.Key)
			}
		}
	}

	if len(result.Extra) > 0 {
		fmt.Fprintf(w, "⚠ Extra keys not in template (%d):\n", len(result.Extra))
		for _, k := range result.Extra {
			fmt.Fprintf(w, "  + %s\n", k)
		}
	}
}

// WriteGeneratedTemplate writes a generated template to w in .env format.
func WriteGeneratedTemplate(w io.Writer, tmpl map[string]string, keys []string) {
	fmt.Fprintln(w, "# Generated .env template")
	fmt.Fprintln(w, "# Fill in values before use")
	fmt.Fprintln(w)
	for _, k := range keys {
		v := tmpl[k]
		fmt.Fprintf(w, "%s=%s\n", k, v)
	}
}
