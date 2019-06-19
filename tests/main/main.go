package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strings"

	"golang.org/x/tools/go/loader"
)

func main3() {
	fset := token.NewFileSet()

	content, err := ioutil.ReadFile("../tpl/g.go")
	if err != nil {
		panic(err) // parse error
	}

	// Parse the input string, []byte, or io.Reader,
	// recording position information in fset.
	// ParseFile returns an *ast.File, a syntax tree.
	f, err := parser.ParseFile(fset, "../tpl/g.go", content, 0)
	if err != nil {
		panic(err) // parse error
	}

	// Run type checker
	info := types.Info{Types: make(map[ast.Expr]types.TypeAndValue)}

	// A Config controls various options of the type checker.
	// The defaults work fine except for one setting:
	// we must specify how to deal with imports.
	conf := types.Config{Importer: importer.Default()}

	// Type-check the package containing only file f.
	// Check returns a *types.Package.
	_, err = conf.Check("mypkg", fset, []*ast.File{f}, &info)
	if err != nil {
		panic(err) // parse error
	}

	for s, v := range info.Types {
		fmt.Println(s, ":", v.Type.String())
	}

	// Inspect variable types in f()
	for _, varDecl := range f.Decls[2].(*ast.FuncDecl).Body.List {
		fmt.Println(reflect.TypeOf(varDecl))
		switch v := varDecl.(type) {
		case *ast.RangeStmt:
			for _, i := range v.Body.List {
				e := i.(*ast.ExprStmt)
				c := e.X.(*ast.CallExpr)
				fmt.Println(c)
				fmt.Println(reflect.TypeOf(c.Fun))
				switch t2 := c.Fun.(type) {
				case *ast.SelectorExpr:
					if strings.HasPrefix(t2.Sel.Name, "HTMLEscape") {
						fmt.Println(t2.Sel.Name)
						arg := c.Args[0].(*ast.SelectorExpr)
						fmt.Println("-------")
						fmt.Println(arg.Sel)
					}
				}
			}
			// switch t2 := v.Fun.(type) {
			// case *ast.SelectorExpr:
			// 	// if t2.X == "gorazor" && strings.HasPrefix(t2.Sel.Name, "HTMLEscape") {
			// 	if strings.HasPrefix(t2.Sel.Name, "HTMLEscape") {
			// 		fmt.Println(v.Args[0])
			// 	}
			// }
		}

		// value := varDecl.(*ast.DeclStmt).Decl.(*ast.GenDecl).Specs[0].(*ast.ValueSpec)

		// pos := fset.Position(value.Type.Pos())
		// typ := info.Types[value.Type].Type

		// fmt.Println(pos, "basic:", typ.String())
	}
}

func main2() {

	filename := "../tpl/g.go"
	var conf loader.Config

	// Parse the specified files and create an ad hoc package with path "foo".
	// All files must have the same 'package' declaration.
	conf.CreateFromFilenames("../tpl/", filename)
	prog, err := conf.Load()
	if err != nil {
		panic(err)
	}
	printProgram(prog)

}

func printProgram(prog *loader.Program) {
	// Created packages are the initial packages specified by a call
	// to CreateFromFilenames or CreateFromFiles.
	var names []string
	for _, info := range prog.Created {
		names = append(names, info.Pkg.Path())
	}
	fmt.Printf("created: %s\n", names)

	// Imported packages are the initial packages specified by a
	// call to Import or ImportWithTests.
	names = nil
	for _, info := range prog.Imported {
		if strings.Contains(info.Pkg.Path(), "internal") {
			continue // skip, to reduce fragility
		}
		names = append(names, info.Pkg.Path())
	}
	sort.Strings(names)
	fmt.Printf("imported: %s\n", names)

	// InitialPackages contains the union of created and imported.
	names = nil
	for _, info := range prog.InitialPackages() {
		names = append(names, info.Pkg.Path())
	}
	sort.Strings(names)
	fmt.Printf("initial: %s\n", names)

	// AllPackages contains all initial packages and their dependencies.
	names = nil
	for pkg := range prog.AllPackages {
		names = append(names, pkg.Path())
	}
	sort.Strings(names)
	fmt.Printf("all: %s\n", names)
}

func main() {

	filename := "../tpl/g.go"
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing file: " + filename)
		os.Exit(255)
	}

	info := types.Info{Types: make(map[ast.Expr]types.TypeAndValue)}

	// conf := types.Config{Importer: importer.Default()}
	conf := types.Config{Importer: importer.For("source", nil)}

	_, err = conf.Check("mypkg", fset, []*ast.File{node}, &info)
	if err != nil {
		println("Type check error: " + err.Error())
	}

	// traverse all tokens
	ast.Inspect(node, func(n ast.Node) bool {
		// fmt.Println(reflect.TypeOf(n))
		switch t := n.(type) {
		case *ast.CallExpr:
			switch t2 := t.Fun.(type) {
			case *ast.SelectorExpr:
				if t2.Sel.Name == "HTMLEscape" && t2.X.(*ast.Ident).String() == "gorazor" {
					typ := info.Types[t.Args[0]]
					if typ.Type != nil {
						fmt.Println(t)
						fmt.Println(t.Args[0], typ.Type.String())
					}
				}
			}

		// find variable declarations
		case *ast.TypeSpec:
			// which are public
			if t.Name.IsExported() {
				switch t.Type.(type) {
				// and are interfaces
				case *ast.InterfaceType:
					fmt.Println(t.Name.Name)
				}
			}
		}
		return true
	})
}
