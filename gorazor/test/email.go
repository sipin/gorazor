package cases

import (
	"bytes"
	"github.com/sipin/gorazor/gorazor"
	"io"
	"strings"
)

func Email() string {
	var _b strings.Builder
	RenderEmail(&_b)
	return _b.String()
}

func RenderEmail(_buffer io.StringWriter) {
	_buffer.WriteString("<span>rememberingsteve@apple.com ")
	_buffer.WriteString(gorazor.HTMLEscape(username))
	_buffer.WriteString("</span>")

}
