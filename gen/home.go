package now

import (
	"bytes"
	"tpl/layout"
)

func Home(totalMessage int) string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n\n<p>hello</p>")
	return layout.Base(_buffer.String())
}
