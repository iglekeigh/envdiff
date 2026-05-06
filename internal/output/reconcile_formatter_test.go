package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/output"
	"github.com/user/envdiff/internal/reconcile"
)

func TestWriteReconcileResult_SummaryLine(t *testing.T) {
	result := reconcile.Result{
		Env:     map[string]string{"FOO": "bar"},
		Added:   1,
		Kept:    2,
		Overwritten: 0,
		Conflicts: nil,
	}
	var buf bytes.Buffer
	if err := output.WriteReconcileResult(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "1 added, 2 kept, 0 overwritten") {
		t.Errorf("expected summary line, got: %s", buf.String())
	}
}

func TestWriteReconcileResult_NoConflicts_OmitsConflictSection(t *testing.T) {
	result := reconcile.Result{
		Env:       map[string]string{"A": "1"},
		Conflicts: nil,
	}
	var buf bytes.Buffer
	_ = output.WriteReconcileResult(&buf, result)
	if strings.Contains(buf.String(), "Conflicts resolved") {
		t.Errorf("expected no conflicts section, got: %s", buf.String())
	}
}

func TestWriteReconcileResult_WithConflicts_PrintsConflicts(t *testing.T) {
	result := reconcile.Result{
		Env: map[string]string{"KEY": "base_val"},
		Conflicts: map[string]reconcile.Conflict{
			"KEY": {Base: "base_val", Other: "other_val", Resolved: "base_val"},
		},
		Overwritten: 1,
	}
	var buf bytes.Buffer
	_ = output.WriteReconcileResult(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "Conflicts resolved (1)") {
		t.Errorf("expected conflict header, got: %s", out)
	}
	if !strings.Contains(out, "KEY") {
		t.Errorf("expected KEY in conflict output, got: %s", out)
	}
	if !strings.Contains(out, "base_val") {
		t.Errorf("expected base_val in conflict output, got: %s", out)
	}
}

func TestWriteReconcileResult_EnvSectionSorted(t *testing.T) {
	result := reconcile.Result{
		Env: map[string]string{
			"ZEBRA": "z",
			"ALPHA": "a",
			"MANGO": "m",
		},
	}
	var buf bytes.Buffer
	_ = output.WriteReconcileResult(&buf, result)
	out := buf.String()
	alpha := strings.Index(out, "ALPHA")
	mango := strings.Index(out, "MANGO")
	zebra := strings.Index(out, "ZEBRA")
	if !(alpha < mango && mango < zebra) {
		t.Errorf("expected sorted env output, got: %s", out)
	}
}
