package razorcore

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var testdata = "testdata"

func TestCap(t *testing.T) {
	if Capitalize("") != "" {
		t.Error()
	}
	if Capitalize("hello") != "Hello" {
		t.Error()
	}
	if Capitalize("0hello") != "0hello" {
		t.Error()
	}
}

func TestWidget(t *testing.T) {
	w := Widget{
		Label:       "Username",
		Value:       "test",
		Name:        "username",
		PlaceHolder: "Enter username",
		Type:        "text",
		ErrorMsg:    "invalid",
	}
	if w.Label != "Username" || w.Value != "test" || w.Name != "username" {
		t.Errorf("unexpected widget values: %+v", w)
	}
}

func TestLayoutCache(t *testing.T) {
	cache := NewLayoutCache()
	cache.Set("hello", []string{"this", "is", "good"})
	cache.Set("world", []string{"funny"})
	if len(cache.Get("hello")) != 3 {
		t.Error()
	}
	if len(cache.Get("world")) != 1 {
		t.Error()
	}
	if len(cache.Get("NO")) != 0 {
		t.Error()
	}
	// test Get/Set with nil cache
	var nilCache *LayoutCache
	if len(nilCache.Get("hello")) != 0 {
		t.Error("expected 0 args for nil cache")
	}
	nilCache.Set("hello", []string{"test"}) // should not panic
}

func TestLexer1(t *testing.T) {
	text := "case do func var switch"
	lex := &Lexer{text, Tests}
	res, err := lex.Scan()
	if err != nil {
		t.Error(err)
	}
	if len(res) != 10 {
		t.Error("token number")
	}
	for i, x := range res {
		if i%2 == 0 && x.Type != tkKeyword {
			t.Error("KEYWORD", x)
		}
	}
}
func TestLexer2(t *testing.T) {
	text := "case casex do do3 func func_ var var+ "
	lex := &Lexer{text, Tests}
	res, err := lex.Scan()
	if err != nil {
		t.Error(err)
	}
	if len(res) != 18 {
		t.Error(err)
	}
	for i, x := range res {
		if i == 0 || i == 4 || i == 8 || i == 12 || i == 14 {
			if x.Type != tkKeyword {
				t.Error("KEYWORD")
			}
		} else if x.Type == tkKeyword {
			t.Error("Should NOT KEYWORD", x)
		}
	}
}

func TestDebug(t *testing.T) {
	casedir, _ := filepath.Abs(filepath.Dir("./cases/"))
	outdir, _ := filepath.Abs(filepath.Dir("./" + testdata + "/"))
	option := Option{}
	option.IsDebug = true
	GenFile(context.Background(), casedir+"/var.gohtml", outdir+"/_var.gohtml", option)
}

func TestGenerate(t *testing.T) {
	casedir, _ := filepath.Abs(filepath.Dir("./cases/"))
	sap := string(filepath.Separator)

	visit := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() { // regular file
			if strings.HasPrefix(filepath.Base(path), ".") {
				return nil
			}
			name := strings.Replace(path, ".gohtml", ".go", 1)
			cmp := strings.ReplaceAll(name, sap+"cases"+sap, sap+testdata+sap)
			dirname := filepath.Dir(cmp)
			log := filepath.Join(dirname, "_"+filepath.Base(cmp))
			if !exists(dirname) {
				os.MkdirAll(dirname, 0755)
			}
			option := Option{QuickMode: true}

			if strings.HasSuffix(path, "panic.gohtml") {
				defer func() {
					if r := recover(); r != nil {
						panicMsg := fmt.Sprint(r)
						if !strings.HasPrefix(panicMsg, "failed to format template") ||
							!strings.Contains(panicMsg, ">>>> func Panic(totalMessage i nt) string {") {
							t.Error("panic.gohtml test failed")
						}
					}
				}()
				GenFile(context.Background(), path, log, option)
			} else {
				GenFile(context.Background(), path, log, option)
			}

			if !exists(cmp) {
				t.Error("No cmp:", cmp)
			} else if !exists(log) {
				t.Error("No log:", log)
			} else {
				//compare the log file and cmp file
				_cmp, _e1 := os.ReadFile(cmp)
				_log, _e2 := os.ReadFile(log)
				if _e1 != nil || _e2 != nil {
					t.Error("Reading")
				} else if string(_cmp) != string(_log) {
					t.Error("MISMATCH:", log, cmp)
				} else {
					t.Log("PASS")
				}
			}
		}
		return nil
	}
	err := filepath.Walk(casedir, visit)
	if err != nil {
		t.Error("walk")
	}
}

