package razorcore

import (
	"context"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// GorazorNamespace defines util pkg namespace used in template
var GorazorNamespace = `gorazor "github.com/sipin/gorazor/runtime"`


// ------------------------------ Compiler ------------------------------ //
const (
	CMKP = iota
	CBLK
	CSTAT
)

var execDir string

func init() {
	// make sure running in source directory
	_, filename, _, _ := runtime.Caller(0)
	execDir = path.Dir(filename) + "/"
}

func getValStr(e interface{}) string {
	switch v := e.(type) {
	case *Ast:
		return v.TagName
	case Token:
		if !(v.Type == tkAt || v.Type == tkAtColon) {
			return v.Text
		}
		return ""
	default:
		panic(e)
	}
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func compactTag(tag string) string {
	var result strings.Builder
	inQuote := byte(0)
	n := len(tag)
	i := 0
	for i < n {
		c := tag[i]
		if inQuote != 0 {
			if c == inQuote {
				inQuote = 0
			}
			result.WriteByte(c)
			i++
		} else if c == '"' || c == '\'' {
			inQuote = c
			result.WriteByte(c)
			i++
		} else if isWhitespace(c) {
			for i < n && isWhitespace(tag[i]) {
				i++
			}
			if i < n && tag[i] != '>' && tag[i] != '/' {
				result.WriteByte(' ')
			}
		} else {
			result.WriteByte(c)
			i++
		}
	}
	return result.String()
}

func compactHTML(html string) string {
	var result strings.Builder
	inPre := false
	inTagName := ""

	i := 0
	n := len(html)
	for i < n {
		if html[i] == '<' {
			if inPre {
				closingPrefix := "</" + inTagName
				isActualClosing := false
				if len(html)-i >= len(closingPrefix) && strings.EqualFold(html[i:i+len(closingPrefix)], closingPrefix) {
					afterIdx := i + len(closingPrefix)
					if afterIdx == len(html) || html[afterIdx] == '>' || isWhitespace(html[afterIdx]) || html[afterIdx] == '/' {
						isActualClosing = true
					}
				}
				if !isActualClosing {
					result.WriteByte(html[i])
					i++
					continue
				}
			}

			// Scan the entire tag up to '>'
			j := i
			inQuote := byte(0)
			for j < n {
				c := html[j]
				if inQuote != 0 {
					if c == inQuote {
						inQuote = 0
					}
				} else if c == '"' || c == '\'' {
					inQuote = c
				} else if c == '>' {
					j++
					break
				}
				j++
			}

			tagContent := html[i:j]

			// Now parse the tag name from tagContent
			if len(tagContent) > 1 {
				tagName := ""
				isClosing := false
				k := 1
				if tagContent[k] == '/' {
					isClosing = true
					k++
				}
				startTagName := k
				for k < len(tagContent) && tagContent[k] != '>' && tagContent[k] != ' ' && tagContent[k] != '\t' && tagContent[k] != '\n' && tagContent[k] != '/' {
					k++
				}
				tagName = strings.ToLower(tagContent[startTagName:k])

				if isClosing {
					if tagName == inTagName {
						inPre = false
						inTagName = ""
					}
				} else {
					if tagName == "pre" || tagName == "textarea" || tagName == "script" || tagName == "style" {
						inPre = true
						inTagName = tagName
					}
				}
			}

			result.WriteString(compactTag(tagContent))
			i = j
			continue
		}

		if inPre {
			result.WriteByte(html[i])
			i++
		} else {
			if isWhitespace(html[i]) {
				// Scan all consecutive whitespaces
				for i < n && isWhitespace(html[i]) {
					i++
				}

				// Check character before and after
				var prevChar byte
				if result.Len() > 0 {
					prev := result.String()
					prevChar = prev[len(prev)-1]
				}
				var nextChar byte
				if i < n {
					nextChar = html[i]
				}

				// Omit whitespace completely if it is between '>' and '<'
				if prevChar == '>' && nextChar == '<' {
					// Omit it
				} else {
					// Otherwise collapse to a single space, but only if we are not at the very start/end of the string
					if result.Len() > 0 && i < n {
						result.WriteByte(' ')
					}
				}
			} else {
				result.WriteByte(html[i])
				i++
			}
		}
	}
	return result.String()
}

// Part represent gorazor template parts
type Part struct {
	ptype int
	value string
	line  int
}

// Compiler generate go code for gorazor template
type Compiler struct {
	ctx         context.Context
	currentLine int
	inputPath   string
	tplPath     string
	ast         *Ast
	buf         string //the final result
	isLayout    bool
	layout      string
	layoutAlias string
	firstBLK    int
	params      []string
	paramNames  []string
	parts       []Part
	imports     map[string]bool
	options     Option
	dir         string
	file        string
}

type compilerError struct {
	err error
}

func (cp *Compiler) errorf(line int, format string, args ...interface{}) {
	if line <= 0 {
		line = cp.currentLine
	}
	msg := fmt.Sprintf(format, args...)
	snippet := ""
	if line > 0 && cp.inputPath != "" {
		if content, err := os.ReadFile(cp.inputPath); err == nil {
			lines := strings.Split(string(content), "\n")
			if line <= len(lines) {
				snippet = fmt.Sprintf("\n\nContext:\n---> %d: %s\n", line, strings.TrimSpace(lines[line-1]))
			}
		}
	}
	detailErr := fmt.Errorf("compilation error in %s:%d: %s%s", cp.inputPath, line, msg, snippet)
	panic(compilerError{err: detailErr})
}

func (cp *Compiler) addPart(part Part) {
	if len(cp.parts) == 0 {
		cp.parts = append(cp.parts, part)
		return
	}
	last := &cp.parts[len(cp.parts)-1]
	if last.ptype == part.ptype {
		last.value += part.value
	} else {
		cp.parts = append(cp.parts, part)
	}
}

func (cp *Compiler) isLayoutSectionPart(p Part) (is bool, val string) {
	if !cp.isLayout {
		return
	}

	if !strings.HasPrefix(p.value, "_buffer.WriteString((") {
		return
	}

	if !strings.HasSuffix(p.value, "))\n") {
		return
	}

	val = p.value[21 : len(p.value)-3]
	for _, name := range cp.paramNames {
		if val == name {
			return true, val
		}
	}

	return
}

func (cp *Compiler) isLayoutSectionTest(p Part) (is bool, val string) {
	if !cp.isLayout {
		return
	}

	line := strings.TrimSpace(p.value)
	line = strings.ReplaceAll(line, " ", "")

	for _, p := range cp.paramNames {
		if line == "if"+p+`==""{` {
			return true, "if " + p + " == nil {\n"
		}
		if line == "if"+p+`!=""{` {
			return true, "if " + p + " != nil {\n"
		}
	}

	return
}

func (cp *Compiler) getLineHint(line int) string {
	if cp.options.NoLineNumber {
		return ""
	}
	return "// Line: " + strconv.Itoa(line) + "\n"
}

func (cp *Compiler) genPart() {
	res := ""

	for _, p := range cp.parts {
		if p.ptype == CMKP && p.value != "" {
			if cp.options.CompactMode {
				p.value = compactHTML(p.value)
			}
			// do some escapings
			for strings.HasSuffix(p.value, "\n") {
				p.value = p.value[:len(p.value)-1]
			}
			if p.value != "" {
				p.value = fmt.Sprintf("%#v", p.value)
				if p.line > 0 {
					res += cp.getLineHint(p.line)
				}

				res += "_buffer.WriteString(" + p.value + ")\n"
			}
		} else if p.ptype == CBLK {
			if ok, val := cp.isLayoutSectionTest(p); ok {
				res += val
			} else {
				res += p.value + "\n"
			}
		} else if ok, val := cp.isLayoutSectionPart(p); ok {
			res += cp.getLineHint(p.line)
			res += val + "(_buffer)\n"
		} else {
			res += p.value
		}
	}
	cp.buf = res
}

func makeCompiler(ctx context.Context, ast *Ast, options Option, input string) *Compiler {
	inputPath, _ := filepath.Abs(input)
	dir := filepath.Base(filepath.Dir(inputPath))
	file := strings.Replace(filepath.Base(input), gzExtension, "", 1)
	if !options.NameNotChange {
		file = Capitalize(file)
	}
	cp := &Compiler{
		ctx:         ctx,
		currentLine: 1,
		ast:         ast,
		buf:         "",
		layout:      "", firstBLK: 0,
		params:      []string{}, parts: []Part{},
		imports:     map[string]bool{},
		options:     options,
		dir:         dir,
		file:        file,
	}

	if dir == "layout" {
		cp.isLayout = true
	}

	cp.inputPath = strings.ReplaceAll(input, "\\", "/")
	tplPath := strings.ReplaceAll(cp.inputPath, execDir, "")
	if filepath.IsAbs(tplPath) {
		if cwd, err := os.Getwd(); err == nil && cwd != "" {
			cwdSlash := strings.ReplaceAll(cwd, "\\", "/") + "/"
			tplPath = strings.TrimPrefix(tplPath, cwdSlash)
		}
	}
	cp.tplPath = strings.TrimPrefix(tplPath, "./")
	return cp
}

func (cp *Compiler) visitBLK(child Token) {
	cp.addPart(Part{CBLK, getValStr(child), child.Line})
}

func (cp *Compiler) visitMKP(child Token) {
	cp.addPart(Part{CMKP, getValStr(child), child.Line})
}

func (cp *Compiler) settleLayout(layoutFunc string) {
	path := cp.layout + "/" + layoutFunc + ".gohtml"

	if !exists(path) && cp.options.TemplateNamespacePrefix != "" {
		path = path[len(cp.options.TemplateNamespacePrefix)+1:]
	}

	if !exists(path) {
		layoutFunc = strings.ToLower(layoutFunc[0:1]) + layoutFunc[1:]
		path = cp.layout + "/" + layoutFunc + ".gohtml"

		if !exists(path) && cp.options.TemplateNamespacePrefix != "" {
			path = path[len(cp.options.TemplateNamespacePrefix)+1:]
		}
	}

	cp.layout = cp.layout + "/" + layoutFunc
	if !exists(path) {
		cp.errorf(1, "Can't find layout: %s [%s]", cp.layout, cp.file)
	}

	if len(cp.options.LayoutCache.Get(cp.layout)) == 0 {
		//TODO, bad for performance
		_cp, err := run(cp.ctx, path, cp.options)
		if err != nil {
			cp.errorf(1, "failed to compile layout %s: %v", path, err)
		}
		cp.options.LayoutCache.Set(cp.layout, _cp.params)
	}
}

// First block contains imports and parameters, specific action for layout,
// NOTE, layout have some conventions.
func (cp *Compiler) visitFirstBLK(blk *Ast) {
	blockContent := cp.extractBlockContent(blk)
	cp.processImports(blockContent)
	layoutFunc := cp.processDeclarations(blockContent)
	cp.finalizeLayout(layoutFunc)
}

// extractBlockContent processes the AST block and returns the generated content
func (cp *Compiler) extractBlockContent(blk *Ast) string {
	pre := cp.buf
	cp.buf = ""
	backup := cp.parts
	cp.parts = []Part{}
	
	cp.visitAst(blk)
	cp.genPart()
	
	content := cp.buf
	cp.buf = pre
	cp.parts = backup
	
	return content
}

// processImports parses and processes Go imports from the block content
func (cp *Compiler) processImports(content string) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", "package main\n"+content, parser.ImportsOnly)
	if err != nil {
		line := 1
		if matches := regexp.MustCompile(`(?:^|:)(\d+):`).FindStringSubmatch(err.Error()); len(matches) > 1 {
			if l, e := strconv.Atoi(matches[1]); e == nil {
				line = l - 1
				if line < 1 {
					line = 1
				}
			}
		}
		cp.errorf(line, "failed to parse imports block: %v", err)
	}
	
	for _, s := range f.Imports {
		importPath := s.Path.Value
		alias := ""
		if s.Name != nil {
			alias = s.Name.Name
			importPath = alias + " " + importPath
		}
		
		cp.imports[importPath] = true
		cp.detectLayoutImport(s.Path.Value, alias)
	}
}

// detectLayoutImport checks if an import path is a layout import and sets cp.layout and cp.layoutAlias
func (cp *Compiler) detectLayoutImport(pathValue, alias string) {
	cleanPath := strings.ReplaceAll(pathValue, "\"", "")
	parts := strings.SplitN(cleanPath, "/", -1)
	if len(parts) >= 1 && parts[len(parts)-1] == "layout" {
		cp.layout = cleanPath
		cp.layoutAlias = alias
	}
}

// processDeclarations processes variable declarations and layout assignments
func (cp *Compiler) processDeclarations(content string) string {
	lines := strings.SplitN(content, "\n", -1)
	var layoutFunc string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		layoutFunc = cp.processDeclarationLine(line, layoutFunc)
	}
	
	return layoutFunc
}

