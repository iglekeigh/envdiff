package reconcile_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/reconcile"
)

func TestReconcile_AddsExtraKeysFromOther(t *testing.T) {
	base := map[string]string{"A": "1"}
	other := map[string]string{"A": "1", "B": "2"}

	res, err := reconcile.Reconcile(base, other, reconcile.PreferBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["B"] != "2" {
		t.Errorf("expected B=2, got %q", res.Merged["B"])
	}
}

func TestReconcile_RetainsBaseOnlyKeys(t *testing.T) {
	base := map[string]string{"A": "1", "C": "3"}
	other := map[string]string{"A": "1"}

	res, err := reconcile.Reconcile(base, other, reconcile.PreferBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["C"] != "3" {
		t.Errorf("expected C=3, got %q", res.Merged["C"])
	}
}

func TestReconcile_PreferBase_KeepsBaseOnConflict(t *testing.T) {
	base := map[string]string{"KEY": "base_val"}
	other := map[string]string{"KEY": "other_val"}

	res, err := reconcile.Reconcile(base, other, reconcile.PreferBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["KEY"] != "base_val" {
		t.Errorf("expected base_val, got %q", res.Merged["KEY"])
	}
}

func TestReconcile_PreferOther_UsesOtherOnConflict(t *testing.T) {
	base := map[string]string{"KEY": "base_val"}
	other := map[string]string{"KEY": "other_val"}

	res, err := reconcile.Reconcile(base, other, reconcile.PreferOther)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["KEY"] != "other_val" {
		t.Errorf("expected other_val, got %q", res.Merged["KEY"])
	}
}

func TestReconcile_ErrorOnConflict_ReturnsError(t *testing.T) {
	base := map[string]string{"KEY": "base_val"}
	other := map[string]string{"KEY": "other_val"}

	_, err := reconcile.Reconcile(base, other, reconcile.ErrorOnConflict)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "KEY") {
		t.Errorf("expected error to mention KEY, got: %v", err)
	}
}

func TestReconcile_WriteWarnings(t *testing.T) {
	base := map[string]string{"A": "1"}
	other := map[string]string{"A": "1", "B": "2"}

	res, _ := reconcile.Reconcile(base, other, reconcile.PreferBase)
	var buf strings.Builder
	res.WriteWarnings(&buf)
	if !strings.Contains(buf.String(), "added missing key") {
		t.Errorf("expected warning about missing key, got: %q", buf.String())
	}
}
