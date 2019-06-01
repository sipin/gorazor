package cases

import (
	"bytes"
	"io"
	"strings"
)

func Quote() string {
	var _b strings.Builder
	RenderQuote(&_b)
	return _b.String()
}

func RenderQuote(_buffer io.StringWriter) {
	_buffer.WriteString("<html>'text'</html>")

}
