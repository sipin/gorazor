package cases

import (
	"bytes"
	"github.com/sipin/gorazor/gorazor"
	"kp/models"
	"tpl/admin/layout"
)

func Sectionbug() string {
	var _buffer bytes.Buffer

	js := func() string {
		var _buffer bytes.Buffer
		for _, jsFile := range ctx.GetJS() {

			_buffer.WriteString("<script src=\"")
			_buffer.WriteString(gorazor.HTMLEscape(jsFile))
			_buffer.WriteString("\"></script>")

		}
		return _buffer.String()
	}

	return layout.Base(_buffer.String(), js())
}
