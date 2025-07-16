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
