package tpl

import (
	"io"

	"github.com/sipin/gorazor/gorazor"
)

type BenchRow struct {
	ID      int
	Message string
	Print   bool
}

func RenderIndex(_buffer io.StringWriter, rows []BenchRow) {
	for _, row := range rows {
		gorazor.HTMLEscape(row.ID)
		_buffer.WriteString(gorazor.HTMLEscape(row.Message))
	}

}
