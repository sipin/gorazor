package main

import (
	"io"
	"strings"
)

func WriteBase(_buffer io.StringWriter, body func(_buffer io.StringWriter), title func(_buffer io.StringWriter)) {

	companyName := "深圳思品科技有限公司"

	_buffer.WriteString("<title>")
	_buffer.WriteString(companyName)
	body(_buffer)
	// _buffer.WriteString(title)
	_buffer.WriteString("<title>")

}

func main() {
	var b strings.Builder
	t := func(_b io.StringWriter) {
		_b.WriteString(" -  bingo")
	}

	WriteBase(&b, t, nil)
	println(b.String())
}
