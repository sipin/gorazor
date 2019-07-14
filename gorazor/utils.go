package gorazor

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
)

// HTMLEscape wraps template.HTMLEscapeString
func HTMLEscape(m interface{}) string {
	switch v := m.(type) {
	case int:
		return strconv.Itoa(v)
	case string:
		return template.HTMLEscapeString(v)
	}

	s := fmt.Sprint(m)
	return template.HTMLEscapeString(s)
}

// HTMLEscInt strconv.Itoa
func HTMLEscInt(m int) string {
	return strconv.Itoa(m)
}

// HTMLEscStr is alias to template.HTMLEscapeString
func HTMLEscStr(m string) string {
	return template.HTMLEscapeString(m)
}

// Itoa wraps strconv.Itoa
func Itoa(obj int) string {
	return strconv.Itoa(obj)
}

// Capitalize change first character to upper
func Capitalize(str string) string {
	if len(str) == 0 {
		return ""
	}
	return strings.ToUpper(str[0:1]) + str[1:]
}
