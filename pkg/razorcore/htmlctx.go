package razorcore

import "strings"

// htmlContext is the kind of place in an HTML document that an expression is
// emitted into. It selects the escaping function used for that expression.
type htmlContext int

const (
	// ctxText is element content, an ordinary attribute value, or anything
	// else where HTML escaping is the right answer.
	ctxText htmlContext = iota
	// ctxURL is a URL-valued attribute (href, src, ...) at a position where
	// the expression starts the URL, so its scheme must be checked.
	ctxURL
	// ctxURLQuery is a URL-valued attribute past the "?", where the
	// expression is a query component and must be percent-encoded.
	ctxURLQuery
	// ctxScript is the raw-text content of a <script> element.
	ctxScript
	// ctxJSAttr is an inline event handler attribute (onclick, onmouseover, ...)
	// where JavaScript escaping followed by HTML entity escaping is required.
	ctxJSAttr
)

// urlAttrs are the attributes whose value is a single URL.
//
// srcset is deliberately absent: it holds a comma-separated list of URLs plus
// descriptors, so scheme-filtering it as one URL would mangle valid markup.
var urlAttrs = map[string]bool{
	"href":       true,
	"src":        true,
	"action":     true,
	"formaction": true,
	"poster":     true,
	"cite":       true,
	"background": true,
	"manifest":   true,
	"longdesc":   true,
	"profile":    true,
	"codebase":   true,
	"ping":       true,
	"data":       true,
	"xlink:href": true,
}

// scanner states
const (
	stText = iota
	stTagOpen
	stTagName
	stBeforeAttrName
	stAttrName
	stAfterAttrName
	stBeforeAttrValue
	stAttrValueQuoted
	stAttrValueUnquoted
	stRawText // inside <script> or <style>, only looking for the end tag
	stBang
	stDocType
	stCommentDash1
	stComment
	stCommentCloseDash1
	stCommentCloseDash2
)

// htmlScanner tracks just enough HTML structure to answer "what kind of place
// is the template writing into right now?".
//
// It is fed the literal markup of a template in document order and queried
// between feeds, at the point where an expression is emitted. It is not a
// conforming HTML tokenizer and does not need to be: it only has to classify
// positions, and it errs toward ctxText, which keeps today's behaviour.
type htmlScanner struct {
	state    int
	quote    byte   // quote character of the attribute value being scanned
	tagName  string // tag currently being scanned or entered
	attrName string // attribute whose value is being scanned
	rawTag   string // "script" or "style" while inside such an element
	closing  bool   // the tag being scanned is a close tag
	// attrVal is the part of the current attribute value seen so far. Only
	// used to tell a URL's start from its query component.
	attrVal strings.Builder
}

func (s *htmlScanner) reset() {
	*s = htmlScanner{}
}

// clone returns a copy of the scanner. strings.Builder may not be copied once
// used, so the accumulated attribute value is carried over as a string.
func (s *htmlScanner) clone() *htmlScanner {
	c := &htmlScanner{
		state:    s.state,
		quote:    s.quote,
		tagName:  s.tagName,
		attrName: s.attrName,
		rawTag:   s.rawTag,
		closing:  s.closing,
	}
	c.attrVal.WriteString(s.attrVal.String())
	return c
}

func isASCIISpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f'
}

// context reports the kind of position the scanner is currently at.
func (s *htmlScanner) context() htmlContext {
	if s.rawTag == "script" {
		return ctxScript
	}
	// A <style> element is a CSS context. Gorazor has no CSS escaper yet, so
	// it keeps the historical HTML escaping rather than pretending to be safe.
	if s.rawTag != "" {
		return ctxText
	}
	// If the scanner is right after "=" (e.g. <a href=@url> or <button onclick=@handler>),
	// this expression is an unquoted attribute value.
	if s.state == stBeforeAttrValue {
		s.state = stAttrValueUnquoted
	}
	if s.state != stAttrValueQuoted && s.state != stAttrValueUnquoted {
		return ctxText
	}
	if s.attrIsEvent() {
		return ctxJSAttr
	}
	if !s.attrIsURL() {
		return ctxText
	}
	// Past a "?" the expression is a query component, not the URL itself.
	val := s.attrVal.String()
	if strings.ContainsRune(val, '?') {
		return ctxURLQuery
	}
	// In RFC 3986, a scheme cannot contain '/' or '#' and only appears at the
	// very start of a URL. If the attribute value seen so far already contains
	// '/' or '#', the expression is in the path, authority, or fragment
	// component of the URL, not the scheme. HTMLEscape is sufficient and avoids
	// false-positive #ZgotmplZ replacements on path segments containing colons
	// (e.g. <a href="/items/@id"> where id is "doc:123").
	if strings.ContainsRune(val, '/') || strings.ContainsRune(val, '#') {
		return ctxText
	}
	return ctxURL
}

// attrIsEvent reports whether the attribute currently being scanned is an inline
// event handler such as onclick, onmouseover, onerror, etc.
func (s *htmlScanner) attrIsEvent() bool {
	return strings.HasPrefix(s.attrName, "on") && len(s.attrName) > 2
}

// notifyExp informs the scanner that an expression was emitted. If the scanner
// is currently waiting for an attribute value (stBeforeAttrValue), the expression
// serves as the unquoted attribute value.
func (s *htmlScanner) notifyExp() {
	if s.state == stBeforeAttrValue {
		s.state = stAttrValueUnquoted
	}
}

