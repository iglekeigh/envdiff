package output

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/redact"
)

// WriteRedactedEnv writes a redacted view of an env map to the given writer.
// Sensitive values are replaced with the redactor's placeholder, while
// non-sensitive values are printed as-is.
func WriteRedactedEnv(w io.Writer, env map[string]string, r *redact.Redactor) error {
	redacted := r.Apply(env)

	keys := make([]string, 0, len(redacted))
	for k := range redacted {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		_, err := fmt.Fprintf(w, "%s=%s\n", k, redacted[k])
		if err != nil {
			return fmt.Errorf("write redacted env: %w", err)
		}
	}
	return nil
}

// WriteRedactSummary writes a summary of which keys were redacted.
func WriteRedactSummary(w io.Writer, env map[string]string, r *redact.Redactor) error {
	var redactedKeys []string
	for k := range env {
		if r.IsSensitive(k) {
			redactedKeys = append(redactedKeys, k)
		}
	}
	sort.Strings(redactedKeys)

	if len(redactedKeys) == 0 {
		_, err := fmt.Fprintln(w, "No sensitive keys detected.")
		return err
	}

	_, err := fmt.Fprintf(w, "Redacted %d sensitive key(s):\n", len(redactedKeys))
	if err != nil {
		return err
	}
	for _, k := range redactedKeys {
		_, err := fmt.Fprintf(w, "  - %s\n", k)
		if err != nil {
			return err
		}
	}
	return nil
}
