package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/template"
)

func TestWriteTemplateResult_NoIssues(t *testing.T) {
	result := template.Result{Missing: nil, Extra: nil}
	var buf bytes.Buffer
	WriteTemplateResult(&buf, result)

	if !strings.Contains(buf.String(), "✔") {
		t.Errorf("expected pass indicator, got: %s", buf.String())
	}
}

func TestWriteTemplateResult_MissingKeys(t *testing.T) {
	result := template.Result{
		Missing: []template.MissingKey{
			{Key: "DB_HOST", Comment: "database host"},
			{Key: "DB_PORT", Comment: ""},
		},
	}
	var buf bytes.Buffer
	WriteTemplateResult(&buf, result)
	out := buf.String()

	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in output")
	}
	if !strings.Contains(out, "database host") {
		t.Error("expected comment 'database host' in output")
	}
	if !strings.Contains(out, "DB_PORT") {
		t.Error("expected DB_PORT in output")
	}
	if !strings.Contains(out, "Missing keys (2)") {
		t.Error("expected missing keys count in output")
	}
}

func TestWriteTemplateResult_ExtraKeys(t *testing.T) {
	result := template.Result{
		Extra: []string{"LEGACY_FLAG"},
	}
	var buf bytes.Buffer
	WriteTemplateResult(&buf, result)
	out := buf.String()

	if !strings.Contains(out, "LEGACY_FLAG") {
		t.Error("expected LEGACY_FLAG in output")
	}
	if !strings.Contains(out, "Extra keys") {
		t.Error("expected extra keys section in output")
	}
}

func TestWriteGeneratedTemplate_ContainsKeys(t *testing.T) {
	tmpl := map[string]string{
		"HOST": "localhost",
		"PORT": "# required",
	}
	keys := []string{"HOST", "PORT"}
	var buf bytes.Buffer
	WriteGeneratedTemplate(&buf, tmpl, keys)
	out := buf.String()

	if !strings.Contains(out, "HOST=localhost") {
		t.Errorf("expected HOST=localhost in output, got: %s", out)
	}
	if !strings.Contains(out, "PORT=# required") {
		t.Errorf("expected PORT=# required in output, got: %s", out)
	}
	if !strings.Contains(out, "Generated") {
		t.Error("expected header comment in generated template")
	}
}
