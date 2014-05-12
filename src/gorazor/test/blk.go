package cases

import (
	"bytes"
)

func Blk() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n\n")

	return _buffer.String()
}
