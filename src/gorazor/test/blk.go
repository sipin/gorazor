package cases

import (
	"bytes"
	"gorazor"
)

func Blk() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n\n")

	return _buffer.String()
}
