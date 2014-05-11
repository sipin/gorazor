package gorazor

import (
	"html/template"
	"strconv"
	"unicode"
)

func HTMLEscape(obj string) string {
	return template.HTMLEscapeString(obj)
}

func Itoa(obj int) string {
	return strconv.Itoa(obj)
}

func Capitalize(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}
