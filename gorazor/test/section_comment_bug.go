package cases

import (
	"bytes"
	"io"
	"strings"
	"tpl/layout"
)

func Section_comment_bug() string {
	var _b strings.Builder
	WriteSection_comment_bug(&_b)
	return _b.String()
}

func WriteSection_comment_bug(_buffer io.StringWriter) {
	_buffer.WriteString("\n\n<a>\n    <!-- comment -->\n</a>")

	side := func() string {
		var _buffer bytes.Buffer

		_buffer.WriteString("<!-- comment -->\n    plain text")
		return _buffer.String()
	}

	return layout.Base(_buffer.String(), side())
}
