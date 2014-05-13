package gorazor

import (
	"fmt"
	"regexp"
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

type Option map[string]interface{}

type TokenMatch struct {
	Type  int
	Text  string
	Regex *regexp.Regexp
}

func rec(reg string) *regexp.Regexp {
	return regexp.MustCompile("^" + reg)
}

// The order is important
var Tests = []TokenMatch{
	TokenMatch{EMAIL, "EMAIL", rec(`([a-zA-Z0-9.%]+@[a-zA-Z0-9.\-]+\.(?:ca|co\.uk|com|edu|net|org))\\b`)},
	TokenMatch{AT_STAR_OPEN, "AT_STAR_OPEN", rec(`@\*`)},
	TokenMatch{AT_STAR_CLOSE, "AT_STAR_CLOSE", rec(`(\*@)`)},
	TokenMatch{AT_COLON, "AT_COLON", rec(`(@\:)`)},
	TokenMatch{AT, "AT", rec(`(@)`)},
	TokenMatch{PAREN_OPEN, "PAREN_OPEN", rec(`(\()`)},
	TokenMatch{PAREN_CLOSE, "PAREN_CLOSE", rec(`(\))`)},
	TokenMatch{HARD_PAREN_OPEN, "HARD_PAREN_OPEN", rec(`(\[)`)},
	TokenMatch{HARD_PAREN_CLOSE, "HARD_PAREN_CLOSE", rec(`(\])`)},
	TokenMatch{BRACE_OPEN, "BRACE_OPEN", rec(`(\{)`)},
	TokenMatch{BRACE_CLOSE, "BRACE_CLOSE", rec(`(\})`)},
	TokenMatch{TEXT_TAG_OPEN, "TEXT_TAG_OPEN", rec(`(<text>)`)},
	TokenMatch{TEXT_TAG_CLOSE, "TEXT_TAG_CLOSE", rec(`(<\/text>)`)},
	TokenMatch{HTML_TAG_OPEN, "HTML_TAG_OPEN", rec(`(<[a-zA-Z@]+?[^>]*?["a-zA-Z]*>)`)},
	TokenMatch{HTML_TAG_CLOSE, "HTML_TAG_CLOSE", rec(`(</[^>@]+?>)`)},
	TokenMatch{HTML_TAG_VOID_CLOSE, "HTML_TAG_VOID_CLOSE", rec(`(\/\s*>)`)},
	TokenMatch{PERIOD, "PERIOD", rec(`(\.)`)},
	TokenMatch{NEWLINE, "NEWLINE", rec(`(\n)`)},
	TokenMatch{WHITESPACE, "WHITESPACE", rec(`(\s)`)},
	//TokenMatch{FUNCTION, "FUNCTION", rec(`(function)([^\d\w])`)},
	TokenMatch{KEYWORD, "KEYWORD", rec(`(case|do|else|section|for|func|goto|if|return|switch|var|with)([^\d\w])`)},
	TokenMatch{IDENTIFIER, "IDENTIFIER", rec(`([_$a-zA-Z][_$a-zA-Z0-9]*)`)}, //need verify
	TokenMatch{FORWARD_SLASH, "FORWARD_SLASH", rec(`(\/)`)},
	TokenMatch{OPERATOR, "OPERATOR", rec(`(===|!==|==|!==|>>>|<<|>>|>=|<=|>|<|\+|-|\/|\*|\^|%|\:|\?)`)},
	TokenMatch{ASSIGN_OPERATOR, "ASSIGN_OPERATOR", rec(`(\|=|\^=|&=|>>>=|>>=|<<=|-=|\+=|%=|\/=|\*=|=)`)},
	TokenMatch{LOGICAL, "LOGICAL", rec(`(&&|\|\||&|\||\^)`)},
	TokenMatch{ESCAPED_QUOTE, "ESCAPED_QUOTE", rec(`(\\+['\"])`)},
	TokenMatch{BACKSLASH, "BACKSLASH", rec(`(\\)`)},
	TokenMatch{DOUBLE_QUOTE, "DOUBLE_QUOTE", rec(`(")`)},
	TokenMatch{SINGLE_QUOTE, "SINGLE_QUOTE", rec(`(')`)},
	TokenMatch{NUMERIC_CONTENT, "NUMERIC_CONTENT", rec(`([0-9]+)`)},
	TokenMatch{CONTENT, "CONTENT", rec(`([^\s})@.]+?)`)},
}

type Token struct {
	Text    string
	TypeStr string
	Type    int
	Line    int
	Pos     int
}

func (token Token) P() {
	textStr := strings.Replace(token.Text, "\n", "\\n", -1)
	textStr = strings.Replace(textStr, "\t", "\\t", -1)
	fmt.Printf("Token: %-20s Location:(%-2d %-2d) Value: %s\n",
		token.TypeStr, token.Line, token.Pos, textStr)
}

type Lexer struct {
	Text    string
	Matches []TokenMatch
}

func lineAndPos(src string, pos int) (int, int) {
	lines := strings.Count(src[:pos], "\n")
	p := pos - strings.LastIndex(src[:pos], "\n")
	return lines + 1, p
}

// Why we need this: Go's regexp DO NOT support lookahead assertion
func regRemoveTail(text string, regs []string) string {
	res := text
	for _, reg := range regs {
		regc, err := regexp.Compile(reg)
		if err != nil {
			panic(err)
		}
		found := regc.FindIndex([]byte(res))
		if found != nil {
			res = res[:found[0]] //BUG?
		}
	}
	return res
}

func tagClean(text string) string {
	regs := []string{
		`([a-zA-Z0-9.%]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4})\b`,
		`(@)`,
		`(\/\s*>)`}
	return regRemoveTail(text, regs)
}

func keyClean(text string) string {
	regs := []string{`(\s|\W)`}
	return regRemoveTail(text, regs)
}

func (lexer *Lexer) Scan() ([]Token, error) {
	pos := 0
	toks := []Token{}
	text := strings.Replace(lexer.Text, "\r\n", "\n", -1)
	text = strings.Replace(lexer.Text, "\r", "\n", -1)
	text += "\n"
	for pos < len(text) {
		left := text[pos:]
		match := false
		length := 0
		for _, m := range lexer.Matches {
			found := m.Regex.FindIndex([]byte(left))
			if found != nil {
				match = true
				line, pos := lineAndPos(text, pos)
				tokenVal := left[found[0]:found[1]]

				if m.Type == HTML_TAG_OPEN {
					tokenVal = tagClean(tokenVal)
				} else if m.Type == KEYWORD {
					tokenVal = keyClean(tokenVal)
				}

				tok := Token{tokenVal, m.Text, m.Type, line, pos}
				toks = append(toks, tok)
				length = len(tokenVal)
				break
			}
		}
		if !match {
			err_line, err_pos := lineAndPos(text, pos)
			return toks, fmt.Errorf("%d:%d: Illegal character: %s",
				err_line, err_pos, string(text[pos]))
		}
		pos += length
	}
	return toks, nil
}
