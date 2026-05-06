// Package schema provides tools for defining and enforcing a schema
// against .env files, ensuring required keys are present and values
// match expected types or patterns.
package schema

import (
	"fmt"
	"regexp"
)

// FieldType represents the expected type of an env value.
type FieldType string

const (
	TypeString  FieldType = "string"
	TypeInt     FieldType = "int"
	TypeBool    FieldType = "bool"
	TypeURL     FieldType = "url"
)

// Field defines a single schema entry.
type Field struct {
	Key      string
	Type     FieldType
	Required bool
	Pattern  *regexp.Regexp // optional additional pattern constraint
}

// Schema is a collection of field definitions.
type Schema struct {
	Fields []Field
}

// Violation represents a schema validation failure.
type Violation struct {
	Key     string
	Message string
}

func (v Violation) Error() string {
	return fmt.Sprintf("schema violation [%s]: %s", v.Key, v.Message)
}

var (
	reInt  = regexp.MustCompile(`^-?\d+$`)
	reBool = regexp.MustCompile(`^(true|false|1|0|yes|no)$`)
	reURL  = regexp.MustCompile(`^https?://`)
)

// Validate checks env against the schema and returns any violations.
func Validate(s Schema, env map[string]string) []Violation {
	var violations []Violation

	for _, field := range s.Fields {
		val, ok := env[field.Key]
		if !ok {
			if field.Required {
				violations = append(violations, Violation{
					Key:     field.Key,
					Message: "required key is missing",
				})
			}
			continue
		}

		if v := checkType(field, val); v != nil {
			violations = append(violations, *v)
			continue
		}

		if field.Pattern != nil && !field.Pattern.MatchString(val) {
			violations = append(violations, Violation{
				Key:     field.Key,
				Message: fmt.Sprintf("value %q does not match pattern %s", val, field.Pattern),
			})
		}
	}

	return violations
}

func checkType(field Field, val string) *Violation {
	switch field.Type {
	case TypeInt:
		if !reInt.MatchString(val) {
			return &Violation{Key: field.Key, Message: fmt.Sprintf("expected int, got %q", val)}
		}
	case TypeBool:
		if !reBool.MatchString(val) {
			return &Violation{Key: field.Key, Message: fmt.Sprintf("expected bool, got %q", val)}
		}
	case TypeURL:
		if !reURL.MatchString(val) {
			return &Violation{Key: field.Key, Message: fmt.Sprintf("expected URL, got %q", val)}
		}
	}
	return nil
}
