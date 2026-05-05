package output_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/output"
	"github.com/user/envdiff/internal/validate"
)

func TestWriteViolations_NoViolations_PrintsPass(t *testing.T) {
	var buf strings.Builder
	n, err := output.WriteViolations(&buf, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 violations reported, got %d", n)
	}
	if !strings.Contains(buf.String(), "passed") {
		t.Errorf("expected 'passed' in output, got: %q", buf.String())
	}
}

func TestWriteViolations_WithViolations_PrintsCount(t *testing.T) {
	violations := []validate.Violation{
		{Key: "BAD_KEY", Rule: "key-format", Message: "must be uppercase"},
		{Key: "EMPTY", Rule: "no-empty-value", Message: "value must not be empty"},
	}
	var buf strings.Builder
	n, err := output.WriteViolations(&buf, violations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 violations, got %d", n)
	}
	out := buf.String()
	if !strings.Contains(out, "2 violation") {
		t.Errorf("expected violation count in output, got: %q", out)
	}
}

func TestWriteViolations_OutputIsSorted(t *testing.T) {
	violations := []validate.Violation{
		{Key: "Z_KEY", Rule: "key-format", Message: "msg"},
		{Key: "A_KEY", Rule: "key-format", Message: "msg"},
	}
	var buf strings.Builder
	_, err := output.WriteViolations(&buf, violations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	aIdx := strings.Index(out, "A_KEY")
	zIdx := strings.Index(out, "Z_KEY")
	if aIdx > zIdx {
		t.Error("expected A_KEY to appear before Z_KEY in sorted output")
	}
}

func TestWriteViolations_ContainsRuleAndKey(t *testing.T) {
	violations := []validate.Violation{
		{Key: "MY_VAR", Rule: "no-empty-value", Message: "value must not be empty"},
	}
	var buf strings.Builder
	_, _ = output.WriteViolations(&buf, violations)
	out := buf.String()
	if !strings.Contains(out, "MY_VAR") {
		t.Error("expected key in output")
	}
	if !strings.Contains(out, "no-empty-value") {
		t.Error("expected rule name in output")
	}
}
