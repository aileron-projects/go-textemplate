package textemplate

import (
	"bytes"
	"io"
	"regexp"
	"testing"

	"github.com/aileron-projects/go-tester"
)

func TestTagRe(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		value string
		match bool
	}{
		{"{}", false},
		{"{ }", false},
		{"{test}", true},
		{"{ test}", true},
		{"{test }", true},
		{"{ test }", true},
		{"{ te-st }", true},
		{"{ te_st }", true},
		{"{ te st }", false},
		{"{ .test }", false},
		{"{ test. }", false},
		{"{ te..st }", false},
		{"{ te.st }", true},
		{"{ t.e.s.t }", true},
	}
	re := regexp.MustCompile("\\{" + tagRe + "\\}")
	for _, tc := range testCases {
		t.Run(tc.value, func(t *testing.T) {
			matched := re.MatchString(tc.value)
			tester.AssertEqual(t, tc.match, matched)
		})
	}
}

func TestTemplate_bracket(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		start string
		end   string
		tpl   string
		want  string
	}{
		"ok {}":       {`{`, `}`, `{ value }`, "test"},
		"ok {{}}":     {`{{`, `}}`, `{{ value }}`, "test"},
		"ok []":       {`[`, `]`, `[ value ]`, "test"},
		"ok [[]]":     {`[[`, `]]`, `[[ value ]]`, "test"},
		"ok %%":       {`%`, `%`, `% value %`, "test"},
		"ok %%%%":     {`%%`, `%%`, `%% value %%`, "test"},
		"ok [}":       {`[`, `}`, `[ value }`, "test"},
		"escape {}":   {`{`, `}`, `\{{ value }\}`, "{test}"},
		"escape {{}}": {`{{`, `}}`, `\{\{{{ value }}\}\}`, "{{test}}"},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tpl := New(tc.tpl, tc.start, tc.end)
			got := tpl.ExecuteString(map[string]any{"value": "test"})
			tester.AssertEqual(t, tc.want, got)
		})
	}
}

func TestTemplate_map(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		tpl  string
		want string
		val  map[string]any
	}{
		"start with tag": {`{{value}} bar`, `foo bar`, map[string]any{"value": "foo"}},
		"end with tag":   {`foo {{value}}`, `foo bar`, map[string]any{"value": "bar"}},
		"middle tag":     {`foo {{value}} baz`, `foo bar baz`, map[string]any{"value": "bar"}},
		"multiple tags":  {`{{f}} {{b}} {{z}}`, `foo bar baz`, map[string]any{"f": "foo", "b": "bar", "z": "baz"}},
		"no tags":        {`foo bar`, `foo bar`, map[string]any{}},
		"escape bracket": {`\{\{value\}\}`, `{{value}}`, map[string]any{"value": "***"}},
		"escape bra":     {`\{\{{{value}}`, `{{foo`, map[string]any{"value": "foo"}},
		"escape cket":    {`{{value}}\}\}`, `foo}}`, map[string]any{"value": "foo"}},
	}

	t.Run("Execute", func(t *testing.T) {
		for name, tc := range testCases {
			t.Run(name, func(t *testing.T) {
				tpl := New(tc.tpl, "{{", "}}")
				got := tpl.Execute(tc.val)
				tester.AssertEqual(t, tc.want, string(got))
			})
		}
	})

	t.Run("ExecuteString", func(t *testing.T) {
		for name, tc := range testCases {
			t.Run(name, func(t *testing.T) {
				tpl := New(tc.tpl, "{{", "}}")
				got := tpl.ExecuteString(tc.val)
				tester.AssertEqual(t, tc.want, got)
			})
		}
	})

	t.Run("ExecuteWriter", func(t *testing.T) {
		for name, tc := range testCases {
			t.Run(name, func(t *testing.T) {
				var b bytes.Buffer
				b.WriteString(">>")
				tpl := New(tc.tpl, "{{", "}}")
				err := tpl.ExecuteWriter(&b, tc.val)
				tester.AssertEqual(t, nil, err)
				tester.AssertEqual(t, ">>"+tc.want, b.String())
			})
		}
	})
}

