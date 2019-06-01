package cases

import (
	"Tpl"
	"bytes"
	"github.com/sipin/gorazor/gorazor"
	"io"
	"strings"
)

func Bug42() string {
	var _b strings.Builder
	RenderBug42(&_b)
	return _b.String()
}

func RenderBug42(_buffer io.StringWriter) {
	_buffer.WriteString("\n<div class=\"container\">\n    ")
	_buffer.WriteString(gorazor.HTMLEscape((Tpl.TplBread([]string{"选择邮寄方式"}))))
	_buffer.WriteString("\n</div>")

}
