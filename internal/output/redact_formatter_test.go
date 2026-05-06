package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/output"
	"github.com/user/envdiff/internal/redact"
)

func TestWriteRedactedEnv_RedactsSensitiveKeys(t *testing.T) {
	env := map[string]string{
		"API_KEY":  "supersecret",
		"APP_NAME": "myapp",
	}
	r := redact.New()
	var buf bytes.Buffer
	err := output.WriteRedactedEnv(&buf, env, r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "supersecret") {
		t.Error("expected sensitive value to be redacted")
	}
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Error("expected non-sensitive value to be preserved")
	}
}

func TestWriteRedactedEnv_OutputIsSorted(t *testing.T) {
	env := map[string]string{
		"Z_VAR": "z",
		"A_VAR": "a",
		"M_VAR": "m",
	}
	r := redact.New()
	var buf bytes.Buffer
	_ = output.WriteRedactedEnv(&buf, env, r)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A_VAR") {
		t.Errorf("expected first line to start with A_VAR, got %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Z_VAR") {
		t.Errorf("expected last line to start with Z_VAR, got %s", lines[2])
	}
}

func TestWriteRedactSummary_NoSensitiveKeys(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
	}
	r := redact.New()
	var buf bytes.Buffer
	err := output.WriteRedactSummary(&buf, env, r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No sensitive keys") {
		t.Errorf("expected no-sensitive message, got: %s", buf.String())
	}
}

func TestWriteRedactSummary_WithSensitiveKeys(t *testing.T) {
	env := map[string]string{
		"SECRET_TOKEN": "abc",
		"DB_PASSWORD":  "pass",
		"APP_NAME":     "myapp",
	}
	r := redact.New()
	var buf bytes.Buffer
	err := output.WriteRedactSummary(&buf, env, r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Redacted 2 sensitive key(s)") {
		t.Errorf("expected redact count in summary, got: %s", out)
	}
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Error("expected DB_PASSWORD in summary")
	}
	if !strings.Contains(out, "SECRET_TOKEN") {
		t.Error("expected SECRET_TOKEN in summary")
	}
}
