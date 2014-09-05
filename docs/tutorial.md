

## Hello world

Gorazor is a translator from `gohtml` to `go`. For every `gohtml` file will translated into a Go program with a function declared, which will return a `string` value as HTML output.

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

Let's use framework [web.go](github.com/hoisie/web) as example, the `Hello world` example in web.go is main.go:
```shell
mkdir src
export GOPATH=$PWD
go get github.com/hoisie/web
```

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

We make a new directory named `tpl` in project dir, and write a index.gohtml in it.

```html
<p>This is Index</p>
```

and then use : `gorazor tpl src/tpl` will generate go files into `src/tpl`.
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

## Section
