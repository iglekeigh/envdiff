package validate_test

import (
	"testing"

	"github.com/user/envdiff/internal/validate"
)

func TestValidate_ValidEnv_NoViolations(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"API_KEY":      "abc123",
	}
	rules := validate.DefaultRules()
	violations := validate.Validate(env, rules)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d: %v", len(violations), violations)
	}
}

func TestValidate_InvalidKeyFormat_ReturnsViolation(t *testing.T) {
	env := map[string]string{
		"invalid-key": "value",
	}
	rules := validate.DefaultRules()
	violations := validate.Validate(env, rules)
	if len(violations) == 0 {
		t.Fatal("expected violations for invalid key format")
	}
	found := false
	for _, v := range violations {
		if v.Rule == "key-format" {
			found = true
		}
	}
	if !found {
		t.Error("expected key-format violation")
	}
}

func TestValidate_EmptyValue_ReturnsViolation(t *testing.T) {
	env := map[string]string{
		"VALID_KEY": "",
	}
	rules := validate.DefaultRules()
	violations := validate.Validate(env, rules)
	if len(violations) == 0 {
		t.Fatal("expected violation for empty value")
	}
	if violations[0].Rule != "no-empty-value" {
		t.Errorf("expected no-empty-value rule, got %q", violations[0].Rule)
	}
}

func TestValidate_WhitespaceKey_ReturnsViolation(t *testing.T) {
	env := map[string]string{
		"KEY WITH SPACE": "value",
	}
	rules := validate.DefaultRules()
	violations := validate.Validate(env, rules)
	if len(violations) == 0 {
		t.Fatal("expected violation for whitespace in key")
	}
}

func TestViolation_Error_FormatsMessage(t *testing.T) {
	v := validate.Violation{
		Key:     "BAD KEY",
		Rule:    "no-whitespace-key",
		Message: "key must not contain whitespace",
	}
	msg := v.Error()
	if msg == "" {
		t.Error("expected non-empty error message")
	}
}

func TestValidate_CustomRule(t *testing.T) {
	env := map[string]string{
		"SECRET": "tooshort",
	}
	rules := []validate.Rule{
		{
			Name:    "min-length",
			Message: "value must be at least 16 characters",
			Check: func(_, value string) bool {
				return len(value) >= 16
			},
		},
	}
	violations := validate.Validate(env, rules)
	if len(violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(violations))
	}
}