func TestTemplate_func(t *testing.T) {
	t.Parallel()
	foo := []byte("foo")
	bar := []byte("bar")

	testCases := map[string]struct {
		tpl  string
		want string
		tf   func(string) ([]byte, bool)
	}{
		"value found":     {`{{value}}`, `foo`, func(string) ([]byte, bool) { return foo, true }},
		"value not found": {`{{value}}`, `{{value}}`, func(string) ([]byte, bool) { return nil, false }},
		"start with tag":  {`{{value}} bar`, `foo bar`, func(string) ([]byte, bool) { return foo, true }},
		"end with tag":    {`foo {{value}}`, `foo bar`, func(string) ([]byte, bool) { return bar, true }},
		"middle tag":      {`foo {{value}} baz`, `foo bar baz`, func(string) ([]byte, bool) { return bar, true }},
		"nil value":       {`{{value}}`, ``, func(string) ([]byte, bool) { return nil, true }},
	}

	t.Run("ExecuteFunc", func(t *testing.T) {
		for name, tc := range testCases {
			t.Run(name, func(t *testing.T) {
				tpl := New(tc.tpl, "{{", "}}")
				got := tpl.ExecuteFunc(tc.tf)
				tester.AssertEqual(t, tc.want, string(got))
			})
		}
	})

	t.Run("ExecuteFuncString", func(t *testing.T) {
		for name, tc := range testCases {
			t.Run(name, func(t *testing.T) {
				tpl := New(tc.tpl, "{{", "}}")
				got := tpl.ExecuteFuncString(tc.tf)
				tester.AssertEqual(t, tc.want, got)
			})
		}
	})

	t.Run("ExecuteWriterFunc", func(t *testing.T) {
		for name, tc := range testCases {
			t.Run(name, func(t *testing.T) {
				var b bytes.Buffer
				b.WriteString(">>")
				tpl := New(tc.tpl, "{{", "}}")
				err := tpl.ExecuteWriterFunc(&b, tc.tf)
				tester.AssertEqual(t, nil, err)
				tester.AssertEqual(t, ">>"+tc.want, b.String())
			})
		}
	})
}

func TestTemplate_WithDefaults(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		want string
		tv   map[string]any // tag value
		val  map[string]any
	}{
		"nil":        {`{{value}}`, nil, nil},
		"not found":  {`{{value}}`, map[string]any{"unused": "***"}, nil},
		"found":      {`foo`, map[string]any{"value": "foo"}, nil},
		"overridden": {`bar`, map[string]any{"value": "foo"}, map[string]any{"value": "bar"}},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tpl := New(`{{value}}`, "{{", "}}")
			tpl.WithDefaults(tc.tv)
			got := tpl.ExecuteString(tc.val)
			tester.AssertEqual(t, tc.want, string(got))
		})
	}
}
func TestTemplate_WithNotFound(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		want string
		nf   func(io.Writer, string) error // notFound handler
	}{
		"default": {`{{foo.bar}}`, nil},
		"custom handler": {`NotFound:foo.bar`, func(w io.Writer, s string) error {
			_, err := w.Write([]byte("NotFound:" + s))
			return err
		}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tpl := New(`{{foo.bar}}`, "{{", "}}")
			tpl.WithNotFound(tc.nf)
			got := tpl.ExecuteString(nil)
			tester.AssertEqual(t, tc.want, string(got))
		})
	}
}

func TestTemplate_Panics(t *testing.T) {
	t.Parallel()
	t.Run("fixed value write", func(t *testing.T) {
		tpl := New(`test`, "{{", "}}")
		w := tester.MaxErrorWriter(3)
		err := tpl.ExecuteWriter(w, nil)
		tester.AssertEqual(t, tester.ErrMaxWritten, err)
		tester.AssertEqual(t, w.String(), "tes")
	})
	t.Run("tag value write", func(t *testing.T) {
		tpl := New(`{{value}}`, "{{", "}}")
		w := tester.MaxErrorWriter(3)
		err := tpl.ExecuteWriter(w, map[string]any{"value": "test"})
		tester.AssertEqual(t, tester.ErrMaxWritten, err)
		tester.AssertEqual(t, w.String(), "tes")
	})
	t.Run("tag func write", func(t *testing.T) {
		tpl := New(`{{value}}`, "{{", "}}")
		w := tester.MaxErrorWriter(3)
		err := tpl.ExecuteWriterFunc(w, func(tag string) (value []byte, found bool) {
			return []byte("test"), true
		})
		tester.AssertEqual(t, tester.ErrMaxWritten, err)
		tester.AssertEqual(t, w.String(), "tes")
	})
	t.Run("default value write", func(t *testing.T) {
		tpl := New(`{{value}}`, "{{", "}}")
		tpl.WithDefaults(map[string]any{"value": "test"})
		w := tester.MaxErrorWriter(3)
		err := tpl.ExecuteWriter(w, nil)
		tester.AssertEqual(t, tester.ErrMaxWritten, err)
		tester.AssertEqual(t, w.String(), "tes")
	})
	t.Run("not found write", func(t *testing.T) {
		tpl := New(`{{value}}`, "{{", "}}")
		w := tester.MaxErrorWriter(3)
		err := tpl.ExecuteWriter(w, nil)
		tester.AssertEqual(t, tester.ErrMaxWritten, err)
		tester.AssertEqual(t, w.String(), "{{v")
	})
}
