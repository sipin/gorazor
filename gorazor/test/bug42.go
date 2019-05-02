package cases

import (
	"Tpl"
	"bytes"
	"github.com/sipin/gorazor/gorazor"
)

func Bug42() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n<div class=\"container\">\n    ")
	_buffer.WriteString(gorazor.HTMLEscape((Tpl.TplBread([]string{"选择邮寄方式"}))))
	_buffer.WriteString("\n</div>")

	return _buffer.String()
}
