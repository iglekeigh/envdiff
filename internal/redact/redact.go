// Package redact provides utilities for redacting sensitive values
// in environment variable maps before display or output.
package redact

import "strings"

const RedactedValue = "***REDACTED***"

// DefaultSensitiveKeys contains common patterns for sensitive environment variable names.
var DefaultSensitiveKeys = []string{
	"PASSWORD",
	"SECRET",
	"TOKEN",
	"KEY",
	"PRIVATE",
	"CREDENTIAL",
	"AUTH",
	"API_KEY",
	"ACCESS_KEY",
	"CERT",
}

// Redactor holds configuration for redacting sensitive values.
type Redactor struct {
	SensitivePatterns []string
}

// New creates a Redactor with the default sensitive key patterns.
func New() *Redactor {
	return &Redactor{
		SensitivePatterns: DefaultSensitiveKeys,
	}
}

// NewWithPatterns creates a Redactor with custom sensitive key patterns.
func NewWithPatterns(patterns []string) *Redactor {
	return &Redactor{
		SensitivePatterns: patterns,
	}
}

// IsSensitive returns true if the key matches any sensitive pattern.
func (r *Redactor) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, pattern := range r.SensitivePatterns {
		if strings.Contains(upper, strings.ToUpper(pattern)) {
			return true
		}
	}
	return false
}

// Apply returns a new map with sensitive values replaced by RedactedValue.
func (r *Redactor) Apply(env map[string]string) map[string]string {
	result := make(map[string]string, len(env))
	for k, v := range env {
		if r.IsSensitive(k) {
			result[k] = RedactedValue
		} else {
			result[k] = v
		}
	}
	return result
}
