package importer

import (
	"encoding/json"
	"fmt"
)

// parseJSON reads a flat JSON object and returns string key-value pairs.
// Non-string values are skipped and recorded in the skipped slice.
func parseJSON(data []byte) (map[string]string, []string, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, nil, fmt.Errorf("json: unmarshal failed: %w", err)
	}

	env := make(map[string]string, len(raw))
	var skipped []string

	for k, v := range raw {
		switch val := v.(type) {
		case string:
			env[k] = val
		case float64:
			env[k] = fmt.Sprintf("%g", val)
		case bool:
			if val {
				env[k] = "true"
			} else {
				env[k] = "false"
			}
		case nil:
			env[k] = ""
		default:
			skipped = append(skipped, k)
		}
	}

	return env, skipped, nil
}
