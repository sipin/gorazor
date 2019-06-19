package gorazor

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"html/template"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	rgxSyntaxError = regexp.MustCompile(`(\d+):\d+: `)
)

// FormatBuffer format go code, panic when code is invalid
func FormatBuffer(code string) string {
	buf := bytes.NewBufferString(code)
	output, err := format.Source(buf.Bytes())
	if err == nil {
		return string(output)
	}

	matches := rgxSyntaxError.FindStringSubmatch(err.Error())
	if matches == nil {
		panic(errors.New("failed to format template"))
	}

	lineNum, _ := strconv.Atoi(matches[1])
	scanner := bufio.NewScanner(buf)
	errBuf := &bytes.Buffer{}
	line := 1
	for ; scanner.Scan(); line++ {
		if delta := line - lineNum; delta < -5 || delta > 5 {
			continue
		}

		if line == lineNum {
			errBuf.WriteString(">>>> ")
		} else {
			fmt.Fprintf(errBuf, "% 4d ", line)
		}
		errBuf.Write(scanner.Bytes())
		errBuf.WriteByte('\n')
	}

	panic("failed to format template\n\n" + string(errBuf.Bytes()))

	return ""
}

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

// HTMLEscapeInt strconv.Itoa
func HTMLEscapeInt(m int) string {
	return strconv.Itoa(m)
}

// HTMLEscapeStr is alias to template.HTMLEscapeString
func HTMLEscapeStr(m string) string {
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
