package util

import "encoding/json"

// MergeJsonChanges merges typed changes into raw JSON while preserving unknown fields.
// Supports both top-level maps and arrays.
func MergeJsonChanges(data []byte, changes any) ([]byte, error) {
	var original any
	if err := json.Unmarshal(data, &original); err != nil {
		return nil, err
	}

	changedBytes, err := json.Marshal(changes)
	if err != nil {
		return nil, err
	}

	var changed any
	if err := json.Unmarshal(changedBytes, &changed); err != nil {
		return nil, err
	}

	merged := mergeAny(original, changed, "")
	return json.Marshal(merged)
}

func mergeAny(original, changes any, key string) any {
	// Special rule: "list" is always replaced, never merged
	if key == "list" {
		return changes
	}

	switch o := original.(type) {

	case map[string]any:
		if c, ok := changes.(map[string]any); ok {
			for k, v := range c {
				if ov, exists := o[k]; exists {
					o[k] = mergeAny(ov, v, k)
				} else {
					o[k] = v
				}
			}
			return o
		}
		return changes

	case []any:
		if c, ok := changes.([]any); ok {
			minLen := len(o)
			if len(c) < minLen {
				minLen = len(c)
			}
			for i := 0; i < minLen; i++ {
				o[i] = mergeAny(o[i], c[i], "")
			}
			if len(c) > len(o) {
				o = append(o, c[len(o):]...)
			}
			return o
		}
		return changes

	default:
		// primitives
		return changes
	}
}




