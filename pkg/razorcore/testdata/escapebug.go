// This file is generated by gorazor 1.2.2
// DON'T modified manually
// Should edit source file and re-generate: cases/escapebug.gohtml

package cases

import (
	"io"
	"strings"
)

// Escapebug generates cases/escapebug.gohtml
func Escapebug() string {
	var _b strings.Builder
	RenderEscapebug(&_b)
	return _b.String()
}

// RenderEscapebug render cases/escapebug.gohtml
func RenderEscapebug(_buffer io.StringWriter) {
	// Line: 1
	_buffer.WriteString("<script type=\"text/javascript\">console.log(\"\\n\");</script>")

}