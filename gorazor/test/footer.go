package cases

import (
	"bytes"
	"io"
	"strings"
)

func Footer() string {
	var _b strings.Builder
	RenderFooter(&_b)
	return _b.String()
}

func RenderFooter(_buffer io.StringWriter) {
	_buffer.WriteString("<div>copyright 2014</div>")

}
