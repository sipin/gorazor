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
	if strings.Count(out, "#ZgotmplZ") != 2 {
		t.Errorf("expected both URL attributes to be filtered:\n%s", out)
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
	if strings.Count(out, "张三") != 3 {
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
