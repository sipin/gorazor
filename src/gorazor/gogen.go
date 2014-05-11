package gorazor

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//------------------------------ Compiler ------------------------------ //
type Compiler struct {
	ast      *Ast
	buf      string
	firstBLK int
	params   []string
	imports  map[string]bool
	options  Option
}

func getValStr(e interface{}) string {
	switch v := e.(type) {
	case *Ast:
		return v.TagName
	case Token:
		return v.Text
	default:
		panic(e)
	}
}

func (cp *Compiler) visitMKP(child interface{}, ast *Ast) {
	v := strings.Replace(getValStr(child), "\n", "\\n", -1)
	cp.buf += "MKP(" + v + ")MKP"
}

// First block contains imports and parameters, specific action for layout,
// NOTE, layout have some conventions.
func (cp *Compiler) visitFirstBLK(blk *Ast) {
	pre := cp.buf
	cp.buf = ""
	first := ""
	cp.visitAst(blk)
	first, cp.buf = cp.buf, pre
	first = cp.cleanUp(first)
	isImport := false

	lines := bytes.SplitN([]byte(first), []byte("\n"), -1)
	for _, l := range lines {
		l = bytes.TrimSpace(l)
		if bytes.HasPrefix(l, []byte("import")) {
			isImport = true
			continue
		}
		if bytes.Compare(l, []byte(")")) == 0 {
			isImport = false
			continue
		}

		if isImport {
			parts := bytes.SplitN([]byte(l), []byte("/"), -1)
			if len(parts) >= 2 && bytes.Compare(parts[len(parts)-2], []byte("layout")) == 0 {
				lay := string(parts[len(parts)-1])
				lay = lay[:len(lay)-1]
				dir := strings.Replace(string(l), lay, "", 1)
				cp.imports[dir] = true
			} else {
				cp.imports[string(l)] = true
			}
		} else if bytes.HasPrefix(l, []byte("var")) {
			vname := string(l)[4:]
			cp.params = append(cp.params, vname)
		}
	}
}

func (cp *Compiler) visitBLK(child interface{}, ast *Ast) {
	cp.buf += "BLK(" + getValStr(child) + ")BLK"
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
	if false { //TODO, options
		if ppNotExp && idx == 0 && isHomo {
			if val == "helper" || val == "raw" { //TODO
				start += "("
			} else {
				start += "gorazor.HTMLEscape("
			}
		}
		if ppNotExp && idx == ppChildCnt-1 && isHomo {
			end += ")"
		}
	}

	if ppNotExp && idx == 0 {
		start = "_buffer.WriteString(" + start
	}
	if ppNotExp && idx == ppChildCnt-1 {
		end += ")\n"
	}

	if val == "raw" {
		cp.buf += start + end
	} else {
		cp.buf += start + val + end
	}
}

