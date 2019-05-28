package cases

import (
	"bytes"
	"hello"
	"huhu"
	"io"
	"now"
	"strconv"
	"strings"
	"this"
)

func Import() string {
	var _b strings.Builder
	WriteImport(&_b)
	return _b.String()
}

func WriteImport(_buffer io.StringWriter) {
	_buffer.WriteString("\n\n\n<p>hello</p>")

}
