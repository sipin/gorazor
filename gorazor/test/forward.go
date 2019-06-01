package cases

import (
	"bytes"
	"cases/layout"
	"io"
	"strings"
)

func Forward(content string, err string) string {
	var _b strings.Builder
	WriteForward(&_b, content, err)
	return _b.String()
}

func WriteForward(_buffer io.StringWriter, content string, err string) {

	_body := func(_buffer io.StringWriter) {

		//hello word
		/* hello this */

	}

	return layout.Base(_buffer, body, nil, nil)
}
