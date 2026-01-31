package yqlib

import (
	"bufio"
	"bytes"
	"container/list"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

func TestNodeInfoPrinter_PrintResults(t *testing.T) {
	// Create a simple CandidateNode
	node := &CandidateNode{
		Kind:        ScalarNode,
		Style:       DoubleQuotedStyle,
		Tag:         "!!str",
		Value:       "hello world",
		Line:        5,
		Column:      7,
		HeadComment: "head",
		LineComment: "line",
		FootComment: "foot",
		Anchor:      "anchor",
	}
	listNodes := list.New()
	listNodes.PushBack(node)

	var output bytes.Buffer
	writer := bufio.NewWriter(&output)
	printer := NewNodeInfoPrinter(NewSinglePrinterWriter(writer))
	err := printer.PrintResults(listNodes)
	writer.Flush()
	if err != nil {
		t.Fatalf("PrintResults error: %v", err)
	}

	outStr := output.String()
	// Check for key NodeInfo fields in YAML output using substring checks
	test.AssertResult(t, true, strings.Contains(outStr, "kind: ScalarNode"))
	test.AssertResult(t, true, strings.Contains(outStr, "style: DoubleQuotedStyle"))
	test.AssertResult(t, true, strings.Contains(outStr, "tag: '!!str'"))
	test.AssertResult(t, true, strings.Contains(outStr, "value: hello world"))
	test.AssertResult(t, true, strings.Contains(outStr, "line: 5"))
	test.AssertResult(t, true, strings.Contains(outStr, "column: 7"))
	test.AssertResult(t, true, strings.Contains(outStr, "headComment: head"))
	test.AssertResult(t, true, strings.Contains(outStr, "lineComment: line"))
	test.AssertResult(t, true, strings.Contains(outStr, "footComment: foot"))
	test.AssertResult(t, true, strings.Contains(outStr, "anchor: anchor"))
}

func TestNodeInfoPrinter_PrintedAnything_True(t *testing.T) {
	node := &CandidateNode{
		Kind:  ScalarNode,
		Tag:   "!!str",
		Value: "test",
	}
	listNodes := list.New()
	listNodes.PushBack(node)

	var output bytes.Buffer
	writer := bufio.NewWriter(&output)
	printer := NewNodeInfoPrinter(NewSinglePrinterWriter(writer))

	// Before printing, should be false
	test.AssertResult(t, false, printer.PrintedAnything())

	err := printer.PrintResults(listNodes)
	writer.Flush()
	if err != nil {
		t.Fatalf("PrintResults error: %v", err)
	}

	// After printing, should be true
	test.AssertResult(t, true, printer.PrintedAnything())
}

func TestNodeInfoPrinter_PrintedAnything_False(t *testing.T) {
	listNodes := list.New()

	var output bytes.Buffer
	writer := bufio.NewWriter(&output)
	printer := NewNodeInfoPrinter(NewSinglePrinterWriter(writer))

	err := printer.PrintResults(listNodes)
	writer.Flush()
	if err != nil {
		t.Fatalf("PrintResults error: %v", err)
	}

	// No nodes printed, should still be false
	test.AssertResult(t, false, printer.PrintedAnything())
}

func TestNodeInfoPrinter_SetNulSepOutput(_ *testing.T) {
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)
	printer := NewNodeInfoPrinter(NewSinglePrinterWriter(writer))

	// Should not panic or error
	printer.SetNulSepOutput(true)
	printer.SetNulSepOutput(false)
}

func TestNodeInfoPrinter_SetAppendix(t *testing.T) {
	node := &CandidateNode{
		Kind:  ScalarNode,
		Tag:   "!!str",
		Value: "test",
	}
	listNodes := list.New()
	listNodes.PushBack(node)

	var output bytes.Buffer
	writer := bufio.NewWriter(&output)
	printer := NewNodeInfoPrinter(NewSinglePrinterWriter(writer))

	appendixText := "This is appendix text\n"
	appendixReader := strings.NewReader(appendixText)
	printer.SetAppendix(appendixReader)

	err := printer.PrintResults(listNodes)
	writer.Flush()
	if err != nil {
		t.Fatalf("PrintResults error: %v", err)
	}

	outStr := output.String()
	test.AssertResult(t, true, strings.Contains(outStr, "test"))
	test.AssertResult(t, true, strings.Contains(outStr, appendixText))
}

func TestNodeInfoPrinter_MultipleNodes(t *testing.T) {
	node1 := &CandidateNode{
		Kind:  ScalarNode,
		Tag:   "!!str",
		Value: "first",
	}
	node2 := &CandidateNode{
		Kind:  ScalarNode,
		Tag:   "!!str",
		Value: "second",
	}
	listNodes := list.New()
	listNodes.PushBack(node1)
	listNodes.PushBack(node2)

	var output bytes.Buffer
	writer := bufio.NewWriter(&output)
	printer := NewNodeInfoPrinter(NewSinglePrinterWriter(writer))

	err := printer.PrintResults(listNodes)
	writer.Flush()
	if err != nil {
		t.Fatalf("PrintResults error: %v", err)
	}

	outStr := output.String()
	test.AssertResult(t, true, strings.Contains(outStr, "value: first"))
	test.AssertResult(t, true, strings.Contains(outStr, "value: second"))
}

func TestNodeInfoPrinter_SequenceNode(t *testing.T) {
	node := &CandidateNode{
		Kind:  SequenceNode,
		Tag:   "!!seq",
		Style: FlowStyle,
	}
	listNodes := list.New()
	listNodes.PushBack(node)

	var output bytes.Buffer
	writer := bufio.NewWriter(&output)
	printer := NewNodeInfoPrinter(NewSinglePrinterWriter(writer))

	err := printer.PrintResults(listNodes)
	writer.Flush()
	if err != nil {
		t.Fatalf("PrintResults error: %v", err)
	}

	outStr := output.String()
	test.AssertResult(t, true, strings.Contains(outStr, "kind: SequenceNode"))
	test.AssertResult(t, true, strings.Contains(outStr, "tag: '!!seq'"))
	test.AssertResult(t, true, strings.Contains(outStr, "style: FlowStyle"))
}

func TestNodeInfoPrinter_MappingNode(t *testing.T) {
	node := &CandidateNode{
		Kind: MappingNode,
		Tag:  "!!map",
	}
	listNodes := list.New()
	listNodes.PushBack(node)

	var output bytes.Buffer
	writer := bufio.NewWriter(&output)
	printer := NewNodeInfoPrinter(NewSinglePrinterWriter(writer))

	err := printer.PrintResults(listNodes)
	writer.Flush()
	if err != nil {
		t.Fatalf("PrintResults error: %v", err)
	}

	outStr := output.String()
	test.AssertResult(t, true, strings.Contains(outStr, "kind: MappingNode"))
	test.AssertResult(t, true, strings.Contains(outStr, "tag: '!!map'"))
}

func TestNodeInfoPrinter_EmptyList(t *testing.T) {
	listNodes := list.New()

	var output bytes.Buffer
	writer := bufio.NewWriter(&output)
	printer := NewNodeInfoPrinter(NewSinglePrinterWriter(writer))

	err := printer.PrintResults(listNodes)
	writer.Flush()
	if err != nil {
		t.Fatalf("PrintResults error: %v", err)
	}

	test.AssertResult(t, "", output.String())
	test.AssertResult(t, false, printer.PrintedAnything())
}
