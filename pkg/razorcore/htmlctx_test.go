package razorcore

import "testing"

func (c htmlContext) String() string {
	switch c {
	case ctxText:
		return "ctxText"
	case ctxURL:
		return "ctxURL"
	case ctxURLQuery:
		return "ctxURLQuery"
	case ctxScript:
		return "ctxScript"
	}
	return "unknown"
}

// TestHTMLScannerContext feeds the scanner the markup that would precede an
// expression and checks the context it reports at that point.
func TestHTMLScannerContext(t *testing.T) {
	cases := []struct {
		name   string
		markup string
		want   htmlContext
	}{
		// Element content
		{"empty", "", ctxText},
		{"plain text", "<p>hello ", ctxText},
		{"after closed tag", "<div></div>", ctxText},
		{"between tags", "<ul><li>", ctxText},

		// URL attributes
		{"href double quoted", `<a href="`, ctxURL},
		{"href single quoted", `<a href='`, ctxURL},
		{"href unquoted", `<a href=`, ctxText}, // value not started yet
		{"href unquoted started", `<a href=x`, ctxURL},
		{"src", `<img src="`, ctxURL},
		{"uppercase attr", `<a HREF="`, ctxURL},
		{"uppercase tag", `<A href="`, ctxURL},
		{"form action", `<form action="`, ctxURL},
		{"button formaction", `<button formaction="`, ctxURL},
		{"video poster", `<video poster="`, ctxURL},
		{"blockquote cite", `<blockquote cite="`, ctxURL},
		{"script src", `<script src="`, ctxURL},
		{"href after other attrs", `<a class="x" id="y" href="`, ctxURL},
		{"href before other attrs", `<a href="`, ctxURL},
		{"attr with spaces around eq", `<a href = "`, ctxURL},
		{"href with prefix", `<a href="/base/`, ctxURL},

		// Non-URL attributes stay on HTML escaping
		{"class attr", `<div class="`, ctxText},
		{"alt attr", `<img src="x" alt="`, ctxText},
		{"title attr", `<a href="x" title="`, ctxText},
		{"data-* attr", `<div data-url="`, ctxText},
		{"srcset excluded", `<img srcset="`, ctxText},
		{"style attr", `<div style="`, ctxText},

		// "data" is a URL only on <object>
		{"object data", `<object data="`, ctxURL},
		{"div data", `<div data="`, ctxText},

		// Query position inside a URL attribute
		{"query after ?", `<a href="/search?q=`, ctxURLQuery},
		{"query second param", `<a href="/s?a=1&b=`, ctxURLQuery},
		{"fragment is not query", `<a href="/p#`, ctxURL},
		{"question mark outside attr", `<p>what? <a href="`, ctxURL},

		// Script raw text
		{"in script", `<script>var a = "`, ctxScript},
		{"script with attrs", `<script type="text/javascript">x = "`, ctxScript},
		{"script uppercase", `<SCRIPT>var a = "`, ctxScript},
		{"after script closed", `<script>x</script><p>`, ctxText},
		{"after script closed uppercase", `<script>x</SCRIPT><p>`, ctxText},
		{"script closed with space", `<script>x</script ><p>`, ctxText},

		// A "<" inside a script is content, not a tag. These are the cases
		// that a naive tag scanner corrupts.
		{"less than in script", `<script>if (i < 10) { x = "`, ctxScript},
		{"less than no space", `<script>if (i<10) { x = "`, ctxScript},
		{"comparison chain", `<script>if (a<b && c>d) { x = "`, ctxScript},
		{"html string in script", `<script>var h = "<div>"; var y = "`, ctxScript},
		{"arrow function", `<script>const f = (a) => a<1 ? "`, ctxScript},

		// Style is raw text too, but has no CSS escaper yet, so it stays
		// on HTML escaping rather than pretending to be safe.
		{"in style", `<style>.a { color: `, ctxText},
		{"after style closed", `<style>.a{}</style><p>`, ctxText},
		{"href after style", `<style>.a{}</style><a href="`, ctxURL},

		// Self-closing and void elements must not swallow what follows.
		{"after self closing", `<br/><a href="`, ctxURL},
		{"after void element", `<hr><a href="`, ctxURL},
		{"self closing with attr", `<input type="text"/><a href="`, ctxURL},

		// Attribute values must not leak past their closing quote.
		{"after href closed", `<a href="/x">`, ctxText},
		{"after href closed single", `<a href='/x'>`, ctxText},
		{"quote inside other attr", `<a title="a > b" href="`, ctxURL},
		{"gt inside quoted attr", `<a title=">" href="`, ctxURL},
		{"unquoted then href", `<a class=foo href="`, ctxURL},

		// Bare "<" in ordinary text should not start a tag.
		{"lt in text", `<p>a < b and `, ctxText},
		{"lt then href", `<p>a < b</p><a href="`, ctxURL},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := &htmlScanner{}
			s.feed(tc.markup)
			if got := s.context(); got != tc.want {
				t.Errorf("after %q: context = %v, want %v", tc.markup, got, tc.want)
			}
		})
	}
}

