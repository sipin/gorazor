package cases

import (
	"bytes"
	"io"
	"strings"
)

func Double_quote() string {
	var _b strings.Builder
	RenderDouble_quote(&_b)
	return _b.String()
}

func RenderDouble_quote(_buffer io.StringWriter) {
	_buffer.WriteString("<meta charset=\"utf-8\" />")

}
