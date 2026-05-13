package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/clone"
)

func TestWriteCloneResult_UnchangedKeys(t *testing.T) {
	res := &clone.Result{
		Env:    map[string]string{"FOO": "1"},
		Mapped: map[string]string{"FOO": "FOO"},
	}
	var buf bytes.Buffer
	WriteCloneResult(&buf, res)
	out := buf.String()
	if !strings.Contains(out, "Cloned 1 key(s)") {
		t.Errorf("expected summary line, got: %s", out)
	}
	if !strings.Contains(out, "FOO (unchanged)") {
		t.Errorf("expected unchanged label, got: %s", out)
	}
}

func TestWriteCloneResult_RenamedKeys(t *testing.T) {
	res := &clone.Result{
		Env:    map[string]string{"DB_HOST": "localhost"},
		Mapped: map[string]string{"HOST": "DB_HOST"},
	}
	var buf bytes.Buffer
	WriteCloneResult(&buf, res)
	out := buf.String()
	if !strings.Contains(out, "HOST -> DB_HOST") {
		t.Errorf("expected rename arrow, got: %s", out)
	}
}

func TestWriteCloneResult_WithSkipped(t *testing.T) {
	res := &clone.Result{
		Env:     map[string]string{"HOST": "localhost"},
		Mapped:  map[string]string{"DB_HOST": "HOST"},
		Skipped: []string{"APP_PORT"},
	}
	var buf bytes.Buffer
	WriteCloneResult(&buf, res)
	out := buf.String()
	if !strings.Contains(out, "Skipped 1 key(s)") {
		t.Errorf("expected skipped section, got: %s", out)
	}
	if !strings.Contains(out, "APP_PORT") {
		t.Errorf("expected APP_PORT in skipped, got: %s", out)
	}
}

func TestWriteClonedEnv_OutputIsSorted(t *testing.T) {
	env := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	var buf bytes.Buffer
	WriteClonedEnv(&buf, env)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected A_KEY first, got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Z_KEY") {
		t.Errorf("expected Z_KEY last, got: %s", lines[2])
	}
}

func TestWriteClonedEnv_FormatsKeyValue(t *testing.T) {
	env := map[string]string{"HOST": "localhost"}
	var buf bytes.Buffer
	WriteClonedEnv(&buf, env)
	if !strings.Contains(buf.String(), "HOST=localhost") {
		t.Errorf("expected KEY=VALUE format, got: %s", buf.String())
	}
}
