package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sipin/gorazor/tests/data"
	"github.com/sipin/gorazor/tests/tpl"
)

func BenchmarkRazorTemplate1(b *testing.B) {
	benchmarkRazorTemplate(b, 1)
}

func BenchmarkRazorTemplate10(b *testing.B) {
	benchmarkRazorTemplate(b, 10)
}

func BenchmarkRazorTemplate100(b *testing.B) {
	benchmarkRazorTemplate(b, 100)
}

func benchmarkRazorTemplate(b *testing.B, rowsCount int) {
	rows := getBenchRows(rowsCount)
	b.RunParallel(func(pb *testing.PB) {
		// bb := quicktemplate.AcquireByteBuffer()
		// for pb.Next() {
		// 	templates.WriteBenchPage(bb, rows)
		// 	bb.Reset()
		// }
		// quicktemplate.ReleaseByteBuffer(bb)
		for pb.Next() {
			var sb strings.Builder
			tpl.RenderIndex(&sb, rows)
		}
	})
}

func getBenchRows(n int) []data.BenchRow {
	rows := make([]data.BenchRow, n)
	for i := 0; i < n; i++ {
		rows[i] = data.BenchRow{
			ID:      i,
			Message: fmt.Sprintf("message %d", i),
			Print:   ((i & 1) == 0),
		}
	}
	return rows
}
