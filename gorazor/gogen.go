package gorazor

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/fsnotify.v1"
)

var GorazorNamespace = `"github.com/sipin/gorazor/gorazor"`

//------------------------------ Compiler ------------------------------ //
const (
	CMKP = iota
	CBLK
	CSTAT
)

func getValStr(e interface{}) string {
	switch v := e.(type) {
	case *Ast:
		return v.TagName
	case Token:
		if !(v.Type == AT || v.Type == AT_COLON) {
			return v.Text
		}
		return ""
	default:
		panic(e)
	}
}

type Part struct {
	ptype int
	value string
}

type Compiler struct {
	ast      *Ast
	buf      string //the final result
	layout   string
	firstBLK int
	params   []string
	parts    []Part
	imports  map[string]bool
	options  Option
	dir      string
	file     string
}

func (self *Compiler) addPart(part Part) {
	if len(self.parts) == 0 {
		self.parts = append(self.parts, part)
		return
	}
	last := &self.parts[len(self.parts)-1]
	if last.ptype == part.ptype {
		last.value += part.value
	} else {
		self.parts = append(self.parts, part)
	}
}

func (self *Compiler) genPart() {
	res := ""

	for _, p := range self.parts {
		if p.ptype == CMKP && p.value != "" {
			// do some escapings
			for strings.HasSuffix(p.value, "\n") {
				p.value = p.value[:len(p.value)-1]
			}
			p.value = fmt.Sprintf("%#v", p.value)
			p.value = p.value[1 : len(p.value)-1]
			p.value = strings.Replace(p.value, `\t`, "\t", -1)

			if p.value != `\n` && p.value != "" {
				res += "_buffer.WriteString(\"" + p.value + "\")\n"
			}
		} else if p.ptype == CBLK {
			res += p.value + "\n"
		} else {
			res += p.value
		}
	}
	self.buf = res
}

func makeCompiler(ast *Ast, options Option, input string) *Compiler {
	dir := filepath.Base(filepath.Dir(input))
	file := strings.Replace(filepath.Base(input), gz_extension, "", 1)
	if options["NameNotChange"] == nil {
		file = Capitalize(file)
	}
	return &Compiler{ast: ast, buf: "",
		layout: "", firstBLK: 0,
		params: []string{}, parts: []Part{},
		imports: map[string]bool{},
		options: options,
		dir:     dir,
		file:    file,
	}
}

func (cp *Compiler) visitBLK(child interface{}, ast *Ast) {
	cp.addPart(Part{CBLK, getValStr(child)})
}

func (cp *Compiler) visitMKP(child interface{}, ast *Ast) {

	cp.addPart(Part{CMKP, getValStr(child)})
}

// First block contains imports and parameters, specific action for layout,
// NOTE, layout have some conventions.
func (cp *Compiler) visitFirstBLK(blk *Ast) {
	pre := cp.buf
	cp.buf = ""
	first := ""
	backup := cp.parts
	cp.parts = []Part{}
	cp.visitAst(blk)
	cp.genPart()
	first, cp.buf = cp.buf, pre
	cp.parts = backup

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", "package main\n"+first, parser.ImportsOnly)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		for _, s := range f.Imports {
			v := s.Path.Value
			if s.Name != nil {
				v = s.Name.Name + " " + v
			}
			parts := strings.SplitN(v, "/", -1)
			if len(parts) >= 2 && parts[len(parts)-2] == "layout" {
				cp.layout = strings.Replace(v, "\"", "", -1)
				dir := strings.Join(parts[0:len(parts)-1], "/") + "\""
				cp.imports[dir] = true
			} else {
				cp.imports[v] = true
			}
		}
	}

	lines := strings.SplitN(first, "\n", -1)
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "var") {
			vname := l[4:]
			if strings.HasSuffix(l, "gorazor.Widget") {
				cp.imports[GorazorNamespace] = true
				cp.params = append(cp.params, vname[:len(vname)-14]+"gorazor.Widget")
			} else {
				cp.params = append(cp.params, vname)
			}
		}
	}
	if cp.layout != "" {
		path := cp.layout + ".gohtml"
		if exists(path) && len(LayOutArgs(path)) == 0 {
			//TODO, bad for performance
			_cp, err := run(path, cp.options)
			if err != nil {
				panic(err)
			}
			SetLayout(cp.layout, _cp.params)
		}
	}
}

