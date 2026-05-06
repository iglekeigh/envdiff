package schema

import (
	"regexp"
	"testing"
)

func TestValidate_AllRequiredPresent_NoViolations(t *testing.T) {
	s := Schema{Fields: []Field{
		{Key: "APP_NAME", Type: TypeString, Required: true},
		{Key: "PORT", Type: TypeInt, Required: true},
	}}
	env := map[string]string{"APP_NAME": "myapp", "PORT": "8080"}
	v := Validate(s, env)
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %v", v)
	}
}

func TestValidate_MissingRequiredKey_ReturnsViolation(t *testing.T) {
	s := Schema{Fields: []Field{
		{Key: "DATABASE_URL", Type: TypeURL, Required: true},
	}}
	v := Validate(s, map[string]string{})
	if len(v) != 1 || v[0].Key != "DATABASE_URL" {
		t.Fatalf("expected violation for DATABASE_URL, got %v", v)
	}
}

func TestValidate_OptionalMissingKey_NoViolation(t *testing.T) {
	s := Schema{Fields: []Field{
		{Key: "OPTIONAL_KEY", Type: TypeString, Required: false},
	}}
	v := Validate(s, map[string]string{})
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %v", v)
	}
}

func TestValidate_InvalidInt_ReturnsViolation(t *testing.T) {
	s := Schema{Fields: []Field{
		{Key: "PORT", Type: TypeInt, Required: true},
	}}
	v := Validate(s, map[string]string{"PORT": "abc"})
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %v", v)
	}
}

func TestValidate_ValidBool_NoViolation(t *testing.T) {
	s := Schema{Fields: []Field{
		{Key: "DEBUG", Type: TypeBool, Required: true},
	}}
	for _, val := range []string{"true", "false", "1", "0", "yes", "no"} {
		v := Validate(s, map[string]string{"DEBUG": val})
		if len(v) != 0 {
			t.Errorf("expected no violations for %q, got %v", val, v)
		}
	}
}

func TestValidate_InvalidURL_ReturnsViolation(t *testing.T) {
	s := Schema{Fields: []Field{
		{Key: "API_URL", Type: TypeURL, Required: true},
	}}
	v := Validate(s, map[string]string{"API_URL": "not-a-url"})
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %v", v)
	}
}

func TestValidate_PatternMismatch_ReturnsViolation(t *testing.T) {
	s := Schema{Fields: []Field{
		{Key: "REGION", Type: TypeString, Required: true, Pattern: regexp.MustCompile(`^us-`)},
	}}
	v := Validate(s, map[string]string{"REGION": "eu-west-1"})
	if len(v) != 1 {
		t.Fatalf("expected 1 violation for pattern mismatch, got %v", v)
	}
}

func TestViolation_Error_FormatsMessage(t *testing.T) {
	v := Violation{Key: "FOO", Message: "required key is missing"}
	got := v.Error()
	if got != "schema violation [FOO]: required key is missing" {
		t.Errorf("unexpected error string: %q", got)
	}
}
