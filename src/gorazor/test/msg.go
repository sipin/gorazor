package cases

import (
	"bytes"
	"gorazor"
)

func Msg() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("<h4>Hello ")
	_buffer.WriteString(username)
	_buffer.WriteString("</h4>\n\n<div>")
	_buffer.WriteString((u.Intro))
	_buffer.WriteString("</div>\n</div>\n\n")

	return _buffer.String()
}
