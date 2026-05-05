package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/output"
)

func makeResult(entries []diff.Entry) *diff.Result {
	return &diff.Result{Entries: entries}
}

func TestFormatter_Text_Match(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatText)
	r := makeResult([]diff.Entry{
		{Key: "FOO", BaseValue: "bar", OtherValue: "bar", Status: diff.StatusMatch},
	})
	if err := f.Write(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "  FOO=bar") {
		t.Errorf("expected match line, got: %q", buf.String())
	}
}

func TestFormatter_Text_Missing(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatText)
	r := makeResult([]diff.Entry{
		{Key: "SECRET", BaseValue: "abc", Status: diff.StatusMissing},
	})
	_ = f.Write(r)
	if !strings.Contains(buf.String(), "- SECRET=abc") {
		t.Errorf("expected missing line, got: %q", buf.String())
	}
}

func TestFormatter_Text_Extra(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatText)
	r := makeResult([]diff.Entry{
		{Key: "NEW_KEY", OtherValue: "xyz", Status: diff.StatusExtra},
	})
	_ = f.Write(r)
	if !strings.Contains(buf.String(), "+ NEW_KEY=xyz") {
		t.Errorf("expected extra line, got: %q", buf.String())
	}
}

func TestFormatter_Text_Changed(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatText)
	r := makeResult([]diff.Entry{
		{Key: "PORT", BaseValue: "8080", OtherValue: "9090", Status: diff.StatusChanged},
	})
	_ = f.Write(r)
	if !strings.Contains(buf.String(), "~ PORT: 8080 -> 9090") {
		t.Errorf("expected changed line, got: %q", buf.String())
	}
}

func TestFormatter_JSON_Output(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatJSON)
	r := makeResult([]diff.Entry{
		{Key: "FOO", BaseValue: "bar", OtherValue: "bar", Status: diff.StatusMatch},
	})
	if err := f.Write(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "\"Key\"") && !strings.Contains(out, "\"key\"") {
		t.Errorf("expected JSON output with key field, got: %q", out)
	}
}
