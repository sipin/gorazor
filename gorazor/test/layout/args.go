package layout

import (
	"bytes"
	"strconv"
	"zfw/models"
)

func Args(objs ...*models.Widget) string {
	var _buffer bytes.Buffer

	size := strconv.Itoa(12 / len(objs))

	return _buffer.String()
}
