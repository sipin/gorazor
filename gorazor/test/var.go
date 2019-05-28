package cases

import (
	"bytes"
	"io"
	"strings"
)

func Var(totalMessage int) string {
	var _b strings.Builder
	WriteVar(&_b, totalMessage)
	return _b.String()
}

func WriteVar(_buffer io.StringWriter, totalMessage int) {

}
