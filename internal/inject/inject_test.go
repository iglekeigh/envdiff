package inject_test

import (
	"os"
	"testing"

	"github.com/your/envdiff/internal/inject"
)

func TestIntoMap_NoOptions_InjectsAll(t *testing.T) {
	target := map[string]string{"EXISTING": "old"}
	source := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result, err := inject.IntoMap(target, source, inject.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Injected) != 2 {
		t.Errorf("expected 2 injected, got %d", len(result.Injected))
	}
	if target["FOO"] != "bar" || target["BAZ"] != "qux" {
		t.Error("expected keys to be injected into target")
	}
	if target["EXISTING"] != "old" {
		t.Error("existing key should be preserved")
	}
}

func TestIntoMap_Overwrite_ReplacesExisting(t *testing.T) {
	target := map[string]string{"FOO": "old"}
	source := map[string]string{"FOO": "new"}

	result, err := inject.IntoMap(target, source, inject.Options{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Overwrote) != 1 || result.Overwrote[0] != "FOO" {
		t.Errorf("expected FOO in Overwrote, got %v", result.Overwrote)
	}
	if target["FOO"] != "new" {
		t.Errorf("expected FOO=new, got %q", target["FOO"])
	}
}

func TestIntoMap_NoOverwrite_SkipsExisting(t *testing.T) {
	target := map[string]string{"FOO": "old"}
	source := map[string]string{"FOO": "new"}

	result, err := inject.IntoMap(target, source, inject.Options{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(result.Skipped))
	}
	if target["FOO"] != "old" {
		t.Error("existing key should not be overwritten")
	}
}

func TestIntoMap_PrefixFilter_OnlyMatchingKeys(t *testing.T) {
	target := map[string]string{}
	source := map[string]string{"APP_FOO": "1", "APP_BAR": "2", "OTHER": "3"}

	result, err := inject.IntoMap(target, source, inject.Options{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Injected) != 2 {
		t.Errorf("expected 2 injected, got %d", len(result.Injected))
	}
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(result.Skipped))
	}
	if _, ok := target["OTHER"]; ok {
		t.Error("OTHER should not be injected")
	}
}

func TestIntoMap_StripPrefix_RemovesPrefixFromKey(t *testing.T) {
	target := map[string]string{}
	source := map[string]string{"APP_FOO": "bar"}

	_, err := inject.IntoMap(target, source, inject.Options{Prefix: "APP_", StripPrefix: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if target["FOO"] != "bar" {
		t.Errorf("expected FOO=bar after strip, got %v", target)
	}
	if _, ok := target["APP_FOO"]; ok {
		t.Error("APP_FOO should not exist after strip")
	}
}

func TestIntoOS_InjectsIntoProcessEnv(t *testing.T) {
	os.Unsetenv("ENVDIFF_TEST_INJECT_KEY")
	t.Cleanup(func() { os.Unsetenv("ENVDIFF_TEST_INJECT_KEY") })

	source := map[string]string{"ENVDIFF_TEST_INJECT_KEY": "hello"}
	result, err := inject.IntoOS(source, inject.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Injected) != 1 {
		t.Errorf("expected 1 injected, got %d", len(result.Injected))
	}
	if got := os.Getenv("ENVDIFF_TEST_INJECT_KEY"); got != "hello" {
		t.Errorf("expected hello, got %q", got)
	}
}
