package cases

import (
	"bytes"
	"io"
	"strings"
)

func Footer() string {
	var _b strings.Builder
	WriteFooter(&_b)
	return _b.String()
}

func WriteFooter(_buffer io.StringWriter) {
	_buffer.WriteString("<div>copyright 2014</div>")

}
