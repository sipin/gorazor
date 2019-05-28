package cases

import (
	"bytes"
	"cases/layout"
	"github.com/sipin/gorazor/gorazor"
	"io"
	. "kp/models"
	"strings"
	"tpl/helper"
)

func Argsbug(totalMessage int, u *User) string {
	var _b strings.Builder
	WriteArgsbug(&_b, totalMessage, u)
	return _b.String()
}

func WriteArgsbug(_buffer io.StringWriter, totalMessage int, u *User) {

	messages := []string{}

	_buffer.WriteString("\n\n<p>")
	_buffer.WriteString(gorazor.HTMLEscape(gorazor.Itoa(args(messages...))))
	_buffer.WriteString("</p>")

	return layout.Args(_buffer.String())
}
