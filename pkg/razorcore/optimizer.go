package razorcore

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

	type replacement struct {
		offset  int
		oldName string
		newName string
	}
	var replacements []replacement

	// traverse all tokens
	ast.Inspect(node, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.CallExpr:
			switch t2 := t.Fun.(type) {
			case *ast.SelectorExpr:
				ident, ok := t2.X.(*ast.Ident)
				if !ok || ident.Name != "gorazor" || len(t.Args) == 0 {
					return true
				}
				typ := info.Types[t.Args[0]]
				if typ.Type == nil {
					return true
				}
				argType := typ.Type.String()
				var newName string
				switch t2.Sel.Name {
				case "HTMLEscape":
					if argType == "int" {
						newName = "HTMLEscInt"
					} else if argType == "string" {
						newName = "HTMLEscStr"
					}
				case "URLEscape":
					if argType == "string" {
						newName = "URLEscStr"
					}
				case "URLQueryEscape":
					if argType == "string" {
						newName = "URLQueryEscStr"
					}
				case "JSEscape":
					if argType == "string" {
						newName = "JSEscStr"
					}
				case "JSAttrEscape":
					if argType == "string" {
						newName = "JSAttrEscStr"
					}
				}
				if newName != "" {
					offset := fset.Position(t2.Sel.Pos()).Offset
					replacements = append(replacements, replacement{
						offset:  offset,
						oldName: t2.Sel.Name,
						newName: newName,
					})
				}
			}
		}
		return true
	})

	// Replace in reverse order so any offset changes would not affect earlier replacements.
	for i := len(replacements) - 1; i >= 0; i-- {
		r := replacements[i]
		if r.offset >= 0 && r.offset+len(r.oldName) <= len(content) && content[r.offset:r.offset+len(r.oldName)] == r.oldName {
			content = content[:r.offset] + r.newName + content[r.offset+len(r.oldName):]
		}
	}

	return len(replacements) > 0, content
}