// processDeclarationLine processes a single line of declarations
func (cp *Compiler) processDeclarationLine(line, currentLayoutFunc string) string {
	switch {
	case strings.HasPrefix(line, "var"):
		return cp.processVariableDeclaration(line, currentLayoutFunc)
	case strings.HasPrefix(line, "isLayout"):
		cp.isLayout = strings.HasSuffix(line, "true")
		return currentLayoutFunc
	case strings.HasPrefix(line, "layout:=") || strings.HasPrefix(line, "layout :="):
		return cp.processLayoutAssignment(line)
	default:
		return currentLayoutFunc
	}
}

// processVariableDeclaration processes variable declarations (var statements)
func (cp *Compiler) processVariableDeclaration(line, currentLayoutFunc string) string {
	vname := line[4:] // Remove "var " prefix
	
	switch {
	case strings.HasSuffix(line, "gorazor.Widget"):
		cp.processWidgetVariable(vname)
		return currentLayoutFunc
	case strings.HasPrefix(vname, "layout"):
		return cp.extractLayoutFunctionName(vname)
	default:
		cp.processRegularVariable(vname)
		return currentLayoutFunc
	}
}

// processWidgetVariable handles gorazor.Widget variable declarations
func (cp *Compiler) processWidgetVariable(vname string) {
	cp.imports[GorazorNamespace] = true
	cp.params = append(cp.params, vname[:len(vname)-14]+"gorazor.Widget")
	name := strings.SplitN(vname, " ", 2)[0]
	cp.paramNames = append(cp.paramNames, name)
}

