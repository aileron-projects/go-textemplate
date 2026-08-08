package textemplate_test

import (
	"bytes"
	"testing"
	ttpl "text/template"

	"github.com/aileron-projects/go-textemplate"
)

var template = `
nil = {{ .nil }}
string = {{ .string }}
int = {{ .int }}
int8 = {{ .int8 }}
int16 = {{ .int16 }}
int32 = {{ .int32 }}
int64 = {{ .int64 }}
uint = {{ .uint }}
uint8 = {{ .uint8 }}
uint16 = {{ .uint16 }}
uint32 = {{ .uint32 }}
uint64 = {{ .uint64 }}
float32 = {{ .float32 }}
float64 = {{ .float64 }}
complex64 = {{ .complex64 }}
complex128 = {{ .complex128 }}
`

var values = map[string]any{
	"nil":        nil,
	"string":     "foo",
	"int":        123,
	"int8":       int8(123),
	"int16":      int16(123),
	"int32":      int32(123),
	"int64":      int64(123),
	"uint":       uint(123),
	"uint8":      uint8(123),
	"uint16":     uint16(123),
	"uint32":     uint32(123),
	"uint64":     uint64(123),
	"float32":    float32(1.141592653589),
	"float64":    float64(1.141592653589),
	"complex64":  complex64(123 + 456i),
	"complex128": complex128(123 + 456i),
}

func BenchmarkTextemplate(b *testing.B) {
	tpl := textemplate.New(template, "{{", "}}")

	b.ResetTimer()
	for b.Loop() {
		tpl.Execute(values)
	}
}

func BenchmarkTextTemplate(b *testing.B) {
	// Standard package text/template.
	tpl, _ := ttpl.New("").Parse(template)

	b.ResetTimer()
	for b.Loop() {
		var b bytes.Buffer
		tpl.Execute(&b, values)
	}
}
