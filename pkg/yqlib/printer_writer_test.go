package yqlib

import (
	"container/list"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"testing"
)

type failingPrinterEncoder struct{}

func (failingPrinterEncoder) Encode(io.Writer, *CandidateNode) error {
	return errors.New("encode failed")
}

func (failingPrinterEncoder) PrintDocumentSeparator(io.Writer) error {
	return nil
}

func (failingPrinterEncoder) PrintLeadingContent(io.Writer, string) error {
	return nil
}

func (failingPrinterEncoder) CanHandleAliases() bool {
	return true
}

func TestMultiPrintWriterClosesSplitFiles(t *testing.T) {
	InitExpressionParser()
	t.Chdir(t.TempDir())

	expression, err := ExpressionParser.ParseExpression(`$index | tostring`)
	if err != nil {
		t.Fatal(err)
	}
	printer := NewPrinter(
		NewYamlEncoder(NewDefaultYamlPreferences()),
		NewMultiPrinterWriter(expression, YamlFormat),
	)

	matches := list.New()
	for i := range 40 {
		matches.PushBack(&CandidateNode{
			Kind:  ScalarNode,
			Tag:   "!!str",
			Value: fmt.Sprintf("item-%d", i),
		})
	}

	oldGCPercent := debug.SetGCPercent(-1)
	t.Cleanup(func() { debug.SetGCPercent(oldGCPercent) })
	before := countOpenFileDescriptors()
	if err := printer.PrintResults(matches); err != nil {
		t.Fatal(err)
	}
	after := countOpenFileDescriptors()
	if after > before {
		t.Fatalf("open file descriptors grew from %d to %d", before, after)
	}

	for i := range 40 {
		filename := filepath.Join(".", fmt.Sprintf("%d.yml", i))
		content, err := os.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}
		if string(content) != fmt.Sprintf("item-%d\n", i) {
			t.Fatalf("%s contains %q", filename, content)
		}
	}
}

func TestMultiPrintWriterClosesSplitFilesOnEncodeError(t *testing.T) {
	InitExpressionParser()
	t.Chdir(t.TempDir())

	expression, err := ExpressionParser.ParseExpression(`$index | tostring`)
	if err != nil {
		t.Fatal(err)
	}
	printer := NewPrinter(
		failingPrinterEncoder{},
		NewMultiPrinterWriter(expression, YamlFormat),
	)
	matches := list.New()
	matches.PushBack(&CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "item"})

	oldGCPercent := debug.SetGCPercent(-1)
	t.Cleanup(func() { debug.SetGCPercent(oldGCPercent) })
	before := countOpenFileDescriptors()
	if err := printer.PrintResults(matches); err == nil {
		t.Fatal("expected encoding error")
	}
	after := countOpenFileDescriptors()
	if after > before {
		t.Fatalf("open file descriptors grew from %d to %d", before, after)
	}
}