// processRegularVariable handles regular variable declarations
func (cp *Compiler) processRegularVariable(vname string) {
	cp.params = append(cp.params, vname)
	name := strings.SplitN(vname, " ", 2)[0]
	cp.paramNames = append(cp.paramNames, name)
}

// processLayoutAssignment processes layout := assignments
func (cp *Compiler) processLayoutAssignment(line string) string {
	vname := strings.TrimSpace(strings.Split(line, ":=")[1])
	return cp.extractLayoutFunctionName(vname)
}

// extractLayoutFunctionName extracts function name from layout variable assignments
func (cp *Compiler) extractLayoutFunctionName(vname string) string {
	funcName := strings.SplitN(vname, ".", -1)
	return funcName[len(funcName)-1]
}

// finalizeLayout completes layout processing if layout was detected
func (cp *Compiler) finalizeLayout(layoutFunc string) {
	if cp.layout != "" {
		cp.settleLayout(layoutFunc)
	}
}

func (cp *Compiler) isExpNeedEscape(val string) (needEsape bool) {
	switch {
	case val == "helper" || val == "html" || val == "raw":
		return false
	case cp.dir == "layout":
		for _, param := range cp.params {
			if strings.HasPrefix(param, val+" ") {
				return false
			}
		}
	}
	return true
}

