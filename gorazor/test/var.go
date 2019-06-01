package cases

import (
	"bytes"
	"io"
	"strings"
)

func Var(totalMessage int) string {
	var _b strings.Builder
	RenderVar(&_b, totalMessage)
	return _b.String()
}

func RenderVar(_buffer io.StringWriter, totalMessage int) {

}
