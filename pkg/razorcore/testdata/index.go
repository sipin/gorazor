// This file is generated by gorazor 1.2.2
// DON'T modified manually
// Should edit source file and re-generate: cases/index.gohtml

package cases

import (
	"cases/layout"
	"github.com/sipin/gorazor/gorazor"
	"io"
	"kp/models"
	"strings"
)

// Index generates cases/index.gohtml
func Index(users []*models.User, total int, limit int, offset int) string {
	var _b strings.Builder
	RenderIndex(&_b, users, total, limit, offset)
	return _b.String()
}

// RenderIndex render cases/index.gohtml
func RenderIndex(_buffer io.StringWriter, users []*models.User, total int, limit int, offset int) {

	_body := func(_buffer io.StringWriter) {
		// Line: 12
		_buffer.WriteString("\n\n<h2 class=\"sub-header\">用户总数：")
		// Line: 14
		_buffer.WriteString(gorazor.HTMLEscape(gorazor.Itoa(total)))
		// Line: 14
		_buffer.WriteString("</h2>\n<div class=\"table-responsive\">\n\t<table class=\"table table-striped\">\n\t\t<thead>\n\t\t\t<tr>\n\t\t\t\t<th>名字</th>\n\t\t\t\t<th>电邮</th>\n\t\t\t\t<th>编辑</th>\n\t\t\t</tr>\n\t\t</thead>\n\t\t<tbody>\n\t\t\t")
		for _, u := range users {

			// Line: 26
			_buffer.WriteString("<tr>\n\t\t\t\t<td>")
			// Line: 27
			_buffer.WriteString(gorazor.HTMLEscape(u.Name))
			// Line: 27
			_buffer.WriteString("</td>\n\t\t\t\t<td>")
			// Line: 28
			_buffer.WriteString(gorazor.HTMLEscape(u.Email))
			// Line: 28
			_buffer.WriteString("</td>\n\t\t\t\t<td><a href=\"/admin/user/edit?id=")
			// Line: 29
			_buffer.WriteString(gorazor.HTMLEscape(u.ID.Hex()))
			// Line: 29
			_buffer.WriteString("\">编辑</a></td>\n\t\t\t</tr>")

		}
		// Line: 31
		_buffer.WriteString("\n\t\t</tbody>\n\t</table>\n</div>")

	}

	_js := func(_buffer io.StringWriter) {

	}

	_title := func(_buffer io.StringWriter) {

		// Line: 41
		_buffer.WriteString("用户管理")

	}

	layout.RenderBase(_buffer, _body, _title, _js)
}