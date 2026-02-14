package marshal

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"sort"
	"strconv"
)

// Writer serializes Go types to Ruby Marshal format.
type Writer struct {
	buf        *bytes.Buffer
	symbols    map[string]int      // symbol -> index for symlinks
	objects    map[interface{}]int // object -> index for object links
	objCounter int
}

// NewWriter creates a new Marshal writer.
func NewWriter() *Writer {
	return &Writer{
		buf:        &bytes.Buffer{},
		symbols:    make(map[string]int),
		objects:    make(map[interface{}]int),
		objCounter: 0,
	}
}

// Write serializes data to Ruby Marshal format.
func Write(data interface{}) ([]byte, error) {
	w := NewWriter()
	// Write Marshal version header (4.8)
	w.buf.WriteByte(4)
	w.buf.WriteByte(8)

	if err := w.writeValue(data); err != nil {
		return nil, err
	}
	return w.buf.Bytes(), nil
}

// writeValue writes a single value to the marshal stream.
func (w *Writer) writeValue(v interface{}) error {
	if v == nil {
		w.buf.WriteByte(TypeNil)
		return nil
	}

	switch val := v.(type) {
	case bool:
		if val {
			w.buf.WriteByte(TypeTrue)
		} else {
			w.buf.WriteByte(TypeFalse)
		}
		return nil

	case int:
		w.buf.WriteByte(TypeFixnum)
		w.writeInt(val)
		return nil

	case int64:
		w.buf.WriteByte(TypeFixnum)
		w.writeInt(int(val))
		return nil

	case float64:
		return w.writeFloat(val)

	case string:
		return w.writeString(val)

	case []interface{}:
		return w.writeArray(val)

	case map[string]interface{}:
		return w.writeHash(val)

	case *RubyObject:
		return w.writeObject(val)

	case *RubyStruct:
		return w.writeStruct(val)

	case *RubyTable:
		return w.writeTable(val)

	case *RubyColor:
		return w.writeColor(val)

	case *RubyTone:
		return w.writeTone(val)

	default:
		return fmt.Errorf("marshal: unsupported type %T", v)
	}
}

// writeInt writes a packed integer in Ruby's format.
func (w *Writer) writeInt(n int) {
	if n == 0 {
		w.buf.WriteByte(0)
		return
	}

	if n > 0 && n < 123 {
		w.buf.WriteByte(byte(n + 5))
		return
	}

	if n < 0 && n > -124 {
		w.buf.WriteByte(byte(n - 5))
		return
	}

	// Multi-byte encoding
	if n > 0 {
		if n <= 0xFF {
			w.buf.WriteByte(1)
			w.buf.WriteByte(byte(n))
		} else if n <= 0xFFFF {
			w.buf.WriteByte(2)
			w.buf.WriteByte(byte(n))
			w.buf.WriteByte(byte(n >> 8))
		} else if n <= 0xFFFFFF {
			w.buf.WriteByte(3)
			w.buf.WriteByte(byte(n))
			w.buf.WriteByte(byte(n >> 8))
			w.buf.WriteByte(byte(n >> 16))
		} else {
			w.buf.WriteByte(4)
			w.buf.WriteByte(byte(n))
			w.buf.WriteByte(byte(n >> 8))
			w.buf.WriteByte(byte(n >> 16))
			w.buf.WriteByte(byte(n >> 24))
		}
	} else {
		// Negative numbers
		un := uint32(n)
		if n >= -0x100 {
			w.buf.WriteByte(0xFF) // -1
			w.buf.WriteByte(byte(un))
		} else if n >= -0x10000 {
			w.buf.WriteByte(0xFE) // -2
			w.buf.WriteByte(byte(un))
			w.buf.WriteByte(byte(un >> 8))
		} else if n >= -0x1000000 {
			w.buf.WriteByte(0xFD) // -3
			w.buf.WriteByte(byte(un))
			w.buf.WriteByte(byte(un >> 8))
			w.buf.WriteByte(byte(un >> 16))
		} else {
			w.buf.WriteByte(0xFC) // -4
			w.buf.WriteByte(byte(un))
			w.buf.WriteByte(byte(un >> 8))
			w.buf.WriteByte(byte(un >> 16))
			w.buf.WriteByte(byte(un >> 24))
		}
	}
}

// writeFloat writes a float as a Ruby float (string representation).
func (w *Writer) writeFloat(f float64) error {
	w.buf.WriteByte(TypeFloat)
	var str string
	if math.IsInf(f, 1) {
		str = "inf"
	} else if math.IsInf(f, -1) {
		str = "-inf"
	} else if math.IsNaN(f) {
		str = "nan"
	} else {
		str = strconv.FormatFloat(f, 'g', -1, 64)
	}
	w.writeRawString(str)
	w.objCounter++
	return nil
}

// writeRawString writes a length-prefixed string (without type marker).
func (w *Writer) writeRawString(s string) {
	w.writeInt(len(s))
	w.buf.WriteString(s)
}

// writeString writes a string with IVAR encoding wrapper for UTF-8.
func (w *Writer) writeString(s string) error {
	// Write as IVAR-wrapped string with UTF-8 encoding
	w.buf.WriteByte(TypeIvar)
	w.buf.WriteByte(TypeString)
	w.writeRawString(s)
	w.objCounter++

	// Write encoding ivar: 1 ivar, :E (encoding flag) => true (UTF-8)
	w.writeInt(1)
	w.writeSymbol("E")
	w.buf.WriteByte(TypeTrue)

	return nil
}

