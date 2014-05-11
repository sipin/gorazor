package cases

import (
	"bytes"
	"gorazor"
)

func Base() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("<div>")
	_buffer.WriteString(gorazor.HTMLEscape(body))
	_buffer.WriteString("</div>\n<div>")
	_buffer.WriteString(gorazor.HTMLEscape(side))
	_buffer.WriteString("</div>\n\n")

	return _buffer.String()
}
