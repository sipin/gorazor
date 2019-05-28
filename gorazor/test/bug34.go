package cases

import (
	"bytes"
	"io"
	"strings"
)

func Bug34() string {
	var _b strings.Builder
	WriteBug34(&_b)
	return _b.String()
}

func WriteBug34(_buffer io.StringWriter) {
	_buffer.WriteString("value=\\\"<?= h(aabasdf\\Admin\\Document::$asdf) ?>\\\"/>\\n")

}
