package gorazor

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
)

func optimize(filename string, pkgname string, content string) (optimized bool, result string) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, content, parser.ParseComments)
	if err != nil {
		fmt.Println("Parsing error when optimize file: " + err.Error())
		return false, content
	}

	info := types.Info{Types: make(map[ast.Expr]types.TypeAndValue)}

	conf := types.Config{Importer: importer.For("source", nil)}

	_, err = conf.Check(pkgname, fset, []*ast.File{node}, &info)
	if err != nil {
		fmt.Println("Type check error when optimize file: " + err.Error())
		return false, content
	}

	replacedCount := 0
	var intPos []token.Pos
	var strPos []token.Pos

	// traverse all tokens
	ast.Inspect(node, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.CallExpr:
			switch t2 := t.Fun.(type) {
			case *ast.SelectorExpr:
				if t2.Sel.Name == "HTMLEscape" && t2.X.(*ast.Ident).String() == "gorazor" {
					typ := info.Types[t.Args[0]]
					if typ.Type != nil {
						if typ.Type.String() == "int" {
							intPos = append(intPos, t2.Pos())
						} else if typ.Type.String() == "string" {
							strPos = append(strPos, t2.Pos())
						}
					}
				}
			}
		}
		return true
	})

	replacedCount = len(intPos) + len(strPos)

	for _, pos := range intPos {
		pos += 7
		content = content[0:pos] + "HTMLEscInt" + content[pos+10:]
	}

	for _, pos := range strPos {
		pos += 7
		content = content[0:pos] + "HTMLEscStr" + content[pos+10:]
	}

	return replacedCount > 0, content
}
