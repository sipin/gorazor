package gorazor

import (
	"testing"
	"path/filepath"
	"strings"
	"os"
	"io/ioutil"
	_ "fmt"
)

func TestCap(t *testing.T) {
	if Capitalize("hello") != "Hello" {
		t.Error()
	}
	if Capitalize("0hello") != "0hello" {
		t.Error()
	}
}

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
	cmpsdir, _ := filepath.Abs(filepath.Dir("./cmps/"))

        visit := func(path string, info os.FileInfo, err error) error {
                if !info.IsDir() {
                        name := strings.Replace(filepath.Base(path), ".gohtml", ".go", 1)
                        cmp := filepath.Join(cmpsdir, name)
                        log := filepath.Join(cmpsdir, "_" + name)
			option := Option{}
                	GenFile(path, log, option )
			if !exists(cmp) || !exists(log) {
				t.Error("No Log")
			} else {
				//compare the log file and cmp file
				_cmp, _e1 := ioutil.ReadFile(cmp)
				_log, _e2 := ioutil.ReadFile(log)
				if _e1 != nil || _e2 != nil {
					t.Error("Reading")
				} else if string(_cmp) != string(_log) {
					t.Error("MISMATCH:", cmp, log)
				} else {
					t.Log("PASS")
				}
			}
                }
                return nil
        }
        err := filepath.Walk(casedir, visit)
	if err !=  nil {
		t.Error("walk")
	}
}
