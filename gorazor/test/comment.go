package cases

import (
	"bytes"
)

func Comment() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n\n\n\n<p>hello </p>")

	hello

	return _buffer.String()
}
