package audit_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/audit"
)

func TestNew_SetsSource(t *testing.T) {
	l := audit.New("prod")
	if l == nil {
		t.Fatal("expected non-nil log")
	}
	if l.Len() != 0 {
		t.Errorf("expected empty log, got %d entries", l.Len())
	}
}

func TestRecord_AppendsEntry(t *testing.T) {
	l := audit.New("test")
	l.Record("DB_PASS", audit.ActionAdded, "", "secret")
	if l.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", l.Len())
	}
	entries := l.Entries()
	e := entries[0]
	if e.Key != "DB_PASS" {
		t.Errorf("expected key DB_PASS, got %s", e.Key)
	}
	if e.Action != audit.ActionAdded {
		t.Errorf("expected action added, got %s", e.Action)
	}
	if e.NewValue != "secret" {
		t.Errorf("expected new value 'secret', got %s", e.NewValue)
	}
	if e.Source != "test" {
		t.Errorf("expected source 'test', got %s", e.Source)
	}
}

func TestEntries_ReturnsCopy(t *testing.T) {
	l := audit.New("src")
	l.Record("KEY", audit.ActionChanged, "old", "new")
	a := l.Entries()
	a[0].Key = "MUTATED"
	b := l.Entries()
	if b[0].Key == "MUTATED" {
		t.Error("Entries() should return a copy, not a reference")
	}
}

func TestEntry_String_Added(t *testing.T) {
	l := audit.New("ci")
	l.Record("API_KEY", audit.ActionAdded, "", "xyz")
	s := l.Entries()[0].String()
	if !strings.Contains(s, "ADD") || !strings.Contains(s, "API_KEY") {
		t.Errorf("unexpected string format: %s", s)
	}
}

func TestEntry_String_Removed(t *testing.T) {
	l := audit.New("ci")
	l.Record("OLD_KEY", audit.ActionRemoved, "val", "")
	s := l.Entries()[0].String()
	if !strings.Contains(s, "REMOVE") || !strings.Contains(s, "OLD_KEY") {
		t.Errorf("unexpected string format: %s", s)
	}
}

func TestEntry_String_Changed(t *testing.T) {
	l := audit.New("ci")
	l.Record("PORT", audit.ActionChanged, "3000", "8080")
	s := l.Entries()[0].String()
	if !strings.Contains(s, "CHANGE") || !strings.Contains(s, "PORT") {
		t.Errorf("unexpected string format: %s", s)
	}
}

func TestEntry_String_Redacted(t *testing.T) {
	l := audit.New("ci")
	l.Record("SECRET", audit.ActionRedacted, "", "")
	s := l.Entries()[0].String()
	if !strings.Contains(s, "REDACT") || !strings.Contains(s, "SECRET") {
		t.Errorf("unexpected string format: %s", s)
	}
}
