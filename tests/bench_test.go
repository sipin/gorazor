package tests

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"
	"text/template"
	"unsafe"

	"github.com/sipin/gorazor/tests/data"
	"github.com/sipin/gorazor/tests/tpl"
	"github.com/valyala/quicktemplate"
)

var htmltpl = template.Must(template.ParseFiles("tpl/bench.tpl"))

func init() {
	// make sure that both template engines generate the same result
	rows := getBenchRows(3)

	bb1 := &quicktemplate.ByteBuffer{}
	if err := htmltpl.Execute(bb1, rows); err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	bb := quicktemplate.AcquireByteBuffer()
	q := &quickStringWriter{}
	q.bb = bb
	tpl.RenderIndex(q, rows)

	if !bytes.Equal(bb1.B, bb.B) {
		log.Fatalf("results mismatch:\n%q\n%q", bb1, bb)
		quicktemplate.ReleaseByteBuffer(bb)
	}
	quicktemplate.ReleaseByteBuffer(bb)
}

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
		for pb.Next() {
			var sb strings.Builder
			tpl.RenderIndex(&sb, rows)
		}
	})
}

func BenchmarkRazorQuickTemplate1(b *testing.B) {
	benchmarkRazorQuickTemplate(b, 1)
}

func BenchmarkRazorQuickTemplate10(b *testing.B) {
	benchmarkRazorQuickTemplate(b, 10)
}

func BenchmarkRazorQuickTemplate100(b *testing.B) {
	benchmarkRazorQuickTemplate(b, 100)
}

func benchmarkRazorQuickTemplate(b *testing.B, rowsCount int) {
	rows := getBenchRows(rowsCount)

	b.RunParallel(func(pb *testing.PB) {
		bb := quicktemplate.AcquireByteBuffer()
		var q quickStringWriter
		q.bb = bb
		for pb.Next() {
			tpl.RenderIndex(&q, rows)
			bb.Reset()
		}
		quicktemplate.ReleaseByteBuffer(bb)
	})
}

func unsafeStrToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

type quickStringWriter struct {
	bb *quicktemplate.ByteBuffer
}

func (q *quickStringWriter) WriteString(s string) (i int, e error) {
	return q.bb.Write(unsafeStrToBytes(s))
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
