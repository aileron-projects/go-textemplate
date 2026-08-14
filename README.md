<!-- markdownlint-disable MD033 MD041 -->

<div align="center">

[![Release](https://img.shields.io/github/v/release/aileron-projects/go-textemplate?sort=semver)](https://github.com/aileron-projects/go-textemplate/releases)
[![Reference](https://pkg.go.dev/badge/github.com/aileron-projects/go-textemplate.svg)](https://pkg.go.dev/github.com/aileron-projects/go-textemplate)
[![DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/aileron-projects/go-textemplate)
[![Test](https://github.com/aileron-projects/go-textemplate/actions/workflows/test.yaml/badge.svg)](https://github.com/aileron-projects/go/actions/workflows/test.yaml)

[![Insights](https://badgen.net/badge/Insights/open%2Fsource%2Finsights/cyan)](https://deps.dev/go/github.com%2Faileron-projects%go-textemplate)
[![Insights](https://badgen.net/badge/Insights/OSS%2FInsight/orange)](https://ossinsight.io/analyze/aileron-projects/go-textemplate)

</div>

# go-textemplate

**Simple and fast text template engine for Go.**

## Features

- Fast and lightweight.
- Tagged value embedding like `Hello {{ world }}`
- Custom bracket.
- Various type support.

## Usages

### Basic Usage

- Allowed tag pattern is `[0-9a-zA-Z_\-]+(?:\.[A-Za-z0-9_\-]+)*`.
  - OK: `foo`, `foo.bar`, `foo.bar.baz`, `foo-bar`, `foo_bar`
  - NG: `.foo`, `foo.` Tag must not start or end with dots.
  - NG: `foo..bar` Dots must not be consecutive.
- Dots in tags are only used for accessing nested map.
- Brackets are customizable.
  - `{{`, `}}` are very popular (`{{ foo }}`).
  - Other patterns are also allowed like `% foo %`, `& foo &`, `^ foo $` etc.
- Escape brackets if necessary.
  - For example, `\{\{` for `{{`.
  - Escape expression is determined by [regexp#QuoteMeta](https://pkg.go.dev/regexp#QuoteMeta).
- Tags are output as it is if corresponding value was not found.

```go
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
```

### Nested map value

Nested maps of `map[string]any` are allowed.
Use dots to access internal maps and values.

```go
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
```

### Values by TagValueFunc

TagValueFunc `func(tag string) []byte` can provide tag values.

```go
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
```

### Default values

Default tag-value pairs can be bounded to a template with `WithDefaults` method.

```go
defaults := map[string]any{"value": "world!!"}
values := map[string]any{"value": "Go!!"}

template := `Hello {{ value }}`
tpl := textemplate.New(template, "{{", "}}")
tpl.WithDefaults(defaults) // Register default values.

fmt.Println(tpl.ExecuteString(nil))    // Hello world!!
fmt.Println(tpl.ExecuteString(values)) // Hello Go!!
```

### Escape brackets

```go
template := `\{\{ "Hello": "{{ value }}" \}\}`
values := map[string]any{
    "value": "world!!",
}

tpl := textemplate.New(template, "{{", "}}")
result := tpl.Execute(values)

fmt.Println(string(result))
// Output:
// {{ "Hello": "world!!" }}
```

## Docs & Examples

- GoDoc: <https://pkg.go.dev/github.com/aileron-projects/go-textemplate>
- Examples: [example_test.go](./example_test.go)
- Benchmarks: [benchmark_test.go](./benchmark_test.go)

## Benchmarks

### Simple string only  template

Template:

```go
template := `{{scheme}}://{{host}}{{path}}?foo={{foo}}&bar={{bar}}#{{fragment}}`
values := map[string]any{
    "scheme":   "https",
    "host":     "godoc.org",
    "path":     "/github.com/aileron-projects/go-textemplate",
    "foo":      "dummy",
    "bar":      "dummy",
    "fragment": "section-readme",
}
```

Results:

```txt
                                    Iteration          Speed  Heap size    Alloc count
                                      -------   ------------   --------   ------------
BenchmarkFmtFprintf-8                 5278726    357.0 ns/op   144 B/op    2 allocs/op
BenchmarkStringsReplaceAll-8          1416745    785.9 ns/op   528 B/op    6 allocs/op
BenchmarkTextTemplate-8                847612   2360.0 ns/op   560 B/op   17 allocs/op
BenchmarkFasttemplate-8               1400848    731.0 ns/op   344 B/op    9 allocs/op
BenchmarkTextemplateExecute-8         2457674    467.1 ns/op   256 B/op    4 allocs/op ★
BenchmarkTextemplateExecuteString-8   1848076    580.1 ns/op   352 B/op    5 allocs/op ★
BenchmarkTextemplateExecuteWriter-8   2002717    505.3 ns/op   336 B/op    5 allocs/op ★
```

### Various types template

See [./benchmark_test.go](./benchmark_test.go).

```txt
                       Iteration        Speed        Heap     Allocation
                         -------   ----------   ---------   ------------
BenchmarkTextemplate-8    822277   1359 ns/op    496 B/op    3 allocs/op ★ This library
BenchmarkTextTemplate-8   301384   4457 ns/op   1136 B/op   38 allocs/op ★ Std text/template
```

## References

- [text template package](https://pkg.go.dev/text/template)
- [valyala/fasttemplate](https://github.com/valyala/fasttemplate)
