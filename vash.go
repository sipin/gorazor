package main

import (
	"os"
	"fmt"
	"regexp"
	"bufio"
	"strings"
)

const (
	UNDEF = iota
	AT
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
        HTML_TAG_OPEN
	HTML_TAG_CLOSE
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
        TokenMatch{HTML_TAG_OPEN, `(<[a-zA-Z@]+?[^>]*?["a-zA-Z]*>)`},
	TokenMatch{HTML_TAG_CLOSE, `(<\/[^>@\\b]+?>)`},
	TokenMatch{HTML_TAG_VOID_CLOSE, `(\/\s*>)`},
        TokenMatch{PERIOD, `(\.)`},
	TokenMatch{NEWLINE, `(\n)`},
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

var TokenStr = []string{
	"UNDEF",
        "AT",
        "ASSIGN_OPERATOR",
        "AT_COLON",
        "AT_STAR_CLOSE",
        "AT_STAR_OPEN",
        "BACKSLASH",
        "BRACE_CLOSE",
        "BRACE_OPEN",
        "CONTENT",
        "EMAIL",
        "ESCAPED_QUOTE",
        "FORWARD_SLASH",
        "FUNCTION",
        "HARD_PAREN_CLOSE",
        "HARD_PAREN_OPEN",
        "HTML_TAG_CLOSE",
        "THML_TAG_OPEN",
        "HTML_TAG_VOID_CLOSE",
        "IDENTIFIER",
        "KEYWORD",
        "LOGICAL",
        "NEWLINE",
        "NUMERIC_CONTENT",
        "OPERATOR",
        "PAREN_CLOSE",
        "PAREN_OPEN",
        "PERIOD",
        "SINGLE_QUOTE",
        "DOUBLE_QUOTE",
        "TEXT_TAG_CLOSE",
        "TEXT_TAG_OPEN",
        "WHITESPACE"}

type Token struct {
	Text  string
	Type  int
	Line  int
	Pos   int
}

func (token Token)P() {
	typeStr := TokenStr[token.Type]
	textStr := token.Text
	if textStr == "\n" {
		textStr = "\\n"
	}
	fmt.Printf("Token: %-20s Location:(%-2d %-2d) Value: %s\n",
		    typeStr, token.Line, token.Pos, textStr)
}

type Lexer struct {
	Text    string
        Cache   [](*regexp.Regexp)
	Types   []int
}

func LineAndPos(src string, pos int) (int, int) {
	lines := strings.Count(src[:pos], "\n")
	p := pos - strings.LastIndex(src[:pos], "\n")
	return lines, p
}

func (lexer *Lexer) TagOpen(text string) (string) {
        regs := []string {
		`([a-zA-Z0-9.%]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4})\b`,
		`(@)`,
                `(\/\s*>)`}
	res := text
	for _, reg := range regs {
		regc, err := regexp.Compile(reg)
		if err != nil {
			panic(err)
		}
		found := regc.FindIndex([]byte(text))
		if found != nil {
			res = res[found[1]:]
		}
	}
	return res
}


func (lexer *Lexer) Scan() ([]Token, error) {
	pos := 0
	toks := []Token{}
	for _, test := range Tests {
		reg, err := regexp.Compile("^" + test.Regex)
		if err != nil {
			panic(err)
		}
		lexer.Cache = append(lexer.Cache, reg)
		lexer.Types = append(lexer.Types, test.Type)
	}

        text := strings.Replace(lexer.Text, "\r\n", "\n", -1)
	text = strings.Replace(lexer.Text, "\r", "\n", -1)
        for pos < len(text) {
		left := text[pos:]
		match := false
		length := 0
		for idx, regexp := range lexer.Cache {
			found := regexp.FindIndex([]byte(left))
			if found != nil {
				match = true
				line, pos := LineAndPos(text, pos)
				tokenVal := left[found[0]:found[1]]
				toType := lexer.Types[idx]
				if toType == HTML_TAG_OPEN {
					tokenVal = lexer.TagOpen(tokenVal)
				}
				tok := Token{tokenVal, toType, line, pos}
				toks = append(toks, tok)
				length = len(tokenVal)
				break
			}
		}
		if !match {
			err_line, err_pos := LineAndPos(text, pos)
			return toks, fmt.Errorf("%d:%d: Illegal character: %s",
				     err_line, err_pos, string(text[pos]))
		}
		pos += length
	}
	return toks, nil
}

type Ast struct {
	Parent *Ast
	Children []*Ast
	Value    string
}

func (ast *Ast) AddChild(child *Ast) {
	ast.Children = append(ast.Children, child)
}

type Parser struct {
	root       Ast
	tokens     []Token
	preTokens  []Token
	inComment  bool
        curr       Token
}

func (parser *Parser)Run() (err error) {
	if(parser == nil) {
		return
	}
	parser.curr = Token{"UNDEF", UNDEF, 0, 0}
	for {
		parser.preTokens = append(parser.preTokens, parser.curr)
		if (len(parser.tokens) == 0) {
			break
		}
		parser.curr = parser.tokens[0]
		parser.tokens = parser.tokens[1:]
		fmt.Println("now: ", parser.curr)
	}
	return nil
}

// func test() {
// 	buffer := "casex case"
// 	lex := &Lexer{ buffer }
// 	res, _ := Lex(lex, buffer)
// 	fmt.Println(res)
// }

func main() {
	buf := ""
	file, _ := os.Open("./now/demo.gohtml")
	reader := bufio.NewReader(file)
	for {
		byte, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		buf += byte
	}
	lex := &Lexer{buf, []*regexp.Regexp{}, []int{}}
	res, err := lex.Scan()
	fmt.Println("buf:", buf)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	for _, elem := range res {
		elem.P()
		//fmt.Println(elem)
	}

}
