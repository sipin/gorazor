package cases

import (
	"bytes"
)

func Quote() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("<html>'text'</html>")

	return _buffer.String()
}
