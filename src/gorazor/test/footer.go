package cases

import (
	"bytes"
	"gorazor"
)

func Footer() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("<div>copyright 2014</div>\n")

	return _buffer.String()
}
