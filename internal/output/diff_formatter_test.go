package output_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/output"
)

func makeDiffResult(entries []diff.Entry) diff.Result {
	return diff.Result{Entries: func() []diff.Entry { return entries }}
}

func TestWriteDiffSummary_NoDifferences(t *testing.T) {
	var buf strings.Builder
	result := diff.NewResult(nil)
	output.WriteDiffSummary(&buf, result)
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-differences message, got: %q", buf.String())
	}
}

func TestWriteDiffSummary_MissingKey(t *testing.T) {
	var buf strings.Builder
	base := map[string]string{"FOO": "bar"}
	other := map[string]string{}
	result := diff.Compare(base, other)
	output.WriteDiffSummary(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "MISSING") {
		t.Errorf("expected MISSING label, got: %q", out)
	}
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected key FOO in output, got: %q", out)
	}
}

func TestWriteDiffSummary_ExtraKey(t *testing.T) {
	var buf strings.Builder
	base := map[string]string{}
	other := map[string]string{"NEW_KEY": "value"}
	result := diff.Compare(base, other)
	output.WriteDiffSummary(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "EXTRA") {
		t.Errorf("expected EXTRA label, got: %q", out)
	}
	if !strings.Contains(out, "NEW_KEY") {
		t.Errorf("expected key NEW_KEY in output, got: %q", out)
	}
}

func TestWriteDiffSummary_ChangedValue(t *testing.T) {
	var buf strings.Builder
	base := map[string]string{"HOST": "localhost"}
	other := map[string]string{"HOST": "production.example.com"}
	result := diff.Compare(base, other)
	output.WriteDiffSummary(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "CHANGED") {
		t.Errorf("expected CHANGED label, got: %q", out)
	}
	if !strings.Contains(out, "HOST") {
		t.Errorf("expected key HOST in output, got: %q", out)
	}
}

func TestWriteDiffSummary_SummaryLine(t *testing.T) {
	var buf strings.Builder
	base := map[string]string{"A": "1", "B": "old"}
	other := map[string]string{"B": "new", "C": "3"}
	result := diff.Compare(base, other)
	output.WriteDiffSummary(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "Summary:") {
		t.Errorf("expected Summary line, got: %q", out)
	}
	if !strings.Contains(out, "1 missing") {
		t.Errorf("expected 1 missing in summary, got: %q", out)
	}
	if !strings.Contains(out, "1 extra") {
		t.Errorf("expected 1 extra in summary, got: %q", out)
	}
	if !strings.Contains(out, "1 changed") {
		t.Errorf("expected 1 changed in summary, got: %q", out)
	}
}
