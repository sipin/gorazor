package cases

import (
	"bytes"
)

func Double_quote() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("<meta charset=\"utf-8\" />\n")

	return _buffer.String()
}
