package textemplate

import (
	"bytes"
	"io"
	"maps"
	"regexp"
	"strings"
)

// tagRe is the tag regular expression.
const tagRe = ` *[0-9a-zA-Z_\-]+(?:\.[A-Za-z0-9_\-]+)* *`

// TagFunc returns tag value.
type TagFunc func(tag string) (value []byte, found bool)

// segment is the template segment.
// It represents a fixed value or a tag.
type segment struct {
	isTag bool   // isTag if true, this segment is a tag to be resolved.
	tag   string // tag is the tag name.
	value []byte // value is, if not tag, the fixed value for this segment.
}

// New returns a new instance of the Template.
// Allowed tag name pattern is `[0-9a-zA-Z_\-]+(?:\.[A-Za-z0-9_\-]+)*`.
func New(tpl string, start, end string) *Template {
	qStart := regexp.QuoteMeta(start)
	qEnd := regexp.QuoteMeta(end)
	re := regexp.MustCompile(qStart + tagRe + qEnd)

	unescape := func(s string) []byte {
		s = strings.ReplaceAll(s, qStart, start)
		s = strings.ReplaceAll(s, qEnd, end)
		return []byte(s)
	}

	segs := []*segment{}
	bufSize := 0 // Initial estimated buffer size.
	pos := 0
	for _, ids := range re.FindAllStringIndex(tpl, -1) {
		tagStart, tagEnd := ids[0], ids[1]
		fixedVal := unescape(tpl[pos:tagStart])
		tagVal := strings.TrimSpace(tpl[tagStart+len(start) : tagEnd-len(end)])
		bufSize += len(fixedVal) + 16 // Occupy initially 16 bytes for tag value.
		segs = append(segs, &segment{value: fixedVal})
		segs = append(segs, &segment{isTag: true, tag: tagVal})
		pos = tagEnd
	}
	if pos < len(tpl) {
		segs = append(segs, &segment{value: unescape(tpl[pos:])})
	}

	return &Template{
		tagStart: start,
		tagEnd:   end,
		segs:     segs,
		bufSize:  bufSize,
	}
}

// Template is the text template engine.
// Use [New] to instantiate a new Template.
// Template supports following types.
//
//   - string
//   - []byte
//   - bool
//   - int, int8, int16, int32, int64
//   - uint, uint8, uint16, uint32, uint64
//   - float32, float64
//   - complex64, complex128
//   - interface{ String() string }
//   - interface{ Bytes() []bytes }
//   - interface{ Append([]byte) []bytes }
//   - func() string
//   - func() []byte
//   - others : fallback to fmt.Sprint
type Template struct {
	tagStart string         // tagStart is the tag start marker.
	tagEnd   string         // tagEnd is the tag end marker.
	segs     []*segment     // segs is the segment list.
	bufSize  int            // bufSize is the initial buffer size.
	defaults map[string]any // defaultVals are the default tag values.

	// notFound is called when no tag values were found.
	notFound func(io.Writer, string) error
}

// WithDefaults asscociates default values to the template.
// WithDefaults can be called multiple times.
// Currently there is no way to removed the registered value.
// It replaces the existing value if duplicate keys were provided.
func (t *Template) WithDefaults(values map[string]any) {
	if len(values) == 0 {
		return
	}
	if t.defaults == nil {
		t.defaults = make(map[string]any, len(values))
	}
	maps.Copy(t.defaults, values)
}

// WithNotFound registeres not found hadler to the template.
// notFound will be called as fallback when no tag values were found.
// It replaces the existing notFound handler when called multiple times.
// Use WithNotFound(nil) to remove the registered handler.
func (t *Template) WithNotFound(notFound func(io.Writer, string) error) {
	t.notFound = notFound
}

// Execute executes the template with given values.
// See [Template] for supported data types.
func (t *Template) Execute(m map[string]any) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, t.bufSize))
	err := t.execute(buf, m, nil)
	return buf.Bytes(), err

}

// MustExecute calls [Execute].
// It panics if an error was returned from [Execute].
func (t *Template) MustExecute(m map[string]any) []byte {
	b, err := t.Execute(m)
	mustNil(err)
	return b

}

// ExecuteString executes the template with given values.
// See [Template] for supported data types.
func (t *Template) ExecuteString(m map[string]any) (string, error) {
	buf := bytes.NewBuffer(make([]byte, 0, t.bufSize))
	err := t.execute(buf, m, nil)
	return buf.String(), err
}

// MustExecuteString calls [ExecuteString].
// It panics if an error was returned from [ExecuteString].
func (t *Template) MustExecuteString(m map[string]any) string {
	s, err := t.ExecuteString(m)
	mustNil(err)
	return s
}

// ExecuteFunc executes the template with tag function.
func (t *Template) ExecuteFunc(tf TagFunc) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, t.bufSize))
	err := t.execute(buf, nil, tf)
	return buf.Bytes(), err
}

// MustExecuteFunc calls [ExecuteFunc].
// It panics if an error was returned from [ExecuteFunc].
func (t *Template) MustExecuteFunc(tf TagFunc) []byte {
	b, err := t.ExecuteFunc(tf)
	mustNil(err)
	return b
}

// ExecuteStringFunc executes the template with tag function.
func (t *Template) ExecuteFuncString(tf TagFunc) (string, error) {
	buf := bytes.NewBuffer(make([]byte, 0, t.bufSize))
	err := t.execute(buf, nil, tf)
	return buf.String(), err
}

// MustExecuteFuncString calls [ExecuteFuncString].
// It panics if an error was returned from [ExecuteFuncString].
func (t *Template) MustExecuteFuncString(tf TagFunc) string {
	s, err := t.ExecuteFuncString(tf)
	mustNil(err)
	return s
}

// ExecuteWriter executes the template and writes result to w.
// Note that the w.Wirte will be called multiple times.
// See [Template] for supported data types.
func (t *Template) ExecuteWriter(w io.Writer, m map[string]any) error {
	return t.execute(w, m, nil)
}

// ExecuteWriterFunc executes the template with tag function.
// Note that the w.Wirte will be called multiple times.
func (t *Template) ExecuteWriterFunc(w io.Writer, tf TagFunc) error {
	return t.execute(w, nil, tf)
}

func (t *Template) execute(w io.Writer, m map[string]any, tf TagFunc) (err error) {
	defer func() { recover() }() // Recover from panics.

	var found bool
	buf := make([]byte, 0, 32) // Reusable buffer.
	for _, s := range t.segs {
		if !s.isTag {
			_, err = w.Write(s.value)
			mustNil(err)
			continue
		}
		if m != nil {
			if buf, found = appendTagValue(buf[:0], m, s.tag); found {
				_, err = w.Write(buf)
				mustNil(err)
				continue
			}
		}
		if tf != nil {
			if b, found := tf(s.tag); found {
				_, err = w.Write(b)
				mustNil(err)
				continue
			}
		}
		if t.defaults != nil {
			if buf, found = appendTagValue(buf[:0], t.defaults, s.tag); found {
				_, err = w.Write(buf)
				mustNil(err)
				continue
			}
		}
		if t.notFound != nil {
			err = t.notFound(w, s.tag)
			mustNil(err)
		} else {
			_, err = w.Write([]byte(t.tagStart + s.tag + t.tagEnd))
			mustNil(err)
		}
	}
	return nil
}

func mustNil(err error) {
	if err != nil {
		panic(err)
	}
}
