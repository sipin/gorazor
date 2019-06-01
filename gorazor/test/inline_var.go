package cases

import (
	"bytes"
	"github.com/sipin/gorazor/gorazor"
	"github.com/sunfmin/gorazortests/models"
	"io"
	"strings"
)

func Inline_var() string {
	var _b strings.Builder
	RenderInline_var(&_b)
	return _b.String()
}

func RenderInline_var(_buffer io.StringWriter) {
	_buffer.WriteString("\n\n<body>")
	_buffer.WriteString(gorazor.HTMLEscape(Hello("Felix Sun", "h1", 30, &models.Author{"Van", 20}, 10)))
	_buffer.WriteString("\n</body>")

}
