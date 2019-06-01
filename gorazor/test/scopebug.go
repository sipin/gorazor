package cases

import (
	"bytes"
	"dm"
	"io"
	"strings"
	"zfw/models"
	. "zfw/tplhelper"
)

func Scopebug(obj *models.Widget) string {
	var _b strings.Builder
	RenderScopebug(&_b, obj)
	return _b.String()
}

func RenderScopebug(_buffer io.StringWriter, obj *models.Widget) {

	if 1 == 2 {
	} else {
		values := []int{}
		for _, v := range values {
			if v, ok := v.(type); ok {

				_buffer.WriteString("<a>\n\t\t\t\t\t")
				for _, v := range values {
				}
				_buffer.WriteString("\n\t\t\t\t</a>")

			} else {

			}
		}
	}

}
