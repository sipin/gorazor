package cases

import (
	"bytes"
)

func Menu() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("<ul class=\"nav nav-sidebar\">\n\t<li role=\"presentation\" class=\"dropdown-header\">用户管理</li>\n\t<li><a href=\"/admin/user\">查看用户</a></li>\n\t<li><a href=\"/admin/user/create\">添加用户</a></li>\n\t<li role=\"presentation\" class=\"divider\"></li>\n\t<li role=\"presentation\" class=\"dropdown-header\">公文管理</li>\n\t<li><a href=\"#\">收文管理</a></li>\n\t<li><a href=\"#\">收文登记</a></li>\n\t<li><a href=\"#\">发送公文</a></li>\n\t<li><a href=\"#\">发文管理</a></li>\n\t<li><a href=\"#\">发文登记</a></li>\n</ul>\n<ul class=\"nav nav-sidebar\">\n\t<li><a href=\"\">领导审批</a></li>\n\t<li><a href=\"\">流程监控</a></li>\n\t<li role=\"presentation\" class=\"divider\"></li>\n\t<li role=\"presentation\" class=\"dropdown-header\">其它</li>\n\t<li><a href=\"\">添加日程</a></li>\n\t<li><a href=\"\">公共通讯录</a></li>\n\t<li><a href=\"\">添加联系人</a></li>\n\t<li><a href=\"\">投票</a></li>\n</ul>")

	return _buffer.String()
}
