package cases

import (
	"bytes"
	"github.com/sipin/gorazor/gorazor"
)

func Keyword() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("BLK(<span>rememberingsteve@apple.com ")
	_buffer.WriteString(gorazor.HTMLEscape(username))
	_buffer.WriteString("</span>)BLK")

	return _buffer.String()
}
