package textemplate_test

import (
	"fmt"
	"io"

	textemplate "github.com/aileron-projects/go-textemplate"
)

func Example() {
	template := `Hello {{ value }}!!`
	tpl := textemplate.New(template, "{{", "}}")
	val := map[string]any{
		"value": "world",
	}
	fmt.Println(tpl.ExecuteString(val))
	// Output:
	// Hello world!!
}

func ExampleTemplate_basic() {
	template := `
nil = {{ nil }}
string = {{ string }}
int = {{ int }}
int8 = {{ int8 }}
int16 = {{ int16 }}
int32 = {{ int32 }}
int64 = {{ int64 }}
uint = {{ uint }}
uint8 = {{ uint8 }}
uint16 = {{ uint16 }}
uint32 = {{ uint32 }}
uint64 = {{ uint64 }}
float32 = {{ float32 }}
float64 = {{ float64 }}
complex64 = {{ complex64 }}
complex128 = {{ complex128 }}
struct = {{ struct }}
NotFound = {{ not-found }}
`

	values := map[string]any{
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
		"struct":     struct{ a, b string }{"foo", "bar"},
	}

	tpl := textemplate.New(template, "{{", "}}")
	result := tpl.Execute(values)

	fmt.Println(string(result))
	// Output:
	// nil = <nil>
	// string = foo
	// int = 123
	// int8 = 123
	// int16 = 123
	// int32 = 123
	// int64 = 123
	// uint = 123
	// uint8 = 123
	// uint16 = 123
	// uint32 = 123
	// uint64 = 123
	// float32 = 1.1415926
	// float64 = 1.141592653589
	// complex64 = (123+456i)
	// complex128 = (123+456i)
	// struct = {foo bar}
	// NotFound = {{not-found}}
}

func ExampleTemplate_nestmap() {
	template := `foo.bar.baz = {{ foo.bar.baz }}`
	values := map[string]any{
		"foo": map[string]any{
			"bar": map[string]any{
				"baz": "FOO.BAR.BAZ",
			},
		},
		"foo.bar.baz": "this field has no effect",
	}

	tpl := textemplate.New(template, "{{", "}}")
	result := tpl.Execute(values)

	fmt.Println(string(result))
	// Output:
	// foo.bar.baz = FOO.BAR.BAZ
}

func ExampleTemplate_escape() {
	template := `\{\{ "Hello": "{{ value }}" \}\}`
	values := map[string]any{
		"value": "world!!",
	}

	tpl := textemplate.New(template, "{{", "}}")
	result := tpl.Execute(values)

	fmt.Println(string(result))
	// Output:
	// {{ "Hello": "world!!" }}
}

func ExampleTagValueFunc() {
	template := `Hello {{ value }}`
	values := map[string]any{
		"value": textemplate.TagValueFunc(func(tag string) []byte {
			return []byte("world!!")
		}),
	}

	tpl := textemplate.New(template, "{{", "}}")
	result := tpl.Execute(values)

	fmt.Println(string(result))
	// Output:
	// Hello world!!
}

func ExampleTemplate_WithDefaults() {
	defaults := map[string]any{"value": "world!!"}
	values := map[string]any{"value": "Go!!"}

	template := `Hello {{ value }}`
	tpl := textemplate.New(template, "{{", "}}")
	tpl.WithDefaults(defaults) // Register default values.

	fmt.Println(tpl.ExecuteString(nil))    // Hello world!!
	fmt.Println(tpl.ExecuteString(values)) // Hello Go!!
	// Output:
	// Hello world!!
	// Hello Go!!
}

func ExampleTemplate_WithNotFound() {
	values := map[string]any{"value": "Go!!"}

	template := `Hello {{ value }}`
	tpl := textemplate.New(template, "{{", "}}")
	tpl.WithNotFound(func(w io.Writer, tag string) error { // Register not found handler.
		_, err := w.Write([]byte("#NotFound:" + tag))
		return err
	})

	fmt.Println(tpl.ExecuteString(nil))    // Hello #NotFound:value
	fmt.Println(tpl.ExecuteString(values)) // Hello Go!!
	// Output:
	// Hello #NotFound:value
	// Hello Go!!
}
