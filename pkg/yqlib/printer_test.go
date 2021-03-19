package yqlib

import (
	"bufio"
	"bytes"
	"container/list"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var multiDocSample = `a: banana
---
a: apple
---
a: coconut
`

func nodeToList(candidate *CandidateNode) *list.List {
	elMap := list.New()
	elMap.PushBack(candidate)
	return elMap
}

func TestPrinterMultipleDocsInSequence(t *testing.T) {
	var output bytes.Buffer
	var writer = bufio.NewWriter(&output)
	printer := NewPrinter(writer, false, true, false, 2, true)

	inputs, err := readDocuments(strings.NewReader(multiDocSample), "sample.yml", 0)
	if err != nil {
		panic(err)
	}

	el := inputs.Front()
	sample1 := nodeToList(el.Value.(*CandidateNode))

	el = el.Next()
	sample2 := nodeToList(el.Value.(*CandidateNode))

	el = el.Next()
	sample3 := nodeToList(el.Value.(*CandidateNode))

	err = printer.PrintResults(sample1)
	if err != nil {
		panic(err)
	}

	err = printer.PrintResults(sample2)
	if err != nil {
		panic(err)
	}

	err = printer.PrintResults(sample3)
	if err != nil {
		panic(err)
	}

	writer.Flush()
	test.AssertResult(t, multiDocSample, output.String())
}

func TestPrinterMultipleFilesInSequence(t *testing.T) {
	var output bytes.Buffer
	var writer = bufio.NewWriter(&output)
	printer := NewPrinter(writer, false, true, false, 2, true)

	inputs, err := readDocuments(strings.NewReader(multiDocSample), "sample.yml", 0)
	if err != nil {
		panic(err)
	}

	el := inputs.Front()
	elNode := el.Value.(*CandidateNode)
	elNode.Document = 0
	elNode.FileIndex = 0
	sample1 := nodeToList(elNode)

	el = el.Next()
	elNode = el.Value.(*CandidateNode)
	elNode.Document = 0
	elNode.FileIndex = 1
	sample2 := nodeToList(elNode)

	el = el.Next()
	elNode = el.Value.(*CandidateNode)
	elNode.Document = 0
	elNode.FileIndex = 2
	sample3 := nodeToList(elNode)

	err = printer.PrintResults(sample1)
	if err != nil {
		panic(err)
	}

	err = printer.PrintResults(sample2)
	if err != nil {
		panic(err)
	}

	err = printer.PrintResults(sample3)
	if err != nil {
		panic(err)
	}

	writer.Flush()
	test.AssertResult(t, multiDocSample, output.String())
}

func TestPrinterMultipleDocsInSinglePrint(t *testing.T) {
	var output bytes.Buffer
	var writer = bufio.NewWriter(&output)
	printer := NewPrinter(writer, false, true, false, 2, true)

	inputs, err := readDocuments(strings.NewReader(multiDocSample), "sample.yml", 0)
	if err != nil {
		panic(err)
	}

	err = printer.PrintResults(inputs)
	if err != nil {
		panic(err)
	}

	writer.Flush()
	test.AssertResult(t, multiDocSample, output.String())
}

func TestPrinterMultipleDocsJson(t *testing.T) {
	var output bytes.Buffer
	var writer = bufio.NewWriter(&output)
	// note printDocSeparators is true, it should still not print document separators
	// when outputing JSON.
	printer := NewPrinter(writer, true, true, false, 0, true)

	inputs, err := readDocuments(strings.NewReader(multiDocSample), "sample.yml", 0)
	if err != nil {
		panic(err)
	}

	err = printer.PrintResults(inputs)
	if err != nil {
		panic(err)
	}

	expected := `{"a":"banana"}
{"a":"apple"}
{"a":"coconut"}
`

	writer.Flush()
	test.AssertResult(t, expected, output.String())
}
