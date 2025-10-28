## Warning

> The gorazor syntax mentioned in this article has been modified in version 2.0
> This document was translated with a LLM

# Intro 
Not all scenarios are suitable for frontend-backend separation. There will still be situations where the server needs to return rendered pages — especially for pages requiring SEO.
The built-in template engine in Go is not easy to use and has often been criticized — for example [here](http://weibo.com/1729408273/B1ZU32ynE) or [here](http://www.yinwang.org/blog-cn/2014/04/18/golang).
I once tried to use [mustache](http://mustache.github.io) to alleviate this issue. Unfortunately, the [Go implementation](https://github.com/hoisie/mustache) provided on the official mustache website is almost half-broken — it doesn’t even support dot notation. So, I quickly [forked a version](https://github.com/Wuvist/mustache/commits/master) and added support for dot notation and lambda functions.
It worked, but this implementation of mustache relies heavily on reflection, which made me very concerned about performance. Furthermore, mustache’s main advantage is language interoperability, so its feature set is quite limited — it’s not a particularly convenient templating system.
Since I hadn’t yet started full development, I wanted to spend some time clearing these technical obstacles in advance.
How do we fix Go’s template weakness? Implement or port a new one.
I’ve always been fascinated by web template engines, having tried countless varieties across multiple languages.
Two in particular left a strong impression on me:

* Django
* Razor

## Django template engine

The Django template engine is powerful and completely separates logic from presentation. I used to love this “pure” template style.
I once believed that templates should never contain code — that absolute separation of code and template was essential. Developers couldn’t be trusted, and technical barriers should prevent them from doing something foolish.
But as my understanding evolved, I grew increasingly tired of “purity.” I came to believe that developers *can* be trusted — that they deserve choice. It’s not about prohibition, but about discipline.
In the end, it’s a matter of skill and maturity. When I was a novice, I wanted hard restrictions to prevent myself from making mistakes. But now that I’m not a novice anymore, restrictions only frustrate me. When I choose to break a rule, I have a good reason — and I know exactly what I’m doing.

For example, let’s say we need to format a date/time within a template. With Django, you’d have to first implement a date formatting function in Python, then register it as a filter or tag. But now, I see that as unnecessary overhead.

> “Entities must not be multiplied beyond necessity.”

Formatting a date is a business requirement — the formatting function itself is necessary. But defining filters/tags just to expose that function to the template adds redundant complexity.
If the time function’s interface changes, all those extra layers must change too. This increase in entities multiplies complexity.
Layer upon layer — when the bottom layer changes, everything above must change too. That’s annoying. When ASP.NET first came out, Web Forms used a similar multi-tier model, and it turned out to be painful in practice. That’s why Microsoft later switched to ASP.NET MVC.
Fewer layers are better; fewer entities are better. While adding an abstraction layer *can* solve many problems, not every problem *needs* one.
Templates like Django’s require inventing a whole new template language for logic control.
That template language *is* an abstraction layer. Go’s built-in template engine works the same way.
Whether through reflection or interpretation, these template languages are eventually mapped to the host language. Anything expressible in the template language can also be expressed directly in the host language.

So why not use the native language directly? Why force developers to learn another mini-language?
Is this abstraction layer of the template language really worth it?

## Razor template engine

ASP.NET MVC’s default Razor template is remarkably clean. No need for symbols like `<?, <%, {{` that only add keystrokes and visual noise.
Logic control? Just use the native language — for example:

```go
@if(totalMessage == 1) {
    <p>@Userame has 1 message</p>
} else {
    <p>@Userame has @totalMessage messages</p>
}
```

It only needs an `@` to embed code. The compiler automatically determines what’s code and what’s template.
When I first saw Razor, I worried it might confuse code and markup. But after using it, I almost never encountered that issue. In fact, most developers who’ve used Razor [say it’s delightful](http://www.zhihu.com/question/19973649).
It’s been years since I’ve done .NET work, but Razor remains, in my opinion, the most elegantly designed template system — even after trying alternatives like Jade in Node.js.
**Port Razor to Go! Port Razor to Go! Port Razor to Go!**
That thought wouldn’t leave my head. But I kept suppressing it — startup life is about delivering business value, not tech for tech’s sake!
Still, Razor is *so simple*! And Microsoft’s .NET implementation is [fully open source](https://github.com/ASP-NET-MVC/aspnetwebstack/tree/master/src/System.Web.Razor). It already supports multiple languages. I’d just need to modify it to support Go — maybe a day’s work? But I haven’t written C# in years… and I’d have to install Visual/Xamarin Studio… worth it?
So, as usual — “When in doubt, ask Google.” That’s how I discovered [vash](https://github.com/kirbysayshi/vash): a JavaScript implementation of Razor. Only ~2000 lines of code. Hacking it to output Go should absolutely be a one-day task!
And that’s how [gorazor](https://github.com/sipin/gorazor) was born.
Building the **proof of concept** was quick; maturing it takes time, but that can evolve alongside business development.
I didn’t even read vash’s full implementation at first — I just found its compiler generator and started hacking Go syntax in. If it worked, I moved on. *Quick & dirty*, but effective — as long as the design is clear.


# gorazor

## Code Generation

The first design decision: **code generation**.
Templates are just templates — `gorazor` compiles them into Go source code. The app imports and calls these generated Go files directly, not the templates.
So when templates change, you must recompile. In Go, that’s not a problem — Go is compiled anyway, and it’s blazing fast.
Use a file watcher: when a template file changes, kill the process, recompile, restart. Often, this all happens between `save + S / alt + tab / F5`.
Go may need compilation, but with this workflow, developing in Go feels just as smooth as hot-reload in Python.

## Binding to Functions
The second design decision: generated templates are **functions**, not **classes**.
A template `msg/inbox.gohtml` compiles to an `Inbox` function inside the `msg` package — not an `Inbox` struct.
So in the app, you can simply:

```go
import tpl/msg
msg.Inbox(...)
```

Instead of:

```go
inbox := &msg.Inbox{...}
inbox.render()
```

The generated code uses the template’s directory as its namespace and the filename as its function name.
This approach makes template nesting in gorazor naturally seamless.

## Template Declaration
Go is statically typed and compiled — this is not a burden, but a strength.
Using something like `page.ViewData["XXX"]` (a generic map/dictionary) to pass data into templates is bad practice.
Templates should explicitly declare what data types they accept — letting the compiler catch mistakes early.
In gorazor, the **first code block** must declare the template’s imports and parameters, e.g.:

```go
@{
	import (
		. "myapp/models"
	)
	var totalMessage int
	var u *User
}
```

This means the template accepts two parameters: an `int` and a `*User` from `myapp/models`.
If no data/model is needed, just leave the first block empty.
These declarations are converted to function parameters:

```go
func Inbox(totalMessage int, u *User) string {
   ....
}
```
Thus, editors/IDEs can offer autocompletion and type hints — no more guessing which parameters a template expects.

## Section / Layout

Function-based templates already solve nesting — making gorazor as powerful as mustache.
Through function nesting, we can also simulate `section` and `layout` functionality.
However, having native support for `layout` and `section` makes life easier.
So, what syntax to use? After some thought, I decided to use a **magic import**, e.g.:

```go
@{
	import (
		. "myapp/models"
		"tpl/layout/base"
	)
	var totalMessage int
	var u *User
}
```

The namespace segment `layout` receives special treatment. `tpl/layout/base` is automatically resolved to `tpl/layout.Base()`. (Yes, Go’s capitalization rules are annoying!)
I considered adding a keyword like:

```go
@{
	import (
		. "myapp/models"
	)
	Layout "tpl/layout/base.gohtml"
	var totalMessage int
	var u *User
}
```

But I ultimately rejected that idea.

1. I didn’t want to introduce a new keyword like `Layout` — it’s not part of Go’s syntax.
2. I didn’t want to refer to template files by filesystem paths like `tpl/layout/base.gohtml`. Multiple template sets could exist, sharing layouts — better to handle everything through Go namespaces.

Layouts can also import other layouts, supporting multi-level nesting naturally via function calls.
For performance, I considered making generated functions write to a `writer` instead of returning `string`. But since layouts need to build sections into strings first, this offers no clear advantage.
Layouts must also be declared, for example:

```go
@{
	var body string
	var title string
	var sideMenu string
}
<!DOCTYPE html>
<html>
	<head>
	<meta charset="utf-8" />
	@title
	</head>
<body>
	<div>@body</div>
	<div>@sideMenu</div>
</body>
</html>
```

Layout templates accept parameters as `string`s, in order. The **first must be `body`**.
When a page template specifies a layout via import, we assume it can’t directly read that layout’s source. But the page template can define `sections`:

```go
@{
	import (
		. "kp/models"
		"tpl/layout/base"
	)
	var totalMessage int
	var u *User
}

@if(totalMessage == 1) {
	<p>@Userame has 1 message</p>
} else {
	<p>@Userame has @totalMessage messages</p>
}

@section title {
	<title>@u.Name's Homepage</title>
}

@section sideMenu {

}
```

Any content outside a `section` is treated as the `body`. Section order determines the argument order passed to the layout — though this can be error-prone.
Gorazor could be smarter — it can infer the layout’s file path from its namespace, parse it, and automatically handle defaults.
For instance, if the layout defines a section that isn’t provided by the page, it could use a default value:

```go
@{
	var body string
	var title string
	var sideMenu string
}
<!DOCTYPE html>
<html>
	<head>
	@title
	</head>
<body>
	<div>@body</div>
	@if sideMenu == "" {
		<div>Default side menu</div>
	} else {
		<div>@sideMenu</div>
	}
</body>
</html>
```

*P.S.* Handling `@title` can be tricky — writing `<title>@title</title>` in the layout makes the page side more cumbersome.

# Using GoRazor
Just get the package 
```bash
go get github.com/sipin/gorazor
```
and have a look at the [tutorial](https://github.com/sipin/gorazor/blob/main/docs/tutorial.md) or the [example](https://github.com/sipin/gorazor/tree/main/examples) 
