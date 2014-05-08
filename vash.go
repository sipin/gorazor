package main

import (
	"os"
	"io"
	"fmt"
	"regexp"
	"bytes"
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
	Text  string
	Regex *regexp.Regexp
}

func rec(reg string) (*regexp.Regexp) {
	res, err := regexp.Compile("^" + reg)
	if err != nil {
		panic(err)
	}
	return res
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
        TokenMatch{HTML_TAG_CLOSE, "HTML_TAG_CLOSE", rec(`(<\/[^>@\\b]+?>)`)},
        TokenMatch{HTML_TAG_VOID_CLOSE, "HTML_TAG_VOID_CLOSE", rec(`(\/\s*>)`)},
        TokenMatch{PERIOD, "PERIOD", rec(`(\.)`)},
        TokenMatch{NEWLINE, "NEWLINE", rec(`(\n)`)},
        TokenMatch{WHITESPACE, "WHITESPACE", rec(`(\s)`)},
        TokenMatch{FUNCTION, "FUNCTION", rec(`(function)([\D\W])`)},
        TokenMatch{KEYWORD, "KEYWORD", rec(`(case|do|else|section|for|func|goto|if|return|switch|try|var|while|with)([\D\W])`)},
        TokenMatch{IDENTIFIER, "IDENTIFIER", rec(`([_$a-zA-Z][_$a-zA-Z0-9]*)`)}, //need verify
        TokenMatch{FORWARD_SLASH, "FORWARD_SLASH", rec(`(\/)`)},
        TokenMatch{OPERATOR, "OPERATOR", rec(`(===|!==|==|!==|>>>|<<|>>|>=|<=|>|<|\+|-|\/|\*|\^|%|\:|\?)`)},
        TokenMatch{ASSIGN_OPERATOR, "ASSIGN_OPERATOR", rec(`(\|=|\^=|&=|>>>=|>>=|<<=|-=|\+=|%=|\/=|\*=|=)`)},
        TokenMatch{LOGICAL, "LOGICAL", rec(`(&&|\|\||&|\||\^)`)},
        TokenMatch{ESCAPED_QUOTE, "ESCAPED_QUOTE", rec(`(\\+['\"])`)},
        TokenMatch{BACKSLASH, "BACKSLASH", rec(`(\\)`)},
        TokenMatch{DOUBLE_QUOTE, "DOUBLE_QUOTE", rec(`(\\")`)},
        TokenMatch{SINGLE_QUOTE, "SINGLE_QUOTE", rec(`(\')`)},
        TokenMatch{NUMERIC_CONTENT, "NUMERIC_CONTENT", rec(`([0-9]+)`)},
        TokenMatch{CONTENT, "CONTENT", rec(`([^\s})@.]+?)`)},
}

type Token struct {
	Text string
	TypeStr string
	Type int
	Line int
	Pos  int
}

func (token Token) P() {
	textStr := token.Text
	if textStr == "\n" {
		textStr = "\\n"
	}
	fmt.Printf("Token: %-20s Location:(%-2d %-2d) Value: %s\n",
		token.TypeStr, token.Line, token.Pos, textStr)
}

type Lexer struct {
	Text  string
	Matches []TokenMatch
}

func LineAndPos(src string, pos int) (int, int) {
	lines := strings.Count(src[:pos], "\n")
	p := pos - strings.LastIndex(src[:pos], "\n")
	return lines, p
}

func TagOpen(text string) string {
	regs := []string{
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
			res = res[found[1]:] //BUG?
		}
	}
	return res
}

