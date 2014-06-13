package cases

import (
	"bytes"
	"strconv"
	"zfw/models"
)

func Slashbug(objs ...*models.Widget) string {
	var _buffer bytes.Buffer

	size := strconv.Itoa(12 / len(objs))

	return _buffer.String()
}
