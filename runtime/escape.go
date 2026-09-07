package runtime

import (
	"fmt"
	"html/template"
	"net/url"
	"strconv"
	"strings"
)

// FilteredURL replaces a URL whose scheme is not in safeURLSchemes.
// It mirrors html/template's "#ZgotmplZ" failsafe: it is deliberately
// harmless and greppable, so an unexpected value shows up in the output
// instead of silently becoming an executable URL.
const FilteredURL = "#ZgotmplZ"

// safeURLSchemes lists the URL schemes that may appear in a URL-valued
// attribute such as href or src. Everything else - most importantly
// javascript:, vbscript: and data: - is replaced by FilteredURL.
//
// This is an allow list on purpose. A deny list is not defensible here:
// browsers accept a wide and quirky range of scheme spellings, so anything
// not positively known to be inert must be rejected.
var safeURLSchemes = map[string]bool{
	"http":   true,
	"https":  true,
	"mailto": true,
	"tel":    true,
	"ftp":    true,
	"ftps":   true,
}

func toString(m interface{}) string {
	switch v := m.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	}
	return fmt.Sprint(m)
}

// filterURL returns s unchanged when it is a relative URL or carries a scheme
// from safeURLSchemes, and FilteredURL otherwise.
//
// When looking for the scheme, ASCII whitespace and control characters are
// skipped, because browsers ignore them while resolving a scheme: "java\tscript:"
// and " javascript:" both execute. Detecting the scheme on the raw string only
// would let either through.
func filterURL(s string) string {
	var scheme strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == ':':
			if scheme.Len() == 0 {
				// A leading ":" is not a scheme, but it is odd enough
				// that it is not worth trusting.
				return FilteredURL
			}
			if safeURLSchemes[strings.ToLower(scheme.String())] {
				return s
			}
			return FilteredURL
		case c == '/' || c == '?' || c == '#':
			// Reached the path, query or fragment before any ":", so
			// this is a relative URL and there is no scheme to check.
			return s
		case c <= 0x20 || c == 0x7f:
			// Ignored by browsers when resolving the scheme.
			continue
		default:
			scheme.WriteByte(c)
		}
	}
	// No ":" at all - a relative URL.
	return s
}

// URLEscape escapes m for use as a whole URL in an attribute such as
// href or src. Unsafe schemes are replaced by FilteredURL; use @raw() to
// opt out of filtering.
func URLEscape(m interface{}) string {
	return URLEscStr(toString(m))
}

// URLEscStr is URLEscape for a value already known to be a string.
func URLEscStr(s string) string {
	return template.HTMLEscapeString(filterURL(s))
}

// URLQueryEscape escapes m for use inside the query string of a URL
// attribute, as in href="/search?q=@query". The value is percent-encoded so
// that it cannot inject extra query parameters.
func URLQueryEscape(m interface{}) string {
	return URLQueryEscStr(toString(m))
}

// URLQueryEscStr is URLQueryEscape for a value already known to be a string.
func URLQueryEscStr(s string) string {
	return template.HTMLEscapeString(url.QueryEscape(s))
}

// JSEscape escapes m for use inside a JavaScript string literal in a
// <script> block.
//
// HTMLEscape must not be used here: a <script> element holds raw text, so a
// browser never decodes HTML entities inside it and an escaped quote would
// reach the script as the literal text "&#34;".
func JSEscape(m interface{}) string {
	return JSEscStr(toString(m))
}

// JSEscStr is JSEscape for a value already known to be a string.
func JSEscStr(s string) string {
	// template.JSEscapeString covers quotes, backslashes, the characters
	// that could close the script element, and the line terminators that
	// would break a string literal. It leaves the backtick alone, which
	// matters inside a template literal, so escape that separately.
	//
	// The order is safe: JSEscapeString has already doubled every existing
	// backslash, so the backslash introduced below can only be read as the
	// start of this new escape sequence.
	return strings.ReplaceAll(template.JSEscapeString(s), "`", "\\u0060")
}

// JSAttrEscape escapes m for use inside an inline event handler attribute
// (such as onclick="handleClick('@name')"). It applies JavaScript string
// escaping followed by HTML entity escaping so that neither the JS string literal
// nor the HTML attribute delimiter can be broken out of.
func JSAttrEscape(m interface{}) string {
	return JSAttrEscStr(toString(m))
}

// JSAttrEscStr is JSAttrEscape for a value already known to be a string.
func JSAttrEscStr(s string) string {
	return template.HTMLEscapeString(JSEscStr(s))
}
