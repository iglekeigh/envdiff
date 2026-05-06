package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/output"
)

func TestWriteEncryptedEnv_OutputIsSorted(t *testing.T) {
	env := map[string]string{
		"ZEBRA": "enc:abc123",
		"ALPHA": "enc:xyz789",
		"MIDDLE": "enc:def456",
	}

	var buf bytes.Buffer
	output.WriteEncryptedEnv(&buf, env)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "ALPHA=") {
		t.Errorf("expected first line to start with ALPHA=, got %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "ZEBRA=") {
		t.Errorf("expected last line to start with ZEBRA=, got %s", lines[2])
	}
}

func TestWriteEncryptedEnv_ContainsValues(t *testing.T) {
	env := map[string]string{
		"API_KEY": "enc:supersecret",
	}

	var buf bytes.Buffer
	output.WriteEncryptedEnv(&buf, env)

	if !strings.Contains(buf.String(), "API_KEY=enc:supersecret") {
		t.Errorf("expected output to contain API_KEY=enc:supersecret, got: %s", buf.String())
	}
}

func TestWriteEncryptSummary_AllSucceeded(t *testing.T) {
	var buf bytes.Buffer
	output.WriteEncryptSummary(&buf, 5, 5, nil)

	out := buf.String()
	if !strings.Contains(out, "Encrypted 5/5 keys") {
		t.Errorf("unexpected output: %s", out)
	}
	if strings.Contains(out, "Failed") {
		t.Errorf("should not mention failed keys when none failed")
	}
}

func TestWriteEncryptSummary_WithFailures(t *testing.T) {
	var buf bytes.Buffer
	output.WriteEncryptSummary(&buf, 5, 3, []string{"BAD_KEY", "ANOTHER"})

	out := buf.String()
	if !strings.Contains(out, "Encrypted 3/5 keys") {
		t.Errorf("expected count line, got: %s", out)
	}
	if !strings.Contains(out, "Failed keys (2)") {
		t.Errorf("expected failed section, got: %s", out)
	}
	if !strings.Contains(out, "ANOTHER") || !strings.Contains(out, "BAD_KEY") {
		t.Errorf("expected failed key names in output, got: %s", out)
	}
}

func TestWriteDecryptSummary_NoFailures(t *testing.T) {
	var buf bytes.Buffer
	output.WriteDecryptSummary(&buf, 4, 4, nil)

	out := buf.String()
	if !strings.Contains(out, "Decrypted 4/4 keys") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestWriteDecryptSummary_FailedKeysSorted(t *testing.T) {
	var buf bytes.Buffer
	output.WriteDecryptSummary(&buf, 3, 1, []string{"Z_KEY", "A_KEY"})

	out := buf.String()
	idxA := strings.Index(out, "A_KEY")
	idxZ := strings.Index(out, "Z_KEY")
	if idxA == -1 || idxZ == -1 {
		t.Fatalf("expected both keys in output, got: %s", out)
	}
	if idxA > idxZ {
		t.Errorf("expected A_KEY before Z_KEY in sorted output")
	}
}
