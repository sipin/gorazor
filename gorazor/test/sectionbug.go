package cases

import (
	"bytes"
	"cases/layout"
	"github.com/sipin/gorazor/gorazor"
	"io"
	"kp/models"
	"strings"
)

func Sectionbug() string {
	var _b strings.Builder
	WriteSectionbug(&_b)
	return _b.String()
}

func WriteSectionbug(_buffer io.StringWriter) {

	_body := func(_buffer io.StringWriter) {

	}

	js := func(_buffer io.StringWriter) {
		for _, jsFile := range ctx.GetJS() {

			_buffer.WriteString("<script src=\"")
			_buffer.WriteString(gorazor.HTMLEscape(jsFile))
			_buffer.WriteString("\"></script>")

		}

	}

	return layout.Base(_buffer, body, nil, js)
}
