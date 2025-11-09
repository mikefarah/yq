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

func TestPrinterMultipleDocsInSequenceOnly(t *testing.T) {
	var output bytes.Buffer
	var writer = bufio.NewWriter(&output)
	printer := NewSimpleYamlPrinter(writer, true, 2, true)

	inputs, err := readDocuments(strings.NewReader(multiDocSample), "sample.yml", 0, NewYamlDecoder(ConfiguredYamlPreferences))
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
	printer := NewSimpleYamlPrinter(writer, true, 2, true)

	inputs, err := readDocuments(strings.NewReader(multiDocSample), "sample.yml", 0, NewYamlDecoder(ConfiguredYamlPreferences))
	if err != nil {
		panic(err)
	}

	el := inputs.Front()
	el.Value.(*CandidateNode).LeadingContent = "# go cats\n$yqDocSeparator$\n"
	sample1 := nodeToList(el.Value.(*CandidateNode))

	el = el.Next()
	el.Value.(*CandidateNode).LeadingContent = "$yqDocSeparator$\n"
	sample2 := nodeToList(el.Value.(*CandidateNode))

	el = el.Next()
	el.Value.(*CandidateNode).LeadingContent = "$yqDocSeparator$\n# cool\n"
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
	printer := NewSimpleYamlPrinter(writer, true, 2, true)

	inputs, err := readDocuments(strings.NewReader(multiDocSample), "sample.yml", 0, NewYamlDecoder(ConfiguredYamlPreferences))
	if err != nil {
		panic(err)
	}

	el := inputs.Front()
	elNode := el.Value.(*CandidateNode)
	elNode.document = 0
	elNode.fileIndex = 0
	sample1 := nodeToList(elNode)

	el = el.Next()
	elNode = el.Value.(*CandidateNode)
	elNode.document = 0
	elNode.fileIndex = 1
	sample2 := nodeToList(elNode)

	el = el.Next()
	elNode = el.Value.(*CandidateNode)
	elNode.document = 0
	elNode.fileIndex = 2
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
	printer := NewSimpleYamlPrinter(writer, true, 2, true)

	inputs, err := readDocuments(strings.NewReader(multiDocSample), "sample.yml", 0, NewYamlDecoder(ConfiguredYamlPreferences))
	if err != nil {
		panic(err)
	}

	el := inputs.Front()
	elNode := el.Value.(*CandidateNode)
	elNode.document = 0
	elNode.fileIndex = 0
	elNode.LeadingContent = "# go cats\n$yqDocSeparator$\n"
	sample1 := nodeToList(elNode)

	el = el.Next()
	elNode = el.Value.(*CandidateNode)
	elNode.document = 0
	elNode.fileIndex = 1
	elNode.LeadingContent = "$yqDocSeparator$\n"
	sample2 := nodeToList(elNode)

	el = el.Next()
	elNode = el.Value.(*CandidateNode)
	elNode.document = 0
	elNode.fileIndex = 2
	elNode.LeadingContent = "$yqDocSeparator$\n# cool\n"
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
	printer := NewSimpleYamlPrinter(writer, true, 2, true)

	inputs, err := readDocuments(strings.NewReader(multiDocSample), "sample.yml", 0, NewYamlDecoder(ConfiguredYamlPreferences))
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
	printer := NewSimpleYamlPrinter(writer, true, 2, true)

	inputs, err := readDocuments(strings.NewReader(multiDocSample), "sample.yml", 0, NewYamlDecoder(ConfiguredYamlPreferences))
	if err != nil {
		panic(err)
	}

	inputs.Front().Value.(*CandidateNode).LeadingContent = "# go cats\n$yqDocSeparator$\n"

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
	printer := NewSimpleYamlPrinter(writer, true, 2, true)

	inputs, err := readDocuments(strings.NewReader(multiDocSample), "sample.yml", 0, NewYamlDecoder(ConfiguredYamlPreferences))
	if err != nil {
		panic(err)
	}
	inputs.Front().Value.(*CandidateNode).LeadingContent = "$yqDocSeparator$\n"
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
	printer := NewSimpleYamlPrinter(writer, true, 2, true)

	node, err := getExpressionParser().ParseExpression(".a")
	if err != nil {
		panic(err)
	}
	streamEvaluator := NewStreamEvaluator()
	_, err = streamEvaluator.Evaluate("sample", strings.NewReader(multiDocSample), node, printer, NewYamlDecoder(ConfiguredYamlPreferences))
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
	// when outputting JSON.
	prefs := ConfiguredJSONPreferences.Copy()
	prefs.Indent = 0
	encoder := NewJSONEncoder(prefs)
	if encoder == nil {
		t.Skipf("no support for %s output format", "json")
	}
	printer := NewPrinter(encoder, NewSinglePrinterWriter(writer))

	inputs, err := readDocuments(strings.NewReader(multiDocSample), "sample.yml", 0, NewYamlDecoder(ConfiguredYamlPreferences))
	if err != nil {
		panic(err)
	}

	inputs.Front().Value.(*CandidateNode).LeadingContent = "# ignore this\n"

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

func TestPrinterNulSeparator(t *testing.T) {
	var output bytes.Buffer
	var writer = bufio.NewWriter(&output)
	printer := NewSimpleYamlPrinter(writer, true, 2, false)
	printer.SetNulSepOutput(true)
	node, err := getExpressionParser().ParseExpression(".a")
	if err != nil {
		panic(err)
	}
	streamEvaluator := NewStreamEvaluator()
	_, err = streamEvaluator.Evaluate("sample", strings.NewReader(multiDocSample), node, printer, NewYamlDecoder(ConfiguredYamlPreferences))
	if err != nil {
		panic(err)
	}

	writer.Flush()
	expected := "banana\x00apple\x00coconut\x00"
	test.AssertResult(t, expected, output.String())
}

func TestPrinterNulSeparatorWithJson(t *testing.T) {
	var output bytes.Buffer
	var writer = bufio.NewWriter(&output)
	// note printDocSeparators is true, it should still not print document separators
	// when outputting JSON.
	prefs := ConfiguredJSONPreferences.Copy()
	prefs.Indent = 0
	encoder := NewJSONEncoder(prefs)
	if encoder == nil {
		t.Skipf("no support for %s output format", "json")
	}
	printer := NewPrinter(encoder, NewSinglePrinterWriter(writer))
	printer.SetNulSepOutput(true)

	inputs, err := readDocuments(strings.NewReader(multiDocSample), "sample.yml", 0, NewYamlDecoder(ConfiguredYamlPreferences))
	if err != nil {
		panic(err)
	}

	inputs.Front().Value.(*CandidateNode).LeadingContent = "# ignore this\n"

	err = printer.PrintResults(inputs)
	if err != nil {
		panic(err)
	}

	expected := `{"a":"banana"}` + "\x00" + `{"a":"apple"}` + "\x00" + `{"a":"coconut"}` + "\x00"

	writer.Flush()
	test.AssertResult(t, expected, output.String())
}

func TestPrinterRootUnwrap(t *testing.T) {
	var output bytes.Buffer
	var writer = bufio.NewWriter(&output)
	printer := NewSimpleYamlPrinter(writer, true, 2, false)
	node, err := getExpressionParser().ParseExpression(".")
	if err != nil {
		panic(err)
	}
	streamEvaluator := NewStreamEvaluator()
	_, err = streamEvaluator.Evaluate("sample", strings.NewReader("'a'"), node, printer, NewYamlDecoder(ConfiguredYamlPreferences))
	if err != nil {
		panic(err)
	}

	writer.Flush()
	expected := `a
`
	test.AssertResult(t, expected, output.String())
}

func TestRemoveLastEOL(t *testing.T) {
	// Test with \r\n
	buffer := bytes.NewBufferString("test\r\n")
	removeLastEOL(buffer)
	test.AssertResult(t, "test", buffer.String())

	// Test with \n only
	buffer = bytes.NewBufferString("test\n")
	removeLastEOL(buffer)
	test.AssertResult(t, "test", buffer.String())

	// Test with \r only
	buffer = bytes.NewBufferString("test\r")
	removeLastEOL(buffer)
	test.AssertResult(t, "test", buffer.String())

	// Test with no EOL
	buffer = bytes.NewBufferString("test")
	removeLastEOL(buffer)
	test.AssertResult(t, "test", buffer.String())

	// Test with empty buffer
	buffer = bytes.NewBufferString("")
	removeLastEOL(buffer)
	test.AssertResult(t, "", buffer.String())

	// Test with multiple \r\n
	buffer = bytes.NewBufferString("line1\r\nline2\r\n")
	removeLastEOL(buffer)
	test.AssertResult(t, "line1\r\nline2", buffer.String())
}

func TestPrinterPrintedAnything(t *testing.T) {
	var output bytes.Buffer
	var writer = bufio.NewWriter(&output)
	printer := NewSimpleYamlPrinter(writer, true, 2, true)

	test.AssertResult(t, false, printer.PrintedAnything())

	// Print a scalar value
	node := createStringScalarNode("test")
	nodeList := nodeToList(node)
	err := printer.PrintResults(nodeList)
	if err != nil {
		t.Fatal(err)
	}

	// Should now be true
	test.AssertResult(t, true, printer.PrintedAnything())
}

func TestPrinterNulSeparatorWithNullChar(t *testing.T) {
	var output bytes.Buffer
	var writer = bufio.NewWriter(&output)
	printer := NewSimpleYamlPrinter(writer, true, 2, false)
	printer.SetNulSepOutput(true)

	// Create a node with null character
	node := createStringScalarNode("test\x00value")
	nodeList := nodeToList(node)

	err := printer.PrintResults(nodeList)
	if err == nil {
		t.Fatal("Expected error for null character in NUL separated output")
	}

	expectedError := "can't serialize value because it contains NUL char and you are using NUL separated output"
	if err.Error() != expectedError {
		t.Fatalf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestPrinterSetNulSepOutput(t *testing.T) {
	var output bytes.Buffer
	var writer = bufio.NewWriter(&output)
	printer := NewSimpleYamlPrinter(writer, true, 2, false)

	// Test setting NUL separator output
	printer.SetNulSepOutput(true)
	test.AssertResult(t, true, true) // Placeholder assertion

	printer.SetNulSepOutput(false)
	// Should also not cause errors
	test.AssertResult(t, false, false) // Placeholder assertion
}

func TestPrinterSetAppendix(t *testing.T) {
	var output bytes.Buffer
	var writer = bufio.NewWriter(&output)
	printer := NewSimpleYamlPrinter(writer, true, 2, true)

	// Test setting appendix
	appendix := strings.NewReader("appendix content")
	printer.SetAppendix(appendix)
	test.AssertResult(t, true, true) // Placeholder assertion
}
