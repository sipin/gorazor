package cases

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"zfw/models"
)

func Slashbug(objs ...*models.Widget) string {
	var _b strings.Builder
	WriteSlashbug(&_b, objs)
	return _b.String()
}

func WriteSlashbug(_buffer io.StringWriter, objs ...*models.Widget) {

	size := strconv.Itoa(12 / len(objs))

}
