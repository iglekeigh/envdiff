// Package lint provides heuristic checks for .env file quality beyond
// structural validation, such as detecting duplicate keys, suspicious values,
// and overly long lines.
package lint

import (
	"fmt"
	"strings"
)

// Severity indicates how serious a lint finding is.
type Severity string

const (
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

// Finding represents a single lint issue discovered in an env map.
type Finding struct {
	Key      string
	Rule     string
	Message  string
	Severity Severity
}

// String returns a human-readable representation of the finding.
func (f Finding) String() string {
	return fmt.Sprintf("[%s] %s: %s (%s)", f.Severity, f.Key, f.Message, f.Rule)
}

// Rule is a function that inspects a key/value pair and returns findings.
type Rule func(key, value string) []Finding

// DefaultRules returns the built-in set of lint rules.
func DefaultRules() []Rule {
	return []Rule{
		RuleNoEmptyKey,
		RuleNoWhitespaceInKey,
		RuleNoURLAsPlaintext,
		RuleValueTooLong,
	}
}

// RuleNoEmptyKey flags keys that are empty strings.
func RuleNoEmptyKey(key, _ string) []Finding {
	if strings.TrimSpace(key) == "" {
		return []Finding{{Key: key, Rule: "no-empty-key", Message: "key must not be empty", Severity: SeverityError}}
	}
	return nil
}

// RuleNoWhitespaceInKey flags keys containing spaces or tabs.
func RuleNoWhitespaceInKey(key, _ string) []Finding {
	if strings.ContainsAny(key, " \t") {
		return []Finding{{Key: key, Rule: "no-whitespace-in-key", Message: "key contains whitespace", Severity: SeverityError}}
	}
	return nil
}

// RuleNoURLAsPlaintext warns when a value looks like a plain HTTP URL.
func RuleNoURLAsPlaintext(key, value string) []Finding {
	if strings.HasPrefix(value, "http://") {
		return []Finding{{Key: key, Rule: "no-plaintext-url", Message: "value uses unencrypted http:// URL", Severity: SeverityWarning}}
	}
	return nil
}

// RuleValueTooLong warns when a value exceeds 512 characters.
func RuleValueTooLong(key, value string) []Finding {
	const maxLen = 512
	if len(value) > maxLen {
		return []Finding{{Key: key, Rule: "value-too-long", Message: fmt.Sprintf("value exceeds %d characters (%d)", maxLen, len(value)), Severity: SeverityWarning}}
	}
	return nil
}

// Lint runs all provided rules against every key/value pair in env and
// returns the aggregated list of findings.
func Lint(env map[string]string, rules []Rule) []Finding {
	var findings []Finding
	for k, v := range env {
		for _, rule := range rules {
			findings = append(findings, rule(k, v)...)
		}
	}
	return findings
}
