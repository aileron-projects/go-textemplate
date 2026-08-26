package textemplate

import (
	"testing"

	"github.com/aileron-projects/go-tester"
)

type testStringer struct {
	s string
}

func (t *testStringer) String() string {
	return t.s
}

type testByter struct {
	s string
}

func (t *testByter) Bytes() []byte {
	return []byte(t.s)
}

type testAppender struct {
	s string
}

func (t *testAppender) Append(dst []byte) []byte {
	return append(dst, []byte(t.s)...)
}

func TestAppendTagValue(t *testing.T) {
	t.Parallel()
	t.Run("value", func(t *testing.T) {
		testCases := map[string]struct {
			m      map[string]any
			tag    string
			expect string
			found  bool
		}{
			"not found":            {map[string]any{"key": nil}, "KEY", "#", false},
			"nil map":              {nil, "key", "#", false},
			"nil":                  {map[string]any{"key": nil}, "key", "#<nil>", true},
			"string":               {map[string]any{"key": "test"}, "key", "#test", true},
			"[]byte":               {map[string]any{"key": []byte("test")}, "key", "#test", true},
			"bool":                 {map[string]any{"key": true}, "key", "#true", true},
			"int":                  {map[string]any{"key": 123}, "key", "#123", true},
			"int8":                 {map[string]any{"key": int8(123)}, "key", "#123", true},
			"int16":                {map[string]any{"key": int16(123)}, "key", "#123", true},
			"int32":                {map[string]any{"key": int32(123)}, "key", "#123", true},
			"int64":                {map[string]any{"key": int64(123)}, "key", "#123", true},
			"uint":                 {map[string]any{"key": uint(123)}, "key", "#123", true},
			"uint8":                {map[string]any{"key": uint8(123)}, "key", "#123", true},
			"uint16":               {map[string]any{"key": uint16(123)}, "key", "#123", true},
			"uint32":               {map[string]any{"key": uint32(123)}, "key", "#123", true},
			"uint64":               {map[string]any{"key": uint64(123)}, "key", "#123", true},
			"complex64":            {map[string]any{"key": complex64(1 + 2i)}, "key", "#(1+2i)", true},
			"complex128":           {map[string]any{"key": complex128(1 - 2i)}, "key", "#(1-2i)", true},
			"stringer":             {map[string]any{"key": &testStringer{"test"}}, "key", "#test", true},
			"byter":                {map[string]any{"key": &testByter{"test"}}, "key", "#test", true},
			"appender":             {map[string]any{"key": &testAppender{"test"}}, "key", "#test", true},
			"tagfunc":              {map[string]any{"key": TagValueFunc(func(string) []byte { return []byte("test") })}, "key", "#test", true},
			"struct":               {map[string]any{"key": struct{ x string }{"test"}}, "key", "#{test}", true},
			"nested map":           {map[string]any{"foo": map[string]any{"bar": 123}}, "foo.bar", "#123", true},
			"nested map not found": {map[string]any{"foo": map[string]any{"baz": 123}}, "foo.bar", "#", false},
			"nested map is nil":    {map[string]any{"foo": nil}, "foo.bar", "#", false},
		}

		for name, tc := range testCases {
			t.Run(name, func(t *testing.T) {
				dst := []byte("#")
				got, found := appendTagValue(dst, tc.m, tc.tag)
				tester.AssertEqual(t, found, tc.found)
				tester.AssertEqual(t, tc.expect, string(got))
			})
		}
	})
}
