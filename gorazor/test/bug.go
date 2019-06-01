package cases

import (
	"bytes"
	"io"
	"strings"
)

func Bug() string {
	var _b strings.Builder
	RenderBug(&_b)
	return _b.String()
}

func RenderBug(_buffer io.StringWriter) {
	_buffer.WriteString("<html>\n  <head>\n    <title>Title</title>\n  </head>\n\n  <body>\n  Body\n  </body>\n</html>")

}
