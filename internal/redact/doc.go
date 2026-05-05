// Package redact provides functionality to identify and redact sensitive
// environment variable values before they are displayed or written to output.
//
// Sensitive keys are identified by matching against a configurable list of
// patterns (e.g., "PASSWORD", "SECRET", "TOKEN"). Matched values are replaced
// with a placeholder string to prevent accidental exposure of credentials.
//
// Usage:
//
//	r := redact.New()
//	safeEnv := r.Apply(envMap)
//
// Custom patterns can be provided via NewWithPatterns.
package redact
