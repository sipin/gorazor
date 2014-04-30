package gorazor

import (
	"html/template"
	"strconv"
)

func HTMLEscape(obj interface{}) string {
	switch obj := obj.(type) {
	case int:
		return strconv.Itoa(obj)
	case string:
		return template.HTMLEscapeString(obj)
	}
	return obj.(string)
}
