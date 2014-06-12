package gorazor

import (
	"fmt"
	"html/template"
	"os"
	"strconv"
	"strings"
	"time"
)

func HTMLEscape(m interface{}) string {
	s := fmt.Sprint(m)
	return template.HTMLEscapeString(s)
}

func StrTime(timestamp int64, format string) string {
	return time.Unix(timestamp, 0).Format(format)
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

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
