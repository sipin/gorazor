// This file is generated by gorazor 2.0.1
// DON'T modified manually
// Should edit source file and re-generate: cases/bug8.gohtml

package cases

import (
	"github.com/sipin/gorazor/pkg/gorazor"
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