func (cp *Compiler) visitExp(child interface{}, parent *Ast, idx int, isHomo bool) {
	start := ""
	end := ""
	ppNotExp := true
	ppChildCnt := len(parent.Children)
	if parent.Parent != nil && parent.Parent.Mode == EXP {
		ppNotExp = false
	}
	val := getValStr(child)

	if ppNotExp && idx == 0 && isHomo {
		if cp.isExpNeedEscape(val) {
			start += "gorazor.HTMLEscape("
			cp.imports[GorazorNamespace] = true
		} else {
			start += "("
		}
	}
	if ppNotExp && idx == ppChildCnt-1 && isHomo {
		end += ")"
	}

	lineHint := ""
	lineNumber := 0
	if ppNotExp && idx == 0 {
		if token, ok := child.(Token); ok {
			lineNumber = token.Line
			lineHint = cp.getLineHint(token.Line)
		}
		start = "_buffer.WriteString(" + start
	}
	if ppNotExp && idx == ppChildCnt-1 {
		end += ")\n"
	}

	if val == "raw" {
		cp.addPart(Part{CSTAT, lineHint + start + end, lineNumber})
	} else {
		p := Part{CSTAT, start + val + end, lineNumber}
		if ok, _ := cp.isLayoutSectionPart(p); !ok {
			p.value = lineHint + p.value
		}
		cp.addPart(p)
	}
}

