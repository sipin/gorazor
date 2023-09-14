// This file is generated by gorazor 1.2.2
// DON'T modified manually
// Should edit source file and re-generate: cases/menu_layout.gohtml

package cases

import (
	"cases/layout"
	"github.com/sipin/gorazor/gorazor"
	"io"
	"kp/models"
	"strings"
)

// Menu_layout generates cases/menu_layout.gohtml
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

// RenderMenu_layout render cases/menu_layout.gohtml
func RenderMenu_layout(_buffer io.StringWriter, body func(_buffer io.StringWriter), title func(_buffer io.StringWriter), menu func(_buffer io.StringWriter)) {

	_body := func(_buffer io.StringWriter) {
		// Line: 13
		_buffer.WriteString("\n\n<div id=\"body\">\n    <div id=\"menu\">")
		// Line: 16
		_buffer.WriteString(gorazor.HTMLEscape(menu))
		// Line: 16
		_buffer.WriteString("</div>\n    <div id=\"content\">")
		// Line: 17
		_buffer.WriteString(gorazor.HTMLEscape(body))
		// Line: 17
		_buffer.WriteString("</div>\n</div>")

	}

	_title := func(_buffer io.StringWriter) {

		// Line: 21
		_buffer.WriteString(gorazor.HTMLEscape(title))

	}

	layout.RenderBase(_buffer, _body, _title, nil)
}