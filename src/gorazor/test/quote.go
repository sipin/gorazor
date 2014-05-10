package cases

import (
	"bytes"
	"gorazor"
)

func Quote() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("<html>'text'</html>\n\n")

	return _buffer.String()
}
