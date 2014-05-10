package cases

import (
	"bytes"
	"gorazor"
)

func Header() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("<div>Page Header</div>\n")

	return _buffer.String()
}
