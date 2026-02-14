// Package marshal provides Ruby Marshal format parsing for RPG Maker VX Ace .rvdata2 files.
package marshal

// Marshal format type constants
const (
	TypeNil         = '0'
	TypeTrue        = 'T'
	TypeFalse       = 'F'
	TypeFixnum      = 'i'
	TypeBignum      = 'l'
	TypeFloat       = 'f'
	TypeString      = '"'
	TypeRegexp      = '/'
	TypeArray       = '['
	TypeHash        = '{'
	TypeHashDef     = '}'
	TypeObject      = 'o'
	TypeUserDef     = 'u'
	TypeUserMarshal = 'U'
	TypeSymbol      = ':'
	TypeSymlink     = ';'
	TypeIvar        = 'I'
	TypeLink        = '@'
	TypeData        = 'd'
	TypeClass       = 'c'
	TypeModule      = 'm'
	TypeStruct      = 'S'
	TypeExtended    = 'e'
)

// RubyObject represents a parsed Ruby object with its class name and instance variables.
type RubyObject struct {
	Class      string
	Properties map[string]interface{}
}

// RubyStruct represents a parsed Ruby Struct.
type RubyStruct struct {
	Name    string
	Members map[string]interface{}
}

// RubyTable represents RPG Maker's Table class (tile data, etc).
type RubyTable struct {
	Dim   int
	XSize int
	YSize int
	ZSize int
	Data  []int16
}

// RubyColor represents RPG Maker's Color class.
type RubyColor struct {
	Red   float64
	Green float64
	Blue  float64
	Alpha float64
}

// RubyTone represents RPG Maker's Tone class.
type RubyTone struct {
	Red   float64
	Green float64
	Blue  float64
	Gray  float64
}

// GetString safely extracts a string from a property value.
func (o *RubyObject) GetString(key string) string {
	if v, ok := o.Properties[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetInt safely extracts an int from a property value.
func (o *RubyObject) GetInt(key string) int {
	if v, ok := o.Properties[key]; ok {
		switch n := v.(type) {
		case int:
			return n
		case int64:
			return int(n)
		case float64:
			return int(n)
		}
	}
	return 0
}

// GetBool safely extracts a bool from a property value.
func (o *RubyObject) GetBool(key string) bool {
	if v, ok := o.Properties[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// GetArray safely extracts an array from a property value.
func (o *RubyObject) GetArray(key string) []interface{} {
	if v, ok := o.Properties[key]; ok {
		if arr, ok := v.([]interface{}); ok {
			return arr
		}
	}
	return nil
}

// GetObject safely extracts a nested RubyObject from a property value.
func (o *RubyObject) GetObject(key string) *RubyObject {
	if v, ok := o.Properties[key]; ok {
		if obj, ok := v.(*RubyObject); ok {
			return obj
		}
	}
	return nil
}

// GetMap safely extracts a map from a property value.
func (o *RubyObject) GetMap(key string) map[string]interface{} {
	if v, ok := o.Properties[key]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			return m
		}
	}
	return nil
}

// SetString sets a string property value.
func (o *RubyObject) SetString(key string, value string) {
	if o.Properties == nil {
		o.Properties = make(map[string]interface{})
	}
	o.Properties[key] = value
}
