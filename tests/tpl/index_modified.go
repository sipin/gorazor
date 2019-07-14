// This file is manually modified to ensure exact output as quicktemplate
// Source file for re-generate: tests/tpl/index.gohtml

package tpl

import (
	"io"
	"strings"

	"github.com/sipin/gorazor/gorazor"
	"github.com/sipin/gorazor/tests/data"
)

// IndexModified generates tests/tpl/index.gohtml
func IndexModified(rows []data.BenchRow) string {
	var _b strings.Builder
	RenderIndex(&_b, rows)
	return _b.String()
}

// RenderIndexModified render tests/tpl/index.gohtml
func RenderIndexModified(_buffer io.StringWriter, rows []data.BenchRow) {
	_buffer.WriteString("<html>\n\t<head><title>test</title></head>\n\t<body>\n\t\t<ul>\n\t\t\n\t\t\t")
	for _, row := range rows {
		if row.Print {

			_buffer.WriteString("\n\t\t\t\t<li>ID=")
			_buffer.WriteString(gorazor.HTMLEscInt(row.ID))
			_buffer.WriteString(", Message=")
			_buffer.WriteString(gorazor.HTMLEscStr(row.Message))
			_buffer.WriteString("</li>\n\t\t\t\n\t\t\n\t\t")

		} else {
			_buffer.WriteString("\t\n\t\t\n\t\t\t")
		}
	}
	_buffer.WriteString("</ul>\n\t</body>\n</html>\n")

}
