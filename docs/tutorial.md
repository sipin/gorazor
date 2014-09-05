

## Hello world

Gorazor is a translator from `gohtml` to `go`. For every `gohtml` file will be translated into a Go program with a function declared, which will return a `string` value as HTML output.

For example:

```html
<p>Hello world</p>
```

will be translated into:

```go
package demo

import (
	"bytes"
)

func Hello() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("<p>Hello world</p>")

	return _buffer.String()
}
```

Note: put hello.gohtml in a directory, the directory name will be used as package name in Go program.

## Routes

Let's use framework [web.go](github.com/hoisie/web) as example,

Firstly let's install web.go as below:
```shell
mkdir src
export GOPATH=$PWD
go get github.com/hoisie/web
```

the `Hello world` example in web.go is main.go:

```go
package main

import (
    "github.com/hoisie/web"
)

func hello(val string) string {
    return "hello " + val
}

func main() {
    web.Get("/(.*)", hello)
    web.Run("0.0.0.0:9999")
}

```

use command: `go run src/main.go` to start web server, and localhost:9999 will ready for use. For more details please have a look at: [web.go toturial](http://webgo.io/).

Then we make a new directory named `tpl` in project dir, and write an `index.gohtml` in it.

```html
<p>This is Index</p>
```

and then use : `gorazor tpl src/tpl` will generate `Go` files into `src/tpl`.
and then modify main.go:

```go
package main

import (
	"tpl"
	"github.com/hoisie/web"
)

func init_web() {
	web.Get("/index", tpl.Index)
}

func hello(val string) string {
	return "hello " + val
}

func main() {
	init_web()
	web.Get("/(.*)", hello)
	web.Run("0.0.0.0:9999")
}

```

## Code sections

In `gohtml` you may insert `Go` code snippet, like this:
```html
@{
  import (
    "time"
  )
}

<p>This is Index</p>

@{
   t := time.Now()
   StrTime := t.Format("2006-01-02 15:04:05")
   <p>Time now is:  @StrTime </p>
}
```

For more details syntax please refer to [Web programming using the Razor syntax](http://www.asp.net/web-pages/tutorials/basics/2-introduction-to-asp-net-web-programming-using-the-razor-syntax).
And you may also add `javascript` code in `gohtml`, where ctx is `var ctx *web.Context`, details please refer to [sipin/web](http://github.com/sipin/web).


```javascript
@section js {
ctx.AddJS("/assets/js/moment.js")
ctx.AddJS("/assets/js/bootstrap-datetimepicker.js")
<script type="text/javascript">
  jQuery(document).ready(function($) {
    $(".datetimepicker").datetimepicker({
      format: "YYYY-MM-DD HH:mm:ss",
    });
  });
</script>
ctx.AddJS("/assets/js/bootstrap-multiselect.js")
<script>
  $(document).ready(function() {
    $('.multiselect').multiselect({
      enableFiltering: true,
      buttonWidth: '170px',
      maxHeight: 200
    });
  });
</script>
}
```
