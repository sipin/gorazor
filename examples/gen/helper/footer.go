package helper

import (
	"bytes"
)

func Footer() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("<div>copyright 2014</div>")

	return _buffer.String()
}
