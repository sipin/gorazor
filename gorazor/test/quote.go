package cases

import (
	"bytes"
	"io"
	"strings"
)

func Quote() string {
	var _b strings.Builder
	WriteQuote(&_b)
	return _b.String()
}

func WriteQuote(_buffer io.StringWriter) {
	_buffer.WriteString("<html>'text'</html>")

}
