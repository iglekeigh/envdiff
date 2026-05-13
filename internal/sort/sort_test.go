package sort

import (
	"testing"
)

func TestApply_Alpha_DefaultOrder(t *testing.T) {
	env := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	keys, out := Apply(env, DefaultOptions)

	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "APPLE" || keys[1] != "MANGO" || keys[2] != "ZEBRA" {
		t.Errorf("unexpected order: %v", keys)
	}
	if out["APPLE"] != "2" {
		t.Errorf("map value mismatch for APPLE")
	}
}

func TestApply_Alpha_Descending(t *testing.T) {
	env := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	keys, _ := Apply(env, Options{Strategy: StrategyAlpha, Descending: true})

	if keys[0] != "ZEBRA" || keys[2] != "APPLE" {
		t.Errorf("expected descending alpha order, got %v", keys)
	}
}

func TestApply_Length_SortsByKeyLength(t *testing.T) {
	env := map[string]string{"AB": "1", "A": "2", "ABC": "3"}
	keys, _ := Apply(env, Options{Strategy: StrategyLength})

	if keys[0] != "A" || keys[1] != "AB" || keys[2] != "ABC" {
		t.Errorf("expected length order, got %v", keys)
	}
}

func TestApply_Length_TieBreakAlpha(t *testing.T) {
	env := map[string]string{"BB": "1", "AA": "2"}
	keys, _ := Apply(env, Options{Strategy: StrategyLength})

	if keys[0] != "AA" || keys[1] != "BB" {
		t.Errorf("expected alpha tie-break within same length, got %v", keys)
	}
}

func TestApply_Group_ClustersByPrefix(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_PORT": "5432",
		"APP_NAME": "envdiff",
	}
	keys, _ := Apply(env, Options{Strategy: StrategyGroup})

	// APP group should come before DB group.
	appDone := false
	for _, k := range keys {
		if k == "APP_NAME" || k == "APP_PORT" {
			if appDone {
				t.Errorf("APP keys should be contiguous before DB keys; got %v", keys)
			}
		}
		if k == "DB_HOST" || k == "DB_PORT" {
			appDone = true
		}
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	_, out := Apply(env, DefaultOptions)
	out["KEY"] = "mutated"

	if env["KEY"] == "mutated" {
		t.Error("Apply mutated the original map")
	}
}

func TestApply_EmptyMap_ReturnsEmpty(t *testing.T) {
	keys, out := Apply(map[string]string{}, DefaultOptions)
	if len(keys) != 0 || len(out) != 0 {
		t.Error("expected empty results for empty input")
	}
}
