package rpgmaker

import (
	"bytes"
	"encoding/json"
	"htpatcher/internal/util"
)

// Trait represents a character/item trait in RPG Maker
type Trait struct {
	Code   int                        `json:"code"`
	DataId int                        `json:"dataId"`
	Value  float64                    `json:"value"`
	Extras map[string]json.RawMessage `json:"-"`
}

func (t *Trait) UnmarshalJSON(data []byte) error {
	type Alias Trait
	aux := (*Alias)(t)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	t.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(t))
	return nil
}

func (t Trait) MarshalJSON() ([]byte, error) {
	type Alias Trait
	return util.MarshalWithExtras((*Alias)(&t), t.Extras)
}

// EventCommand represents a command in an event or troop page
type EventCommand struct {
	Code       int   `json:"code"`
	Indent     int   `json:"indent"`
	Parameters []any `json:"parameters"`
}

// UnmarshalJSON custom unmarshals EventCommand to preserve object key order in Parameters.
func (e *EventCommand) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))

	// Read opening brace
	if _, err := dec.Token(); err != nil {
		return err
	}

	// Read key-value pairs
	for dec.More() {
		keyToken, err := dec.Token()
		if err != nil {
			return err
		}
		key := keyToken.(string)

		switch key {
		case "code":
			var code int
			if err := dec.Decode(&code); err != nil {
				return err
			}
			e.Code = code

		case "indent":
			var indent int
			if err := dec.Decode(&indent); err != nil {
				return err
			}
			e.Indent = indent

		case "parameters":
			params, err := parseParametersArray(dec)
			if err != nil {
				return err
			}
			e.Parameters = params

		default:
			// Skip unknown fields
			var skip any
			if err := dec.Decode(&skip); err != nil {
				return err
			}
		}
	}

	// Read closing brace
	if _, err := dec.Token(); err != nil {
		return err
	}

	return nil
}

// parseParametersArray parses the parameters array, converting objects to *util.OrderedMap.
func parseParametersArray(dec *json.Decoder) ([]any, error) {
	// Read opening bracket
	token, err := dec.Token()
	if err != nil {
		return nil, err
	}
	if delim, ok := token.(json.Delim); !ok || delim != '[' {
		return []any{}, nil
	}

	params := []any{}
	for dec.More() {
		value, err := parseParameterValue(dec)
		if err != nil {
			return nil, err
		}
		params = append(params, value)
	}

	// Read closing bracket
	if _, err := dec.Token(); err != nil {
		return nil, err
	}

	return params, nil
}

// parseParameterValue parses a single parameter value, converting objects to *util.OrderedMap.
func parseParameterValue(dec *json.Decoder) (any, error) {
	token, err := dec.Token()
	if err != nil {
		return nil, err
	}

	switch t := token.(type) {
	case json.Delim:
		switch t {
		case '{':
			// Parse object as OrderedMap
			om := util.NewOrderedMap()
			for dec.More() {
				keyToken, err := dec.Token()
				if err != nil {
					return nil, err
				}
				key := keyToken.(string)

				value, err := parseParameterValue(dec)
				if err != nil {
					return nil, err
				}

				om.Set(key, value)
			}
			// Consume closing brace
			if _, err := dec.Token(); err != nil {
				return nil, err
			}
			return om, nil

		case '[':
			// Parse array recursively
			arr := []any{}
			for dec.More() {
				value, err := parseParameterValue(dec)
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

	return nil, nil
}
