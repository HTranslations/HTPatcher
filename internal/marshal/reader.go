package marshal

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Reader parses Ruby Marshal format data.
type Reader struct {
	data    []byte
	pos     int
	symbols []string
	objects []interface{}
}

// NewReader creates a new Marshal reader from raw bytes.
func NewReader(data []byte) *Reader {
	return &Reader{
		data:    data,
		pos:     0,
		symbols: make([]string, 0),
		objects: make([]interface{}, 0),
	}
}

// Parse reads Marshal data and returns the parsed value.
// This is the main entry point for parsing .rvdata2 files.
func Parse(data []byte) (interface{}, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("marshal: data too short")
	}
	if data[0] != 4 || data[1] != 8 {
		return nil, fmt.Errorf("marshal: unsupported version %d.%d (expected 4.8)", data[0], data[1])
	}

	reader := NewReader(data[2:])
	return reader.Read()
}

func (r *Reader) readByte() (byte, error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	b := r.data[r.pos]
	r.pos++
	return b, nil
}

func (r *Reader) readBytes(n int) ([]byte, error) {
	if r.pos+n > len(r.data) {
		return nil, io.EOF
	}
	bytes := r.data[r.pos : r.pos+n]
	r.pos += n
	return bytes, nil
}

// readInt reads a packed integer (Ruby's special format).
func (r *Reader) readInt() (int, error) {
	b, err := r.readByte()
	if err != nil {
		return 0, err
	}
	c := int(int8(b))

	if c == 0 {
		return 0, nil
	}
	if c > 0 {
		if c <= 4 {
			result := 0
			for i := 0; i < c; i++ {
				b, err := r.readByte()
				if err != nil {
					return 0, err
				}
				result |= int(b) << (8 * i)
			}
			return result, nil
		}
		return c - 5, nil
	}
	if c >= -4 {
		result := -1
		for i := 0; i < -c; i++ {
			b, err := r.readByte()
			if err != nil {
				return 0, err
			}
			result &= ^(0xff << (8 * i))
			result |= int(b) << (8 * i)
		}
		return result, nil
	}
	return c + 5, nil
}

