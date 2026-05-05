package redact_test

import (
	"testing"

	"github.com/user/envdiff/internal/redact"
)

func TestIsSensitive_MatchesDefaultPatterns(t *testing.T) {
	r := redact.New()
	sensitiveKeys := []string{
		"DB_PASSWORD",
		"API_SECRET",
		"AUTH_TOKEN",
		"PRIVATE_KEY",
		"AWS_ACCESS_KEY",
	}
	for _, key := range sensitiveKeys {
		if !r.IsSensitive(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}
}

func TestIsSensitive_NotSensitive(t *testing.T) {
	r := redact.New()
	safeKeys := []string{
		"APP_NAME",
		"PORT",
		"DEBUG",
		"LOG_LEVEL",
	}
	for _, key := range safeKeys {
		if r.IsSensitive(key) {
			t.Errorf("expected %q to NOT be sensitive", key)
		}
	}
}

func TestApply_RedactsSensitiveValues(t *testing.T) {
	r := redact.New()
	env := map[string]string{
		"DB_PASSWORD": "supersecret",
		"APP_NAME":    "myapp",
		"API_TOKEN":   "abc123",
	}
	result := r.Apply(env)
	if result["DB_PASSWORD"] != redact.RedactedValue {
		t.Errorf("expected DB_PASSWORD to be redacted")
	}
	if result["API_TOKEN"] != redact.RedactedValue {
		t.Errorf("expected API_TOKEN to be redacted")
	}
	if result["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME to remain unchanged")
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	r := redact.New()
	env := map[string]string{
		"DB_PASSWORD": "supersecret",
	}
	_ = r.Apply(env)
	if env["DB_PASSWORD"] != "supersecret" {
		t.Errorf("original map should not be mutated")
	}
}

func TestNewWithPatterns_CustomPatterns(t *testing.T) {
	r := redact.NewWithPatterns([]string{"INTERNAL"})
	if !r.IsSensitive("INTERNAL_URL") {
		t.Errorf("expected INTERNAL_URL to be sensitive with custom pattern")
	}
	if r.IsSensitive("DB_PASSWORD") {
		t.Errorf("expected DB_PASSWORD to NOT be sensitive with custom pattern")
	}
}
