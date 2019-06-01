package layout

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"zfw/models"
)

func Args(objs ...*models.Widget) string {
	var _b strings.Builder

	_objs := func(_buffer io.StringWriter) {
		_buffer.WriteString(objs)
	}

	WriteArgs(_b, _objs)
	return _b.String()
}

func WriteArgs(_buffer io.StringWriter, objs ...*models.Widget) {

	size := strconv.Itoa(12 / len(objs))

}
