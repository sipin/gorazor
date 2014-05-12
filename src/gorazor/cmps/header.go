package test

import (
	"bytes"
)

func Header() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("<div>Page Header</div>")
	return _buffer.String()
}
