package cases

import (
	"bytes"
	"cases/layout"
	"github.com/sipin/gorazor/gorazor"
	"io"
	. "kp/models"
	"strings"
	"tpl/helper"
)

func End(totalMessage int, u *User) string {
	var _b strings.Builder
	WriteEnd(&_b, totalMessage, u)
	return _b.String()
}

func WriteEnd(_buffer io.StringWriter, totalMessage int, u *User) {

	_body := func(_buffer io.StringWriter) {
		_buffer.WriteString((helper.Header()))
		_buffer.WriteString((helper.Msg(u)))
		for i := 0; i < 2; i++ {
			if totalMessage > 0 {
				if totalMessage == 1 {

					_buffer.WriteString("<p>")
					_buffer.WriteString(gorazor.HTMLEscape(u.Name))
					_buffer.WriteString(" has 1 message</p>")

				} else {

					_buffer.WriteString("<p>")
					_buffer.WriteString(gorazor.HTMLEscape(u.Name))
					_buffer.WriteString(" has ")
					_buffer.WriteString(gorazor.HTMLEscape(gorazor.Itoa(totalMessage)))
					_buffer.WriteString(" messages</p>")

				}
			} else {

				_buffer.WriteString("<p>")
				_buffer.WriteString(gorazor.HTMLEscape(u.Name))
				_buffer.WriteString(" has no messages</p>")

			}
		}

		for i := 0; i < 2; i++ {
			if totalMessage > 0 {
				if totalMessage == 1 {

					_buffer.WriteString("<p>")
					_buffer.WriteString(gorazor.HTMLEscape(u.Name))
					_buffer.WriteString(" has 1 message</p>")

				} else {

					_buffer.WriteString("<p>")
					_buffer.WriteString(gorazor.HTMLEscape(u.Name))
					_buffer.WriteString(" has ")
					_buffer.WriteString(gorazor.HTMLEscape(gorazor.Itoa(totalMessage)))
					_buffer.WriteString(" messages</p>")

				}
			} else {

				_buffer.WriteString("<p>")
				_buffer.WriteString(gorazor.HTMLEscape(u.Name))
				_buffer.WriteString(" has no messages</p>")

			}
		}

		switch totalMessage {
		case 1:

			_buffer.WriteString("<p>")
			_buffer.WriteString(gorazor.HTMLEscape(u.Name))
			_buffer.WriteString(" has 1  message</p>")

		case 2:

			_buffer.WriteString("<p>")
			_buffer.WriteString(gorazor.HTMLEscape(u.Name))
			_buffer.WriteString(" has 2 messages</p>")

		default:

			_buffer.WriteString("<p>")
			_buffer.WriteString(gorazor.HTMLEscape(u.Name))
			_buffer.WriteString(" has no messages</p>")

		}

		_buffer.WriteString((helper.Footer()))

	}

	title := func(_buffer io.StringWriter) {

		_buffer.WriteString("<title>")
		_buffer.WriteString(gorazor.HTMLEscape(u.Name))
		_buffer.WriteString("'s homepage</title>")

	}

	side := func(_buffer io.StringWriter) {

	}

	return layout.Base(_buffer.String(), title(), "")
}
