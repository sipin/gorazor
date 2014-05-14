package cases

import (
	"bytes"
	"github.com/sipin/gorazor/gorazor"
	"kp/models"
	"tpl/admin/layout"
)

func Edit(u *models.User) string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n<div style=\"width: 500px\">\n<form role=\"form\">\n  <div class=\"form-group\">\n    <label for=\"exampleInputEmail1\">名字</label>\n    <input type=\"email\" class=\"form-control\" id=\"exampleInputEmail1\" placeholder=\"Enter email\" value=\"")
	_buffer.WriteString(gorazor.HTMLEscape(u.Name))
	_buffer.WriteString("\">\n  </div>\n  <div class=\"form-group\">\n    <label for=\"exampleInputPassword1\">电邮</label>\n    <input type=\"email\" class=\"form-control\" id=\"exampleInputPassword1\" placeholder=\"电邮\" value=\"")
	_buffer.WriteString(gorazor.HTMLEscape(u.Email))
	_buffer.WriteString("\">\n  </div>\n  <button type=\"submit\" class=\"btn btn-primary\">保存</button>\n  <a href=\"/admin/user\" class=\"btn btn-default pull-right\">返回</a>\n</form>\n</div>")

	title := func() string {
		var _buffer bytes.Buffer

		_buffer.WriteString("用户管理")

		return _buffer.String()
	}

	return layout.Base(_buffer.String(), title())
}
