package cases

import (
	"bytes"
	"github.com/sipin/gorazor/gorazor"
	. "kp/models"
)

func Msg(u *User) string {
	var _buffer bytes.Buffer

	getName := func(u *User) string {
		return "(" + u.Name + ")"
	}

	var username string
	if u.Email != "" {
		username = getName(u) + "(" + u.Email + ")"
	}
	_buffer.WriteString("\n<div class=\"welcome\">\n<h4>Hello ")
	_buffer.WriteString(gorazor.HTMLEscape(username))
	_buffer.WriteString("</h4>\n\n<div>")
	_buffer.WriteString((u.Intro))
	_buffer.WriteString("</div>\n</div>")

	return _buffer.String()
}
