package cases

import (
	"bytes"
	"io"
	"strings"
)

func Brace_bug() string {
	var _b strings.Builder
	WriteBrace_bug(&_b)
	return _b.String()
}

func WriteBrace_bug(_buffer io.StringWriter) {

	isActive := func(name string) {
		if active == name {

			_buffer.WriteString("<li class=\"active\">\n        ")
		} else {

			_buffer.WriteString("<li>\n        ")
		}
	}

}
