// Package ctxcases holds a template compiled by gorazor, so that the escaping
// the compiler chose can be checked by actually rendering it.
//
// Regenerate with: gorazor pkg/razorcore/ctxcases/escaping.gohtml pkg/razorcore/ctxcases/escaping.go
package ctxcases

import (
	"strings"
	"testing"
)

// TestRenderedOutputIsSafe renders the template with hostile input and checks
// the resulting document, rather than the generated source. It is the end of
// the chain the other tests cover in pieces.
func TestRenderedOutputIsSafe(t *testing.T) {
	out := Escaping(
		"javascript:alert(1)",
		"guest&role=admin",
		`</script><script>alert(1)</script>`,
	)

	// The javascript: URL must not survive into href or src.
	if strings.Contains(out, "javascript:alert(1)") {
		t.Errorf("javascript: URL reached the output:\n%s", out)
	}
	if strings.Count(out, "#ZgotmplZ") != 3 {
		t.Errorf("expected every URL attribute to be filtered:\n%s", out)
	}

	// The query value must not be able to add a parameter of its own.
	if strings.Contains(out, "q=guest&role=admin") {
		t.Errorf("query value injected a parameter:\n%s", out)
	}
	if !strings.Contains(out, "q=guest%26role%3Dadmin") {
		t.Errorf("query value not percent-encoded:\n%s", out)
	}

	// The interpolated value must not be able to close its string literal or
	// the script element. Only that line is checked: the surrounding literal
	// JavaScript legitimately contains "<" comparisons.
	assign := ""
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "var n = ") {
			assign = line
			break
		}
	}
	if assign == "" {
		t.Fatalf("assignment line missing:\n%s", out)
	}
	want := `var n = "\u003C/script\u003E\u003Cscript\u003Ealert(1)\u003C/script\u003E";`
	if assign != want {
		t.Errorf("value not neutralised in the script context:\n got: %s\nwant: %s", assign, want)
	}
	if strings.ContainsAny(assign[len(`var n = "`):len(assign)-2], `<>`) {
		t.Errorf("interpolated value kept a raw angle bracket: %s", assign)
	}

	// Only one script element: the payload did not open a second one.
	if strings.Count(out, "<script>") != 1 {
		t.Errorf("payload opened another script element:\n%s", out)
	}
}

// TestRenderedOutputKeepsValidInput checks that the escaping is not so eager
// that it breaks ordinary values.
func TestRenderedOutputKeepsValidInput(t *testing.T) {
	out := Escaping("https://example.com/a?b=1", "hello world", "张三")

	if !strings.Contains(out, `href="https://example.com/a?b=1"`) {
		t.Errorf("valid https URL was altered:\n%s", out)
	}
	if !strings.Contains(out, "q=hello+world") {
		t.Errorf("query value not encoded as expected:\n%s", out)
	}
	// Non-ASCII content must survive both the HTML and the JS context.
	if strings.Count(out, "张三") != 4 {
		t.Errorf("non-ASCII text was mangled:\n%s", out)
	}
	if strings.Contains(out, "#ZgotmplZ") {
		t.Errorf("valid URL was filtered:\n%s", out)
	}
}

// TestRenderedScriptStaysIntact checks that literal JavaScript in the template
// is untouched: the "<" comparisons are the case a naive scanner corrupts.
func TestRenderedScriptStaysIntact(t *testing.T) {
	out := Escaping("/x", "y", "z")

	for _, want := range []string{"i < 10", "i < 5", "</script>"} {
		if !strings.Contains(out, want) {
			t.Errorf("script literal %q missing from output:\n%s", want, out)
		}
	}
}

// TestUnquotedAttrCannotBreakOut renders the unquoted attributes with a payload
// that would otherwise end the attribute value and start a handler of its own.
func TestUnquotedAttrCannotBreakOut(t *testing.T) {
	out := Escaping("/x onmouseover=alert(1)", "q", "y onmouseover=alert(2)")

	// Only the unquoted tags matter here: the same payload also lands in a
	// quoted href and in element content, where quoting and HTML escaping
	// already contain it.
	for _, val := range unquotedValues(t, out) {
		if strings.ContainsAny(val, " \t\n\f\r") {
			t.Errorf("value ended the attribute and started another: %q", val)
		}
		if strings.Contains(val, "=") {
			t.Errorf("value introduced a second attribute: %q", val)
		}
	}

	// The value survives, encoded, as a single attribute value.
	if !strings.Contains(out, "href=/x&#32;onmouseover&#61;alert(1)>") {
		t.Errorf("unquoted href not encoded as expected:\n%s", out)
	}
	if !strings.Contains(out, "class=y&#32;onmouseover&#61;alert(2)>") {
		t.Errorf("unquoted class not encoded as expected:\n%s", out)
	}
}

// TestUnquotedAttrStillFiltersScheme checks the wrapper did not displace the
// scheme filtering underneath it.
func TestUnquotedAttrStillFiltersScheme(t *testing.T) {
	out := Escaping("javascript:alert(1)", "q", "n")
	if strings.Contains(out, "javascript:alert(1)") {
		t.Errorf("javascript: URL survived in an unquoted attribute:\n%s", out)
	}
	if strings.Count(out, "#ZgotmplZ") != 3 {
		t.Errorf("expected all three URL attributes filtered:\n%s", out)
	}
}

// TestUnquotedAttrLeavesOrdinaryValuesAlone checks the encoding does not fire
// on values that need nothing.
func TestUnquotedAttrLeavesOrdinaryValuesAlone(t *testing.T) {
	out := Escaping("/page/1", "q", "btn")
	if !strings.Contains(out, "href=/page/1>") {
		t.Errorf("ordinary unquoted href was altered:\n%s", out)
	}
	if !strings.Contains(out, "class=btn>") {
		t.Errorf("ordinary unquoted class was altered:\n%s", out)
	}
	if strings.Contains(out, "&#32;") {
		t.Errorf("nothing should have been encoded:\n%s", out)
	}
}

// unquotedValues returns the values of the two attributes the template writes
// without quotes around them.
func unquotedValues(t *testing.T, out string) []string {
	t.Helper()

	var values []string
	for _, prefix := range []string{"<a href=", "<div class="} {
		found := false
		for i := 0; ; {
			j := strings.Index(out[i:], prefix)
			if j < 0 {
				break
			}
			at := i + j + len(prefix)
			i = at
			// Skip the quoted attribute that shares this prefix.
			if at < len(out) && (out[at] == '"' || out[at] == '\'') {
				continue
			}
			end := strings.IndexByte(out[at:], '>')
			if end < 0 {
				t.Fatalf("unterminated tag after %q:\n%s", prefix, out)
			}
			values = append(values, out[at:at+end])
			found = true
			break
		}
		if !found {
			t.Fatalf("no unquoted %q found in output:\n%s", prefix, out)
		}
	}
	return values
}
