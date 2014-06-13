package cases

import (
	"bytes"
	"cases/layout"
	"github.com/sipin/gorazor/gorazor"
	. "kp/models"
	"tpl/helper"
)

func Argsbug(totalMessage int, u *User) string {
	var _buffer bytes.Buffer

	messages := []string{}
	_buffer.WriteString("\n\n<p>")
	_buffer.WriteString(gorazor.HTMLEscape(gorazor.Itoa(args(messages...))))
	_buffer.WriteString("</p>")

	return layout.Args(_buffer.String())
}
