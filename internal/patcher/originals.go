package patcher

import (
	"encoding/json"
)

const htOriginalField = "ht_original"

func setHTOriginal(extras *map[string]json.RawMessage, path string, original string) {
	if extras == nil {
		return
	}
	if *extras == nil {
		*extras = make(map[string]json.RawMessage)
	}

	originals := make(map[string]string)
	if raw, ok := (*extras)[htOriginalField]; ok {
		_ = json.Unmarshal(raw, &originals)
	}
	originals[path] = original

	raw, err := json.Marshal(originals)
	if err != nil {
		return
	}
	(*extras)[htOriginalField] = raw
}

func originalValueString(value any) string {
	if text, ok := value.(string); ok {
		return text
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(raw)
}
