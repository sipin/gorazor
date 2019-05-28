package cases

import (
	"bytes"
	"cases/layout"
	"github.com/sipin/gorazor/gorazor"
	"io"
	"kp/models"
	"strings"
)

func Index(users []*models.User, total int, limit int, offset int) string {
	var _b strings.Builder
	WriteIndex(&_b, users, total, limit, offset)
	return _b.String()
}

func WriteIndex(_buffer io.StringWriter, users []*models.User, total int, limit int, offset int) {
	_buffer.WriteString("\n\n<h2 class=\"sub-header\">用户总数：")
	_buffer.WriteString(gorazor.HTMLEscape(gorazor.Itoa(total)))
	_buffer.WriteString("</h2>\n<div class=\"table-responsive\">\n\t<table class=\"table table-striped\">\n\t\t<thead>\n\t\t\t<tr>\n\t\t\t\t<th>名字</th>\n\t\t\t\t<th>电邮</th>\n\t\t\t\t<th>编辑</th>\n\t\t\t</tr>\n\t\t</thead>\n\t\t<tbody>\n\t\t\t")
	for _, u := range users {

		_buffer.WriteString("<tr>\n\t\t\t\t<td>")
		_buffer.WriteString(gorazor.HTMLEscape(u.Name))
		_buffer.WriteString("</td>\n\t\t\t\t<td>")
		_buffer.WriteString(gorazor.HTMLEscape(u.Email))
		_buffer.WriteString("</td>\n\t\t\t\t<td><a href=\"/admin/user/edit?id=")
		_buffer.WriteString(gorazor.HTMLEscape(u.ID.Hex()))
		_buffer.WriteString("\">编辑</a></td>\n\t\t\t</tr>")

	}
	_buffer.WriteString("\n\t\t</tbody>\n\t</table>\n</div>")

	js := func() string {
		var _buffer bytes.Buffer
		return _buffer.String()
	}

	title := func() string {
		var _buffer bytes.Buffer

		_buffer.WriteString("用户管理")

		return _buffer.String()
	}

	return layout.Base(_buffer.String(), title(), js())
}
