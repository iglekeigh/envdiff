package output

import (
	"fmt"
	"io"
	"sort"
)

// WriteEncryptedEnv writes the encrypted environment map to w in KEY=VALUE format.
// Encrypted values are opaque ciphertext strings.
func WriteEncryptedEnv(w io.Writer, env map[string]string) {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(w, "%s=%s\n", k, env[k])
	}
}

// WriteEncryptSummary writes a human-readable summary of the encryption operation.
func WriteEncryptSummary(w io.Writer, total, encrypted int, failed []string) {
	fmt.Fprintf(w, "Encrypted %d/%d keys\n", encrypted, total)

	if len(failed) > 0 {
		sort.Strings(failed)
		fmt.Fprintf(w, "Failed keys (%d):\n", len(failed))
		for _, k := range failed {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}
}

// WriteDecryptSummary writes a human-readable summary of the decryption operation.
func WriteDecryptSummary(w io.Writer, total, decrypted int, failed []string) {
	fmt.Fprintf(w, "Decrypted %d/%d keys\n", decrypted, total)

	if len(failed) > 0 {
		sort.Strings(failed)
		fmt.Fprintf(w, "Failed keys (%d):\n", len(failed))
		for _, k := range failed {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}
}
