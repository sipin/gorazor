package cases

import (
	"bytes"
	"cases/layout/args"
	"github.com/sipin/gorazor/gorazor"
	"io"
	. "kp/models"
	"strings"
	"tpl/helper"
)

func Argsbug(totalMessage int, u *User) string {
	var _b strings.Builder
	RenderArgsbug(&_b, totalMessage, u)
	return _b.String()
}

func RenderArgsbug(_buffer io.StringWriter, totalMessage int, u *User) {

	messages := []string{}

	_buffer.WriteString("\n\n<p>")
	_buffer.WriteString(gorazor.HTMLEscape(gorazor.Itoa(args(messages...))))
	_buffer.WriteString("</p>")

}
