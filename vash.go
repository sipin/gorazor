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
        TokenMatch{IDENTIFIER, `([_$a-zA-Z][_$a-zA-Z0-9]*)`}, //need verify
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

func TagOpen(text string) (string) {
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
			res = res[found[1]:] //BUG?
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
					tokenVal = TagOpen(tokenVal)
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

//------------------------------ Parser ------------------------------//
const (
	PRG = iota
	MKP
	BLK
	EXP
)

type Ast struct {
	Parent     *Ast
	Children   []interface{}
	Mode       int
	TagName    string
}

func (ast *Ast) addChild(child interface{}) {
	ast.Children = append(ast.Children, child)
}

func (ast *Ast) popChild() {
	l := len(ast.Children)
	if l == 0 {
		return
	}
	ast.Children = ast.Children[:l-1]
}

func (ast *Ast) root() (*Ast) {
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

func(ast *Ast)  beget(mode int, tag string) (*Ast) {
	child := &Ast{ast, []interface{}{}, mode, tag}
	ast.addChild(child)
	return child
}

func (ast *Ast) closest(mode int, tag string) (*Ast) {
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

type Parser struct {
	ast        *Ast
	tokens     []Token
	preTokens  []Token
        curr       Token
        inComment  bool
	saveTextTag bool
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
                return "", err
        }
        found := regc.FindIndex([]byte(text))
        if found != nil {
                return text[found[0]:found[1]], nil
        }
        return "", nil
}

func (parser *Parser) handleMKP(token Token) {
	next  := parser.peekToken(0)
	//nnext := parser.peekToken(1)
	switch token.Type {
	case AT_STAR_OPEN:
		break
	case AT:
		if next != nil {
			switch next.Type {
			case PAREN_OPEN:
			case IDENTIFIER:
				if len(parser.ast.Children) == 0 {
					parser.ast = parser.ast.Parent //BUG
					parser.ast.popChild() //remove empty MKP block
				}
				parser.ast = parser.ast.beget(EXP, "")
				break
			case KEYWORD:
			case FUNCTION:
                        case BRACE_OPEN:      //BLK
				if len(parser.ast.Children) == 0 {
					parser.ast = parser.ast.Parent
					parser.ast.popChild()
				}
				parser.ast = parser.ast.beget(BLK, "")
				break
			case AT:
			case AT_COLON:
				//we want to keep the token, but remove it's special meanning
				next.Type = CONTENT //BUG, modify from a pointer, work?
				parser.ast.addChild(parser.nextToken())
				break
			default:
				parser.ast.addChild(parser.nextToken())
				break
			}
		}
		break
	case TEXT_TAG_OPEN:
	case HTML_TAG_OPEN:
                tagName, _ := regMatch(`/^<([^\/ >]+)/`, token.Text)
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
		break
	case TEXT_TAG_CLOSE:
	case HTML_TAG_CLOSE:
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
		break
	case HTML_TAG_VOID_CLOSE:
		parser.ast.addChild(token)
		parser.ast = parser.ast.Parent
		break
	case BACKSLASH:
		token.Text += "\\"
		parser.ast.addChild(token)
		break
	default:
		parser.ast.addChild(token)
		break
	}
}

func (parser *Parser) handleBLK(token Token) {
	next := parser.peekToken(0)
	switch token.Type {
	case AT:
		if next.Type != AT && (parser.inComment) {
			parser.deferToken(token)
			parser.ast = parser.ast.beget(MKP, "")
		} else {
			next.Type = CONTENT
			parser.ast.addChild(next)
			parser.skipToken()
		}
		break

	case AT_STAR_OPEN:
		//TODO
		break
	case AT_COLON:
		//TODO subparsre
		break
	case TEXT_TAG_OPEN:
	case TEXT_TAG_CLOSE:
	case HTML_TAG_OPEN:
	case HTML_TAG_CLOSE:
		parser.ast = parser.ast.beget(MKP, "")
		parser.deferToken(token)
		break

	case FORWARD_SLASH:
	case SINGLE_QUOTE:
	case DOUBLE_QUOTE:
		if token.Type == FORWARD_SLASH && next != nil && next.Type == FORWARD_SLASH {
			parser.inComment = true
		}
		if !parser.inComment {
			//TODO
		} else {
			parser.ast.addChild(token)
		}
		break
	case NEWLINE:
		if parser.inComment {
			parser.inComment = false
		}
		parser.ast.addChild(token)
		break

	case BRACE_OPEN:
	case PAREN_OPEN:
		//TODO
	default:
		parser.ast.addChild(token)
		break
	}
}


func (parser *Parser) handleEXP(token Token) {
	switch token.Type {
	case KEYWORD:
	case FUNCTION:
		parser.ast = parser.ast.beget(BLK, "")
		parser.deferToken(token)
		break
	case WHITESPACE:
	case LOGICAL:
	case ASSIGN_OPERATOR:
	case OPERATOR:
	case NUMERIC_CONTENT:
		if parser.ast.Parent != nil && parser.ast.Parent.Mode == EXP {
			parser.ast.addChild(token)
		} else {
			parser.ast = parser.ast.Parent
			parser.deferToken(token)
		}
		break;
	case IDENTIFIER:
		parser.ast.addChild(token)
		break
	case SINGLE_QUOTE:
	case DOUBLE_QUOTE:
		//TODO
		break
	case HARD_PAREN_OPEN:
	case PAREN_OPEN:
		prev := parser.prevToken(0)
		next := parser.peekToken(0) //BUG?
		if token.Type == HARD_PAREN_OPEN && next.Type == HARD_PAREN_CLOSE {
                        // likely just [], which is not likely valid outside of EXP
			parser.deferToken(token)
			parser.ast = parser.ast.Parent
			break
		}
		//TODO subParse
		if (prev != nil && prev.Type == AT) || ( next != nil && next.Type == IDENTIFIER) {
			parser.ast = parser.ast.Parent
		}
		break

	case BRACE_OPEN:
		parser.deferToken(token)
		parser.ast = parser.ast.beget(BLK, "")
		break

	case PERIOD:
		next := parser.peekToken(0)
		if next != nil && ( next.Type == IDENTIFIER ||
			next.Type == KEYWORD ||
			next.Type == FUNCTION ||
			next.Type == PERIOD ||
			(parser.ast.Parent != nil && parser.ast.Parent.Mode == EXP)) {
			parser.ast.addChild(token)
		} else {
			parser.ast = parser.ast.Parent
			parser.deferToken(token)
		}
		break
	default:
		if parser.ast.Parent != nil && parser.ast.Parent.Mode != EXP {
			parser.ast = parser.ast.Parent
			parser.deferToken(token)
		} else {
			parser.ast.addChild(token)
		}
		break
	}
}

func (parser *Parser) advanceUntilNot(tokenType int) ([]Token) {
	res := []Token{}
	for idx, token := range parser.tokens {
		if token.Type == tokenType {
			res = append(res, token)
		} else {
			break;
			parser.tokens = parser.tokens[idx:] //BUG?
		}
	}
	return res
}

func (parser *Parser) Run() (err error) {
	if(parser == nil) {
		return
	}
	parser.ast.Mode = PRG
	parser.curr = Token{"UNDEF", UNDEF, 0, 0}
	for {
		parser.preTokens = append(parser.preTokens, parser.curr)
		if (len(parser.tokens) == 0) {
			break
		}
		parser.curr = parser.nextToken()
		if(parser.ast.Mode == PRG) {
			parser.ast = parser.ast.beget(MKP, "")
		}
		parser.curr.P()
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
	buf := ""
	file, _ := os.Open("./tpl/home.gohtml")
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

	parser := &Parser{&Ast{}, res, []Token{}, Token{}, false, false}
	err = parser.Run()

}
