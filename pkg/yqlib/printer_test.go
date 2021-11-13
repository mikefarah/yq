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

var multiDocSampleLeadingExpected = `# go cats
---
a: banana
---
a: apple
---
# cool
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
	printer := NewPrinterWithSingleWriter(writer, YamlOutputFormat, true, false, 2, true)

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

func TestPrinterMultipleDocsInSequenceWithLeadingContent(t *testing.T) {
	var output bytes.Buffer
	var writer = bufio.NewWriter(&output)
	printer := NewPrinterWithSingleWriter(writer, YamlOutputFormat, true, false, 2, true)

	inputs, err := readDocuments(strings.NewReader(multiDocSample), "sample.yml", 0)
	if err != nil {
		panic(err)
	}

	el := inputs.Front()
	el.Value.(*CandidateNode).LeadingContent = "$yqLeadingContent$\n# go cats\n$yqDocSeperator$\n"
	sample1 := nodeToList(el.Value.(*CandidateNode))

	el = el.Next()
	el.Value.(*CandidateNode).LeadingContent = "$yqLeadingContent$\n$yqDocSeperator$\n"
	sample2 := nodeToList(el.Value.(*CandidateNode))

	el = el.Next()
	el.Value.(*CandidateNode).LeadingContent = "$yqLeadingContent$\n$yqDocSeperator$\n# cool\n"
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

	test.AssertResult(t, multiDocSampleLeadingExpected, output.String())
}

func TestPrinterMultipleFilesInSequence(t *testing.T) {
	var output bytes.Buffer
	var writer = bufio.NewWriter(&output)
	printer := NewPrinterWithSingleWriter(writer, YamlOutputFormat, true, false, 2, true)

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

func TestPrinterMultipleFilesInSequenceWithLeadingContent(t *testing.T) {
	var output bytes.Buffer
	var writer = bufio.NewWriter(&output)
	printer := NewPrinterWithSingleWriter(writer, YamlOutputFormat, true, false, 2, true)

	inputs, err := readDocuments(strings.NewReader(multiDocSample), "sample.yml", 0)
	if err != nil {
		panic(err)
	}

	el := inputs.Front()
	elNode := el.Value.(*CandidateNode)
	elNode.Document = 0
	elNode.FileIndex = 0
	elNode.LeadingContent = "$yqLeadingContent$\n# go cats\n$yqDocSeperator$\n"
	sample1 := nodeToList(elNode)

	el = el.Next()
	elNode = el.Value.(*CandidateNode)
	elNode.Document = 0
	elNode.FileIndex = 1
	elNode.LeadingContent = "$yqLeadingContent$\n$yqDocSeperator$\n"
	sample2 := nodeToList(elNode)

	el = el.Next()
	elNode = el.Value.(*CandidateNode)
	elNode.Document = 0
	elNode.FileIndex = 2
	elNode.LeadingContent = "$yqLeadingContent$\n$yqDocSeperator$\n# cool\n"
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
	test.AssertResult(t, multiDocSampleLeadingExpected, output.String())
}

func TestPrinterMultipleDocsInSinglePrint(t *testing.T) {
	var output bytes.Buffer
	var writer = bufio.NewWriter(&output)
	printer := NewPrinterWithSingleWriter(writer, YamlOutputFormat, true, false, 2, true)

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

func TestPrinterMultipleDocsInSinglePrintWithLeadingDoc(t *testing.T) {
	var output bytes.Buffer
	var writer = bufio.NewWriter(&output)
	printer := NewPrinterWithSingleWriter(writer, YamlOutputFormat, true, false, 2, true)

	inputs, err := readDocuments(strings.NewReader(multiDocSample), "sample.yml", 0)
	if err != nil {
		panic(err)
	}

	inputs.Front().Value.(*CandidateNode).LeadingContent = "$yqLeadingContent$\n# go cats\n$yqDocSeperator$\n"

	err = printer.PrintResults(inputs)
	if err != nil {
		panic(err)
	}

	writer.Flush()
	expected := `# go cats
---
a: banana
---
a: apple
---
a: coconut
`
	test.AssertResult(t, expected, output.String())
}

func TestPrinterMultipleDocsInSinglePrintWithLeadingDocTrailing(t *testing.T) {
	var output bytes.Buffer
	var writer = bufio.NewWriter(&output)
	printer := NewPrinterWithSingleWriter(writer, YamlOutputFormat, true, false, 2, true)

	inputs, err := readDocuments(strings.NewReader(multiDocSample), "sample.yml", 0)
	if err != nil {
		panic(err)
	}
	inputs.Front().Value.(*CandidateNode).LeadingContent = "$yqLeadingContent$\n$yqDocSeperator$\n"
	err = printer.PrintResults(inputs)
	if err != nil {
		panic(err)
	}

	writer.Flush()
	expected := `---
a: banana
---
a: apple
---
a: coconut
`
	test.AssertResult(t, expected, output.String())
}

func TestPrinterScalarWithLeadingCont(t *testing.T) {
	var output bytes.Buffer
	var writer = bufio.NewWriter(&output)
	printer := NewPrinterWithSingleWriter(writer, YamlOutputFormat, true, false, 2, true)

	node, err := NewExpressionParser().ParseExpression(".a")
	if err != nil {
		panic(err)
	}
	streamEvaluator := NewStreamEvaluator()
	_, err = streamEvaluator.Evaluate("sample", strings.NewReader(multiDocSample), node, printer, "# blah\n")
	if err != nil {
		panic(err)
	}

	writer.Flush()
	expected := `banana
---
apple
---
coconut
`
	test.AssertResult(t, expected, output.String())
}

func TestPrinterMultipleDocsJson(t *testing.T) {
	var output bytes.Buffer
	var writer = bufio.NewWriter(&output)
	// note printDocSeparators is true, it should still not print document separators
	// when outputing JSON.
	printer := NewPrinterWithSingleWriter(writer, JsonOutputFormat, true, false, 0, true)

	inputs, err := readDocuments(strings.NewReader(multiDocSample), "sample.yml", 0)
	if err != nil {
		panic(err)
	}

	inputs.Front().Value.(*CandidateNode).LeadingContent = "$yqLeadingContent$\n# ignore this\n"

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
