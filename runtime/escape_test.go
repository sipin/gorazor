package runtime

import (
	"strings"
	"testing"
)

func TestURLEscapeSafeSchemes(t *testing.T) {
	cases := []string{
		"http://example.com/a",
		"https://example.com/a",
		"HTTPS://example.com/a",
		"mailto:someone@example.com",
		"tel:+8613800138000",
		"ftp://example.com/f",
		"ftps://example.com/f",
	}
	for _, in := range cases {
		if got := URLEscStr(in); got != in {
			t.Errorf("URLEscStr(%q) = %q, want it kept as-is", in, got)
		}
	}
}

func TestURLEscapeRelative(t *testing.T) {
	cases := []string{
		"/path/to/page",
		"../up",
		"page.html",
		"?q=1",
		"#frag",
		"//cdn.example.com/x.js", // protocol-relative
		"",
	}
	for _, in := range cases {
		if got := URLEscStr(in); got != in {
			t.Errorf("URLEscStr(%q) = %q, want relative URL kept as-is", in, got)
		}
	}
}

// TestURLEscapeBlocksScriptSchemes covers the reason context-aware escaping
// exists: HTML-escaping alone leaves javascript: URLs clickable.
func TestURLEscapeBlocksScriptSchemes(t *testing.T) {
	cases := []string{
		"javascript:alert(1)",
		"JaVaScRiPt:alert(1)",
		"JAVASCRIPT:alert(1)",
		"vbscript:msgbox(1)",
		"data:text/html;base64,PHNjcmlwdD5hbGVydCgxKTwvc2NyaXB0Pg==",
		"file:///etc/passwd",
		"about:blank",
		"blob:https://example.com/uuid",

		// Browsers ignore whitespace and control characters while
		// resolving the scheme, so these all execute if not filtered.
		" javascript:alert(1)",
		"\tjavascript:alert(1)",
		"\njavascript:alert(1)",
		"\rjavascript:alert(1)",
		"java\tscript:alert(1)",
		"java\nscript:alert(1)",
		"java\x00script:alert(1)",
		"jav\x01ascript:alert(1)",
		"\x0cjavascript:alert(1)",
		"javascript\t:alert(1)",

		// A bare leading colon is not a usable scheme, but it is odd
		// enough that it is not worth trusting either.
		":alert(1)",
	}
	for _, in := range cases {
		if got := URLEscStr(in); got != FilteredURL {
			t.Errorf("URLEscStr(%q) = %q, want %q", in, got, FilteredURL)
		}
	}
}

