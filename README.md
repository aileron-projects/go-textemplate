# go-textemplate

**Simple and fast text template engine for Go.**

<div align="center">

[![GoDoc](https://godoc.org/github.com/aileron-projects/go-textemplate?status.svg)](http://godoc.org/github.com/aileron-projects/go-textemplate)
[![Test](https://github.com/aileron-projects/go-textemplate/actions/workflows/test.yaml/badge.svg?branch=main)](https://github.com/aileron-projects/go-textemplate/actions/workflows/test.yaml?query=branch%3Amain)
[![License](https://img.shields.io/badge/License-Apache%202.0-yellow.svg)](./LICENSE)

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/aileron-projects/go-textemplate)
[![OpenSourceInsight](https://badgen.net/badge/open%2Fsource%2F/insight/cyan)](https://deps.dev/go/github.com%2Faileron-projects%2Fgo-textemplate)
[![OSS Insight](https://badgen.net/badge/OSS/Insight/orange)](https://ossinsight.io/analyze/aileron-projects/go-textemplate)

</div>

## Features

- Fast and lightweight.
- Tagged value embedding like `Hello {{ world }}`
- Custom bracket.
- Various type support.

## Tested Environments

Operating System:

- `Linux`: [ubuntu-latest](https://github.com/actions/runner-images)
- `Windows`: [windows-latest](https://github.com/actions/runner-images)
- `macOS`: [macos-latest](https://github.com/actions/runner-images)

Architecture (Using QEMU on linux):

- x86: `amd64`, `386`
- arm: `arm/v5`, `arm/v6`, `arm/v7`, `arm64`
- risc: `riscv64`, `loong64`
- ppc: `ppc64`, `ppc64le`
- mips: `mips`, `mips64`, `mips64le`, `mipsle`
- ibm: `s390x`

## Release Cycle

- Releases are made as needed.
- [Semantic Versioning](https://semver.org/) `vX.Y.Z` is used.

## License

[Apache-2.0](LICENSE)

## Usage

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

## Build Tags

No build tags defined for this library.

## Enviromental Variables

No environmental variables defined for this library.

## Benchmark

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
