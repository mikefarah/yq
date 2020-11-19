package yqlib

import (
	"bufio"
	"bytes"
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

func TestPrinterMultipleDocsInSequence(t *testing.T) {
	var output bytes.Buffer
	var writer = bufio.NewWriter(&output)
	printer := NewPrinter(writer, false, true, false, 2, true)

	inputs, err := readDocuments(strings.NewReader(multiDocSample), "sample.yml")
	if err != nil {
		panic(err)
	}

	el := inputs.Front()
	sample1 := nodeToMap(el.Value.(*CandidateNode))

	el = el.Next()
	sample2 := nodeToMap(el.Value.(*CandidateNode))

	el = el.Next()
	sample3 := nodeToMap(el.Value.(*CandidateNode))

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

	inputs, err := readDocuments(strings.NewReader(multiDocSample), "sample.yml")
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
	printer := NewPrinter(writer, true, true, false, 0, false)

	inputs, err := readDocuments(strings.NewReader(multiDocSample), "sample.yml")
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
