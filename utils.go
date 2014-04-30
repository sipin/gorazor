package gorazor

import (
	"html/template"
	"strconv"
)

func HTMLEscape(obj string) string {
	return template.HTMLEscapeString(obj)
}

func Itoa(obj int) string {
	return strconv.Itoa(obj)
}
