### 本文中提到的gorazor语法已经在2.0版本中修改
### 本文中提到的gorazor语法已经在2.0版本中修改
### 本文中提到的gorazor语法已经在2.0版本中修改

前后端分离的并不是所有场景都适用，必然还是会有需要服务器返回页面的场景，特别是需要SEO的页面。

GO内置的模板引擎并不好用，时常为人所诟病，比方说[这里](http://weibo.com/1729408273/B1ZU32ynE)或者[这里](http://www.yinwang.org/blog-cn/2014/04/18/golang)。

我曾经试图适用[mustache](http://mustache.github.io)来缓解这一问题，可惜mustache官方网站的给的[go实现](https://github.com/hoisie/mustache)基本也是个半残废，连dot notation都不支持；当时我果断便[fork出个版本](https://github.com/Wuvist/mustache/commits/master)，增加了dot notation以及lambda等支持。

用是能用，但mustache的这个实现用了大量反射，性能方面我是非常担忧的，并且，mustache的主要优势是在于跨语言，因此功能也是极其有限的，始终也不是什么好用的模板。

趁现在还未正式开工，我还是很希望能够投入时间把这些技术障碍给提前清理掉。

如何解决go的模板软肋？实现或者移植一个新的吧。

我一直都对网页模板引擎十分着迷，尝试过N种语言的N * X种模板。

令我印象最深刻的，有这么两个模板：

- Django
- Razor

Django模板很强大，并且彻底杜绝了代码与模板的混合。这种“纯粹”的模板技术，我当年是十分喜欢的。

我曾经相信，模板中绝对不可以混杂任何代码，以确保代码分离，程序员是不可相信，必须从技术上阻止他们去做这种傻事。

但随着对技术理解的转变，我是越来越厌恶“纯粹”；我也越来越相信，程序员是可相信的，应该把选择权留给程序员，相信他们，或者说我自己不会滥用。

这些说到底是人的素质的问题；当年我是菜鸟，我为了避免自己做傻事，我觉得有技术的强制限制很爽。但我现在已经不是菜鸟，限制只会让我觉得不爽。当我想要打破什么限制的时候，我是有充分理由，并且完全清楚自己在做什么的。

举个栗子，需要在模板中的时间进行格式化。

如果使用Django的话，那么必须先在python中实现这么个时间格式化函数，然后再将其定义为filter甚至是标签；但我现在觉得额外的负担。

如非必要，勿增实体。

格式化时间，这是业务需求，所以我必须在代码中实现一个时间格式化函数，这个实体是必须增加的。但是，与其对应绑定去模板的filter / tag定义，则是多余得实体。

时间函数接口改了的话，他们也得跟着改。这是实体增加之后带来的复杂性增加。

一层包一层，底层一改，每层都得跟着改；这是非常讨厌的事情，当年asp.net刚出来的时候，web form采用的便是这种多层嵌套的N-tiers模式，被实践证明它非常麻烦，所以微软后来就改推asp.net MVC了。

层次越少越好，实体越少越好。所有的编程问题都可以通过增加一个抽象层来解决，但不是所有的问题都需要通过增加抽象层来解决。

类似Django这种模板，它必须得引入自己新的一个模板语言，以实现各种逻辑控制。

这个模板语言，便是一个抽象层。Go的内置模板，也是一样。

反射也好，解释也好，这模板语言，最终都是会被映射到原生语言的。它能表达的，原语言也一定可以表达。

何不直接用原生语言来表达呢？何苦要让程序员多学习一种新的模板语言呢？

模板语言的这层抽象，真的划得来嘛？

# Razor

asp.net MVC现在的默认模板Razor非常简约。

完全没有各种 `<?, <%, {{` 这些徒增击键次数与视觉干扰的符号。

逻辑控制？使用原生语言就好了，比方说：
```
	@if(totalMessage == 1) {
		<p>@Userame has 1 message</p>
	} else {
		<p>@Userame has @totalMessage messages</p>
	}
```
它只需要一个 @ 符号来插入代码，后面的 `} else {` 以及括号等等，它的编译器可以自动判断出来是代码还是模板。

我当年第一次看到Razor时非常担心它会把代码跟模板搞混，但用过之后几乎完全没有遇到这方面的问题。实际上，用过razor的程序员，[基本都说爽](http://www.zhihu.com/question/19973649)。

没搞.net很多年了，但Razor一直都是我认为设计最优雅的模板，在我使用过node.js jade那样风格的模板后也还是这么认为。

把razor移植到go上来！把razor移植到go上来！把razor移植到go上来！把razor移植到go上来！

这个念头挥之不去；但我也不断在扼制这个念头，创业，最重要的是业务实现！我需要实现的是业务，不是技术！

但Razor真的很简单啊~微软的.net实现也是[完全开源](https://github.com/ASP-NET-MVC/aspnetwebstack/tree/master/src/System.Web.Razor)的，它本来就支持多语言，我只要修改一下它的实现，增加对go的支持就好嘛~一天可以搞定？但我没写c#很多年啊~还得装visual/xamarin studio~搞还不是不搞？

外事不决问google，结果我又发现了[vash](https://github.com/kirbysayshi/vash): razor的javascript实现，一共就两千余行代码，hack一下输出go是绝对是一天可以搞定啊！

嗯，然后就有了[gorazor](https://github.com/sipin/gorazor).

出原型Proof of Concept是非常快的，成熟当然需要时间，但这些我可以在开发业务的时候，遇到细节问题再慢慢改。

我甚至一开始都没有看vash的实现，直接搜索到它的compiler生成函数，就开始hack去go的语法。可以work，我就继续保持quick & dirty；实现一直都很快，只要想清楚设计。

# gorazor

## 代码生成

我做的第一个设计决定就是使用代码生成。

模板只是模板，`gorazor`负责将模板编译为go代码，项目应用中直接引用生成出来的这些代码，而不是模板文件。

所以，模板改动之后，必须重新编译，这对于go来说不是什么问题，go本来就需要编译，而且go的编译速度无比快。

开个文件watcher，检测到硬盘文件有修改，直接kill掉当前进程，重新编译，然后重启进程。这一切经常可以在我按 `save + S / alt + tab / F5` 之间完成。

go需要编译，但开发go，跟开发python等支持热重载的一样畅快。

## 绑定至函数

第二个设计决定是模板生成出来的代码，是函数，而不是一个类。

`msg/inbox.gohtml`这个模板编译之后，生成的是`msg`package下的Inbox函数，而不是Indox struct。

所以，在应用中，我可以直接 `import tpl/msg`，然后使用 `msg.Inbox(...)`输出html。

而不要:
```go
inbox := &msg.Inbox{...}
inbox.render()
```
生成代码使用模板所在的目录吗作为命名空间，文件名作为函数名。

这样的做法也使得gorazor非常自然的就“自动”实现了模板嵌套。


## 模板声明

go是强类型，需要编译的语言；这不是累赘，而是优势。

模板中使用`page.ViewData["XXX"]`，这样一个字典传递数据，我认为是非常不好的做法。

模板接受什么类型的model，应该明确声明；然后让编译器尽可能的找出各种低级错误。

所以，gorazor强制模板开始的`first code block`必须是用来做声明，比方说：
```
	@{
		import (
			. "myapp/models"
		)
		var totalMessage int
		var u *User
	}
```
这声明的意思是说模板接受两个数据，第一个是一个int，第二个是 myapp/models 命名空间下的User struct指针。

如果模板不需要接受参数/model，那就把first code block放空就是。

模板可以绑定任意类型的任意数量的model，这些声明，实际上就是被转换为函数参数：
```go
	func Inbox(totalMessage int, u *User) string {
	   ....
	}
```
也就是说，当代码生成后，应用中要调用函数输出的时候，编辑器/IDE基本上都能有输入提示的，不用再去查模板究竟需要哪些什么类型的参数。

## Section / Layout

绑定函数解决了模板嵌套的问题，gorazor实际上已经跟mustache一样强大了。

通过函数嵌套，也可以模拟去 `section / layout`来。

但作为web视图引擎，能够直接支持`layout / section`那是会方便很多。

需要使用怎样的语法呢？这个问题我倒是犹豫了一天；暂时决定使用magic word，在模板声明指定：
```
	@{
		import (
			. "myapp/models"
			"tpl/layout/base"
		)
		var totalMessage int
		var u *User
	}
```
`tpl/layout/base` gorazor里面的倒数第二次命名空间`layout`做了特殊处理。

它会把这个import自动改为`tpl/layout`，然后调用`Base`函数（首字母自动转换大写，很讨厌go用首字母来区分公有/私有啊！）去生成layout。

我本来是考虑增加一个Layout的关键字，比方说：
```
	@{
		import (
			. "myapp/models"
		)
		Layout	"tpl/layout/base.gohtml"
		var totalMessage int
		var u *User
	}
```
但最终放弃了。

第一，我不想增加 Layout 关键字，这不是go的语法。

第二，我不希望使用`tpl/layout/base.gohtml`这样的路径去指定模板，虽然模板路径与命名空间路径是有联系的，但路径前缀怎么办？一个项目中，未必就只有一套模板，可能有多套，然后不同模板共享layout。

子模板中，未必就知道layout的模板文件路径；全部通过go程序的命名空间来处理好了，这样才统一。

当然，layout文件也支持中import另一个layout。全部都是函数调用，多层layout嵌套也是自动实现的。

出于性能考虑，我也考虑过生成的函数不是返回string，而是接受writer参数。但这样的话，layout就不好处理，layout中的section始终还是需要先写去临时string，然后再拼接写入writer。

考虑到应用中几乎可能所有的页面都会使用layout，那么接受writer参数，似乎并不能改善性能。

这个先不管，要是真出现性能问题，再分析具体情况好了。

Layout也是需要声明的，比方说：
```
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
layout模板接受的参数必须是string，并且它是有顺序的；第一个还必须是**body**。

前面说，页面模板指定layout时，是通过go命名空间去引用；我假设无法直接获得layout模板文件内容；也就是说，它是不能通过解析layout模板文件去获得：
```
	@{
		var body string
		var title string
		var sideMenu string
	}
```
这个函数参数声明；但是，layout的section信息，在页面模板中也是会有：
```
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
页面中没有指定section的部分，自动视为是body，传递给layout，然后，根据section的声明顺序，认为便是layout接受的参数。感觉这很可能会是坑，比方说，程序员无意中改变了section的顺序。

我觉得gorazor可以做得聪明一些，它可以根据layout的命名空间去“猜测”layout模板文件的真实路径，然后解析一下，做更多的提示、转换、默认处理。

比方说，layout接受某个section，而页面中没有声明，则认为该页面的section为空，或者说，使用默认值，layout文件中可以这么写：
```
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
PS: @title的处理其实也非常讨厌，layout中写成`<title>@title</title>`会让页面那边搞得很麻烦。

# 使用 GoRazor

`go get github.com/sipin/gorazor`
