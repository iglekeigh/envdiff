package rename

import (
	"testing"
)

func TestApply_ExactRename_ChangesKey(t *testing.T) {
	env := map[string]string{"OLD_KEY": "value", "OTHER": "x"}
	rules := []Rule{{From: "OLD_KEY", To: "NEW_KEY", Mode: ModeExact}}
	res, err := Apply(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["NEW_KEY"]; !ok {
		t.Error("expected NEW_KEY to exist")
	}
	if _, ok := res.Env["OLD_KEY"]; ok {
		t.Error("expected OLD_KEY to be removed")
	}
	if len(res.Renamed) != 1 || res.Renamed[0].OldKey != "OLD_KEY" {
		t.Errorf("unexpected renamed list: %v", res.Renamed)
	}
}

func TestApply_PrefixRename_ChangesMatchingKeys(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "DB_HOST": "db"}
	rules := []Rule{{From: "APP_", To: "SVC_", Mode: ModePrefix}}
	res, err := Apply(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["SVC_HOST"] != "localhost" {
		t.Errorf("expected SVC_HOST=localhost, got %q", res.Env["SVC_HOST"])
	}
	if res.Env["SVC_PORT"] != "8080" {
		t.Errorf("expected SVC_PORT=8080, got %q", res.Env["SVC_PORT"])
	}
	if res.Env["DB_HOST"] != "db" {
		t.Error("DB_HOST should be unchanged")
	}
	if len(res.Renamed) != 2 {
		t.Errorf("expected 2 renames, got %d", len(res.Renamed))
	}
}

func TestApply_SuffixRename_ChangesMatchingKeys(t *testing.T) {
	env := map[string]string{"DB_HOST": "db", "CACHE_HOST": "redis", "APP_PORT": "8080"}
	rules := []Rule{{From: "_HOST", To: "_ADDR", Mode: ModeSuffix}}
	res, err := Apply(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["DB_ADDR"]; !ok {
		t.Error("expected DB_ADDR")
	}
	if _, ok := res.Env["CACHE_ADDR"]; !ok {
		t.Error("expected CACHE_ADDR")
	}
	if res.Env["APP_PORT"] != "8080" {
		t.Error("APP_PORT should be unchanged")
	}
}

func TestApply_Conflict_LeavesOriginalIntact(t *testing.T) {
	env := map[string]string{"OLD": "old_val", "NEW": "existing"}
	rules := []Rule{{From: "OLD", To: "NEW", Mode: ModeExact}}
	res, err := Apply(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["NEW"] != "existing" {
		t.Errorf("conflict: NEW should retain original value, got %q", res.Env["NEW"])
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0] != "NEW" {
		t.Errorf("expected conflict on NEW, got %v", res.Conflicts)
	}
}

func TestApply_EmptyFromRule_ReturnsError(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	rules := []Rule{{From: "", To: "NEW", Mode: ModeExact}}
	_, err := Apply(env, rules)
	if err == nil {
		t.Error("expected error for empty From field")
	}
}

func TestApply_NoRules_ReturnsIdenticalEnv(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	res, err := Apply(env, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Renamed) != 0 {
		t.Errorf("expected no renames, got %v", res.Renamed)
	}
	if res.Env["A"] != "1" || res.Env["B"] != "2" {
		t.Error("env should be unchanged")
	}
}
