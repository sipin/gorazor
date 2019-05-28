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
	WriteArgs(&_b, objs)
	return _b.String()
}

func WriteArgs(_buffer io.StringWriter, objs ...*models.Widget) {

	size := strconv.Itoa(12 / len(objs))

}
