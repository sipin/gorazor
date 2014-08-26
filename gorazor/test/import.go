package cases

import (
	"bytes"
	"now"
	"strconv"
)

func Import() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n\n\n<p>hello</p>")

	return _buffer.String()
}
