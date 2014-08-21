package cases

import (
	"bytes"
)

func Section_comment_bug() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n\n<a>\n    <!-- comment -->\n</a>")

	side := func() string {
		var _buffer bytes.Buffer

		_buffer.WriteString("<!-- comment -->\n    plain text")
		return _buffer.String()
	}

	return _buffer.String(), side()
}
