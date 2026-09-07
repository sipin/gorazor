package razorcore

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// compileTemplate compiles body as a template and returns the generated Go
// source, so a test can assert on the escaping the compiler chose.
func compileTemplate(t *testing.T, body string, opt Option) string {
	t.Helper()

	// The generated package name comes from the containing directory, so it
	// has to be a valid Go identifier rather than t.TempDir()'s numeric leaf.
	dir := filepath.Join(t.TempDir(), "tpl")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	in := filepath.Join(dir, "page.gohtml")
	out := filepath.Join(dir, "page.go")
	if err := os.WriteFile(in, []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
	// QuickMode skips the optimizer, which cannot type-check a template
	// compiled outside a module.
	opt.QuickMode = true
	if err := GenFile(context.Background(), in, out, opt); err != nil {
		t.Fatalf("GenFile: %v", err)
	}
	generated, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	return string(generated)
}

const ctxTestDecl = `@{
	var url string
	var name string
	var q string
	var cls string
}
`

// TestContextEscapeCodegen checks that an expression is compiled to the
// escaping function that suits the place it lands in.
func TestContextEscapeCodegen(t *testing.T) {
	cases := []struct {
		name     string
		markup   string
		want     string
		unwanted string
	}{
		{
			name:   "element content",
			markup: `<p>@name</p>`,
			want:   "gorazor.HTMLEscape(name)",
		},
		{
			name:   "ordinary attribute",
			markup: `<div class="@cls">x</div>`,
			want:   "gorazor.HTMLEscape(cls)",
		},
		{
			name:     "href",
			markup:   `<a href="@url">x</a>`,
			want:     "gorazor.URLEscape(url)",
			unwanted: "gorazor.HTMLEscape(url)",
		},
		{
			name:     "img src",
			markup:   `<img src="@url">`,
			want:     "gorazor.URLEscape(url)",
			unwanted: "gorazor.HTMLEscape(url)",
		},
		{
			name:     "form action",
			markup:   `<form action="@url"></form>`,
			want:     "gorazor.URLEscape(url)",
			unwanted: "gorazor.HTMLEscape(url)",
		},
		{
			name:     "query component",
			markup:   `<a href="/search?q=@q">x</a>`,
			want:     "gorazor.URLQueryEscape(q)",
			unwanted: "gorazor.HTMLEscape(q)",
		},
		{
			name:     "script body",
			markup:   `<script>var n = "@name";</script>`,
			want:     "gorazor.JSEscape(name)",
			unwanted: "gorazor.HTMLEscape(name)",
		},
		{
			name:   "alt attribute next to a url one",
			markup: `<img src="@url" alt="@name">`,
			want:   "gorazor.HTMLEscape(name)",
		},
		{
			name:   "style body keeps html escaping",
			markup: `<style>.a { color: @name; }</style>`,
			want:   "gorazor.HTMLEscape(name)",
		},
		{
			name:     "unquoted href",
			markup:   `<a href=@url>x</a>`,
			want:     "gorazor.URLEscape(url)",
			unwanted: "gorazor.HTMLEscape(url)",
		},
		{
			name:     "path segment with colon",
			markup:   `<a href="/items/@name">x</a>`,
			want:     "gorazor.HTMLEscape(name)",
			unwanted: "gorazor.URLEscape(name)",
		},
		{
			name:     "onclick event handler",
			markup:   `<button onclick="handleClick('@name')">click</button>`,
			want:     "gorazor.JSAttrEscape(name)",
			unwanted: "gorazor.HTMLEscape(name)",
		},
		{
			name:     "element following html comment with script",
			markup:   `<!-- <script>alert(1)</script> --><p>@name</p>`,
			want:     "gorazor.HTMLEscape(name)",
			unwanted: "gorazor.JSEscape(name)",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := compileTemplate(t, ctxTestDecl+tc.markup, Option{})
			if !strings.Contains(got, tc.want) {
				t.Errorf("generated code missing %q:\n%s", tc.want, got)
			}
			if tc.unwanted != "" && strings.Contains(got, tc.unwanted) {
				t.Errorf("generated code still contains %q:\n%s", tc.unwanted, got)
			}
		})
	}
}

// TestContextEscapeRecoversAfterScript is the codegen counterpart of the
// scanner test: a "<" used as a comparison inside a script must not leave the
// rest of the document stuck in the script context.
func TestContextEscapeRecoversAfterScript(t *testing.T) {
	body := ctxTestDecl + `<script>
for (var i = 0; i < 10; i++) { if (i < 5) { log("@name"); } }
</script>
<a href="@url">link</a>
<p>@name</p>`

	got := compileTemplate(t, body, Option{})
	for _, want := range []string{
		"gorazor.JSEscape(name)",
		"gorazor.URLEscape(url)",
		"gorazor.HTMLEscape(name)",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("generated code missing %q:\n%s", want, got)
		}
	}
}

