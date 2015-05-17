package cases

import (
	"bytes"
	"cases/layout"
	"github.com/sipin/gorazor/gorazor"
)

func Add(content string, err string) string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n\n<link rel=\"stylesheet\" href=\"/css/bootstrap-datetimepicker.css\">\n\n<style>\n.row {\n\tmargin-top: 10px;\n}\n</style>\n\n<h2>日程登记</h2>\n\n<div class=\"container-fluid\">\n\t<form method=\"POST\" action=\"\">\n\t<div class=\"row\" >\n\t\t<p class=\"bg-danger\">")
	_buffer.WriteString(gorazor.HTMLEscape(err))
	_buffer.WriteString("</p>\n\t</div>\n\n\t<div class=\"row\">\n\t内容:\n\t<input type='text' class=\"form-control\" name=\"content\" value=\"")
	_buffer.WriteString(gorazor.HTMLEscape(content))
	_buffer.WriteString("\"/>\n\t</div>\n\t\n\t<div class=\"row\">\n\t开始时间:\n\t<input type='text' class=\"datetimepicker form-control\" name=\"startTime\"/>\n\t</div>\n\t\n\t<div class=\"row\">\n\t结束时间:\n\t<input type='text' class=\"datetimepicker form-control\" name=\"endTime\"/>\n\t</div>\n\n\t<div class=\"row\">\n\t日程指派:\n\t<select name=\"appoint\">\n\t\t<option>cheney</option>\n\t\t<option>wuvist</option>\n\t</select>\n\t</div>\n\t\n\t<div class=\"row\">\n\t<input style=\"float:right\" type=\"submit\" value=\"保存\" class=\"btn btn-primary\"/>\n\t</div>\n\t</form>\n</div>")

	title := func() string {
		var _buffer bytes.Buffer

		_buffer.WriteString("管理后台 - 添加日程")

		return _buffer.String()
	}

	js := func() string {
		var _buffer bytes.Buffer

		_buffer.WriteString("<script src=\"/js/moment.js\"></script>")

		_buffer.WriteString("<script src=\"/js/bootstrap-datetimepicker.js\"></script>")

		_buffer.WriteString("<script type=\"text/javascript\">\n\t$(function () {\n\t\t$(\".datetimepicker\").datetimepicker({\n\t\t\tformat: \"YYYY-MM-DD HH:mm\",\n\t\t\tdefaultDate: \"2014-05-01 00:00\",\n\t\t})\n\t});\n</script>")

		return _buffer.String()
	}

	return layout.Base(_buffer.String(), title(), js())
}
