package textemplate_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	ttpl "text/template"

	"github.com/aileron-projects/go-textemplate"
	"github.com/valyala/fasttemplate"
)

func BenchmarkFmtFprintf(b *testing.B) {
	template := `%s://%s%s?foo=%s&bar=%s#%s`

	b.ResetTimer()
	for b.Loop() {
		var b bytes.Buffer
		fmt.Fprintf(&b, template, "https", "godoc.org", "/github.com/aileron-projects/go-textemplate", "dummy", "dummy", "section-readme")
	}
}

func BenchmarkStringsReplaceAll(b *testing.B) {
	template := `{{scheme}}://{{host}}{{path}}?foo={{foo}}&bar={{bar}}#{{fragment}}`

	b.ResetTimer()
	for b.Loop() {
		s := template
		s = strings.ReplaceAll(s, "{{scheme}}", "https")
		s = strings.ReplaceAll(s, "{{host}}", "godoc.org")
		s = strings.ReplaceAll(s, "{{path}}", "/github.com/aileron-projects/go-textemplate")
		s = strings.ReplaceAll(s, "{{foo}}", "dummy")
		s = strings.ReplaceAll(s, "{{bar}}", "dummy")
		s = strings.ReplaceAll(s, "{{fragment}}", "section-readme")
	}
}

func BenchmarkTextTemplate(b *testing.B) {
	template := `{{.scheme}}://{{.host}}{{.path}}?foo={{.foo}}&bar={{.bar}}#{{.fragment}}`
	values := map[string]any{
		"scheme":   "https",
		"host":     "godoc.org",
		"path":     "/github.com/aileron-projects/go-textemplate",
		"foo":      "dummy",
		"bar":      "dummy",
		"fragment": "section-readme",
	}

	tpl, _ := ttpl.New("").Parse(template)

	b.ResetTimer()
	for b.Loop() {
		var b bytes.Buffer
		tpl.Execute(&b, values)
	}
}

func BenchmarkFasttemplate(b *testing.B) {
	template := `{{scheme}}://{{host}}{{path}}?foo={{foo}}&bar={{bar}}#{{fragment}}`
	values := map[string]any{
		"scheme":   "https",
		"host":     "godoc.org",
		"path":     "/github.com/aileron-projects/go-textemplate",
		"foo":      "dummy",
		"bar":      "dummy",
		"fragment": "section-readme",
	}

	tpl := fasttemplate.New(template, "{{", "}}")

	b.ResetTimer()
	for b.Loop() {
		var b bytes.Buffer
		tpl.Execute(&b, values)
	}
}

func BenchmarkTextemplateExecute(b *testing.B) {
	template := `{{scheme}}://{{host}}{{path}}?foo={{foo}}&bar={{bar}}#{{fragment}}`
	values := map[string]any{
		"scheme":   "https",
		"host":     "godoc.org",
		"path":     "/github.com/aileron-projects/go-textemplate",
		"foo":      "dummy",
		"bar":      "dummy",
		"fragment": "section-readme",
	}

	tpl := textemplate.New(template, "{{", "}}")

	b.ResetTimer()
	for b.Loop() {
		tpl.Execute(values)
	}
}

func BenchmarkTextemplateExecuteString(b *testing.B) {
	template := `{{scheme}}://{{host}}{{path}}?foo={{foo}}&bar={{bar}}#{{fragment}}`
	values := map[string]any{
		"scheme":   "https",
		"host":     "godoc.org",
		"path":     "/github.com/aileron-projects/go-textemplate",
		"foo":      "dummy",
		"bar":      "dummy",
		"fragment": "section-readme",
	}

	tpl := textemplate.New(template, "{{", "}}")

	b.ResetTimer()
	for b.Loop() {
		tpl.ExecuteString(values)
	}
}

func BenchmarkTextemplateExecuteWriter(b *testing.B) {
	template := `{{scheme}}://{{host}}{{path}}?foo={{foo}}&bar={{bar}}#{{fragment}}`
	values := map[string]any{
		"scheme":   "https",
		"host":     "godoc.org",
		"path":     "/github.com/aileron-projects/go-textemplate",
		"foo":      "dummy",
		"bar":      "dummy",
		"fragment": "section-readme",
	}

	tpl := textemplate.New(template, "{{", "}}")

	b.ResetTimer()
	for b.Loop() {
		var b bytes.Buffer
		tpl.ExecuteWriter(&b, values)
	}
}
