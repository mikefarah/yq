package yqlib

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

// countOpenFileDescriptors returns the number of file descriptors this process
// currently holds open, or -1 if the platform does not expose them.
func countOpenFileDescriptors() int {
	for _, dir := range []string{"/proc/self/fd", "/dev/fd"} {
		// Readdirnames avoids stat-ing each entry, which races with descriptors
		// (including this directory handle) being closed underneath us.
		handle, err := os.Open(dir)
		if err != nil {
			continue
		}
		names, err := handle.Readdirnames(-1)
		safelyCloseFile(handle)
		if err != nil {
			continue
		}
		// discount the directory handle itself
		return len(names) - 1
	}
	return -1
}

func writeSampleFiles(t *testing.T, count int) []string {
	t.Helper()
	dir := t.TempDir()
	filenames := make([]string, count)
	for i := 0; i < count; i++ {
		filename := filepath.Join(dir, "sample-"+strconv.Itoa(i)+".yml")
		if err := os.WriteFile(filename, []byte("a: apple\n"), 0600); err != nil {
			t.Fatalf("failed to write sample file: %v", err)
		}
		filenames[i] = filename
	}
	return filenames
}

func discardingPrinter() Printer {
	return NewPrinter(NewYamlEncoder(ConfiguredYamlPreferences), NewSinglePrinterWriter(bufio.NewWriter(io.Discard)))
}

func assertNoLeakedFileDescriptors(t *testing.T, evaluate func(filenames []string) error) {
	t.Helper()
	InitExpressionParser()

	// os.File finalisers close leaked descriptors on collection, which would
	// let a genuine leak pass unnoticed.
	defer debug.SetGCPercent(debug.SetGCPercent(-1))

	before := countOpenFileDescriptors()
	if before < 0 {
		t.Skip("file descriptors are not observable on this platform")
	}

	if err := evaluate(writeSampleFiles(t, 50)); err != nil {
		t.Fatalf("failed to evaluate files: %v", err)
	}

	after := countOpenFileDescriptors()
	if after > before {
		t.Errorf("expected no additional open file descriptors, had %d before and %d after", before, after)
	}
}

func TestStreamEvaluatorClosesInputFiles(t *testing.T) {
	assertNoLeakedFileDescriptors(t, func(filenames []string) error {
		return NewStreamEvaluator().EvaluateFiles(".a", filenames, discardingPrinter(), NewYamlDecoder(ConfiguredYamlPreferences))
	})
}

func TestAllAtOnceEvaluatorClosesInputFiles(t *testing.T) {
	assertNoLeakedFileDescriptors(t, func(filenames []string) error {
		return NewAllAtOnceEvaluator().EvaluateFiles(".a", filenames, discardingPrinter(), NewYamlDecoder(ConfiguredYamlPreferences))
	})
}

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
