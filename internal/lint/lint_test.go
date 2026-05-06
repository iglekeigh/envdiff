package lint_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/lint"
)

func TestLint_NoFindings_CleanEnv(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
	}
	findings := lint.Lint(env, lint.DefaultRules())
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d: %v", len(findings), findings)
	}
}

func TestRuleNoEmptyKey_FlagsEmptyKey(t *testing.T) {
	findings := lint.RuleNoEmptyKey("", "value")
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Rule != "no-empty-key" {
		t.Errorf("unexpected rule: %s", findings[0].Rule)
	}
	if findings[0].Severity != lint.SeverityError {
		t.Errorf("expected error severity, got %s", findings[0].Severity)
	}
}

func TestRuleNoWhitespaceInKey_FlagsSpaceInKey(t *testing.T) {
	findings := lint.RuleNoWhitespaceInKey("MY KEY", "val")
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Rule != "no-whitespace-in-key" {
		t.Errorf("unexpected rule: %s", findings[0].Rule)
	}
}

func TestRuleNoWhitespaceInKey_CleanKey_NoFindings(t *testing.T) {
	findings := lint.RuleNoWhitespaceInKey("MY_KEY", "val")
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d", len(findings))
	}
}

func TestRuleNoURLAsPlaintext_FlagsHTTP(t *testing.T) {
	findings := lint.RuleNoURLAsPlaintext("API_URL", "http://example.com")
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != lint.SeverityWarning {
		t.Errorf("expected warning severity")
	}
}

func TestRuleNoURLAsPlaintext_HTTPS_NoFindings(t *testing.T) {
	findings := lint.RuleNoURLAsPlaintext("API_URL", "https://example.com")
	if len(findings) != 0 {
		t.Errorf("expected no findings for https URL")
	}
}

func TestRuleValueTooLong_FlagsLongValue(t *testing.T) {
	longVal := strings.Repeat("x", 513)
	findings := lint.RuleValueTooLong("BIG_KEY", longVal)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Rule != "value-too-long" {
		t.Errorf("unexpected rule: %s", findings[0].Rule)
	}
}

func TestRuleValueTooLong_ShortValue_NoFindings(t *testing.T) {
	findings := lint.RuleValueTooLong("KEY", "short")
	if len(findings) != 0 {
		t.Errorf("expected no findings for short value")
	}
}

func TestFinding_String_ContainsFields(t *testing.T) {
	f := lint.Finding{
		Key:      "MY_KEY",
		Rule:     "some-rule",
		Message:  "something is wrong",
		Severity: lint.SeverityWarning,
	}
	s := f.String()
	for _, want := range []string{"MY_KEY", "some-rule", "something is wrong", "warning"} {
		if !strings.Contains(s, want) {
			t.Errorf("expected %q in finding string %q", want, s)
		}
	}
}
