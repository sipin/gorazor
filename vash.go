package main

import (
	"os"
	"fmt"
	"regexp"
	"bufio"
	"strings"
)

const (
	AT = iota
	ASSIGN_OPERATOR
	AT_COLON
	AT_STAR_CLOSE
	AT_STAR_OPEN
	BACKSLASH
	BRACE_CLOSE
	BRACE_OPEN
	CONTENT
	EMAIL
	ESCAPED_QUOTE
	FORWARD_SLASH
	FUNCTION
	HARD_PAREN_CLOSE
	HARD_PAREN_OPEN
	HTML_TAG_CLOSE
	THML_TAG_OPEN
	HTML_TAG_VOID_CLOSE
	IDENTIFIER
	KEYWORD
	LOGICAL
	NEWLINE
	NUMERIC_CONTENT
	OPERATOR
	PAREN_CLOSE
	PAREN_OPEN
	PERIOD
	SINGLE_QUOTE
	DOUBLE_QUOTE
	TEXT_TAG_CLOSE
	TEXT_TAG_OPEN
	WHITESPACE
)

type TokenMatch struct {
	Type  int
	Regex string
}

// The order is important
var Tests = []TokenMatch{
	TokenMatch{EMAIL, `([a-zA-Z0-9.%]+@[a-zA-Z0-9.\-]+\.(?:ca|co\.uk|com|edu|net|org))\\b`},
	TokenMatch{AT_STAR_OPEN, `@\*`},
	TokenMatch{AT_STAR_CLOSE, `(\*@)`},
        TokenMatch{AT_COLON, `(@\:)`},
	TokenMatch{AT, `(@)`},
        TokenMatch{PAREN_OPEN, `(\()`},
        TokenMatch{PAREN_CLOSE, `(\))`},
        TokenMatch{HARD_PAREN_OPEN, `(\[)`},
        TokenMatch{HARD_PAREN_CLOSE, `(\])`},
        TokenMatch{BRACE_OPEN, `(\{)`},
        TokenMatch{BRACE_CLOSE, `(\})`},
        TokenMatch{TEXT_TAG_OPEN, `(<text>)`},
	TokenMatch{TEXT_TAG_CLOSE, `(<\/text>)`},
	TokenMatch{HTML_TAG_CLOSE, `(<\/[^>@\\b]+?>)`},
	TokenMatch{HTML_TAG_VOID_CLOSE, `(\/\s*>)`},
        TokenMatch{PERIOD, `(\.)`},
        TokenMatch{WHITESPACE, `(\s)`},
        TokenMatch{FUNCTION, `(function)([\D\W])`},
        TokenMatch{KEYWORD, `(case|do|else|section|for|func|goto|if|return|switch|try|var|while|with)([\D\W])`},
        TokenMatch{IDENTIFIER, `([_$a-zA-Z\\xA0-\\uFFFF][_$a-zA-Z0-9\\xA0-\\uFFFF]*)`},
        TokenMatch{FORWARD_SLASH, `(\/)`},
        TokenMatch{OPERATOR, `(===|!==|==|!==|>>>|<<|>>|>=|<=|>|<|\+|-|\/|\*|\^|%|\:|\?)`},
	TokenMatch{ASSIGN_OPERATOR, `(\|=|\^=|&=|>>>=|>>=|<<=|-=|\+=|%=|\/=|\*=|=)`},
        TokenMatch{LOGICAL, `(&&|\|\||&|\||\^)`},
        TokenMatch{ESCAPED_QUOTE, `(\\+['\"])`},
	TokenMatch{BACKSLASH, `(\\)`},
        TokenMatch{DOUBLE_QUOTE, `(\\")`},
	TokenMatch{SINGLE_QUOTE, `(\')`},
	TokenMatch{NUMERIC_CONTENT, `([0-9]+)`},
        TokenMatch{CONTENT, `([^\s})@.]+?)`},
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
				fmt.Printf("length: %d found[0]: %d found[1]: %d\n",
					length, found[0], found[1])
				line, pos := LineAndPos(text, pos)
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


func test() {
	buffer := "casex case"
	lex := &Lexer{ buffer }
	res, _ := Lex(lex, buffer)
	fmt.Println(res)
}

func main() {
	buf := ""
	file, _ := os.Open("./tpl/layout/base.gohtml")
	reader := bufio.NewReader(file)
	for {
		byte, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		buf += byte
	}
	lex := &Lexer{ buf}
	res, err := Lex(lex, buf)
	fmt.Println("buf:", buf)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	for _, elem := range res {
		fmt.Println(elem)
	}

	test()
}
