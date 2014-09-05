package tpl

import (
	"bytes"
	"github.com/sipin/gorazor/gorazor"
	"time"
)

func Index() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n\n<p>This is Index</p>")
	{
		t := time.Now()
		StrTime := t.Format("2006-01-02 15:04:05")

		_buffer.WriteString("<p>Time now is:  ")
		_buffer.WriteString(gorazor.HTMLEscape(StrTime))
		_buffer.WriteString(" </p>")

	}

	return _buffer.String()
}
