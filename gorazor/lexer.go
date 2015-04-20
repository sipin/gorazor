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
	COMMENT_TAG_OPEN
	COMMENT_TAG_CLOSE
	WHITESPACE
)

var typeStr = [...]string{
	"UNDEF", "AT", "ASSIGN_OPERATOR", "AT_COLON",
	"AT_STAR_CLOSE", "AT_STAR_OPEN", "BACKSLASH",
	"BRACE_CLOSE", "BRACE_OPEN", "CONTENT",
	"EMAIL", "ESCAPED_QUOTE",
	"HARD_PAREN_CLOSE", "HARD_PAREN_OPEN",
	"HTML_TAG_OPEN", "HTML_TAG_CLOSE", "HTML_TAG_VOID_CLOSE",
	"IDENTIFIER", "KEYWORD", "LOGICAL",
	"NEWLINE", "NUMERIC_CONTENT", "OPERATOR",
	"PAREN_CLOSE", "PAREN_OPEN", "PERIOD",
	"SINGLE_QUOTE", "DOUBLE_QUOTE", "TEXT_TAG_CLOSE",
	"TEXT_TAG_OPEN", "COMMENT_TAG_OPEN", "COMMENT_TAG_CLOSE", "WHITESPACE"}

// Option have following options:
//   Debug bool
//   Watch bool
//   NameNotChange bool
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
	TokenMatch{EMAIL, "EMAIL", rec(`([a-zA-Z0-9.%]+@[a-zA-Z0-9.\-]+\.(?:ca|co\.uk|com|edu|net|org))\b`)},
	TokenMatch{HTML_TAG_OPEN, "HTML_TAG_OPEN", rec(`(<[a-zA-Z@]+?[^>]*?["a-zA-Z]*>)`)},
	TokenMatch{HTML_TAG_CLOSE, "HTML_TAG_CLOSE", rec(`(</[^>@]+?>)`)},
	TokenMatch{HTML_TAG_VOID_CLOSE, "HTML_TAG_VOID_CLOSE", rec(`(\/\s*>)`)},
	TokenMatch{KEYWORD, "KEYWORD", rec(`(case|do|else|section|for|func|goto|if|return|switch|var|with)([^\d\w])`)},
	TokenMatch{IDENTIFIER, "IDENTIFIER", rec(`([_$a-zA-Z][_$a-zA-Z0-9]*(\.\.\.)?)`)}, //need verify
	TokenMatch{OPERATOR, "OPERATOR", rec(`(==|!=|>>|<<|>=|<=|>|<|\+|-|\/|\*|\^|%|\:|\?)`)},
	TokenMatch{ESCAPED_QUOTE, "ESCAPED_QUOTE", rec(`(\\+['\"])`)},
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

// Why we need this: Go's regexp DO NOT support lookahead assertion
func regRemoveTail(text string, regs []string) string {
	res := text
	for _, reg := range regs {
		regc := regexp.MustCompile(reg)
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
	pos := len(text) - 1
	for {
		v := text[pos]
		if (v >= 'a' && v <= 'z') ||
			(v >= 'A' && v <= 'Z') {
			break
		} else {
			pos--
		}
	}
	return text[:pos+1]
}

func peekNext(expect string, text string) bool {
	if strings.HasPrefix(text, expect) {
		return true
	}
	return false
}

func makeToken(val string, tokenType int) Token {
	return Token{val, typeStr[tokenType], tokenType, 0, 0}
}

func (lexer *Lexer) Scan() ([]Token, error) {
	toks := []Token{}
	text := strings.Replace(lexer.Text, "\r\n", "\n", -1)
	text = strings.Replace(lexer.Text, "\r", "\n", -1)
	text += "\n"
	cur, line, pos := 0, 0, 0
	for cur < len(text) {
		val, left := text[cur], text[cur:]
		var tok Token
		switch val {
		case '\n':
			tok = makeToken(string(val), NEWLINE)
		case ' ', '\t', '\f', '\r':
			tok = makeToken(string(val), WHITESPACE)
		case '(':
			tok = makeToken(string(val), PAREN_OPEN)
		case ')':
			tok = makeToken(string(val), PAREN_CLOSE)
		case '[':
			tok = makeToken(string(val), HARD_PAREN_OPEN)
		case ']':
			tok = makeToken(string(val), HARD_PAREN_CLOSE)
		case '{':
			tok = makeToken(string(val), BRACE_OPEN)
		case '}':
			tok = makeToken(string(val), BRACE_CLOSE)
		case '"', '`':
			tok = makeToken(string(val), DOUBLE_QUOTE)
		case '\'':
			tok = makeToken(string(val), SINGLE_QUOTE)
		case '.':
			tok = makeToken(string(val), PERIOD)
		case '@':
			if peekNext(string(':'), left[1:]) {
				tok = makeToken("@:", AT_COLON)
			} else if peekNext(string('*'), left[1:]) {
				tok = makeToken("@*", AT_STAR_OPEN)
			} else {
				tok = makeToken("@", AT)
			}
		default:
			if peekNext("*@", left) {
				tok = makeToken("*@", AT_STAR_CLOSE)
			} else if peekNext("<text>", left) {
				tok = makeToken("<text>", TEXT_TAG_OPEN)
			} else if peekNext("</text>", left) {
				tok = makeToken("</text>", TEXT_TAG_CLOSE)
			} else if peekNext("<!--", left) {
				tok = makeToken("<!--", COMMENT_TAG_OPEN)
			} else if peekNext("-->", left) {
				tok = makeToken("-->", COMMENT_TAG_CLOSE)
			} else {
				//try rec
				match := false
				for _, m := range lexer.Matches {
					found := m.Regex.FindIndex([]byte(left))
					if found != nil {
						match = true
						tokenVal := left[found[0]:found[1]]
						if m.Type == HTML_TAG_OPEN {
							tokenVal = tagClean(tokenVal)
						} else if m.Type == KEYWORD {
							tokenVal = keyClean(tokenVal)
						}
						tok = makeToken(tokenVal, m.Type)
						break
					}
				}
				if !match {
					return toks, fmt.Errorf("%d:%d: Illegal character: %s",
						line, pos, string(text[pos]))
				}
			}
		}
		tok.Line, tok.Pos = line, pos
		toks = append(toks, tok)
		cur += len(tok.Text)
		if tok.Type == NEWLINE {
			line, pos = line+1, 0
		} else {
			pos += len(tok.Text)
		}
	}
	return toks, nil
}