func (r *Reader) readRawString() (string, error) {
	length, err := r.readInt()
	if err != nil {
		return "", err
	}
	bytes, err := r.readBytes(length)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (r *Reader) readSymbol() (string, error) {
	sym, err := r.readRawString()
	if err != nil {
		return "", err
	}
	r.symbols = append(r.symbols, sym)
	return sym, nil
}

func (r *Reader) readSymbolOrLink() (string, error) {
	typeCode, err := r.readByte()
	if err != nil {
		return "", err
	}

	switch typeCode {
	case TypeSymbol:
		return r.readSymbol()
	case TypeSymlink:
		idx, err := r.readInt()
		if err != nil {
			return "", err
		}
		if idx < 0 || idx >= len(r.symbols) {
			return "", fmt.Errorf("marshal: invalid symbol reference %d", idx)
		}
		return r.symbols[idx], nil
	default:
		return "", fmt.Errorf("marshal: expected symbol or symlink, got 0x%02x", typeCode)
	}
}

func (r *Reader) registerObject(obj interface{}) int {
	idx := len(r.objects)
	r.objects = append(r.objects, obj)
	return idx
}

// Read reads the next value from the marshal stream.
func (r *Reader) Read() (interface{}, error) {
	typeCode, err := r.readByte()
	if err != nil {
		return nil, err
	}

	switch typeCode {
	case TypeNil:
		return nil, nil

	case TypeTrue:
		return true, nil

	case TypeFalse:
		return false, nil

	case TypeFixnum:
		return r.readInt()

	case TypeFloat:
		str, err := r.readRawString()
		if err != nil {
			return nil, err
		}
		// Ruby 1.9.2 may append \0 + IEEE 754 binary bytes for float precision.
		// Strip everything from the null byte onward so ParseFloat succeeds.
		if idx := strings.IndexByte(str, 0); idx >= 0 {
			str = str[:idx]
		}
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			r.registerObject(str)
			return str, nil
		}
		r.registerObject(f)
		return f, nil

	case TypeString:
		str, err := r.readRawString()
		if err != nil {
			return nil, err
		}
		// Note: String registration happens in readIvar for strings with encoding
		// Only register here for bare strings (rare)
		r.registerObject(str)
		return str, nil

	case TypeSymbol:
		return r.readSymbol()

	case TypeSymlink:
		idx, err := r.readInt()
		if err != nil {
			return nil, err
		}
		if idx < 0 || idx >= len(r.symbols) {
			return nil, fmt.Errorf("marshal: invalid symbol reference %d", idx)
		}
		return r.symbols[idx], nil

	case TypeArray:
		return r.readArray()

	case TypeHash, TypeHashDef:
		return r.readHash(typeCode == TypeHashDef)

	case TypeObject:
		return r.readObject()

	case TypeIvar:
		return r.readIvar()

	case TypeLink:
		idx, err := r.readInt()
		if err != nil {
			return nil, err
		}
		if idx < 0 || idx >= len(r.objects) {
			return nil, fmt.Errorf("marshal: invalid object reference %d (have %d)", idx, len(r.objects))
		}
		return r.objects[idx], nil

	case TypeUserDef:
		return r.readUserDef()

	case TypeUserMarshal:
		return r.readUserMarshal()

	case TypeClass:
		name, err := r.readRawString()
		if err != nil {
			return nil, err
		}
		obj := &RubyObject{Class: "__class__:" + name, Properties: nil}
		r.registerObject(obj)
		return obj, nil

	case TypeModule:
		name, err := r.readRawString()
		if err != nil {
			return nil, err
		}
		obj := &RubyObject{Class: "__module__:" + name, Properties: nil}
		r.registerObject(obj)
		return obj, nil

	case TypeExtended:
		_, err := r.readSymbolOrLink()
		if err != nil {
			return nil, err
		}
		return r.Read()

	case TypeStruct:
		return r.readStruct()

	case TypeData:
		return r.readData()

	case TypeBignum:
		return r.readBignum()

	case TypeRegexp:
		return r.readRegexp()

	default:
		return nil, fmt.Errorf("marshal: unknown type code 0x%02x at position %d", typeCode, r.pos-1)
	}
}

func (r *Reader) readArray() ([]interface{}, error) {
	length, err := r.readInt()
	if err != nil {
		return nil, err
	}
	arr := make([]interface{}, length)
	objIdx := r.registerObject(arr)
	for i := 0; i < length; i++ {
		arr[i], err = r.Read()
		if err != nil {
			return nil, err
		}
	}
	r.objects[objIdx] = arr
	return arr, nil
}

func (r *Reader) readHash(hasDefault bool) (map[string]interface{}, error) {
	length, err := r.readInt()
	if err != nil {
		return nil, err
	}
	hash := make(map[string]interface{})
	objIdx := r.registerObject(hash)
	for i := 0; i < length; i++ {
		key, err := r.Read()
		if err != nil {
			return nil, err
		}
		value, err := r.Read()
		if err != nil {
			return nil, err
		}
		// Convert key to string
		var keyStr string
		switch k := key.(type) {
		case string:
			keyStr = k
		case int:
			keyStr = strconv.Itoa(k)
		default:
			keyStr = fmt.Sprintf("%v", key)
		}
		hash[keyStr] = value
	}
	if hasDefault {
		_, err = r.Read() // Read and discard default value
		if err != nil {
			return nil, err
		}
	}
	r.objects[objIdx] = hash
	return hash, nil
}

func (r *Reader) readObject() (*RubyObject, error) {
	className, err := r.readSymbolOrLink()
	if err != nil {
		return nil, err
	}

	numVars, err := r.readInt()
	if err != nil {
		return nil, err
	}

	obj := &RubyObject{
		Class:      className,
		Properties: make(map[string]interface{}),
	}
	objIdx := r.registerObject(obj)

	for i := 0; i < numVars; i++ {
		key, err := r.readSymbolOrLink()
		if err != nil {
			return nil, err
		}
		value, err := r.Read()
		if err != nil {
			return nil, err
		}
		// Remove @ prefix from instance variables
		if strings.HasPrefix(key, "@") {
			key = key[1:]
		}
		obj.Properties[key] = value
	}
	r.objects[objIdx] = obj
	return obj, nil
}

