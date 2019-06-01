package cases

import (
	"bytes"
	"io"
	"strings"
)

func Header() string {
	var _b strings.Builder
	RenderHeader(&_b)
	return _b.String()
}

func RenderHeader(_buffer io.StringWriter) {
	_buffer.WriteString("<div>Page Header</div>")

}
