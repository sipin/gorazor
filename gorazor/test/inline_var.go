package cases

import (
	"bytes"
	"github.com/sipin/gorazor/gorazor"
	"github.com/sunfmin/gorazortests/models"
)

func Inline_var() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n\n<body>")
	_buffer.WriteString(gorazor.HTMLEscape(Hello("Felix Sun", "h1", 30, &models.Author{"Van", 20}, 10)))
	_buffer.WriteString("\n</body>")

	return _buffer.String()
}
