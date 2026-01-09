package util

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// OrderedMap is a map that preserves insertion order of keys.
// This is necessary for JSON serialization where key order matters (e.g., RPG Maker plugin parameters).
type OrderedMap struct {
	Keys   []string
	Values map[string]any
}

// NewOrderedMap creates a new empty OrderedMap.
func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		Keys:   []string{},
		Values: make(map[string]any),
	}
}

// Set adds or updates a key-value pair. If the key is new, it's appended to Keys.
func (o *OrderedMap) Set(key string, value any) {
	if _, exists := o.Values[key]; !exists {
		o.Keys = append(o.Keys, key)
	}
	o.Values[key] = value
}

// Get retrieves a value by key.
func (o *OrderedMap) Get(key string) (any, bool) {
	val, ok := o.Values[key]
	return val, ok
}

// MarshalJSON serializes the OrderedMap to JSON with keys in insertion order.
func (o *OrderedMap) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("{")

	for i, key := range o.Keys {
		if i > 0 {
			buf.WriteString(",")
		}

		// Marshal key
		keyBytes, err := json.Marshal(key)
		if err != nil {
			return nil, err
		}
		buf.Write(keyBytes)
		buf.WriteString(":")

		// Marshal value (recursively handles nested *OrderedMap via json.Marshal)
		valBytes, err := json.Marshal(o.Values[key])
		if err != nil {
			return nil, err
		}
		buf.Write(valBytes)
	}

	buf.WriteString("}")
	return buf.Bytes(), nil
}

// UnmarshalJSON deserializes JSON into OrderedMap while preserving key order.
func (o *OrderedMap) UnmarshalJSON(data []byte) error {
	o.Keys = []string{}
	o.Values = make(map[string]any)

	dec := json.NewDecoder(bytes.NewReader(data))

	// Read opening brace
	token, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := token.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expected '{', got %v", token)
	}

	// Read key-value pairs
	for dec.More() {
		// Read key
		token, err := dec.Token()
		if err != nil {
			return err
		}
		key, ok := token.(string)
		if !ok {
			return fmt.Errorf("expected string key, got %v", token)
		}

		// Read value
		value, err := parseJSONValue(dec)
		if err != nil {
			return err
		}

		o.Keys = append(o.Keys, key)
		o.Values[key] = value
	}

	// Read closing brace
	token, err = dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := token.(json.Delim); !ok || delim != '}' {
		return fmt.Errorf("expected '}', got %v", token)
	}

	return nil
}

// parseJSONValue parses the next JSON value from the decoder, preserving object key order.
func parseJSONValue(dec *json.Decoder) (any, error) {
	token, err := dec.Token()
	if err != nil {
		return nil, err
	}

	switch t := token.(type) {
	case json.Delim:
		switch t {
		case '{':
			// Parse object as OrderedMap
			om := NewOrderedMap()
			for dec.More() {
				// Read key
				keyToken, err := dec.Token()
				if err != nil {
					return nil, err
				}
				key, ok := keyToken.(string)
				if !ok {
					return nil, fmt.Errorf("expected string key, got %v", keyToken)
				}

				// Read value recursively
				value, err := parseJSONValue(dec)
				if err != nil {
					return nil, err
				}

				om.Keys = append(om.Keys, key)
				om.Values[key] = value
			}
			// Consume closing brace
			if _, err := dec.Token(); err != nil {
				return nil, err
			}
			return om, nil

		case '[':
			// Parse array
			var arr []any
			for dec.More() {
				value, err := parseJSONValue(dec)
				if err != nil {
					return nil, err
				}
				arr = append(arr, value)
			}
			// Consume closing bracket
			if _, err := dec.Token(); err != nil {
				return nil, err
			}
			return arr, nil
		}

	case string, float64, bool, nil:
		return t, nil
	}

	return nil, fmt.Errorf("unexpected token: %v", token)
}
