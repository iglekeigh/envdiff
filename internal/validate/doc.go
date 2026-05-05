// Package validate provides rule-based validation for .env file entries.
//
// It supports checking key naming conventions, value constraints, and
// custom user-defined rules. Violations are returned as structured values
// that describe which key failed and why.
//
// Example usage:
//
//	rules := validate.DefaultRules()
//	violations := validate.Validate(envMap, rules)
//	for _, v := range violations {
//		fmt.Println(v.Error())
//	}
package validate
