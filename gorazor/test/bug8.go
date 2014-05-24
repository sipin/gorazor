package cases

import (
	"bytes"
	"github.com/sipin/gorazor/gorazor"
)

func Bug8(l *Locale) string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n<span>")
	_buffer.WriteString(gorazor.HTMLEscape(l.T("for")))
	_buffer.WriteString("</span>")

	return _buffer.String()
}