// attrIsURL reports whether the attribute currently being scanned holds a URL.
func (s *htmlScanner) attrIsURL() bool {
	// "data" is a URL only on <object>. Elsewhere it is either a custom
	// attribute or the tail of a data-* name, and filtering it would reject
	// perfectly ordinary values.
	if s.attrName == "data" {
		return s.tagName == "object"
	}
	return urlAttrs[s.attrName]
}

// feed advances the scanner over a run of literal markup.
func (s *htmlScanner) feed(markup string) {
	for i := 0; i < len(markup); i++ {
		s.step(markup[i])
	}
}

func (s *htmlScanner) step(c byte) {
	switch s.state {
	case stText:
		if c == '<' {
			s.state = stTagOpen
			s.tagName = ""
			s.closing = false
		}

	case stRawText:
		// Only "</" + the raw tag's name ends a raw text element, so any
		// other "<" - a JavaScript less-than, say - is just content.
		if c == '<' {
			s.state = stTagOpen
			s.tagName = ""
			s.closing = false
		}

	case stTagOpen:
		switch {
		case c == '/':
			s.closing = true
			s.state = stTagName
		case isTagNameStart(c):
			s.tagName = strings.ToLower(string(c))
			s.state = stTagName
		case c == '!' && s.rawTag == "":
			s.state = stBang
		default:
			// Not a tag after all ("a < b"); fall back to where we were.
			s.state = s.textState()
		}

	case stBang:
		switch c {
		case '-':
			s.state = stCommentDash1
		case '>':
			s.state = stText
		default:
			s.state = stDocType
		}

	case stDocType:
		if c == '>' {
			s.state = stText
		}

	case stCommentDash1:
		switch c {
		case '-':
			s.state = stComment
		case '>':
			s.state = stText
		default:
			s.state = stDocType
		}

	case stComment:
		if c == '-' {
			s.state = stCommentCloseDash1
		}

	case stCommentCloseDash1:
		switch c {
		case '-':
			s.state = stCommentCloseDash2
		default:
			s.state = stComment
		}

	case stCommentCloseDash2:
		switch c {
		case '>':
			s.state = stText
		case '-':
			// stay in stCommentCloseDash2, handles --->
		default:
			s.state = stComment
		}

	case stTagName:
		switch {
		case c == '>':
			s.finishTag()
		case isASCIISpace(c):
			if s.inRawTextMismatch() {
				break
			}
			s.state = stBeforeAttrName
		case c == '/':
			// Keep the name; the '>' will finish the tag.
		default:
			s.tagName += strings.ToLower(string(c))
		}

	case stBeforeAttrName:
		switch {
		case c == '>':
			s.finishTag()
		case isASCIISpace(c) || c == '/':
			// stay
		default:
			s.attrName = strings.ToLower(string(c))
			s.state = stAttrName
		}

	case stAttrName:
		switch {
		case c == '>':
			s.finishTag()
		case c == '=':
			s.state = stBeforeAttrValue
		case isASCIISpace(c):
			s.state = stAfterAttrName
		default:
			s.attrName += strings.ToLower(string(c))
		}

	case stAfterAttrName:
		switch {
		case c == '>':
			s.finishTag()
		case c == '=':
			s.state = stBeforeAttrValue
		case isASCIISpace(c):
			// stay
		default:
			s.attrName = strings.ToLower(string(c))
			s.state = stAttrName
		}

	case stBeforeAttrValue:
		switch {
		case isASCIISpace(c):
			// stay
		case c == '"' || c == '\'':
			s.quote = c
			s.attrVal.Reset()
			s.state = stAttrValueQuoted
		case c == '>':
			s.finishTag()
		default:
			s.quote = 0
			s.attrVal.Reset()
			s.attrVal.WriteByte(c)
			s.state = stAttrValueUnquoted
		}

	case stAttrValueQuoted:
		if c == s.quote {
			s.quote = 0
			s.attrName = ""
			s.state = stBeforeAttrName
		} else {
			s.attrVal.WriteByte(c)
		}

	case stAttrValueUnquoted:
		switch {
		case c == '>':
			s.attrName = ""
			s.finishTag()
		case isASCIISpace(c):
			s.attrName = ""
			s.state = stBeforeAttrName
		default:
			s.attrVal.WriteByte(c)
		}
	}
}

// textState is the state to return to when a '<' turns out not to open a tag.
func (s *htmlScanner) textState() int {
	if s.rawTag != "" {
		return stRawText
	}
	return stText
}

// inRawTextMismatch reports whether the tag being scanned inside a raw text
// element is not that element's close tag, in which case it is ordinary
// content rather than markup.
func (s *htmlScanner) inRawTextMismatch() bool {
	if s.rawTag == "" {
		return false
	}
	if s.closing && s.tagName == s.rawTag {
		return false
	}
	s.state = stRawText
	return true
}

func (s *htmlScanner) finishTag() {
	if s.rawTag != "" {
		if s.closing && s.tagName == s.rawTag {
			s.rawTag = ""
			s.state = stText
		} else {
			// Some other tag written inside <script>; still raw text.
			s.state = stRawText
		}
		s.attrName = ""
		return
	}

	if !s.closing && (s.tagName == "script" || s.tagName == "style") {
		s.rawTag = s.tagName
		s.state = stRawText
	} else {
		s.state = stText
	}
	s.attrName = ""
}

func isTagNameStart(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
