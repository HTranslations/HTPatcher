package util

import (
	"encoding/json"
	"reflect"
	"strings"
)

// GetJSONFieldNames returns the JSON field names for a struct using reflection.
func GetJSONFieldNames(v any) []string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	var fields []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
			name := strings.Split(tag, ",")[0]
			if name != "" {
				fields = append(fields, name)
			}
		}
	}
	return fields
}

// UnmarshalExtras extracts unknown fields from JSON data.
// knownFields are removed, remaining fields returned as extras.
func UnmarshalExtras(data []byte, knownFields []string) (map[string]json.RawMessage, error) {
	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return nil, err
	}
	for _, field := range knownFields {
		delete(rawMap, field)
	}
	if len(rawMap) == 0 {
		return nil, nil
	}
	return rawMap, nil
}

// MarshalWithExtras marshals a value and merges in extra fields.
func MarshalWithExtras(v any, extras map[string]json.RawMessage) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	if len(extras) == 0 {
		return data, nil
	}

	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return nil, err
	}
	for k, v := range extras {
		rawMap[k] = v
	}
	return json.Marshal(rawMap)
}
