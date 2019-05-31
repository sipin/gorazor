```go
package layout

import (
	"bytes"
	"github.com/sipin/gorazor/gorazor"
	"io"
	"strings"
	"tpl/admin/helper"
)

func Base(body string, title string, js string) string {
	var _b strings.Builder
	WriteBase(&_b, body, title, js)
	return _b.String()
}

func WriteBase(_buffer io.StringWriter, body string, title string, js string) {

	companyName := "深圳思品科技有限公司"

	_buffer.WriteString("\n<!DOCTYPE html>\n<html>\n<head>\n\t<meta charset=\"utf-8\" />\n    <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n\t<link rel=\"stylesheet\" href=\"/css/bootstrap.min.css\">\n\t<link rel=\"stylesheet\" href=\"/css/dashboard.css\">\n    <!-- HTML5 shim and Respond.js IE8 support of HTML5 elements and media queries -->\n    <!--[if lt IE 9]>\n      <script src=\"https://oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js\"></script>\n      <script src=\"https://oss.maxcdn.com/libs/respond.js/1.4.2/respond.min.js\"></script>\n    <![endif]-->\n\t<title>")
	_buffer.WriteString((title))

}
```
