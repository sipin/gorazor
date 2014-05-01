package helper

import (
	"bytes"
)

func Header() string {
	var _buffer bytes.Buffer

	_buffer.WriteString("<html>\n<body>\n")
	return _buffer.String()
}
