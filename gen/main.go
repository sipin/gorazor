package tpl

import (
	"bytes"
	"gorazor"
	. "kp/models"
	"tpl/helper"
)

func Main(totalMessage int, u *User) string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n\n")
	_buffer.WriteString((helper.Header()))

	_buffer.WriteString("\n")
	_buffer.WriteString((helper.Msg(u)))

	_buffer.WriteString("\n\n")
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
	_buffer.WriteString("\n\n\n")

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

	_buffer.WriteString("\n\n")

	switch totalMessage {
	case 1:

		_buffer.WriteString("<p>")
		_buffer.WriteString(gorazor.HTMLEscape(u.name))

		_buffer.WriteString(" has 1  message</p>")

	case 2:

		_buffer.WriteString("<p>")
		_buffer.WriteString(gorazor.HTMLEscape(u.name))

		_buffer.WriteString(" has 2 messages</p>")

	default:

		_buffer.WriteString("<p>")
		_buffer.WriteString(gorazor.HTMLEscape(u.name))

		_buffer.WriteString(" has no messages</p>")

	}

	_buffer.WriteString("\n\n")
	_buffer.WriteString((helper.Footer()))
	return _buffer.String()
}
