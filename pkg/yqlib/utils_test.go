package yqlib

import (
	"bufio"
	"bytes"
	"io"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

// plainWriter only implements io.Writer, so io.WriteString must fall back to Write.
type plainWriter struct {
	buf bytes.Buffer
}

func (w *plainWriter) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}

func TestWriteStringToStringWriter(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	test.AssertResult(t, nil, writeString(writer, "hello world"))
	test.AssertResult(t, nil, writer.Flush())
	test.AssertResult(t, "hello world", buf.String())
}

func TestWriteStringToPlainWriter(t *testing.T) {
	writer := &plainWriter{}
	test.AssertResult(t, nil, writeString(writer, "hello world"))
	test.AssertResult(t, "hello world", writer.buf.String())
}

func TestWriteStringDoesNotAllocate(t *testing.T) {
	writer := bufio.NewWriter(io.Discard)
	allocations := testing.AllocsPerRun(100, func() {
		_ = writeString(writer, "hello world")
	})
	test.AssertResult(t, 0.0, allocations)
}
