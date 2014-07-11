package cases

import (
	"bytes"
)

func Escapebug() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("<script type=\"text/javascript\">console.log(\"\\n\");</script>")

	return _buffer.String()
}
