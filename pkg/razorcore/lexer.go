package razorcore

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	tkUnDef = iota
	tkAt
	tkAssignOperator
	tkAtColon
	tkAtStarClose
	tkAtStarOpen
	tkBackslash
	tkBraceClose
	tkBraceOpen
	tkContent
	tkEmail
	tkEscapedQuote
	tkHardParenClose
	tkHardParenOpen
	tkHTMLTagOpen
	tkHTMLTagClose
	tkHTMLTagVoidClose
	tkIdentifier
	tkKeyword
	tkLogical
	tkNewline
	tkNumericContent
	tkOperator
	tkParenClose
	tkParenOpen
	tkPeriod
	tkSingleQuote
	tkDoubleQuote
	tkTextareaTagClose
	tkTextareaTagOpen
	tkCommentTagOpen
	tkCommentTagClose
	tkWhitespace
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

// TokenMatch store matched token
type TokenMatch struct {
	Type  int
	Text  string
	Regex *regexp.Regexp
}

func rec(reg string) *regexp.Regexp {
	return regexp.MustCompile("^" + reg)
}

// Tests stores TokenMatch list, TokenMatch order is important
var Tests = []TokenMatch{
	{tkEmail, "EMAIL", rec(`([a-zA-Z0-9.%]+@[a-zA-Z0-9.\-]+\.(?:ca|co\.uk|com|edu|net|org))\b`)},
	{tkHTMLTagOpen, "HTML_TAG_OPEN", rec(`(<[a-zA-Z@]+?[^>]*?["a-zA-Z]*>)`)},
	{tkHTMLTagClose, "HTML_TAG_CLOSE", rec(`(</[^>@]+?>)`)},
	{tkHTMLTagVoidClose, "HTML_TAG_VOID_CLOSE", rec(`(\/\s*>)`)},
	{tkKeyword, "KEYWORD", rec(`(case|do|else|section|for|func|goto|if|return|switch|var|with)([^\d\w])`)},
	{tkIdentifier, "IDENTIFIER", rec(`([_$a-zA-Z][_$a-zA-Z0-9]*(\.\.\.)?)`)}, //need verify
	{tkOperator, "OPERATOR", rec(`(==|!=|>>|<<|>=|<=|>|<|\+|-|\/|\*|\^|%|\:|\?)`)},
	{tkEscapedQuote, "ESCAPED_QUOTE", rec(`(\\+['\"])`)},
	{tkNumericContent, "NUMERIC_CONTENT", rec(`([0-9]+)`)},
	{tkContent, "CONTENT", rec(`([^\s})@.]+?)`)},
}

// Token represent a token in code
type Token struct {
	Text    string
	TypeStr string
	Type    int
	Line    int
	Pos     int
}

// P for print
func (token Token) P() {
	textStr := strings.Replace(token.Text, "\n", "\\n", -1)
	textStr = strings.Replace(textStr, "\t", "\\t", -1)
	fmt.Printf("Token: %-20s Location:(%-2d %-2d) Value: %s\n",
		token.TypeStr, token.Line, token.Pos, textStr)
}

// Lexer for gorazor
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

var runeMatch = map[byte]int{
	'\n': tkNewline,
	' ':  tkWhitespace,
	'\t': tkWhitespace,
	'\f': tkWhitespace,
	'\r': tkWhitespace,
	'(':  tkParenOpen,
	')':  tkParenClose,
	'[':  tkHardParenOpen,
	']':  tkHardParenClose,
	'{':  tkBraceOpen,
	'}':  tkBraceClose,
	'"':  tkDoubleQuote,
	'`':  tkDoubleQuote,
	'\'': tkSingleQuote,
	'.':  tkPeriod,
}

var peekNextMatch = map[string]int{
	"*@":          tkAtStarClose,
	"<textarea>":  tkTextareaTagOpen,
	"</textarea>": tkTextareaTagClose,
	"<!--":        tkCommentTagOpen,
	"-->":         tkCommentTagClose,
}

func tryPeekNext(text string) (match string, tokVal int, ok bool) {
	for k, v := range peekNextMatch {
		if peekNext(k, text) {
			return k, v, true
		}
	}
	return
}

// Scan return gorazor doc as list for Token
func (lexer *Lexer) Scan() ([]Token, error) {
	toks := []Token{}
	text := strings.Replace(lexer.Text, "\r\n", "\n", -1)
	text += "\n"
	cur, line, pos := 0, 1, 0
	for cur < len(text) {
		val, left := text[cur], text[cur:]
		var tok Token

		if tokVal, ok := runeMatch[val]; ok {
			tok = makeToken(string(val), tokVal)
		} else if val == '@' {
			if peekNext(string(':'), left[1:]) {
				tok = makeToken("@:", tkAtColon)
			} else if peekNext(string('*'), left[1:]) {
				tok = makeToken("@*", tkAtStarOpen)
			} else {
				tok = makeToken("@", tkAt)
			}
		} else {
			if match, tokVal, ok := tryPeekNext(left); ok {
				tok = makeToken(match, tokVal)
			} else {
				//try rec
				match := false
				for _, m := range lexer.Matches {
					found := m.Regex.FindIndex([]byte(left))
					if found != nil {
						match = true
						tokenVal := left[found[0]:found[1]]
						if m.Type == tkHTMLTagOpen {
							tokenVal = tagClean(tokenVal)
						} else if m.Type == tkKeyword {
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
		if tok.Type == tkNewline {
			line, pos = line+1, 0
		} else {
			pos += len(tok.Text)
		}
	}
	return toks, nil
}