func (cp *Compiler) visitAst(ast *Ast) {
	switch ast.Mode {
	case MKP:
		for _, c := range ast.Children {
			if _, ok := c.(Token); ok {
				cp.visitMKP(c, ast)
			} else {
				cp.visitAst(c.(*Ast))
			}
		}
	case BLK:
		if cp.firstBLK == 0 {
			cp.firstBLK = 1
			cp.visitFirstBLK(ast)
		} else {
			for _, c := range ast.Children {
				if _, ok := c.(Token); ok {
					cp.visitBLK(c, ast)
				} else {
					cp.visitAst(c.(*Ast))
				}
			}
		}
	case EXP:
		nonExp := ast.hasNonExp()
		for i, c := range ast.Children {
			if _, ok := c.(Token); ok {
				cp.visitExp(c, ast, i, nonExp)
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

func (cp *Compiler) cleanUp(buf string) string {
	buf = strings.Replace(buf, ")BLKBLK(", "", -1)
	buf = strings.Replace(buf, ")MKPMKP(", "", -1)
	buf = strings.Replace(buf, "MKP(", "_buffer.WriteString(\"", -1)
	buf = strings.Replace(buf, ")MKP", "\")\n", -1)
	buf = strings.Replace(buf, "BLK(", "", -1)
	buf = strings.Replace(buf, ")BLK", "\n", -1)
	return buf
}

// TODO, this is dirty now
func (cp *Compiler) layout() {
	lines := bytes.SplitN([]byte(cp.buf), []byte("\n"), -1)
	out := ""
	insec := false
	for _, l := range lines {
		l = bytes.TrimSpace(l)
		if bytes.HasPrefix(l, []byte("section")) &&
			bytes.HasSuffix(l, []byte("{")) {
			name := string(l)
			name = name[7 : len(name)-1]
			out += "\n " + name + " := func() string {\n"
			out += "var _buffer bytes.Buffer\n"
			insec = true
		} else if insec && bytes.HasSuffix(l, []byte("}")) {
			out += "return _buffer.String()\n}\n"
			insec = false
		} else {
			out += string(l) + "\n"
		}
	}
	cp.buf = out
}

func (cp *Compiler) visit() {
	cp.visitAst(cp.ast)
	cp.buf = cp.cleanUp(cp.buf)

	dir := cp.options["Dir"].(string)
	fun := cp.options["File"].(string)

	cp.imports[`"bytes"`] = true
	cp.imports[`"gorazor"`] = true
	head := "package " + dir + "\n import (\n"
	for k, _ := range cp.imports {
		head += k + "\n"
	}
	head += "\n)\n func " + fun + "("
	for i, p := range cp.params {
		head += p
		if i != len(cp.params)-1 {
			head += ", "
		}
	}
	head += ") string {\n var _buffer bytes.Buffer\n"
	cp.buf = head + cp.buf + "\nreturn _buffer.String()\n}\n"
	cp.layout()
}

func Generate(path string, Options Option) (string, error) {
	buf := bytes.NewBuffer(nil)
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	io.Copy(buf, f)
	f.Close()

	text := string(buf.Bytes())
	lex := &Lexer{text, Tests}

	res, err := lex.Scan()
	if err != nil {
		return "", err
	}

	//DEBUG
	if Options["Debug"] != nil {
		fmt.Println("------------------- TOKEN START -----------------")
		for _, elem := range res {
			elem.P()
		}
		fmt.Println("--------------------- TOKEN END -----------------\n")
	}

	parser := &Parser{&Ast{}, res, []Token{}, false, false, UNK}
	err = parser.Run()

	//DEBUG
	if Options["Debug"] != nil {
		fmt.Println("--------------------- AST START -----------------")

		parser.ast.debug(0, 7)
		fmt.Println("--------------------- AST END -----------------\n")
		if parser.ast.Mode != PRG {
			panic("TYPE")
		}
	}

	cp := &Compiler{ast: parser.ast, buf: "", firstBLK: 0,
		params: []string{}, imports: map[string]bool{},
		options: Options}
	cp.visit()

	if Options["Debug"] != nil {
		fmt.Println(cp.buf)
	}
	return cp.buf, nil
}

//------------------------------ API ------------------------------
const (
	go_extension = ".go"
	gz_extension = ".gohtml"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// Generate from input to output file,
// gofmt will trigger an error if it fails.
func GenFile(input string, output string, options Option) error {
	//fmt.Printf("Processing: %s --> %s\n", input, output)

	//Use to as package name
	options["Dir"] = filepath.Base(filepath.Dir(input))
	options["File"] = strings.Replace(filepath.Base(input), gz_extension, "", 1)
	options["File"] = Capitalize(options["File"].(string))
	res, err := Generate(input, options)
	if err != nil {
		panic(err)
	} else {
		err := ioutil.WriteFile(output, []byte(res), 0644)
		if err != nil {
			panic(err)
		}
		cmd := exec.Command("gofmt", "-w", output)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// Generate from directory to directory, Find all the files with extension
// of .gohtml and generate it into target dir.
func GenFolder(indir string, outdir string) (err error) {
	if !exists(indir) {
		return errors.New("Input directory does not exsits")
	} else {
		if err != nil {
			return err
		}
	}

	//Make it
	if !exists(outdir) {
		os.MkdirAll(outdir, 0775)
	}

	options := Option{}
	incdir_abs, _ := filepath.Abs(indir)
	outdir_abs, _ := filepath.Abs(outdir)

	visit := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			//Just do file with exstension .gohtml
			if !bytes.HasSuffix([]byte(path), []byte(go_extension)) {
				return nil
			}
			//adjust with the abs path, so that we keep the same directory hierarchy
			input, _ := filepath.Abs(path)
			output := strings.Replace(input, incdir_abs, outdir_abs, 1)
			output = strings.Replace(output, gz_extension, go_extension, -1)
			GenFile(path, output, options)
		}
		return nil
	}
	err = filepath.Walk(indir, visit)
	return err
}
