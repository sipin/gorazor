package cases

import (
	"bytes"
	"io"
	"strings"
)

func Header() string {
	var _b strings.Builder
	WriteHeader(&_b)
	return _b.String()
}

func WriteHeader(_buffer io.StringWriter) {
	_buffer.WriteString("<div>Page Header</div>")

}
