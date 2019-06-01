package cases

import (
	"bytes"
	"cases/layout"
	"io"
	"strings"
)

func Section_comment_bug() string {
	var _b strings.Builder
	WriteSection_comment_bug(&_b)
	return _b.String()
}

func WriteSection_comment_bug(_buffer io.StringWriter) {

	_body := func(_buffer io.StringWriter) {
		_buffer.WriteString("\n\n<a>\n    <!-- comment -->\n</a>")

	}

	side := func(_buffer io.StringWriter) {

		_buffer.WriteString("<!-- comment -->\n    plain text")

	}

	return layout.Base(_buffer, body, nil, nil)
}
