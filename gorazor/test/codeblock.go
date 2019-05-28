package cases

import (
	"bytes"
	"io"
	"strings"
)

func Codeblock() string {
	var _b strings.Builder
	WriteCodeblock(&_b)
	return _b.String()
}

func WriteCodeblock(_buffer io.StringWriter) {

}
