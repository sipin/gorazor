package cases

import (
	"bytes"
	"gorazor"
)

func Var(totalMessage int) string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n\n")

	return _buffer.String()
}