func (cp *Compiler) visitExp(child interface{}, parent *Ast, idx int, isHomo bool) {
	start := ""
	end := ""
	ppNotExp := true
	ppChildCnt := len(parent.Children)
	pack := cp.dir
	htmlEsc := cp.options["htmlEscape"]
	if parent.Parent != nil && parent.Parent.Mode == EXP {
		ppNotExp = false
	}
	val := getValStr(child)
	if htmlEsc == nil {
		if ppNotExp && idx == 0 && isHomo {
			needEsape := true
			switch {
			case val == "helper" || val == "html" || val == "raw":
				needEsape = false
			case pack == "layout":
				needEsape = true
				for _, param := range cp.params {
					if strings.HasPrefix(param, val+" ") {
						needEsape = false
						break
					}
				}
			}

			if needEsape {
				start += "gorazor.HTMLEscape("
				cp.imports[GorazorNamespace] = true
			} else {
				start += "("
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

	v := start
	if val == "raw" {
		v += end
	} else {
		v += val + end
	}
	cp.addPart(Part{CSTAT, v})
}

func (cp *Compiler) visitAst(ast *Ast) {
	switch ast.Mode {
	case MKP:
		cp.firstBLK = 1
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
				if _, ok := c.(Token); ok {
					cp.visitBLK(c, ast)
				} else {
					cp.visitAst(c.(*Ast))
				}
			}
		}
	case EXP:
		cp.firstBLK = 1
		nonExp := ast.hasNonExp()
		for i, c := range ast.Children {
			if _, ok := c.(Token); ok {
				cp.visitExp(c, ast, i, !nonExp)
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

// TODO, this is dirty now
func (cp *Compiler) processLayout() {
	lines := strings.SplitN(cp.buf, "\n", -1)
	out := ""
	sections := []string{}
	scope := 0
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "section") && strings.HasSuffix(l, "{") {
			name := l
			name = strings.TrimSpace(name[7 : len(name)-1])
			out += "\n " + name + " := func() string {\n"
			out += "var _buffer bytes.Buffer\n"
			scope = 1
			sections = append(sections, name)
		} else if scope > 0 {
			if strings.HasSuffix(l, "{") {
				scope++
			} else if strings.HasSuffix(l, "}") {
				scope--
			}
			if scope == 0 {
				out += "return _buffer.String()\n}\n"
				scope = 0
			} else {
				out += l + "\n"
			}
		} else {
			out += l + "\n"
		}
	}
	cp.buf = out
	foot := "\nreturn "
	if cp.layout != "" {
		parts := strings.SplitN(cp.layout, "/", -1)
		base := Capitalize(parts[len(parts)-1])
		foot += "layout." + base + "("
	}
	foot += "_buffer.String()"
	args := LayOutArgs(cp.layout)
	if len(args) == 0 {
		for _, sec := range sections {
			foot += ", " + sec + "()"
		}
	} else {
		for idx, arg := range args {
			//body has been done
			if idx == 0 {
				continue
			}
			arg = strings.Replace(arg, "string", "", -1)
			arg = strings.TrimSpace(arg)
			found := false
			for _, sec := range sections {
				if sec == arg {
					found = true
					foot += ", " + sec + "()"
					break
				}
			}
			if !found {
				foot += ", " + `""`
			}
		}
	}
	if cp.layout != "" {
		foot += ")"
	}
	foot += "\n}\n"
	cp.buf += foot
}

func (cp *Compiler) visit() {
	cp.visitAst(cp.ast)
	cp.genPart()

	pack := cp.dir
	fun := cp.file

	cp.imports[`"bytes"`] = true
	head := "package " + pack + "\n import (\n"
	for k, _ := range cp.imports {
		head += k + "\n"
	}
	head += "\n)\n func " + fun + "("
	for idx, p := range cp.params {
		head += p
		if idx != len(cp.params)-1 {
			head += ", "
		}
	}
	head += ") string {\n var _buffer bytes.Buffer\n"
	cp.buf = head + cp.buf
	cp.processLayout()
}

func run(path string, Options Option) (*Compiler, error) {
	content, err := ioutil.ReadFile(path)
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
	if Options["Debug"] != nil {
		fmt.Println("------------------- TOKEN START -----------------")
		for _, elem := range res {
			elem.P()
		}
		fmt.Println("--------------------- TOKEN END -----------------\n")
	}

	parser := &Parser{&Ast{}, nil, res, []Token{}, false, UNK}
	err = parser.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	//DEBUG
	if Options["Debug"] != nil {
		fmt.Println("--------------------- AST START -----------------")
		parser.ast.debug(0, 20)
		fmt.Println("--------------------- AST END -----------------\n")
		if parser.ast.Mode != PRG {
			panic("TYPE")
		}
	}
	cp := makeCompiler(parser.ast, Options, path)
	cp.visit()
	return cp, nil
}

func generate(path string, output string, Options Option) error {
	cp, err := run(path, Options)
	if err != nil || cp == nil {
		panic(err)
	}
	err = ioutil.WriteFile(output, []byte(cp.buf), 0644)
	cmd := exec.Command("gofmt", "-w", output)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("gofmt: ", err)
		return err
	}
	if Options["Debug"] != nil {
		content, _ := ioutil.ReadFile(output)
		fmt.Println(string(content))
	}
	return err
}

func watchDir(input, output string, options Option) error {
	log.Println("Watching dir:", input, output)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	output_path := func(path string) string {
		res := strings.Replace(path, input, output, 1)
		return res
	}

	gen := func(filename string) error {
		outpath := output_path(filename)
		outpath = strings.Replace(outpath, ".gohtml", ".go", 1)
		outdir := filepath.Dir(outpath)
		if !exists(outdir) {
			os.MkdirAll(outdir, 0775)
		}
		err := GenFile(filename, outpath, options)
		if err == nil {
			log.Printf("%s -> %s\n", filename, outpath)
		}
		return err
	}

	visit_gen := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			//Just do file with exstension .gohtml
			if !strings.HasSuffix(path, ".gohtml") {
				return nil
			}
			filename := filepath.Base(path)
			if strings.HasPrefix(filename, ".#") {
				return nil
			}
			err := gen(path)
			if err != nil {
				return err
			}
		}
		return nil
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				filename := event.Name
				if filename == "" {
					//should be a bug for fsnotify
					continue
				}
				if event.Op&fsnotify.Remove != fsnotify.Remove &&
					(event.Op&fsnotify.Write == fsnotify.Write ||
						event.Op&fsnotify.Create == fsnotify.Create) {
					stat, err := os.Stat(filename)
					if err != nil {
						continue
					}
					if stat.IsDir() {
						log.Println("add dir:", filename)
						watcher.Add(filename)
						output := output_path(filename)
						log.Println("mkdir:", output)
						if !exists(output) {
							os.MkdirAll(output, 0755)
							err = filepath.Walk(filename, visit_gen)
							if err != nil {
								done <- true
							}
						}
						continue
					}
					if !strings.HasPrefix(filepath.Base(filename), ".#") &&
						strings.HasSuffix(filename, ".gohtml") {
						gen(filename)
					}
				} else if event.Op&fsnotify.Remove == fsnotify.Remove ||
					event.Op&fsnotify.Rename == fsnotify.Rename {
					output := output_path(filename)
					if exists(output) {
						//shoud be dir
						watcher.Remove(filename)
						os.RemoveAll(output)
						log.Println("remove dir:", output)
					} else if strings.HasSuffix(output, ".gohtml") {
						output = strings.Replace(output, ".gohtml", ".go", 1)
						if exists(output) {
							os.Remove(output)
							log.Println("removing file:", output)
						}
					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
				continue
			}
		}
	}()

	visit := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			watcher.Add(path)
		}
		return nil
	}

	err = filepath.Walk(input, visit)
	err = watcher.Add(input)
	if err != nil {
		log.Fatal(err)
	}
	<-done
	return nil
}
