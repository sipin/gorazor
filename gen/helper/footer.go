package helper

import (
	"bytes"
)

func Footer() string {
	var _buffer bytes.Buffer

	_buffer.WriteString("<div>copyright 2014</div>\n</body>\n</html>")
	return _buffer.String()
}
