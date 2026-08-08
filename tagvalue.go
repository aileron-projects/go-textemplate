package textemplate

import (
	"fmt"
	"strconv"
	"strings"
)

// TagValueFunc returns tag value.
type TagValueFunc func(tag string) []byte

// appendTagValue appends tag value to dst and returns new by slice.
// It returns false if tag was not found in the m.
func appendTagValue(dst []byte, m map[string]any, tag string) (b []byte, found bool) {
	if m == nil {
		return dst, false
	}
	if b, a, found := strings.Cut(tag, "."); found {
		if mm, ok := m[b].(map[string]any); ok {
			return appendTagValue(dst, mm, a)
		}
		return dst, false
	}
	value, ok := m[tag]
	if !ok {
		return dst, false
	}
	switch v := value.(type) {
	case string:
		b = append(dst, []byte(v)...)
	case fmt.Stringer:
		b = append(dst, []byte(v.String())...)
	case []byte:
		b = append(dst, v...)
	case bool:
		b = strconv.AppendBool(dst, v)
	case int:
		b = strconv.AppendInt(dst, int64(v), 10)
	case int8:
		b = strconv.AppendInt(dst, int64(v), 10)
	case int16:
		b = strconv.AppendInt(dst, int64(v), 10)
	case int32:
		b = strconv.AppendInt(dst, int64(v), 10)
	case int64:
		b = strconv.AppendInt(dst, v, 10)
	case float32:
		b = strconv.AppendFloat(dst, float64(v), 'g', -1, 32)
	case float64:
		b = strconv.AppendFloat(dst, float64(v), 'g', -1, 64)
	case uint:
		b = strconv.AppendUint(dst, uint64(v), 10)
	case uint8:
		b = strconv.AppendUint(dst, uint64(v), 10)
	case uint16:
		b = strconv.AppendUint(dst, uint64(v), 10)
	case uint32:
		b = strconv.AppendUint(dst, uint64(v), 10)
	case uint64:
		b = strconv.AppendUint(dst, v, 10)
	case complex64:
		b = appendComplex(dst, complex128(v), 'g', -1, 32)
	case complex128:
		b = appendComplex(dst, complex128(v), 'g', -1, 64)
	case TagValueFunc:
		b = append(dst, v(tag)...)
	default:
		b = fmt.Append(dst, v) // Fallback to "%+v"
	}
	return b, true
}

func appendComplex(dst []byte, c complex128, fmt byte, prec int, bitSize int) []byte {
	dst = append(dst, '(')
	dst = strconv.AppendFloat(dst, float64(real(c)), fmt, prec, bitSize)
	if imag(c) >= 0 {
		dst = append(dst, '+')
	}
	dst = strconv.AppendFloat(dst, float64(imag(c)), fmt, prec, bitSize)
	dst = append(dst, 'i', ')')
	return dst
}
