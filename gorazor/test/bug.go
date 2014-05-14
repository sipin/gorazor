package cases

import (
	"bytes"
)

func Bug() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("<html>\n  <head>\n    <title>Title</title>\n  </head>\n\n  <body>\n  Body\n  </body>\n</html>\n")

	return _buffer.String()
}
