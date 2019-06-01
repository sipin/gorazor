package cases

import (
	"bytes"
	"github.com/sipin/gorazor/gorazor"
	"io"
	"strings"
)

func Bug8(l *Locale) string {
	var _b strings.Builder
	RenderBug8(&_b, l)
	return _b.String()
}

func RenderBug8(_buffer io.StringWriter, l *Locale) {
	_buffer.WriteString("\n<span>")
	_buffer.WriteString(gorazor.HTMLEscape(l.T("for")))
	_buffer.WriteString("</span>")

}
