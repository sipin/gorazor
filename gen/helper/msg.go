package helper

import (
	"bytes"
	"gorazor"
	. "kp/models"
)

func Msg(u *User) string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n<html>\n<body>\n\n")

	msg := u.Name + "(" + u.Email + ")"

	_buffer.WriteString("\n\n<h4>Hello ")
	_buffer.WriteString(gorazor.HTMLEscape(msg))

	_buffer.WriteString("</h4>\n</body>\n</html>\n")
	return _buffer.String()
}
