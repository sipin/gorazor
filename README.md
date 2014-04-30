# GoRazor

GoRazor is the Go port of the razor template engine originated from [asp.net world](http://weblogs.asp.net/scottgu/archive/2010/07/02/introducing-razor.aspx).

This port is essentially a re-port from razor's port in javascript: [vash](https://github.com/kirbysayshi/vash).

It modifies on vast's generation function to emit go code instead of javascript code.

In summay, GoRazor is:

* Consice syntax
* Able to mix go code in template
  * import is supported
  * Call arbitrary funtions
* Take code generation approach, i.e. no reflection overhead
* Strong type template model
* Utils class


# FAQ

TBA