// writeSymbol writes a symbol, using symlink if already seen.
func (w *Writer) writeSymbol(sym string) {
	if idx, ok := w.symbols[sym]; ok {
		w.buf.WriteByte(TypeSymlink)
		w.writeInt(idx)
		return
	}

	w.symbols[sym] = len(w.symbols)
	w.buf.WriteByte(TypeSymbol)
	w.writeRawString(sym)
}

// writeArray writes an array.
func (w *Writer) writeArray(arr []interface{}) error {
	w.buf.WriteByte(TypeArray)
	w.writeInt(len(arr))
	w.objCounter++

	for _, item := range arr {
		if err := w.writeValue(item); err != nil {
			return err
		}
	}
	return nil
}

// writeHash writes a hash/map.
func (w *Writer) writeHash(hash map[string]interface{}) error {
	w.buf.WriteByte(TypeHash)
	w.writeInt(len(hash))
	w.objCounter++

	// Sort keys for consistent output
	keys := make([]string, 0, len(hash))
	for k := range hash {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		// Try to write key as integer if possible
		if intKey, err := strconv.Atoi(k); err == nil {
			w.buf.WriteByte(TypeFixnum)
			w.writeInt(intKey)
		} else {
			if err := w.writeString(k); err != nil {
				return err
			}
		}
		if err := w.writeValue(hash[k]); err != nil {
			return err
		}
	}
	return nil
}

// writeObject writes a RubyObject.
func (w *Writer) writeObject(obj *RubyObject) error {
	// Check for userdata types
	if userData, ok := obj.Properties["__userdata__"]; ok {
		return w.writeUserDef(obj.Class, userData.([]byte))
	}

	w.buf.WriteByte(TypeObject)
	w.writeSymbol(obj.Class)
	w.objCounter++

	// Count properties (excluding internal markers)
	propCount := 0
	for k := range obj.Properties {
		if k[0] != '_' || len(k) < 2 || k[1] != '_' {
			propCount++
		}
	}
	w.writeInt(propCount)

	// Sort keys for consistent output
	keys := make([]string, 0, len(obj.Properties))
	for k := range obj.Properties {
		// Skip internal markers
		if k[0] == '_' && len(k) >= 2 && k[1] == '_' {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		// Write instance variable with @ prefix
		w.writeSymbol("@" + k)
		if err := w.writeValue(obj.Properties[k]); err != nil {
			return err
		}
	}
	return nil
}

// writeStruct writes a RubyStruct.
func (w *Writer) writeStruct(s *RubyStruct) error {
	w.buf.WriteByte(TypeStruct)
	w.writeSymbol(s.Name)
	w.objCounter++

	w.writeInt(len(s.Members))

	// Sort keys for consistent output
	keys := make([]string, 0, len(s.Members))
	for k := range s.Members {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		w.writeSymbol(k)
		if err := w.writeValue(s.Members[k]); err != nil {
			return err
		}
	}
	return nil
}

// writeUserDef writes a userdef type (Table, Color, Tone, etc.).
func (w *Writer) writeUserDef(className string, data []byte) error {
	w.buf.WriteByte(TypeUserDef)
	w.writeSymbol(className)
	w.writeRawString(string(data))
	w.objCounter++
	return nil
}

// writeTable serializes a RubyTable back to userdef format.
func (w *Writer) writeTable(t *RubyTable) error {
	// Calculate data size
	count := len(t.Data)
	data := make([]byte, 20+count*2)

	binary.LittleEndian.PutUint32(data[0:4], uint32(t.Dim))
	binary.LittleEndian.PutUint32(data[4:8], uint32(t.XSize))
	binary.LittleEndian.PutUint32(data[8:12], uint32(t.YSize))
	binary.LittleEndian.PutUint32(data[12:16], uint32(t.ZSize))
	binary.LittleEndian.PutUint32(data[16:20], uint32(count))

	for i, v := range t.Data {
		binary.LittleEndian.PutUint16(data[20+i*2:], uint16(v))
	}

	return w.writeUserDef("Table", data)
}

// writeColor serializes a RubyColor back to userdef format.
func (w *Writer) writeColor(c *RubyColor) error {
	data := make([]byte, 32)
	binary.LittleEndian.PutUint64(data[0:8], math.Float64bits(c.Red))
	binary.LittleEndian.PutUint64(data[8:16], math.Float64bits(c.Green))
	binary.LittleEndian.PutUint64(data[16:24], math.Float64bits(c.Blue))
	binary.LittleEndian.PutUint64(data[24:32], math.Float64bits(c.Alpha))
	return w.writeUserDef("Color", data)
}

// writeTone serializes a RubyTone back to userdef format.
func (w *Writer) writeTone(t *RubyTone) error {
	data := make([]byte, 32)
	binary.LittleEndian.PutUint64(data[0:8], math.Float64bits(t.Red))
	binary.LittleEndian.PutUint64(data[8:16], math.Float64bits(t.Green))
	binary.LittleEndian.PutUint64(data[16:24], math.Float64bits(t.Blue))
	binary.LittleEndian.PutUint64(data[24:32], math.Float64bits(t.Gray))
	return w.writeUserDef("Tone", data)
}