// TestURLEscapeStillHTMLEscapes guards against the filter replacing, rather
// than adding to, the HTML escaping: a quote in a URL must not be able to
// close the attribute it sits in.
func TestURLEscapeStillHTMLEscapes(t *testing.T) {
	cases := map[string]string{
		`http://a.com/?x=1&y=2`:    `http://a.com/?x=1&amp;y=2`,
		`/p"onmouseover="alert(1)`: `/p&#34;onmouseover=&#34;alert(1)`,
		`/p'onmouseover='alert(1)`: `/p&#39;onmouseover=&#39;alert(1)`,
		`/a<b>c`:                   `/a&lt;b&gt;c`,
		`javascript:alert("x")`:    FilteredURL,
	}
	for in, want := range cases {
		if got := URLEscStr(in); got != want {
			t.Errorf("URLEscStr(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestURLQueryEscape(t *testing.T) {
	cases := map[string]string{
		"hello":        "hello",
		"a b":          "a+b",
		"a&admin=1":    "a%26admin%3D1",
		"a=b":          "a%3Db",
		`"><script>`:   "%22%3E%3Cscript%3E",
		"你好":           "%E4%BD%A0%E5%A5%BD",
		"javascript:x": "javascript%3Ax",
	}
	for in, want := range cases {
		if got := URLQueryEscStr(in); got != want {
			t.Errorf("URLQueryEscStr(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestURLQueryEscapeBlocksParamInjection is the point of the query context:
// a value must not be able to add parameters of its own.
func TestURLQueryEscapeBlocksParamInjection(t *testing.T) {
	got := URLQueryEscStr("guest&role=admin")
	if got != "guest%26role%3Dadmin" {
		t.Fatalf("query value injected a parameter: %q", got)
	}
}

func TestJSEscape(t *testing.T) {
	cases := map[string]string{
		"a\"b":       "a\\\"b",
		"a'b":        "a\\'b",
		"x\\y":       "x\\\\y",
		"</script>":  "\\u003C/script\\u003E",
		"a\tb":       "a\\u0009b",
		"a\nb":       "a\\u000Ab",
		"a\rb":       "a\\u000Db",
		"`tmpl`":     "\\u0060tmpl\\u0060",
		"a=b":        "a\\u003Db",
		"a&b":        "a\\u0026b",
		"a\u2028b":   "a\\u2028b",
		"plain text": "plain text",
	}
	for in, want := range cases {
		if got := JSEscStr(in); got != want {
			t.Errorf("JSEscStr(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestJSEscapeBreakout checks the two ways a value inside a <script> string
// literal can escape it: closing the literal, and closing the element.
func TestJSEscapeBreakout(t *testing.T) {
	for _, in := range []string{
		`";alert(1);//`,
		`';alert(1);//`,
		"`;alert(1);//",
		`</script><script>alert(1)</script>`,
		`\";alert(1);//`,
	} {
		got := JSEscStr(in)
		for _, bad := range []string{`"`, `'`, "`", "<", ">"} {
			if containsUnescaped(got, bad) {
				t.Errorf("JSEscStr(%q) = %q leaves an unescaped %s", in, got, bad)
			}
		}
	}
}

// containsUnescaped reports whether s contains tok not preceded by a backslash.
func containsUnescaped(s, tok string) bool {
	for i := 0; i+len(tok) <= len(s); i++ {
		if s[i:i+len(tok)] != tok {
			continue
		}
		if i > 0 && s[i-1] == '\\' {
			continue
		}
		return true
	}
	return false
}

// TestJSEscapePreservesUnicode guards the property that matters for non-ASCII
// content: escaping must not mangle it into \u00XX byte escapes.
func TestJSEscapePreservesUnicode(t *testing.T) {
	for _, in := range []string{"你好，世界", "café", "日本語テキスト", "emoji 🎉"} {
		if got := JSEscStr(in); got != in {
			t.Errorf("JSEscStr(%q) = %q, want it kept as-is", in, got)
		}
	}
}

func TestEscapersAcceptNonStrings(t *testing.T) {
	if got := URLEscape(42); got != "42" {
		t.Errorf("URLEscape(42) = %q, want \"42\"", got)
	}
	if got := JSEscape(42); got != "42" {
		t.Errorf("JSEscape(42) = %q, want \"42\"", got)
	}
	if got := URLQueryEscape(42); got != "42" {
		t.Errorf("URLQueryEscape(42) = %q, want \"42\"", got)
	}
	if got := JSEscape(true); got != "true" {
		t.Errorf("JSEscape(true) = %q, want \"true\"", got)
	}
	// A non-string type must still be escaped, not just formatted.
	if got := JSEscape(struct{ A string }{`"`}); got != `{\"}` {
		t.Errorf("JSEscape(struct) = %q, want %q", got, `{\"}`)
	}
}

func TestToString(t *testing.T) {
	if got := toString("s"); got != "s" {
		t.Errorf("toString(string) = %q", got)
	}
	if got := toString(7); got != "7" {
		t.Errorf("toString(int) = %q", got)
	}
	if got := toString(1.5); got != "1.5" {
		t.Errorf("toString(float) = %q", got)
	}
}

func TestJSAttrEscape(t *testing.T) {
	cases := map[string]string{
		`hello`:              `hello`,
		`"`:                  `\&#34;`,
		`'`:                  `\&#39;`,
		`'); alert(1); //`:   `\&#39;); alert(1); //`,
		`"><script>alert(1)`: `\&#34;\u003E\u003Cscript\u003Ealert(1)`,
		`x\y`:                `x\\y`,
	}
	for in, want := range cases {
		if got := JSAttrEscStr(in); got != want {
			t.Errorf("JSAttrEscStr(%q) = %q, want %q", in, got, want)
		}
	}
	if got := JSAttrEscape(42); got != "42" {
		t.Errorf("JSAttrEscape(42) = %q, want \"42\"", got)
	}
}

func TestNospaceAttr(t *testing.T) {
	cases := map[string]string{
		"plain":            "plain",
		"/path/to/x":       "/path/to/x",
		"a b":              "a&#32;b",
		"a\tb":             "a&#9;b",
		"a\nb":             "a&#10;b",
		"a\fb":             "a&#12;b",
		"a\rb":             "a&#13;b",
		"a`b":              "a&#96;b",
		"a=b":              "a&#61;b",
		"/x onmouseover=y": "/x&#32;onmouseover&#61;y",
		"":                 "",
	}
	for in, want := range cases {
		if got := NospaceAttr(in); got != want {
			t.Errorf("NospaceAttr(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestNospaceAttrDoesNotDoubleEncode checks the property that lets NospaceAttr
// run on an escaper's output: it only touches characters the escapers leave
// alone, so entities they produced pass through untouched.
func TestNospaceAttrDoesNotDoubleEncode(t *testing.T) {
	for _, in := range []string{
		`&amp;`, `&#34;`, `&#39;`, `&lt;`, `&gt;`, `%20`, `a&amp;b`,
	} {
		if got := NospaceAttr(in); got != in {
			t.Errorf("NospaceAttr(%q) = %q, want it unchanged", in, got)
		}
	}
}

// TestNospaceAttrClosesUnquotedBreakout is the reason it exists: escaping for
// the context alone leaves a value able to start an attribute of its own when
// the template wrote the attribute without quotes.
func TestNospaceAttrClosesUnquotedBreakout(t *testing.T) {
	payload := "/x onmouseover=alert(1)"

	if bare := URLEscStr(payload); !strings.Contains(bare, " ") {
		t.Fatalf("setup wrong: URLEscStr already removed the space: %q", bare)
	}

	got := NospaceAttr(URLEscStr(payload))
	if strings.ContainsAny(got, " \t\n\f\r") {
		t.Errorf("value can still end the attribute: %q", got)
	}
	if strings.Contains(got, "=") {
		t.Errorf("value can still start an attribute of its own: %q", got)
	}
}