func TestAdditionalCoverage(t *testing.T) {
	t.Run("exists_non_existent", func(t *testing.T) {
		if exists("non_existent_file_xyz_123") {
			t.Error("expected false")
		}
	})

	t.Run("getValStr_invalid_type", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic")
			}
		}()
		getValStr(123)
	})

	t.Run("getValStr_ast_pointer", func(t *testing.T) {
		ast := &Ast{TagName: "my-tag"}
		if getValStr(ast) != "my-tag" {
			t.Error("expected my-tag")
		}
	})

	t.Run("ast_modestr_default", func(t *testing.T) {
		badAst := &Ast{Mode: 999}
		if badAst.ModeStr() != "UNDEF" {
			t.Error("expected UNDEF")
		}
	})

	t.Run("ast_check_panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				msg := fmt.Sprint(r)
				if msg != "Maximum number of elements exceeded." {
					t.Error("unexpected panic:", msg)
				}
			} else {
				t.Error("expected panic")
			}
		}()
		panicAst := &Ast{}
		panicAst.Children = make([]interface{}, 100000)
		panicAst.check()
	})

	t.Run("ast_pop_empty", func(t *testing.T) {
		emptyAst := &Ast{}
		emptyAst.popChild() // should not panic
		if len(emptyAst.Children) != 0 {
			t.Error("expected 0 children")
		}
	})

	t.Run("ast_debug_and_has_non_exp", func(t *testing.T) {
		root := &Ast{Mode: PRG, TagName: "root"}
		child1 := &Ast{Mode: EXP, TagName: "exp_child"}
		child2 := &Ast{Mode: MKP, TagName: "markup_child"}
		root.addChild(child1)
		root.addChild(child2)
		root.debug(0, 5) // cover debug output
		root.debug(10, 5) // cover depth limit check

		if !root.hasNonExp() {
			t.Error("expected hasNonExp to be true")
		}
		if child1.hasNonExp() {
			t.Error("expected child1 (EXP with no children) hasNonExp to be false")
		}
	})

	t.Run("ast_has_non_exp_recursive_true", func(t *testing.T) {
		root := &Ast{Mode: EXP}
		child := &Ast{Mode: EXP}
		leaf := &Ast{Mode: MKP}
		child.addChild(leaf)
		root.addChild(child)
		if !root.hasNonExp() {
			t.Error("expected true")
		}
	})

	t.Run("ast_has_non_exp_multiple_children_second_match", func(t *testing.T) {
		root := &Ast{Mode: EXP}
		child1 := &Ast{Mode: EXP} // child1 has no non-exp
		child2 := &Ast{Mode: EXP}
		leaf := &Ast{Mode: MKP} // child2 has non-exp
		child2.addChild(leaf)
		root.addChild(child1)
		root.addChild(child2)
		if !root.hasNonExp() {
			t.Error("expected true when second child has non-exp")
		}
	})

	t.Run("isLayoutSection_direct", func(t *testing.T) {
		cp := &Compiler{
			isLayout:   true,
			paramNames: []string{"sidebar", "js"},
		}
		// test falls through on isLayoutSectionPart
		is, _ := cp.isLayoutSectionPart(Part{ptype: EXP, value: "not_a_write_string"})
		if is {
			t.Error("expected false")
		}

		// test match != in isLayoutSectionTest
		is, val := cp.isLayoutSectionTest(Part{ptype: CBLK, value: `if js != "" {`})
		if !is || val != "if js != nil {\n" {
			t.Error("expected true with match, got:", is, val)
		}
	})

	t.Run("settleLayout_prefix_adjust", func(t *testing.T) {
		defer func() {
			recover() // we expect a panic because the layout file still doesn't exist
		}()
		cp := &Compiler{
			layout: "myprefix/layout",
			file:   "somefile",
			options: Option{
				TemplateNamespacePrefix: "myprefix",
			},
		}
		cp.settleLayout("NonExistentLayout")
	})

	t.Run("settleLayout_compilation_error", func(t *testing.T) {
		tmpLayout := filepath.Join(os.TempDir(), "temp_layout_error.gohtml")
		// Create a directory at tmpLayout so exists() is true but ReadFile fails
		err := os.MkdirAll(tmpLayout, 0755)
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpLayout)

		cp := &Compiler{
			layout: filepath.Dir(tmpLayout),
			file:   "somefile",
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic due to compilation error in settleLayout")
			}
		}()

		cp.settleLayout(strings.TrimSuffix(filepath.Base(tmpLayout), ".gohtml"))
	})

	t.Run("run_non_existent", func(t *testing.T) {
		_, err := run(context.Background(), "non_existent_file_xyz_123", Option{})
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("parser_skip_token", func(t *testing.T) {
		parser := &Parser{
			tokens: []Token{
				{Type: tkContent, Text: "hello"},
				{Type: tkContent, Text: "world"},
			},
		}
		parser.skipToken()
		if parser.peekToken(0).Text != "world" {
			t.Error("expected next token to be world")
		}
	})

	t.Run("lexer_illegal_character", func(t *testing.T) {
		lex := &Lexer{"abc", []TokenMatch{}}
		_, err := lex.Scan()
		if err == nil || !strings.Contains(err.Error(), "Illegal character: a") {
			t.Error("expected Illegal character: a, got:", err)
		}

		// Test multiline illegal character reporting: line 2 column 3 should report 'a', not text[pos] which would be ' '
		lex2 := &Lexer{"(\n   abc", []TokenMatch{}}
		_, err2 := lex2.Scan()
		if err2 == nil || !strings.Contains(err2.Error(), "2:3: Illegal character: a") {
			t.Error("expected '2:3: Illegal character: a', got:", err2)
		}
	})

	t.Run("optimize_errors", func(t *testing.T) {
		// Parsing error
		ok, _ := optimize("dummy.go", "main", "package main\ninvalid syntax")
		if ok {
			t.Error("expected optimize to return false on syntax error")
		}

		// Type check error
		ok, _ = optimize("dummy.go", "main", "package main\nfunc main() { x := undefinedVar }")
		if ok {
			t.Error("expected optimize to return false on type checking error")
		}

		// Non-ident selector call (e.g. getP().HTMLEscape()) should not panic
		codeNonIdent := `package main
type P struct{}
func (p P) HTMLEscape(s string) string { return s }
func getP() P { return P{} }
func main() { _ = getP().HTMLEscape("test") }`
		ok, _ = optimize("dummy.go", "main", codeNonIdent)
		if ok {
			t.Error("expected optimize to return false for non-gorazor selector")
		}
	})

	t.Run("layout_args_zero_sections", func(t *testing.T) {
		cache := NewLayoutCache()
		cache.Set("tpl/layout/zero_args", []string{})
		cp := &Compiler{
			layout: "tpl/layout/zero_args",
			options: Option{
				LayoutCache: cache,
			},
		}
		foot := cp.generateFoot([]string{"mysection"})
		if !strings.Contains(foot, ", mysection()") {
			t.Error("expected section invocation in foot, got:", foot)
		}
	})

	t.Run("layout_args_with_string_substring", func(t *testing.T) {
		cache := NewLayoutCache()
		cache.Set("tpl/layout/base", []string{"body string", "queryString string", "stringVal string"})
		cp := &Compiler{
			layout: "tpl/layout/base",
			options: Option{
				LayoutCache: cache,
			},
		}
		foot := cp.generateFoot([]string{"queryString", "stringVal"})
		if !strings.Contains(foot, ", _queryString, _stringVal)") {
			t.Errorf("expected matching for queryString and stringVal, got: %s", foot)
		}
	})

	t.Run("layout_import_alias", func(t *testing.T) {
		tmpDir := filepath.Join(t.TempDir(), "tpl")
		os.MkdirAll(tmpDir, 0755)
		tplFile := filepath.Join(tmpDir, "alias.gohtml")
		outFile := filepath.Join(tmpDir, "alias.go")
		content := `@{
	import (
		share "cases/layout"
	)
	layout := share.base
}
<h1>Hello with alias layout</h1>`
		if err := os.WriteFile(tplFile, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		cache := NewLayoutCache()
		cache.Set("cases/layout/base", []string{"body string", "title string", "js string"})
		err := GenFile(context.Background(), tplFile, outFile, Option{
			LayoutCache: cache,
		})
		if err != nil {
			t.Fatalf("failed to compile template with layout alias: %v", err)
		}
		outBytes, err := os.ReadFile(outFile)
		if err != nil {
			t.Fatal(err)
		}
		outStr := string(outBytes)
		if !strings.Contains(outStr, `share.RenderBase(_buffer, _body`) {
			t.Errorf("expected share.RenderBase in generated code, got:\n%s", outStr)
		}
		if !strings.Contains(outStr, `share "cases/layout"`) {
			t.Errorf("expected share \"cases/layout\" import in generated code, got:\n%s", outStr)
		}
	})

	t.Run("genfolder_errors", func(t *testing.T) {
		err := GenFolder(context.Background(), "non_existent_folder_xyz_123", "out", Option{})
		if err == nil || !strings.Contains(err.Error(), "input directory does not exsit") {
			t.Error("expected folder not exist error, got:", err)
		}

		// GenFolder with non-existent outdir
		baseDir, _ := filepath.Abs(".")
		casedir := filepath.Join(baseDir, "cases")
		tmpOut := filepath.Join(os.TempDir(), "gorazor_test_out_dir_xyz_123")
		defer os.RemoveAll(tmpOut)

		// This should cover outdir creation in GenFolder
		err = GenFolder(context.Background(), filepath.Join(casedir, "layout"), tmpOut, Option{})
		if err != nil {
			t.Error("expected no error in GenFolder, got:", err)
		}
	})

	t.Run("noline_option", func(t *testing.T) {
		tmpOut := filepath.Join(os.TempDir(), "gorazor_noline_out.go")
		defer os.Remove(tmpOut)

		tmpGohtml := filepath.Join(os.TempDir(), "temp_noline_tpl.gohtml")
		defer os.Remove(tmpGohtml)

		err := os.WriteFile(tmpGohtml, []byte(`@{
	var totalMessage int
}
<div>Hello @totalMessage</div>
`), 0644)
		if err != nil {
			t.Fatal(err)
		}

		option := Option{NoLineNumber: true}
		err = GenFile(context.Background(), tmpGohtml, tmpOut, option)
		if err != nil {
			t.Error("expected no error, got:", err)
		}

		content, err := os.ReadFile(tmpOut)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(content), "// [") {
			t.Error("expected no line hint comments in generated file, but found them")
		}
	})

	t.Run("layout_not_found_error", func(t *testing.T) {
		tmpOut := filepath.Join(os.TempDir(), "gorazor_layout_error_out.go")
		defer os.Remove(tmpOut)

		// Create a temporary layout file that references a non-existent layout
		tmpGohtml := filepath.Join(os.TempDir(), "temp_layout_xyz_123.gohtml")
		defer os.Remove(tmpGohtml)

		err := os.WriteFile(tmpGohtml, []byte(`@{
	import (
		"tpl/layout"
	)
	layout := layout.NonExistentLayoutNameXYZ
}
<div>Hello</div>
`), 0644)
		if err != nil {
			t.Fatal(err)
		}

		option := Option{}
		// This should return an error due to missing layout
		err = GenFile(context.Background(), tmpGohtml, tmpOut, option)
		if err == nil {
			t.Error("expected error, got nil")
		} else if !strings.Contains(err.Error(), "Can't find layout:") {
			t.Error("unexpected error message:", err)
		}
	})

	t.Run("genfolder_comprehensive", func(t *testing.T) {
		tmpDir := filepath.Join(os.TempDir(), "gorazor_comprehensive_test_dir_xyz")
		defer os.RemoveAll(tmpDir)

		inDir := filepath.Join(tmpDir, "in")
		outDir := filepath.Join(tmpDir, "out")
		os.MkdirAll(inDir, 0755)

		// 1. Create a valid .gohtml file
		os.WriteFile(filepath.Join(inDir, "valid.gohtml"), []byte("<div>Hello</div>"), 0644)
		// 2. Create a non-.gohtml file (to cover !strings.HasSuffix branch)
		os.WriteFile(filepath.Join(inDir, "dummy.txt"), []byte("ignore me"), 0644)
		// 3. Create a Emacs lock file (to cover strings.HasPrefix(\".#\") branch)
		os.WriteFile(filepath.Join(inDir, ".#lock.gohtml"), []byte("lock file"), 0644)

		// Call GenFolder on the newly created directories (this also covers outDir not existing, i.e. !exists(outDir) in GenFolder)
		err := GenFolder(context.Background(), inDir, outDir, Option{})
		if err != nil {
			t.Fatal("GenFolder failed:", err)
		}

		// Verify files generated
		if !exists(filepath.Join(outDir, "valid.go")) {
			t.Error("expected valid.go to be generated")
		}
		if exists(filepath.Join(outDir, "dummy.go")) || exists(filepath.Join(outDir, ".#lock.go")) {
			t.Error("expected only valid.go to be generated")
		}

		// 4. Cover GenFile's output directory creation
		nonExistentSubDir := filepath.Join(tmpDir, "non_existent_subdir", "output.go")
		// this will call GenFile and hit the !exists(outdir) branch inside GenFile!
		err = GenFile(context.Background(), filepath.Join(inDir, "valid.gohtml"), nonExistentSubDir, Option{})
		if err != nil {
			t.Fatal("GenFile failed to create outdir:", err)
		}
	})

	t.Run("parser_prev_token_nil", func(t *testing.T) {
		parser := &Parser{
			preTokens: []Token{},
		}
		if parser.prevToken(0) != nil {
			t.Error("expected nil")
		}
	})

	t.Run("regMatch_invalid_regex", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic")
			}
		}()
		regMatch("[invalid", "text")
	})

	t.Run("regMatch_no_match", func(t *testing.T) {
		res, err := regMatch("abc", "def")
		if err != nil || res != "" {
			t.Error("expected empty string and no error")
		}
	})

	t.Run("parser_advance_until_escaped_end", func(t *testing.T) {
		parser := &Parser{
			tokens: []Token{
				{Type: tkAt, Text: "@"},
				{Type: tkParenClose, Text: ")"},
				{Type: tkParenClose, Text: ")"},
			},
		}
		tok := Token{Type: tkParenOpen, Text: "("}
		res, err := parser.advanceUntil(tok, tkParenOpen, tkParenClose, -1, tkAt)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 4 {
			t.Error("expected 4 tokens, got:", len(res))
		}
	})

	t.Run("parser_advance_until_unmatched", func(t *testing.T) {
		parser := &Parser{
			tokens: []Token{},
		}
		tok := Token{Type: tkParenOpen, Text: "("}
		_, err := parser.advanceUntil(tok, tkParenOpen, tkParenClose, -1, tkAt)
		if err == nil || !strings.Contains(err.Error(), "Unmatched tag close") {
			t.Error("expected unmatched tag close error, got:", err)
		}
	})

	t.Run("parser_subparse_unmatched", func(t *testing.T) {
		parser := &Parser{
			tokens: []Token{},
		}
		tok := Token{Type: tkParenOpen, Text: "("}
		err := parser.subParse(tok, EXP, true)
		if err == nil || !strings.Contains(err.Error(), "Unmatched tag close") {
			t.Error("expected error, got:", err)
		}
	})

	t.Run("process_imports_error", func(t *testing.T) {
		tmpOut := filepath.Join(os.TempDir(), "gorazor_import_err_out.go")
		defer os.Remove(tmpOut)

		tmpGohtml := filepath.Join(os.TempDir(), "temp_import_err_xyz.gohtml")
		defer os.Remove(tmpGohtml)

		// Create a file with syntax error in the import statement
		err := os.WriteFile(tmpGohtml, []byte(`@{
	import = 123
}
<div>Hello</div>
`), 0644)
		if err != nil {
			t.Fatal(err)
		}

		err = GenFile(context.Background(), tmpGohtml, tmpOut, Option{})
		if err == nil {
			t.Error("expected error due to invalid import block, got nil")
		} else {
			errMsg := err.Error()
			if !strings.Contains(errMsg, "compilation error in") {
				t.Error("expected diagnostic prefix, got:", errMsg)
			}
			if !strings.Contains(errMsg, "temp_import_err_xyz.gohtml:2") {
				t.Error("expected filename and line 2, got:", errMsg)
			}
			if !strings.Contains(errMsg, "Context:") || !strings.Contains(errMsg, "---> 2: import = 123") {
				t.Error("expected context snippet highlighting line 2, got:", errMsg)
			}
		}
	})

	t.Run("generate_foot_no_layout_sections_error", func(t *testing.T) {
		tmpOut := filepath.Join(os.TempDir(), "gorazor_section_err_out.go")
		defer os.Remove(tmpOut)

		tmpGohtml := filepath.Join(os.TempDir(), "temp_section_err_xyz.gohtml")
		defer os.Remove(tmpGohtml)

		// Create a file with a section but no layout declaration
		err := os.WriteFile(tmpGohtml, []byte(`@{
	// No layout assigned
}
@section Header {
	<div>Header</div>
}
`), 0644)
		if err != nil {
			t.Fatal(err)
		}

		err = GenFile(context.Background(), tmpGohtml, tmpOut, Option{})
		if err == nil {
			t.Error("expected error due to sections with no layout, got nil")
		} else {
			errMsg := err.Error()
			if !strings.Contains(errMsg, "compilation error in") {
				t.Error("expected diagnostic prefix, got:", errMsg)
			}
			if !strings.Contains(errMsg, "temp_section_err_xyz.gohtml:1") {
				t.Error("expected filename and line 1, got:", errMsg)
			}
			if !strings.Contains(errMsg, "Context:") || !strings.Contains(errMsg, "---> 1: @{") {
				t.Error("expected context snippet highlighting line 1, got:", errMsg)
			}
		}
	})

	t.Run("genfolder_with_errors", func(t *testing.T) {
		tmpDir := filepath.Join(os.TempDir(), "gorazor_folder_err_test_xyz")
		defer os.RemoveAll(tmpDir)

		inDir := filepath.Join(tmpDir, "in")
		outDir := filepath.Join(tmpDir, "out")
		os.MkdirAll(inDir, 0755)

		// Create a file with syntax error
		os.WriteFile(filepath.Join(inDir, "error.gohtml"), []byte(`@{
	import = 123
}`), 0644)

		err := GenFolder(context.Background(), inDir, outDir, Option{})
		if err == nil {
			t.Error("expected error from GenFolder due to invalid syntax in file, got nil")
		}
	})

	t.Run("genfolder_context_cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // immediately cancel the context

		tmpDir := filepath.Join(os.TempDir(), "gorazor_cancel_test_dir_xyz")
		defer os.RemoveAll(tmpDir)

		inDir := filepath.Join(tmpDir, "in")
		outDir := filepath.Join(tmpDir, "out")
		os.MkdirAll(inDir, 0755)

		os.WriteFile(filepath.Join(inDir, "valid.gohtml"), []byte("<div>Hello</div>"), 0644)

		err := GenFolder(ctx, inDir, outDir, Option{})
		if err == nil {
			t.Error("expected error due to cancelled context, got nil")
		} else if !strings.Contains(err.Error(), "context canceled") {
			t.Error("unexpected error:", err)
		}
	})

	t.Run("genfolder_panic_recovery", func(t *testing.T) {
		tmpDir := t.TempDir()
		inDir := filepath.Join(tmpDir, "in")
		outDir := filepath.Join(tmpDir, "out")
		os.MkdirAll(inDir, 0755)

		// Create a template that will cause FormatBuffer to panic due to invalid Go syntax
		panicContent := `@{
	var totalMessage i nt
}`
		os.WriteFile(filepath.Join(inDir, "panic.gohtml"), []byte(panicContent), 0644)

		err := GenFolder(context.Background(), inDir, outDir, Option{QuickMode: true})
		if err == nil {
			t.Error("expected error from GenFolder due to panic in file, got nil")
		}
		if !strings.Contains(err.Error(), "error processing") {
			t.Errorf("expected error processing prefix, got: %v", err)
		}
	})

	t.Run("html_compact_mode", func(t *testing.T) {
		tmpOut := filepath.Join(os.TempDir(), "gorazor_compact_out.go")
		defer os.Remove(tmpOut)

		tmpGohtml := filepath.Join(os.TempDir(), "temp_compact_xyz.gohtml")
		defer os.Remove(tmpGohtml)

		// Create a file with complex HTML including pre, script, style, and normal spaces/newlines
		rawHTML := `
		<div   class="header"   id="top" >
			<h1>Hello  World!</h1>
		</div>
		
		<pre  class="code-block" >
			line 1
			  line 2
		</pre>

		<script>
			console.log( "test" );
		</script>

		<style>
			body {  margin: 0;  }
		</style>
		`
		err := os.WriteFile(tmpGohtml, []byte(rawHTML), 0644)
		if err != nil {
			t.Fatal(err)
		}

		err = GenFile(context.Background(), tmpGohtml, tmpOut, Option{CompactMode: true, NoLineNumber: true})
		if err != nil {
			t.Fatalf("failed to compile with CompactMode: %v", err)
		}

		content, err := os.ReadFile(tmpOut)
		if err != nil {
			t.Fatal(err)
		}

		contentStr := string(content)

		// The entire markup should be combined into a single, optimized string allocation with compact HTML
		expectedPayload := `_buffer.WriteString("<div class=\"header\" id=\"top\"><h1>Hello World!</h1></div><pre class=\"code-block\">\n\t\t\tline 1\n\t\t\t  line 2\n\t\t</pre><script>\n\t\t\tconsole.log( \"test\" );\n\t\t</script><style>\n\t\t\tbody {  margin: 0;  }\n\t\t</style>")`
		if !strings.Contains(contentStr, expectedPayload) {
			t.Errorf("expected combined compact HTML statement not found in:\n%s", contentStr)
		}

		// Test script with '<' comparison operators inside to ensure they are preserved and tag does not break
		scriptWithComparison := "<script>\nfor (var i = 0; i < 10; i++) {\n    if (i < 5) {\n        console.log(i);\n    }\n}\n</script>"
		compacted := compactHTML(scriptWithComparison)
		if !strings.Contains(compacted, "i < 10") || !strings.Contains(compacted, "i < 5") {
			t.Errorf("compactHTML corrupted '<' inside script: %s", compacted)
		}
		if !strings.HasSuffix(compacted, "</script>") {
			t.Errorf("compactHTML did not properly close script tag: %s", compacted)
		}
	})
}

