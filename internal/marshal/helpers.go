package marshal

import (
	"encoding/binary"
	"math"
)

// ParseTable decodes RPG Maker's Table userdata.
// Table format: dim(4) + xsize(4) + ysize(4) + zsize(4) + count(4) + data(count*2)
func ParseTable(data []byte) *RubyTable {
	if len(data) < 20 {
		return nil
	}

	dim := int(binary.LittleEndian.Uint32(data[0:4]))
	xsize := int(binary.LittleEndian.Uint32(data[4:8]))
	ysize := int(binary.LittleEndian.Uint32(data[8:12]))
	zsize := int(binary.LittleEndian.Uint32(data[12:16]))
	count := int(binary.LittleEndian.Uint32(data[16:20]))

	table := &RubyTable{
		Dim:   dim,
		XSize: xsize,
		YSize: ysize,
		ZSize: zsize,
	}

	// Parse table data if present
	if len(data) >= 20+count*2 {
		table.Data = make([]int16, count)
		for i := 0; i < count; i++ {
			offset := 20 + i*2
			table.Data[i] = int16(binary.LittleEndian.Uint16(data[offset : offset+2]))
		}
	}

	return table
}

// ParseColor decodes RPG Maker's Color userdata (4 doubles: r, g, b, a).
func ParseColor(data []byte) *RubyColor {
	if len(data) < 32 {
		return nil
	}

	return &RubyColor{
		Red:   math.Float64frombits(binary.LittleEndian.Uint64(data[0:8])),
		Green: math.Float64frombits(binary.LittleEndian.Uint64(data[8:16])),
		Blue:  math.Float64frombits(binary.LittleEndian.Uint64(data[16:24])),
		Alpha: math.Float64frombits(binary.LittleEndian.Uint64(data[24:32])),
	}
}

// ParseTone decodes RPG Maker's Tone userdata (4 doubles: r, g, b, gray).
func ParseTone(data []byte) *RubyTone {
	if len(data) < 32 {
		return nil
	}

	return &RubyTone{
		Red:   math.Float64frombits(binary.LittleEndian.Uint64(data[0:8])),
		Green: math.Float64frombits(binary.LittleEndian.Uint64(data[8:16])),
		Blue:  math.Float64frombits(binary.LittleEndian.Uint64(data[16:24])),
		Gray:  math.Float64frombits(binary.LittleEndian.Uint64(data[24:32])),
	}
}

// ToRubyObject attempts to convert an interface{} to *RubyObject.
func ToRubyObject(v interface{}) *RubyObject {
	if obj, ok := v.(*RubyObject); ok {
		return obj
	}
	return nil
}

// ToRubyObjectSlice converts []interface{} to []*RubyObject, skipping nil entries.
func ToRubyObjectSlice(arr []interface{}) []*RubyObject {
	result := make([]*RubyObject, len(arr))
	for i, v := range arr {
		if v != nil {
			result[i] = ToRubyObject(v)
		}
	}
	return result
}

// GetStringFromValue safely extracts a string from an interface{}.
func GetStringFromValue(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// GetIntFromValue safely extracts an int from an interface{}.
func GetIntFromValue(v interface{}) int {
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	}
	return 0
}

// GetBoolFromValue safely extracts a bool from an interface{}.
func GetBoolFromValue(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}