func (lexer *Lexer) Scan() ([]Token, error) {
	pos := 0
	toks := []Token{}
	text := strings.Replace(lexer.Text, "\r\n", "\n", -1)
	text = strings.Replace(lexer.Text, "\r", "\n", -1)
	for pos < len(text) {
		left := text[pos:]
		match := false
		length := 0
		for _, m := range lexer.Matches {
			found := m.Regex.FindIndex([]byte(left))
			if found != nil {
				match = true
				line, pos := LineAndPos(text, pos)
				tokenVal := left[found[0]:found[1]]
				if m.Type == HTML_TAG_OPEN {
					tokenVal = TagOpen(tokenVal)
				}
				tok := Token{tokenVal, m.Text, m.Type, line, pos}
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

//------------------------------ Parser ------------------------------//
const (
	UNK = iota
	PRG
	MKP
	BLK
	EXP
)

var PAIRS = map[int]int{
	AT_STAR_OPEN:    AT_STAR_CLOSE,
	BRACE_OPEN:      BRACE_CLOSE,
	DOUBLE_QUOTE:    DOUBLE_QUOTE,
	HARD_PAREN_OPEN: HARD_PAREN_CLOSE,
	PAREN_OPEN:      PAREN_CLOSE,
	SINGLE_QUOTE:    SINGLE_QUOTE,
	AT_COLON:        NEWLINE,
	FORWARD_SLASH:   FORWARD_SLASH,
}


type Ast struct {
	Parent   *Ast
	Children []interface{}
	Mode     int
	TagName  string
}

func (ast *Ast) ModeStr() string{
	switch ast.Mode {
	case PRG: return "PROGRAM"
	case MKP: return "MARKUP"
	case BLK: return "BLOCK"
	case EXP: return "EXP"
	default: return "UNDEF"
	}
	return "UNDEF"
}

func (ast *Ast) addChild(child interface{}) {
	ast.Children = append(ast.Children, child)
	ast.check()
}

func (ast *Ast) addChildren(children []Token) { //BUG?
	for _, c := range children {
		ast.addChild(c)
	}
}

func (ast *Ast) addAst(_ast *Ast) {
	fmt.Println("add ast:", ast.Mode, PRG)
	if ast.Mode != PRG {
		fmt.Println("hxxxxxxxxxxxxxxx")
		_ast.debug(0)
		ast.addChild(_ast)
	} else {
		for _, c := range _ast.Children {
			ast.addChild(c)
		}
	}
}

func (ast *Ast) popChild() {
	l := len(ast.Children)
	if l == 0 {
		return
	}
	ast.Children = ast.Children[:l-1]
}

func (ast *Ast) root() *Ast {
	p := ast
	pp := ast.Parent
	for {
		if p == pp || pp == nil {
			return p
		}
		b := pp
		pp = p.Parent
		p = b
	}
	return nil
}

func (ast *Ast) check() {
	if len(ast.Children) >= 100000 {
		panic("Maximum number of elements exceeded.")
	}
}

func (ast *Ast) beget(mode int, tag string) *Ast {
	child := &Ast{ast, []interface{}{}, mode, tag}
	ast.addChild(child)
	return child
}

func (ast *Ast) closest(mode int, tag string) *Ast {
	p := ast
	for {
		if p.TagName != tag && p.Parent != nil {
			p = p.Parent
		} else {
			break
		}
	}
	return p
}
func (ast *Ast) debug(depth int) {
	for i := 0; i < depth; i++ {
		fmt.Printf("%c", '-')
	}
	fmt.Printf("{")
	fmt.Printf("TagName: %s Mode: %s Children: %d\n", ast.TagName, ast.ModeStr(), len(ast.Children))
	for _, a := range ast.Children {
		if _, ok := a.(*Ast); ok {
			b := (*Ast)(a.(*Ast))
			b.debug(depth+1)
		} else {
			aa := (Token)(a.(Token))
                        for i := 0; i < depth+1; i++ {
                                fmt.Printf("%c", '-')
			}
			aa.P()
		}
	}
        for i := 0; i < depth; i++ {
                fmt.Printf("%c", '-')
        }
	fmt.Println("}")
}

type Parser struct {
	ast         *Ast
	tokens      []Token
	preTokens   []Token
	curr        Token
	inComment   bool
	saveTextTag bool
	initMode    int
}

func (parser *Parser) prevToken(idx int) (*Token) {
	l := len(parser.preTokens)
	if l < idx + 1 {
		return nil
	}
	return &(parser.preTokens[l - 1 - idx])
}

func (parser *Parser) deferToken(token Token) {
	parser.tokens = append(parser.tokens, token)
}

func (parser *Parser) peekToken(idx int) (*Token) {
        if len(parser.tokens) <= idx  {
		return nil
        }
        return &(parser.tokens[idx])
}

func (parser *Parser) nextToken() (Token) {
        t := parser.peekToken(0)
        if t != nil {
                parser.tokens = parser.tokens[1:]
        }
        return *t
}

func (parser *Parser) skipToken() {
	parser.tokens = parser.tokens[1:]
}

func regMatch(reg string, text string) (string, error) {
	regc, err := regexp.Compile(reg)
	if err != nil {
		panic(err)
		return "", err
	}
	found := regc.FindIndex([]byte(text))
	if found != nil {
		return text[found[0]:found[1]], nil
	}
	return "", nil
}

func (parser *Parser) advanceUntilNot(tokenType int) []Token {
	res := []Token{}
	for {
		t := parser.peekToken(0)
		if t != nil && t.Type == tokenType {
			res = append(res, parser.nextToken())
		} else {
			break
		}
	}
	return res
}

func (parser *Parser) advanceUntil(token Token, start, end, startEsc, endEsc int) []Token {
	var prev *Token = nil
	next := &token
	res := []Token{}
	nstart := 0
	nend := 0
	for {
		if next.Type == start {
			if (prev != nil && prev.Type != startEsc && start != end) || prev == nil {
				nstart++
			} else if start == end && prev.Type != startEsc {
				nend++
			}
		} else if next.Type == end {
			nend++
			if prev != nil && prev.Type == endEsc {
				nend--
			}
		}
		res = append(res, *next)
		if nstart == nend {
			break
		}
		prev = next
		next = parser.peekToken(0)
		if next == nil {
			panic("UNMATCHED")
		}
		parser.nextToken()
	}
	return res
}

func (parser *Parser) subParse(token Token, modeOpen int, includeDelim bool) {
	subTokens := parser.advanceUntil(token, token.Type, PAIRS[token.Type], -1, AT)
	subTokens = subTokens[1:]
	closer := subTokens[len(subTokens)-1]
	subTokens = subTokens[:len(subTokens)-1]
	if !includeDelim {
		parser.ast.addChild(token)

	}
	fmt.Println("fuck now: ", modeOpen)
        _parser := &Parser{&Ast{}, subTokens, []Token{}, Token{}, false, false, modeOpen}
	_parser.Run()
	if includeDelim {
		_parser.ast.Children = append([]interface{}{token}, _parser.ast.Children...)
		_parser.ast.addChild(closer)
	}
	_parser.ast.debug(0)
	parser.ast.addAst(_parser.ast)
	if !includeDelim {
		parser.ast.addChild(closer)
	}
}

func (parser *Parser) handleMKP(token Token) {
	next  := parser.peekToken(0)
	//nnext := parser.peekToken(1)
	switch token.Type {
	case AT_STAR_OPEN:
		parser.advanceUntil(token, AT_STAR_OPEN, AT_STAR_CLOSE, AT, AT)

	case AT:
		if next != nil {
			switch next.Type {
			case PAREN_OPEN, IDENTIFIER:
				if len(parser.ast.Children) == 0 {
					parser.ast = parser.ast.Parent //BUG
					parser.ast.popChild() //remove empty MKP block
				}
				parser.ast = parser.ast.beget(EXP, "")

			case KEYWORD, FUNCTION, BRACE_OPEN: //BLK
				if len(parser.ast.Children) == 0 {
					parser.ast = parser.ast.Parent
					parser.ast.popChild()
				}
				parser.ast = parser.ast.beget(BLK, "")

			case AT, AT_COLON:
				//we want to keep the token, but remove it's special meanning
				next.Type = CONTENT //BUG, modify from a pointer, work?
				parser.ast.addChild(parser.nextToken())
			default:
				parser.ast.addChild(parser.nextToken())
			}
		}

	case TEXT_TAG_OPEN, HTML_TAG_OPEN:
		tagName, _ := regMatch(`(?i)(^<([^\/ >]+))`, token.Text)
		tagName = strings.Replace(tagName, "<", "", -1)
		//TODO
		if parser.ast.TagName != "" {
			parser.ast = parser.ast.beget(MKP, tagName)
		} else {
			parser.ast.TagName = tagName
		}
		if token.Type == HTML_TAG_OPEN || parser.saveTextTag {
			parser.ast.addChild(token)
		}

	case TEXT_TAG_CLOSE, HTML_TAG_CLOSE:
		tagName, _ := regMatch(`(?i)^<\/([^>]+)`, token.Text)
		tagName = strings.Replace(tagName, "</", "", -1)
		//TODO
		opener := parser.ast.closest(MKP, tagName)
		if opener.TagName != tagName { //unmatched
		} else {
			parser.ast = opener
		}
		if token.Type == HTML_TAG_CLOSE || parser.saveTextTag {
			parser.ast.addChild(token)
		}
		if parser.ast.Parent != nil && parser.ast.Parent.Mode == BLK {
			parser.ast = parser.ast.Parent
		}

	case HTML_TAG_VOID_CLOSE:
		parser.ast.addChild(token)
		parser.ast = parser.ast.Parent

	case BACKSLASH:
		token.Text += "\\"
		parser.ast.addChild(token)
	default:
		parser.ast.addChild(token)
	}
}

func (parser *Parser) handleBLK(token Token) {
	next := parser.peekToken(0)
	ty   := token.Type
	switch ty {
	case AT:
		if next.Type != AT && (parser.inComment) {
			parser.deferToken(token)
			parser.ast = parser.ast.beget(MKP, "")
		} else {
			next.Type = CONTENT
			parser.ast.addChild(next)
			parser.skipToken()
		}

	case AT_STAR_OPEN:
        	parser.advanceUntil(token, AT_STAR_OPEN, AT_STAR_CLOSE, AT, AT)

	case AT_COLON:
                //TODO subparsre
                parser.subParse(token, MKP, true)

	case TEXT_TAG_OPEN, TEXT_TAG_CLOSE, HTML_TAG_OPEN, HTML_TAG_CLOSE:
                parser.ast = parser.ast.beget(MKP, "")
		parser.deferToken(token)

	case FORWARD_SLASH, SINGLE_QUOTE, DOUBLE_QUOTE:
                if token.Type == FORWARD_SLASH && next != nil && next.Type == FORWARD_SLASH {
			parser.inComment = true
		}
		if !parser.inComment {
			subTokens := parser.advanceUntil(token, token.Type,
				PAIRS[token.Type],
				BACKSLASH,
				BACKSLASH)
			for idx, _ := range subTokens {
				if subTokens[idx].Type == AT {
					subTokens[idx].Type = CONTENT
				}
			}
			parser.ast.addChildren(subTokens)
		} else {
			parser.ast.addChild(token)
		}

	case NEWLINE:
                if parser.inComment {
			parser.inComment = false
		}
		parser.ast.addChild(token)

	case BRACE_OPEN, PAREN_OPEN:
		subMode := BLK
		if false && token.Type == BRACE_OPEN {  //TODO
			subMode = MKP
		}
		parser.subParse(token, subMode, false)
		subTokens := parser.advanceUntilNot(WHITESPACE)
		next := parser.peekToken(0)
		if next != nil && next.Type != KEYWORD &&
			next.Type != FUNCTION && next.Type != BRACE_OPEN &&
			next.Type != PAREN_OPEN {
			parser.tokens = append(parser.tokens, subTokens...)
			parser.ast = parser.ast.Parent //BUG?
		} else {
			parser.ast.addChildren(subTokens)
		}
	default:
		parser.ast.addChild(token)
	}
}


func (parser *Parser) handleEXP(token Token) {
	switch token.Type {
	case KEYWORD, FUNCTION:
		parser.ast = parser.ast.beget(BLK, "")
		parser.deferToken(token)

	case WHITESPACE, LOGICAL, ASSIGN_OPERATOR, OPERATOR, NUMERIC_CONTENT:
		if parser.ast.Parent != nil && parser.ast.Parent.Mode == EXP {
			parser.ast.addChild(token)
		} else {
			parser.ast = parser.ast.Parent
			parser.deferToken(token)
		}
	case IDENTIFIER:
		parser.ast.addChild(token)

	case SINGLE_QUOTE, DOUBLE_QUOTE:
		//TODO
		if parser.ast.Parent != nil && parser.ast.Parent.Mode == EXP {
			subTokens := parser.advanceUntil(token, token.Type,
				                         PAIRS[token.Type], BACKSLASH, BACKSLASH)
			parser.ast.addChildren(subTokens)
		} else {
			parser.ast = parser.ast.Parent
			parser.deferToken(token)
		}

	case HARD_PAREN_OPEN, PAREN_OPEN:
		prev := parser.prevToken(0)
		next := parser.peekToken(0) //BUG?
		if token.Type == HARD_PAREN_OPEN && next.Type == HARD_PAREN_CLOSE {
			// likely just [], which is not likely valid outside of EXP
			parser.deferToken(token)
			parser.ast = parser.ast.Parent
			break //BUG?
		}
		parser.subParse(token, EXP, false)
		if (prev != nil && prev.Type == AT) || (next != nil && next.Type == IDENTIFIER) {
			parser.ast = parser.ast.Parent
		}

	case BRACE_OPEN:
		parser.deferToken(token)
		parser.ast = parser.ast.beget(BLK, "")

	case PERIOD:
		next := parser.peekToken(0)
		if next != nil && (next.Type == IDENTIFIER ||
			next.Type == KEYWORD ||
			next.Type == FUNCTION ||
			next.Type == PERIOD ||
			(parser.ast.Parent != nil && parser.ast.Parent.Mode == EXP)) {
			parser.ast.addChild(token)
		} else {
			parser.ast = parser.ast.Parent
			parser.deferToken(token)
		}
	default:
		if parser.ast.Parent != nil && parser.ast.Parent.Mode != EXP {
			parser.ast = parser.ast.Parent
			parser.deferToken(token)
		} else {
			parser.ast.addChild(token)
		}
	}
}

func (parser *Parser) Run() (err error) {
	for _, t := range parser.tokens {
		t.P()
	}
        fmt.Println("-----------------------------------------")
	parser.curr = Token{"UNDEF", "UNDEF", UNDEF, 0, 0}
	fmt.Println("ast now: ", parser.ast)
	parser.ast.Mode = PRG
	for {
		parser.preTokens = append(parser.preTokens, parser.curr)
		if len(parser.tokens) == 0 {
			break
		}
		parser.curr = parser.nextToken()
		if parser.ast.Mode == PRG {
			init := parser.initMode
			if init == UNK {
				init = MKP
			}
			parser.ast = parser.ast.beget(init, "")
			if parser.initMode == EXP {
				parser.ast = parser.ast.beget(EXP, "")
			}
		}
		switch parser.ast.Mode {
		case MKP:
			parser.handleMKP(parser.curr)
		case BLK:
			parser.handleBLK(parser.curr)
		case EXP:
			parser.handleEXP(parser.curr)
		}
	}

	parser.ast = parser.ast.root()
	fmt.Println("-----------------------------------------")

	parser.ast.debug(0)
	return nil
}

//------------------------------ Compiler ------------------------------ //

// func test() {
// 	buffer := "casex case"
// 	lex := &Lexer{ buffer }
// 	res, _ := Lex(lex, buffer)
// 	fmt.Println(res)
// }

func main() {
	buf := bytes.NewBuffer(nil)
	f , err := os.Open("./now/codeblock.gohtml")
	//f, err := os.Open("./tpl/home.gohtml")
	if err != nil {
		panic(err)
	}
	io.Copy(buf, f)
	f.Close()
	text := string(buf.Bytes())
	lex := &Lexer{text, Tests}
	res, err := lex.Scan()
	fmt.Println("buf:", text)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	//for _, elem := range res {
	//elem.P()
	//fmt.Println(elem)
	//}

	parser := &Parser{&Ast{}, res, []Token{}, Token{}, false, false, UNK}
	err = parser.Run()

}
