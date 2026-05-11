package trim_test

import (
	"testing"

	"github.com/user/envdiff/internal/trim"
)

func TestApply_DefaultOptions_TrimsKeysAndValues(t *testing.T) {
	env := map[string]string{
		"  KEY  ": "  value  ",
	}
	opts := trim.DefaultOptions()
	res := trim.Apply(env, opts)

	if v, ok := res.Env["KEY"]; !ok || v != "value" {
		t.Errorf("expected KEY=value, got %q=%q", "KEY", v)
	}
	if len(res.TrimmedKeys) != 1 {
		t.Errorf("expected 1 trimmed key, got %d", len(res.TrimmedKeys))
	}
	if len(res.TrimmedValues) != 1 {
		t.Errorf("expected 1 trimmed value, got %d", len(res.TrimmedValues))
	}
}

func TestApply_NoTrimKeys_LeavesKeyIntact(t *testing.T) {
	env := map[string]string{" SPACED ": "val"}
	opts := trim.Options{TrimKeys: false, TrimValues: false}
	res := trim.Apply(env, opts)

	if _, ok := res.Env[" SPACED "]; !ok {
		t.Error("expected key with spaces to be preserved")
	}
	if len(res.TrimmedKeys) != 0 {
		t.Errorf("expected 0 trimmed keys, got %d", len(res.TrimmedKeys))
	}
}

func TestApply_RemoveEmpty_DropsEmptyValues(t *testing.T) {
	env := map[string]string{
		"PRESENT": "hello",
		"EMPTY":   "",
		"SPACES":  "   ",
	}
	opts := trim.Options{TrimKeys: true, TrimValues: true, RemoveEmpty: true}
	res := trim.Apply(env, opts)

	if _, ok := res.Env["EMPTY"]; ok {
		t.Error("expected EMPTY to be removed")
	}
	if _, ok := res.Env["SPACES"]; ok {
		t.Error("expected SPACES (trimmed to empty) to be removed")
	}
	if _, ok := res.Env["PRESENT"]; !ok {
		t.Error("expected PRESENT to be retained")
	}
	if len(res.RemovedKeys) != 2 {
		t.Errorf("expected 2 removed keys, got %d", len(res.RemovedKeys))
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"KEY ": " val "}
	opts := trim.DefaultOptions()
	trim.Apply(env, opts)

	if _, ok := env["KEY "]; !ok {
		t.Error("original map should not be mutated")
	}
}

func TestApply_EmptyInput_ReturnsEmptyResult(t *testing.T) {
	res := trim.Apply(map[string]string{}, trim.DefaultOptions())
	if len(res.Env) != 0 {
		t.Errorf("expected empty env, got %d keys", len(res.Env))
	}
}
