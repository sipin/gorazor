package gorazor

import (
	"html/template"
	"strconv"
	"strings"
)

func HTMLEscape(obj string) string {
	return template.HTMLEscapeString(obj)
}

func Itoa(obj int) string {
	return strconv.Itoa(obj)
}

func Capitalize(str string) string {
	if len(str) == 0 {
		return ""
	}
	return strings.ToUpper(str[0:1]) + str[1:]
}
