// Package validate provides validation utilities for .env file entries.
package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule represents a validation rule applied to env keys or values.
type Rule struct {
	Name    string
	Message string
	Check   func(key, value string) bool
}

// Violation represents a failed validation for a specific key.
type Violation struct {
	Key     string
	Rule    string
	Message string
}

func (v Violation) Error() string {
	return fmt.Sprintf("key %q violated rule %q: %s", v.Key, v.Rule, v.Message)
}

var validKeyPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// DefaultRules returns the standard set of validation rules.
func DefaultRules() []Rule {
	return []Rule{
		{
			Name:    "key-format",
			Message: "key must be uppercase alphanumeric with underscores, starting with a letter",
			Check: func(key, _ string) bool {
				return validKeyPattern.MatchString(key)
			},
		},
		{
			Name:    "no-empty-value",
			Message: "value must not be empty",
			Check: func(_, value string) bool {
				return strings.TrimSpace(value) != ""
			},
		},
		{
			Name:    "no-whitespace-key",
			Message: "key must not contain whitespace",
			Check: func(key, _ string) bool {
				return !strings.ContainsAny(key, " \t")
			},
		},
	}
}

// Validate runs all rules against the provided env map and returns any violations.
func Validate(env map[string]string, rules []Rule) []Violation {
	var violations []Violation
	for key, value := range env {
		for _, rule := range rules {
			if !rule.Check(key, value) {
				violations = append(violations, Violation{
					Key:     key,
					Rule:    rule.Name,
					Message: rule.Message,
				})
			}
		}
	}
	return violations
}
