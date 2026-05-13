package output

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/pin"
)

func TestWritePinResult_ShowsPinnedKeys(t *testing.T) {
	res := pin.PinResult{
		Pinned: map[string]string{"DB_HOST": "localhost", "API_KEY": "secret"},
	}
	var sb strings.Builder
	WritePinResult(&sb, res)
	out := sb.String()
	if !strings.Contains(out, "Pinned keys (2)") {
		t.Errorf("expected pinned keys header, got: %s", out)
	}
	if !strings.Contains(out, "API_KEY=secret") {
		t.Errorf("expected API_KEY in output, got: %s", out)
	}
}

func TestWritePinResult_NoPinnedKeys(t *testing.T) {
	res := pin.PinResult{Pinned: map[string]string{}}
	var sb strings.Builder
	WritePinResult(&sb, res)
	if !strings.Contains(sb.String(), "No keys currently pinned") {
		t.Errorf("expected empty message, got: %s", sb.String())
	}
}

func TestWritePinResult_ShowsSkipped(t *testing.T) {
	res := pin.PinResult{
		Pinned:  map[string]string{},
		Skipped: []string{"MISSING_KEY"},
	}
	var sb strings.Builder
	WritePinResult(&sb, res)
	if !strings.Contains(sb.String(), "MISSING_KEY") {
		t.Errorf("expected skipped key in output, got: %s", sb.String())
	}
}

func TestWritePinResult_ShowsReleased(t *testing.T) {
	res := pin.PinResult{
		Pinned:   map[string]string{},
		Released: []string{"OLD_KEY"},
	}
	var sb strings.Builder
	WritePinResult(&sb, res)
	if !strings.Contains(sb.String(), "Released") || !strings.Contains(sb.String(), "OLD_KEY") {
		t.Errorf("expected released section, got: %s", sb.String())
	}
}

func TestWritePinDrift_NoDrift(t *testing.T) {
	var sb strings.Builder
	WritePinDrift(&sb, nil)
	if !strings.Contains(sb.String(), "OK") {
		t.Errorf("expected OK message, got: %s", sb.String())
	}
}

func TestWritePinDrift_WithDrift(t *testing.T) {
	var sb strings.Builder
	WritePinDrift(&sb, []string{"DB_PASS", "API_SECRET"})
	out := sb.String()
	if !strings.Contains(out, "DRIFT") {
		t.Errorf("expected DRIFT in output, got: %s", out)
	}
	if !strings.Contains(out, "DB_PASS") || !strings.Contains(out, "API_SECRET") {
		t.Errorf("expected drifted keys in output, got: %s", out)
	}
}