// TestContextEscapeMultipleInOneAttribute checks a URL built from two
// expressions: the first starts the URL, the second is a query value.
func TestContextEscapeMultipleInOneAttribute(t *testing.T) {
	got := compileTemplate(t, ctxTestDecl+`<a href="@url?q=@q">x</a>`, Option{})
	if !strings.Contains(got, "gorazor.URLEscape(url)") {
		t.Errorf("URL start not filtered:\n%s", got)
	}
	if !strings.Contains(got, "gorazor.URLQueryEscape(q)") {
		t.Errorf("query value not percent-encoded:\n%s", got)
	}
}

// TestContextEscapeInFlowControl checks that markup inside @if / @for still
// drives the context.
func TestContextEscapeInFlowControl(t *testing.T) {
	body := ctxTestDecl + `@if url != "" {
	<a href="@url">link</a>
}
@for _, s := range []string{"a"} {
	<script>var v = "@s";</script>
}`

	got := compileTemplate(t, body, Option{})
	if !strings.Contains(got, "gorazor.URLEscape(url)") {
		t.Errorf("href inside @if not URL-escaped:\n%s", got)
	}
	if !strings.Contains(got, "gorazor.JSEscape(s)") {
		t.Errorf("script inside @for not JS-escaped:\n%s", got)
	}
}

// TestContextEscapeRawOptsOut checks that @raw() remains the escape hatch: it
// must bypass the URL filter, not be redirected to a different escaper.
func TestContextEscapeRawOptsOut(t *testing.T) {
	got := compileTemplate(t, ctxTestDecl+`<a href="@raw(url)">x</a>`, Option{})
	if strings.Contains(got, "URLEscape") || strings.Contains(got, "HTMLEscape") {
		t.Errorf("@raw() should not be escaped at all:\n%s", got)
	}
	if !strings.Contains(got, "_buffer.WriteString((url))") {
		t.Errorf("expected raw write of url:\n%s", got)
	}
}

// TestDisableContextEscape checks the opt-out flag restores the old behaviour
// for templates that depended on it.
func TestDisableContextEscape(t *testing.T) {
	body := ctxTestDecl + `<a href="@url">x</a><script>var n = "@name";</script>`

	got := compileTemplate(t, body, Option{DisableContextEscape: true})
	if strings.Contains(got, "URLEscape") || strings.Contains(got, "JSEscape") {
		t.Errorf("DisableContextEscape still applied a context escaper:\n%s", got)
	}
	if !strings.Contains(got, "gorazor.HTMLEscape(url)") {
		t.Errorf("expected HTMLEscape for url:\n%s", got)
	}
	if !strings.Contains(got, "gorazor.HTMLEscape(name)") {
		t.Errorf("expected HTMLEscape for name:\n%s", got)
	}
}

// TestContextEscapeOptimizerOptimizesContextEscapers verifies that the optimizer
// rewrites string arguments to context-specific escapers into their zero-boxing
// EscStr variants.
func TestContextEscapeOptimizerOptimizesContextEscapers(t *testing.T) {
	code := `package main

import gorazor "github.com/sipin/gorazor/runtime"

func main() {
	var s string
	_ = gorazor.HTMLEscape(s)
	_ = gorazor.URLEscape(s)
	_ = gorazor.JSEscape(s)
	_ = gorazor.URLQueryEscape(s)
	_ = gorazor.JSAttrEscape(s)
}`

	ok, out := optimize("dummy.go", "main", code)
	if !ok {
		t.Fatalf("expected optimize to succeed, got false")
	}
	for _, want := range []string{"HTMLEscStr", "URLEscStr", "JSEscStr", "URLQueryEscStr", "JSAttrEscStr"} {
		if !strings.Contains(out, want) {
			t.Errorf("optimizer did not rewrite to %s:\n%s", want, out)
		}
	}
	for _, unwant := range []string{"gorazor.HTMLEscape(", "gorazor.URLEscape(", "gorazor.JSEscape(", "gorazor.URLQueryEscape(", "gorazor.JSAttrEscape("} {
		if strings.Contains(out, unwant) {
			t.Errorf("optimizer left unoptimized call %s:\n%s", unwant, out)
		}
	}
}

