// This file is generated by gorazor 2.0.1
// DON'T modified manually
// Should edit source file and re-generate: cases/menu_layout.gohtml

package cases

import (
	"cases/layout"
	"github.com/sipin/gorazor/pkg/gorazor"
	"io"
	"kp/models"
	"strings"
)

func Menu_layout(body string, title string, menu string) string {
	var _b strings.Builder

	_body := func(_buffer io.StringWriter) {
		_buffer.WriteString(body)
	}

	_title := func(_buffer io.StringWriter) {
		_buffer.WriteString(title)
	}

	_menu := func(_buffer io.StringWriter) {
		_buffer.WriteString(menu)
	}

	RenderMenu_layout(&_b, _body, _title, _menu)
	return _b.String()
}

func RenderMenu_layout(_buffer io.StringWriter, body func(_buffer io.StringWriter), title func(_buffer io.StringWriter), menu func(_buffer io.StringWriter)) {

	_body := func(_buffer io.StringWriter) {
		_buffer.WriteString("\n\n<div id=\"body\">\n    <div id=\"menu\">")
		_buffer.WriteString(gorazor.HTMLEscape(menu))
		_buffer.WriteString("</div>\n    <div id=\"content\">")
		_buffer.WriteString(gorazor.HTMLEscape(body))
		_buffer.WriteString("</div>\n</div>")

	}

	_title := func(_buffer io.StringWriter) {

		_buffer.WriteString(gorazor.HTMLEscape(title))

	}

	layout.RenderBase(_buffer, _body, _title, nil)
}
