package cases

import (
	"bytes"
	"tpl/layout"
)

func Section_comment_bug() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n\n<a>\n    <!-- comment -->\n</a>")

	side := func() string {
		var _buffer bytes.Buffer

		_buffer.WriteString("<!-- comment -->\n    plain text")
		return _buffer.String()
	}

	return layout.Base(_buffer.String(), side())
}