// TestContextEscapeOptimizerLeavesNonStringContextEscapersAlone guards against
// rewriting calls with non-string arguments (e.g. interface{}).
func TestContextEscapeOptimizerLeavesNonStringContextEscapersAlone(t *testing.T) {
	code := `package main

import gorazor "github.com/sipin/gorazor/runtime"

func main() {
	var v interface{}
	_ = gorazor.URLEscape(v)
	_ = gorazor.JSEscape(v)
	_ = gorazor.URLQueryEscape(v)
	_ = gorazor.JSAttrEscape(v)
}`

	ok, out := optimize("dummy.go", "main", code)
	if ok {
		t.Fatalf("expected optimizer not to rewrite interface{} args, got true")
	}
	for _, want := range []string{"URLEscape", "JSEscape", "URLQueryEscape", "JSAttrEscape"} {
		if !strings.Contains(out, want) {
			t.Errorf("optimizer rewrote %s away for interface{} argument:\n%s", want, out)
		}
	}
}

// TestUnquotedAttrCodegen checks that an attribute written without quotes gets
// the extra encoding wrapped around its context escaper.
func TestUnquotedAttrCodegen(t *testing.T) {
	cases := []struct {
		name   string
		markup string
		want   string
	}{
		{
			name:   "unquoted url attribute",
			markup: `<a href=@url>x</a>`,
			want:   "gorazor.NospaceAttr(gorazor.URLEscape(url))",
		},
		{
			name:   "unquoted ordinary attribute",
			markup: `<div class=@cls>x</div>`,
			want:   "gorazor.NospaceAttr(gorazor.HTMLEscape(cls))",
		},
		{
			name:   "unquoted event handler",
			markup: `<button onclick=@name>x</button>`,
			want:   "gorazor.NospaceAttr(gorazor.JSAttrEscape(name))",
		},
		{
			name:   "quoted attribute is not wrapped",
			markup: `<a href="@url">x</a>`,
			want:   "gorazor.URLEscape(url)",
		},
		{
			name:   "element content is not wrapped",
			markup: `<p>@name</p>`,
			want:   "gorazor.HTMLEscape(name)",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := compileTemplate(t, ctxTestDecl+tc.markup, Option{})
			if !strings.Contains(got, tc.want) {
				t.Errorf("generated code missing %q:\n%s", tc.want, got)
			}
		})
	}
}

// TestUnquotedAttrNotWrappedWhenDisabled checks the opt-out still produces the
// old single-call shape.
func TestUnquotedAttrNotWrappedWhenDisabled(t *testing.T) {
	got := compileTemplate(t, ctxTestDecl+`<a href=@url>x</a>`, Option{DisableContextEscape: true})
	if strings.Contains(got, "NospaceAttr") {
		t.Errorf("DisableContextEscape still wrapped the value:\n%s", got)
	}
	if !strings.Contains(got, "gorazor.HTMLEscape(url)") {
		t.Errorf("expected plain HTMLEscape:\n%s", got)
	}
}

// TestUnquotedAttrRawOptsOut checks @raw() still bypasses everything, including
// the unquoted-attribute wrapper.
func TestUnquotedAttrRawOptsOut(t *testing.T) {
	got := compileTemplate(t, ctxTestDecl+`<a href=@raw(url)>x</a>`, Option{})
	if strings.Contains(got, "NospaceAttr") || strings.Contains(got, "Escape") {
		t.Errorf("@raw() should not be escaped or wrapped:\n%s", got)
	}
}

// TestUnquotedAttrMultiPartExpression checks the closing text is emitted on the
// last child of an expression that lexes into several tokens.
func TestUnquotedAttrMultiPartExpression(t *testing.T) {
	body := `@{
	var u struct{ Name string }
}
<a href=@u.Name>x</a>`

	got := compileTemplate(t, body, Option{})
	if !strings.Contains(got, "gorazor.NospaceAttr(gorazor.URLEscape(u.Name))") {
		t.Errorf("multi-token expression not wrapped correctly:\n%s", got)
	}
}

// TestOptimizerRewritesInnerCall checks the optimizer still reaches the inner
// escaper through the NospaceAttr wrapper.
func TestOptimizerRewritesInnerCall(t *testing.T) {
	code := `package main

import gorazor "github.com/sipin/gorazor/runtime"

func main() {
	var s string
	_ = gorazor.NospaceAttr(gorazor.URLEscape(s))
	_ = gorazor.NospaceAttr(gorazor.HTMLEscape(s))
}`

	ok, out := optimize("dummy.go", "main", code)
	if !ok {
		t.Skip("optimizer could not type-check in this environment")
	}
	if !strings.Contains(out, "gorazor.NospaceAttr(gorazor.URLEscStr(s))") {
		t.Errorf("inner URLEscape not optimized:\n%s", out)
	}
	if !strings.Contains(out, "gorazor.NospaceAttr(gorazor.HTMLEscStr(s))") {
		t.Errorf("inner HTMLEscape not optimized:\n%s", out)
	}
	if strings.Contains(out, "NospaceAttrStr") {
		t.Errorf("optimizer rewrote the wrapper itself:\n%s", out)
	}
}
