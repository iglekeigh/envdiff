package importer

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return p
}

func TestImport_EnvFile_ParsesKeyValues(t *testing.T) {
	p := writeTempFile(t, ".env", "FOO=bar\nBAZ=qux\n")
	res, err := Import(p, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["FOO"] != "bar" || res.Env["BAZ"] != "qux" {
		t.Errorf("unexpected env: %v", res.Env)
	}
	if res.Format != FormatEnv {
		t.Errorf("expected format env, got %q", res.Format)
	}
}

func TestImport_JSONFile_ParsesStringValues(t *testing.T) {
	p := writeTempFile(t, "config.json", `{"HOST":"localhost","PORT":"8080"}`)
	res, err := Import(p, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["HOST"] != "localhost" || res.Env["PORT"] != "8080" {
		t.Errorf("unexpected env: %v", res.Env)
	}
}

func TestImport_JSONFile_SkipsNestedObjects(t *testing.T) {
	p := writeTempFile(t, "config.json", `{"KEY":"val","NESTED":{"a":1}}`)
	res, err := Import(p, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "NESTED" {
		t.Errorf("expected NESTED in skipped, got %v", res.Skipped)
	}
}

func TestImport_UnsupportedFormat_ReturnsError(t *testing.T) {
	p := writeTempFile(t, "file.yaml", "foo: bar")
	_, err := Import(p, "yaml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestImport_MissingFile_ReturnsError(t *testing.T) {
	_, err := Import("/nonexistent/.env", FormatEnv)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestInferFormat_JSON(t *testing.T) {
	if inferFormat("config.json") != FormatJSON {
		t.Error("expected JSON format for .json extension")
	}
}

func TestInferFormat_Default(t *testing.T) {
	if inferFormat(".env") != FormatEnv {
		t.Error("expected env format for .env extension")
	}
}
