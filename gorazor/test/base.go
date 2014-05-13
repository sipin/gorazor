package cases

import (
	"bytes"
	"gorazor"
	"tpl/admin/helper"
)

func Base(body string, title string) string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n<!DOCTYPE html>\n<html>\n<head>\n	<meta charset=\"utf-8\" />\n    <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n	<link rel=\"stylesheet\" href=\"/css/bootstrap.min.css\">\n	<link rel=\"stylesheet\" href=\"/css/dashboard.css\">\n    <!-- HTML5 shim and Respond.js IE8 support of HTML5 elements and media queries -->\n    <!--[if lt IE 9]>\n      <script src=\"https://oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js\"></script>\n      <script src=\"https://oss.maxcdn.com/libs/respond.js/1.4.2/respond.min.js\"></script>\n    <![endif]-->\n	<title>")
	_buffer.WriteString(gorazor.HTMLEscape(title))
	_buffer.WriteString("</title>\n</head>\n<body>\n    <div class=\"navbar navbar-inverse navbar-fixed-top\" role=\"navigation\">\n      <div class=\"container-fluid\">\n        <div class=\"navbar-header\">\n          <button type=\"button\" class=\"navbar-toggle\" data-toggle=\"collapse\" data-target=\".navbar-collapse\">\n            <span class=\"sr-only\">Toggle navigation</span>\n            <span class=\"icon-bar\"></span>\n            <span class=\"icon-bar\"></span>\n            <span class=\"icon-bar\"></span>\n          </button>\n          <a class=\"navbar-brand\" href=\"#\">广东省政法委信息化平台</a>\n        </div>\n        <div class=\"navbar-collapse collapse\">\n          <ul class=\"nav navbar-nav navbar-right\">\n            <li><a href=\"/admin/setting\">设置</a></li>\n            <li><a href=\"/admin/help\">帮助</a></li>\n            <li><a href=\"/admin/logout\">退出</a></li>\n          </ul>\n          <form class=\"navbar-form navbar-right\">\n            <input type=\"text\" class=\"form-control\" placeholder=\"搜索...\">\n          </form>\n        </div>\n      </div>\n    </div>\n\n    <div class=\"container-fluid\">\n      <div class=\"row\">\n        <div class=\"col-sm-3 col-md-2 sidebar\">\n			")
	_buffer.WriteString((helper.Menu()))
	_buffer.WriteString("\n        </div>\n        <div class=\"col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main\">\n          ")
	_buffer.WriteString(gorazor.HTMLEscape(body))
	_buffer.WriteString("\n        </div>\n      </div>\n    </div>\n	<script src=\"/js/jquery.min.js\"></script>\n	<script src=\"/js/bootstrap.min.js\"></script>\n  </body>\n</html>\n\n")

	return _buffer.String()
}
