package cases

import (
	"bytes"
	"gorazor"
)

func Layout(body string, title string, side string) string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n<!DOCTYPE html>\n<html>\n<head>\n<meta charset=\"utf-8\" />\n")
	_buffer.WriteString(title)
	_buffer.WriteString("\n</head>\n<body>\n<div>")
	_buffer.WriteString(body)
	_buffer.WriteString("</div>\n<div>")
	_buffer.WriteString(side)
	_buffer.WriteString("</div>\n</body>\n</html>\n")

	return _buffer.String()
}
