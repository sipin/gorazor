package tpl

import (
	"bytes"
)

func Index() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n<p>This is Index</p>")

	return _buffer.String()
}
