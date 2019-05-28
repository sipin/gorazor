package cases

import (
	"bytes"
	"github.com/sipin/gorazor/gorazor"
	"io"
	"strings"
)

func Email() string {
	var _b strings.Builder
	WriteEmail(&_b)
	return _b.String()
}

func WriteEmail(_buffer io.StringWriter) {
	_buffer.WriteString("<span>rememberingsteve@apple.com ")
	_buffer.WriteString(gorazor.HTMLEscape(username))
	_buffer.WriteString("</span>")

}