// TestHTMLScannerIncremental checks that feeding markup in pieces, as the
// compiler does when expressions interrupt the literal text, gives the same
// answer as feeding it all at once.
func TestHTMLScannerIncremental(t *testing.T) {
	markup := `<div class="x"><a href="/search?q=`
	whole := &htmlScanner{}
	whole.feed(markup)

	piecewise := &htmlScanner{}
	for i := 0; i < len(markup); i++ {
		piecewise.feed(markup[i : i+1])
	}

	if whole.context() != piecewise.context() {
		t.Errorf("piecewise feed gave %v, whole feed gave %v",
			piecewise.context(), whole.context())
	}
	if whole.context() != ctxURLQuery {
		t.Errorf("context = %v, want ctxURLQuery", whole.context())
	}
}

// TestHTMLScannerMultipleExpressions walks a document the way the compiler
// does, checking the context at each point an expression would be emitted.
func TestHTMLScannerMultipleExpressions(t *testing.T) {
	steps := []struct {
		markup string
		want   htmlContext
	}{
		{`<a href="`, ctxURL},                // @url
		{`?ref=`, ctxURLQuery},               // @ref
		{`" title="`, ctxText},               // @title
		{`">`, ctxText},                      // @label
		{`</a><script>var n = "`, ctxScript}, // @name
		{`"; var m = "`, ctxScript},          // @other
		{`";</script><p>`, ctxText},          // @body
	}

	s := &htmlScanner{}
	for i, step := range steps {
		s.feed(step.markup)
		if got := s.context(); got != step.want {
			t.Errorf("step %d (after %q): context = %v, want %v",
				i, step.markup, got, step.want)
		}
	}
}

func TestHTMLScannerClone(t *testing.T) {
	s := &htmlScanner{}
	s.feed(`<a href="/search?q=`)

	c := s.clone()
	if c.context() != ctxURLQuery {
		t.Errorf("clone context = %v, want ctxURLQuery", c.context())
	}

	// Advancing the clone must not disturb the original.
	c.feed(`"><script>`)
	if c.context() != ctxScript {
		t.Errorf("clone advanced to %v, want ctxScript", c.context())
	}
	if s.context() != ctxURLQuery {
		t.Errorf("original changed to %v, want ctxURLQuery", s.context())
	}
}

func TestHTMLScannerReset(t *testing.T) {
	s := &htmlScanner{}
	s.feed(`<script>var a = "`)
	if s.context() != ctxScript {
		t.Fatalf("setup failed: context = %v", s.context())
	}
	s.reset()
	if s.context() != ctxText {
		t.Errorf("after reset: context = %v, want ctxText", s.context())
	}
}
