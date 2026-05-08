// Package validate provides rule-based validation for .env file entries.
//
// It supports checking key naming conventions, value constraints, and
// custom user-defined rules. Violations are returned as structured values
// that describe which key failed and why.
//
// # Rules
//
// A Rule is a function that inspects a key-value pair and returns a
// Violation if the entry does not meet the rule's requirements, or nil
// if the entry is valid. Rules can be composed and passed to Validate.
//
// # Violations
//
// A Violation describes a single validation failure. It includes the
// offending key, a human-readable message, and an optional severity
// level (e.g. warning vs. error).
//
// # Example usage
//
//	rules := validate.DefaultRules()
//	violations := validate.Validate(envMap, rules)
//	for _, v := range violations {
//		fmt.Println(v.Error())
//	}
package validate
