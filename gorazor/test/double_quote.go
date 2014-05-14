package cases

import (
	"bytes"
)

func Double_quote() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("<meta charset=\"utf-8\" />")

	return _buffer.String()
}
