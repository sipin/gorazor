package cases

import (
	"bytes"
	"io"
	"strings"
)

func Comment() string {
	var _b strings.Builder
	WriteComment(&_b)
	return _b.String()
}

func WriteComment(_buffer io.StringWriter) {
	_buffer.WriteString("\n\n\n\n<p>hello </p>")

	hello

}
