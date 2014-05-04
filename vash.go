package main

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	EMAIL = iota
	AT_STAR_OPEN
	AT_STAR_CLOSE
)

type TokenMatch struct {
	Type  int
	Regex string
}

var Tests = []TokenMatch{
	TokenMatch{EMAIL, "([a-zA-Z0-9.%]+@[a-zA-Z0-9.\\-]+\\.(?:ca|co\\.uk|com|edu|net|org))\\b"},
	TokenMatch{AT_STAR_OPEN, "@\\*"},
	TokenMatch{AT_STAR_CLOSE, "(\\*@)"},
}

type Token struct {
	Text  string
	Type  int
	Occ   int
	Start int
}

type Lexer struct {
	File string
}

func LineAndPos(src string, pos int) (int, int) {
	lines := strings.Count(src[:pos], "\n")
	p := pos - strings.LastIndex(src[:pos], "\n")
	return lines, p
}

func Lex(lexer *Lexer, text string) ([]Token, error) {
	pos := 0
	toks := []Token{}
	cache := []*regexp.Regexp{}
	for pos < len(text) {
		left := text[pos:]
		match := false
		length := 0
		for idx, test := range Tests {
			pattern, id := test.Regex, test.Type
			if len(cache) < idx+1 {
				reg, err := regexp.Compile("^" + pattern)
				if err != nil {
					panic(err)
				}
				cache = append(cache, reg)
			}

			regexp := cache[idx]
			found := regexp.FindIndex([]byte(left))
			if found != nil {
				match = true
				length = found[1] - found[0]
				line, pos := LineAndPos(text, found[0])
				tok := Token{left[found[0]:found[1]], id, line, pos}
				toks = append(toks, tok)
				break
			}
		}
		if !match {
			err_line, err_pos := LineAndPos(text, pos)
			return toks, fmt.Errorf("%d:%d: Illegal character: %s",
				err_line, err_pos, string(text[pos]))
		} else {
			pos += length
		}
	}
	return toks, nil
}

func main() {

}
