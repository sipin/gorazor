// This file is generated by gorazor 1.2.2
// DON'T modified manually
// Should edit source file and re-generate: cases/textarea.gohtml

package cases

import (
	"github.com/sipin/gorazor/gorazor"
	"io"
	"strings"
)

// Textarea generates cases/textarea.gohtml
func Textarea(count int) string {
	var _b strings.Builder
	RenderTextarea(&_b, count)
	return _b.String()
}

// RenderTextarea render cases/textarea.gohtml
func RenderTextarea(_buffer io.StringWriter, count int) {
	// Line: 3
	_buffer.WriteString("\n<html>\n<body>\n<textarea rows=\"4\" cols=\"50\">\n        At w3schools.com ")
	// Line: 7
	_buffer.WriteString(gorazor.HTMLEscape(count))
	// Line: 7
	_buffer.WriteString(" you will learn how to make a website.\n  We offer free tutorials in all web development technologies.\n\n\n\n  At w3schools.com ")
	// Line: 12
	_buffer.WriteString(gorazor.HTMLEscape(count))
	// Line: 12
	_buffer.WriteString(" you will learn\n  how to make a website.\n  We offer free tutorials in all web development technologies.\n\n</body>\n</html>")

}