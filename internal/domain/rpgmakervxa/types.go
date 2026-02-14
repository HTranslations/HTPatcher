// Package rpgmakervxa provides types and utilities for RPG Maker VX Ace game data.
package rpgmakervxa

import (
	"htpatcher/internal/marshal"
)

// VXAEventCommand represents an event command in VX Ace format.
// The command codes are the same as MV/MZ.
type VXAEventCommand struct {
	Code       int
	Indent     int
	Parameters []interface{}
}

// ParseEventCommands parses an array of event commands from Ruby objects.
func ParseEventCommands(arr []interface{}) []*VXAEventCommand {
	if arr == nil {
		return nil
	}

	commands := make([]*VXAEventCommand, 0, len(arr))
	for _, item := range arr {
		obj, ok := item.(*marshal.RubyObject)
		if !ok {
			continue
		}

		cmd := &VXAEventCommand{
			Code:   obj.GetInt("code"),
			Indent: obj.GetInt("indent"),
		}

		// Parse parameters
		if params := obj.GetArray("parameters"); params != nil {
			cmd.Parameters = make([]interface{}, len(params))
			for i, p := range params {
				cmd.Parameters[i] = extractParameterValue(p)
			}
		}

		commands = append(commands, cmd)
	}
	return commands
}

// extractParameterValue extracts the actual value from a parameter.
func extractParameterValue(v interface{}) interface{} {
	switch val := v.(type) {
	case *marshal.RubyObject:
		// For nested objects, return the object itself so we can serialize it back
		return val
	case *marshal.RubyStruct:
		return val
	default:
		return val
	}
}

// BuildEventCommandArray builds a Ruby object array from event commands.
func BuildEventCommandArray(commands []*VXAEventCommand) []interface{} {
	result := make([]interface{}, len(commands))
	for i, cmd := range commands {
		if cmd == nil {
			result[i] = nil
			continue
		}
		obj := &marshal.RubyObject{
			Class: "RPG::EventCommand",
			Properties: map[string]interface{}{
				"code":       cmd.Code,
				"indent":     cmd.Indent,
				"parameters": cmd.Parameters,
			},
		}
		result[i] = obj
	}
	return result
}

// GetSystemTitle extracts the game title from a parsed System.rvdata2.
func GetSystemTitle(raw interface{}) string {
	obj, ok := raw.(*marshal.RubyObject)
	if !ok {
		return ""
	}
	return obj.GetString("game_title")
}

// Helper functions for accessing RubyObject properties

// GetString safely extracts a string from a RubyObject property.
func GetString(obj *marshal.RubyObject, key string) string {
	if obj == nil {
		return ""
	}
	return obj.GetString(key)
}

// SetString safely sets a string property on a RubyObject.
func SetString(obj *marshal.RubyObject, key, value string) {
	if obj == nil || obj.Properties == nil {
		return
	}
	obj.Properties[key] = value
}

// GetInt safely extracts an int from a RubyObject property.
func GetInt(obj *marshal.RubyObject, key string) int {
	if obj == nil {
		return 0
	}
	return obj.GetInt(key)
}

// GetArray safely extracts an array from a RubyObject property.
func GetArray(obj *marshal.RubyObject, key string) []interface{} {
	if obj == nil {
		return nil
	}
	return obj.GetArray(key)
}

// GetObject safely extracts a nested RubyObject from a property.
func GetObject(obj *marshal.RubyObject, key string) *marshal.RubyObject {
	if obj == nil {
		return nil
	}
	return obj.GetObject(key)
}

// GetMap safely extracts a map from a RubyObject property.
func GetMap(obj *marshal.RubyObject, key string) map[string]interface{} {
	if obj == nil {
		return nil
	}
	return obj.GetMap(key)
}

// SetArray sets an array property on a RubyObject.
func SetArray(obj *marshal.RubyObject, key string, value []interface{}) {
	if obj == nil || obj.Properties == nil {
		return
	}
	obj.Properties[key] = value
}

// PatchStringArray patches translatable strings in an array.
func PatchStringArray(arr []interface{}, dictionary map[string]string, getKey func(string) string) {
	for i, v := range arr {
		if s, ok := v.(string); ok {
			if translation, exists := dictionary[getKey(s)]; exists {
				arr[i] = translation
			}
		}
	}
}
