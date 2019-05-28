package cases

import (
	"bytes"
	"io"
	"strings"
)

func Escapebug() string {
	var _b strings.Builder
	WriteEscapebug(&_b)
	return _b.String()
}

func WriteEscapebug(_buffer io.StringWriter) {
	_buffer.WriteString("<script type=\"text/javascript\">console.log(\"\\n\");</script>")

}
