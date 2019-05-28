package cases

import (
	"bytes"
	"github.com/sipin/gorazor/gorazor"
	"io"
	. "kp/models"
	"strings"
)

func Msg(u *User) string {
	var _b strings.Builder
	WriteMsg(&_b, u)
	return _b.String()
}

func WriteMsg(_buffer io.StringWriter, u *User) {

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

}