func (cp *Compiler) visitAstBlk(ast *Ast) {
	if cp.firstBLK == 0 {
		cp.firstBLK = 1
		cp.visitFirstBLK(ast)
	} else {
		remove := false
		if len(ast.Children) >= 2 {
			first := ast.Children[0]
			last := ast.Children[len(ast.Children)-1]
			v1, ok1 := first.(Token)
			v2, ok2 := last.(Token)
			if ok1 && ok2 && v1.Text == "{" && v2.Text == "}" {
				remove = true
			}
		}
		for idx, c := range ast.Children {
			if remove && (idx == 0 || idx == len(ast.Children)-1) {
				continue
			}
			if token, ok := c.(Token); ok {
				cp.currentLine = token.Line
				cp.visitBLK(token)
			} else {
				cp.visitAst(c.(*Ast))
			}
		}
	}
}

func (cp *Compiler) visitAst(ast *Ast) {
	switch ast.Mode {
	case MKP:
		cp.firstBLK = 1
		for _, c := range ast.Children {
			if token, ok := c.(Token); ok {
				cp.currentLine = token.Line
				cp.visitMKP(token)
			} else {
				cp.visitAst(c.(*Ast))
			}
		}
	case BLK:
		cp.visitAstBlk(ast)
	case EXP:
		cp.firstBLK = 1
		nonExp := ast.hasNonExp()
		for i, c := range ast.Children {
			if token, ok := c.(Token); ok {
				cp.currentLine = token.Line
				cp.visitExp(token, ast, i, !nonExp)
			} else {
				cp.visitAst(c.(*Ast))
			}
		}
	case PRG:
		for _, c := range ast.Children {
			cp.visitAst(c.(*Ast))
		}
	}
}

func (cp *Compiler) hasLayout() bool {
	return cp.layout != ""
}

func (cp *Compiler) generateFoot(sections []string) string {
	foot := ""
	if cp.hasLayout() {
		foot += "\n"
		parts := strings.SplitN(cp.layout, "/", -1)
		base := Capitalize(parts[len(parts)-1])
		pkgPrefix := "layout"
		if cp.layoutAlias != "" {
			pkgPrefix = cp.layoutAlias
		}
		foot += pkgPrefix + ".Render" + base + "("
		foot += "_buffer, _body"
	} else if len(sections) > 0 {
		cp.errorf(1, "expect layout for sections: %s", cp.file)
	}

	args := cp.options.LayoutCache.Get(cp.layout)
	if len(args) == 0 {
		for _, sec := range sections {
			foot += ", " + sec + "()"
		}
	} else {
		for _, arg := range args[1:] {
			parts := strings.Fields(arg)
			if len(parts) > 0 {
				arg = parts[0]
			} else {
				arg = strings.TrimSpace(arg)
			}
			found := false
			for _, sec := range sections {
				if sec == arg {
					found = true
					foot += ", _" + sec
					break
				}
			}
			if !found {
				foot += ", " + `nil`
			}
		}
	}
	if cp.layout != "" {
		foot += ")"
	}

	return foot
}

func (cp *Compiler) processLayout() {
	lines := strings.SplitN(cp.buf, "\n", -1)
	out := ""
	sections := []string{}
	scope := 0
	hasBodyClosed := false

	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "section") && strings.HasSuffix(l, "{") {
			if !hasBodyClosed {
				hasBodyClosed = true
				out += "\n}\n"
			}

			name := l
			name = strings.TrimSpace(name[7 : len(name)-1])
			out += "\n _" + name + " := func(_buffer io.StringWriter) {\n"
			scope = 1
			sections = append(sections, name)
		} else if scope > 0 {
			if strings.HasSuffix(l, "{") {
				scope++
			} else if strings.HasSuffix(l, "}") {
				scope--
			}
			if scope == 0 {
				out += "\n}\n"
				scope = 0
			} else {
				out += l + "\n"
			}
		} else {
			out += l + "\n"
		}
	}

	if cp.hasLayout() && !hasBodyClosed {
		out += "\n}\n"
	}

	cp.buf = out

	foot := cp.generateFoot(sections)

	cp.buf += foot
}

