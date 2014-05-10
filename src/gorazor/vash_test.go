package gorazor

import (
	"testing"
	"path/filepath"
	"strings"
	"os"
	_ "syscall"
	"fmt"
)

func TestLexer(t *testing.T) {
	text := "case do func var switch "
	lex := &Lexer{text, Tests}
	res, err := lex.Scan()
	if err != nil { t.Error(err) }
	if len(res) != 10 { t.Error("token number") }
	for i, x := range res {
		if i % 2 == 0 && x.Type != KEYWORD {
			t.Error("KEYWORD", x)
		}
	}

        text = "case casex do do3 func func_ var var+ "
	lex  = &Lexer{text, Tests}
	res, err = lex.Scan()
	if err != nil { t.Error(err) }
	if len(res) != 17 { t.Error(err) }
	for i, x := range res {
		if i == 0 || i == 4 || i == 8 || i == 12 || i == 14 {
			if x.Type != KEYWORD { t.Error("KEYWORD") }
		} else if x.Type == KEYWORD {
			t.Error("Should NOT KEYWORD", x)
		}
	}
}


func TestGenerate(t *testing.T) {
	casedir, _ := filepath.Abs(filepath.Dir("./test/"))
	fmt.Println("case:", casedir)
	cmpsdir, _ := filepath.Abs(filepath.Dir("./cmps/"))
	fmt.Println("cmps:", cmpsdir)

        visit := func(path string, info os.FileInfo, err error) error {
                if !info.IsDir() {
			output := filepath.Join(cmpsdir,
				strings.Replace(filepath.Base(path), ".gohtml", ".log", 1))
			cmp    := filepath.Join(cmpsdir,
				strings.Replace(filepath.Base(path), ".gohtml", ".go", 1))
			GenFile(path, output)
			if !exists(cmp) || !exists(output) {
				t.Error("MISMATCH")
			} else {
				//TODO: compare
				t.Log("PASS")
			}
                }
                return nil
        }
        err := filepath.Walk(casedir, visit)
	if err !=  nil {
		t.Error("walk")
	}
}
