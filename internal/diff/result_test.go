package diff_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func TestResult_HasDifferences_True(t *testing.T) {
	r := &diff.Result{
		Entries: []diff.Entry{
			{Key: "FOO", Status: diff.StatusMatch, BaseVal: "bar"},
			{Key: "BAZ", Status: diff.StatusChanged, BaseVal: "old", OtherVal: "new"},
		},
	}
	if !r.HasDifferences() {
		t.Error("expected HasDifferences to return true")
	}
}

func TestResult_HasDifferences_False(t *testing.T) {
	r := &diff.Result{
		Entries: []diff.Entry{
			{Key: "FOO", Status: diff.StatusMatch, BaseVal: "bar"},
		},
	}
	if r.HasDifferences() {
		t.Error("expected HasDifferences to return false")
	}
}

func TestResult_Filter(t *testing.T) {
	r := &diff.Result{
		Entries: []diff.Entry{
			{Key: "A", Status: diff.StatusMatch},
			{Key: "B", Status: diff.StatusMissing},
			{Key: "C", Status: diff.StatusExtra},
		},
	}
	missing := r.Filter(diff.StatusMissing)
	if len(missing.Entries) != 1 || missing.Entries[0].Key != "B" {
		t.Errorf("expected 1 missing entry with key B, got %+v", missing.Entries)
	}
}

func TestResult_WriteTo_Output(t *testing.T) {
	r := &diff.Result{
		Entries: []diff.Entry{
			{Key: "APP", Status: diff.StatusMatch, BaseVal: "myapp"},
			{Key: "PORT", Status: diff.StatusChanged, BaseVal: "3000", OtherVal: "4000"},
			{Key: "OLD", Status: diff.StatusMissing, BaseVal: "gone"},
			{Key: "NEW", Status: diff.StatusExtra, OtherVal: "added"},
		},
	}
	var buf bytes.Buffer
	r.WriteTo(&buf)
	out := buf.String()

	if !strings.Contains(out, "  APP=myapp") {
		t.Errorf("expected match line for APP, got:\n%s", out)
	}
	if !strings.Contains(out, "~ PORT: 3000 -> 4000") {
		t.Errorf("expected changed line for PORT, got:\n%s", out)
	}
	if !strings.Contains(out, "- OLD=gone") {
		t.Errorf("expected missing line for OLD, got:\n%s", out)
	}
	if !strings.Contains(out, "+ NEW=added") {
		t.Errorf("expected extra line for NEW, got:\n%s", out)
	}
}
