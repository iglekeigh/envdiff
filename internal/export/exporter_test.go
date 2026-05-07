package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExport_EnvFormat_WritesKeyValues(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.env")
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result, err := Export(env, path, FormatEnv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Count != 2 {
		t.Errorf("expected count 2, got %d", result.Count)
	}

	content, _ := os.ReadFile(path)
	if !strings.Contains(string(content), "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got: %s", content)
	}
}

func TestExport_JSONFormat_WritesValidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.json")
	env := map[string]string{"KEY": "value"}

	_, err := Export(env, path, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, _ := os.ReadFile(path)
	var parsed map[string]string
	if err := json.Unmarshal(content, &parsed); err != nil {
		t.Errorf("output is not valid JSON: %v", err)
	}
	if parsed["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %v", parsed)
	}
}

func TestExport_ShellFormat_WritesExportStatements(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.sh")
	env := map[string]string{"MY_VAR": "hello"}

	result, err := Export(env, path, FormatShell)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, _ := os.ReadFile(path)
	if !strings.Contains(string(content), "export MY_VAR=") {
		t.Errorf("expected export statement, got: %s", content)
	}
	if len(result.Skipped) != 0 {
		t.Errorf("expected no skipped keys, got %v", result.Skipped)
	}
}

func TestExport_ShellFormat_SkipsInvalidKeys(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.sh")
	env := map[string]string{"invalid-key": "val", "VALID": "ok"}

	result, err := Export(env, path, FormatShell)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Skipped) != 1 || result.Skipped[0] != "invalid-key" {
		t.Errorf("expected invalid-key skipped, got %v", result.Skipped)
	}
}

func TestExport_InferFormat_FromExtension(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	env := map[string]string{"A": "1"}

	_, err := Export(env, path, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, _ := os.ReadFile(path)
	var parsed map[string]string
	if err := json.Unmarshal(content, &parsed); err != nil {
		t.Errorf("expected JSON output for .json extension: %v", err)
	}
}

func TestExport_UnsupportedFormat_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.xml")
	env := map[string]string{"X": "y"}

	_, err := Export(env, path, Format("xml"))
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}