func (cp *Compiler) getLayoutOverload() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`
	// %s generates %s
	func %s(%s) string {
		var _b strings.Builder

	`, cp.file, cp.tplPath, cp.file, strings.Join(cp.params, ", ")))

	var funcNames []string
	for _, name := range cp.paramNames {
		b.WriteString(fmt.Sprintf(`
		_%s := func(_buffer io.StringWriter) {
			_buffer.WriteString(%s)
		}
		`, name, name))
		funcNames = append(funcNames, "_"+name)
	}

	b.WriteString(fmt.Sprintf(`
		Render%s(&_b, %s)
		return _b.String()
	}

	`, cp.file, strings.Join(funcNames, ", ")))
	return b.String()
}

func (cp *Compiler) visit() {
	cp.visitAst(cp.ast)
	cp.genPart()

	pack := cp.dir
	fun := cp.file

	cp.imports[`"io"`] = true
	cp.imports[`"strings"`] = true

	head := fmt.Sprintf(`// This file is generated by gorazor %s
// DON'T modified manually
// Should edit source file and re-generate: %s

`, VERSION, cp.tplPath)

	head += "package " + pack + "\n import (\n"
	for k := range cp.imports {
		head += k + "\n"
	}

	funcArgs := strings.Join(cp.params, ", ")

	head += "\n)"

	if cp.isLayout {
		head += cp.getLayoutOverload()
		head += fmt.Sprintf(`
	// Render%s render %s
	`, fun, cp.tplPath)

		head += "func Render" + fun + "(_buffer io.StringWriter, " +
			strings.ReplaceAll(funcArgs, " string", " func(_buffer io.StringWriter)") + ") {\n"
	} else {
		head += fmt.Sprintf(`
	// %s generates %s
	func %s(%s) string {
		var _b strings.Builder
		Render%s(&_b, %s)
		return _b.String()
	}

	`, fun, cp.tplPath, fun, funcArgs, fun, strings.Join(cp.paramNames, ", "))

		head += fmt.Sprintf(`
	// Render%s render %s
	`, fun, cp.tplPath)

		head += "func Render" + fun + "(_buffer io.StringWriter, " + funcArgs + ") {\n"
	}

	if cp.hasLayout() {
		head += "\n_body := func(_buffer io.StringWriter) {\n"
	}

	cp.buf = head + cp.buf
	cp.processLayout()
	foot := "\n}\n"
	cp.buf += foot
}

func run(ctx context.Context, path string, Options Option) (*Compiler, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	text := string(content)
	lex := &Lexer{text, Tests}

	res, err := lex.Scan()
	if err != nil {
		return nil, err
	}

	//DEBUG
	if Options.IsDebug {
		fmt.Println("------------------- TOKEN START -----------------")
		for _, elem := range res {
			elem.P()
		}
		fmt.Println("--------------------- TOKEN END -----------------")
	}

	parser := &Parser{&Ast{}, nil, res, []Token{}, false, UNK}
	err = parser.Run()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}

	//DEBUG
	if Options.IsDebug {
		fmt.Println("--------------------- AST START -----------------")
		parser.ast.debug(0, 20)
		fmt.Println("--------------------- AST END -----------------")
		if parser.ast.Mode != PRG {
			panic("TYPE")
		}
	}
	cp := makeCompiler(ctx, parser.ast, Options, path)
	var compileErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				if ce, ok := r.(compilerError); ok {
					compileErr = ce.err
				} else {
					panic(r)
				}
			}
		}()
		cp.visit()
	}()
	if compileErr != nil {
		return nil, compileErr
	}
	return cp, nil
}

func generate(ctx context.Context, path string, output string, Options Option) error {
	cp, err := run(ctx, path, Options)
	if err != nil {
		return err
	}

	code := FormatBuffer(cp.buf)
	if !Options.QuickMode {
		_, code = optimize(output, cp.dir, code)
	}

	return os.WriteFile(output, []byte(code), 0644)
}
