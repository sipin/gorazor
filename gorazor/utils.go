package gorazor

import (
	"html/template"
	"os"
	"strconv"
	"strings"
	"time"
	"fmt"
)

func HTMLEscape(obj string) string {
	return template.HTMLEscapeString(obj)
}

func Itoa(obj int) string {
	return strconv.Itoa(obj)
}

func StrTime(timestamp int64, format string) string {
	return fmt.Sprintf("%s", time.Unix(timestamp, 0).Format(format))
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