func (r *Reader) readIvar() (interface{}, error) {
	value, err := r.Read()
	if err != nil {
		return nil, err
	}

	numVars, err := r.readInt()
	if err != nil {
		return nil, err
	}

	// Read and discard instance variables (usually encoding info)
	for i := 0; i < numVars; i++ {
		_, err := r.readSymbolOrLink()
		if err != nil {
			return nil, err
		}
		_, err = r.Read()
		if err != nil {
			return nil, err
		}
	}

	return value, nil
}

func (r *Reader) readUserDef() (interface{}, error) {
	className, err := r.readSymbolOrLink()
	if err != nil {
		return nil, err
	}
	data, err := r.readRawString()
	if err != nil {
		return nil, err
	}
	rawBytes := []byte(data)

	// Handle known RPG Maker types
	switch className {
	case "Table":
		table := ParseTable(rawBytes)
		if table != nil {
			r.registerObject(table)
			return table, nil
		}
	case "Color":
		color := ParseColor(rawBytes)
		if color != nil {
			r.registerObject(color)
			return color, nil
		}
	case "Tone":
		tone := ParseTone(rawBytes)
		if tone != nil {
			r.registerObject(tone)
			return tone, nil
		}
	}

	// Unknown userdef, return as RubyObject with raw data note
	obj := &RubyObject{
		Class: className,
		Properties: map[string]interface{}{
			"__userdata__":      rawBytes,
			"__userdata_size__": len(rawBytes),
		},
	}
	r.registerObject(obj)
	return obj, nil
}

func (r *Reader) readUserMarshal() (interface{}, error) {
	className, err := r.readSymbolOrLink()
	if err != nil {
		return nil, err
	}
	value, err := r.Read()
	if err != nil {
		return nil, err
	}
	obj := &RubyObject{
		Class: className,
		Properties: map[string]interface{}{
			"__value__": value,
		},
	}
	r.registerObject(obj)
	return obj, nil
}

func (r *Reader) readStruct() (*RubyStruct, error) {
	className, err := r.readSymbolOrLink()
	if err != nil {
		return nil, err
	}
	numMembers, err := r.readInt()
	if err != nil {
		return nil, err
	}
	s := &RubyStruct{
		Name:    className,
		Members: make(map[string]interface{}),
	}
	r.registerObject(s)
	for i := 0; i < numMembers; i++ {
		key, err := r.readSymbolOrLink()
		if err != nil {
			return nil, err
		}
		value, err := r.Read()
		if err != nil {
			return nil, err
		}
		s.Members[key] = value
	}
	return s, nil
}

func (r *Reader) readData() (*RubyObject, error) {
	className, err := r.readSymbolOrLink()
	if err != nil {
		return nil, err
	}
	value, err := r.Read()
	if err != nil {
		return nil, err
	}
	obj := &RubyObject{
		Class: className,
		Properties: map[string]interface{}{
			"__datavalue__": value,
		},
	}
	r.registerObject(obj)
	return obj, nil
}

func (r *Reader) readBignum() (int64, error) {
	sign, err := r.readByte()
	if err != nil {
		return 0, err
	}
	numShorts, err := r.readInt()
	if err != nil {
		return 0, err
	}
	bytes, err := r.readBytes(numShorts * 2)
	if err != nil {
		return 0, err
	}
	result := int64(0)
	for i := 0; i < len(bytes); i++ {
		result |= int64(bytes[i]) << (8 * i)
	}
	if sign == '-' {
		result = -result
	}
	r.registerObject(result)
	return result, nil
}

func (r *Reader) readRegexp() (*RubyObject, error) {
	pattern, err := r.readRawString()
	if err != nil {
		return nil, err
	}
	flags, err := r.readByte()
	if err != nil {
		return nil, err
	}
	obj := &RubyObject{
		Class: "Regexp",
		Properties: map[string]interface{}{
			"pattern": pattern,
			"flags":   int(flags),
		},
	}
	r.registerObject(obj)
	return obj, nil
}
