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
//   - TagValueFunc
//   - others : fallback to fmt.Sprint
type Template struct {
	tagStart string         // tagStart is the tag start marker.
	tagEnd   string         // tagEnd is the tag end marker.
	segs     []*segment     // segs is the segment list.
	bufSize  int            // bufSize is the initial buffer size.
	defaults map[string]any // defaultVals are the default values.
}

// WithDefaults bounds tag-value pairs used by default to the template.
func (t *Template) WithDefaults(values map[string]any) {
	if len(values) == 0 {
		return
	}
	if t.defaults == nil {
		t.defaults = make(map[string]any, len(values))
	}
	maps.Copy(t.defaults, values)
}

// Execute executes the template and returns result.
// See [Template] for supported data types.
func (t *Template) Execute(m map[string]any) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, t.bufSize))
	err := t.execute(buf, m, nil)
	mustNil(err)
	return buf.Bytes()

}

// ExecuteString executes the template and returns result.
// See [Template] for supported data types.
func (t *Template) ExecuteString(m map[string]any) string {
	buf := bytes.NewBuffer(make([]byte, 0, t.bufSize))
	err := t.execute(buf, m, nil)
	mustNil(err)
	return buf.String()
}

// ExecuteFunc executes the template with tag function.
func (t *Template) ExecuteFunc(tf TagFunc) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, t.bufSize))
	err := t.execute(buf, nil, tf)
	mustNil(err)
	return buf.Bytes()
}

// ExecuteStringFunc executes the template with tag function.
func (t *Template) ExecuteFuncString(tf TagFunc) string {
	buf := bytes.NewBuffer(make([]byte, 0, t.bufSize))
	err := t.execute(buf, nil, tf)
	mustNil(err)
	return buf.String()
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
		// Tag value was not found.
		// Write tags as it is so it can be replaced later.
		_, err = w.Write([]byte(t.tagStart + s.tag + t.tagEnd))
		mustNil(err)
	}
	return nil
}

func mustNil(err error) {
	if err != nil {
		panic(err)
	}
}
