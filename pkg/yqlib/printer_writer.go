package yqlib

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type PrinterWriter interface {
	GetWriter(node *CandidateNode) (*bufio.Writer, error)
}

type singlePrinterWriter struct {
	bufferedWriter *bufio.Writer
}

func NewSinglePrinterWriter(writer io.Writer) PrinterWriter {
	return &singlePrinterWriter{
		bufferedWriter: bufio.NewWriter(writer),
	}
}

func (sp *singlePrinterWriter) GetWriter(node *CandidateNode) (*bufio.Writer, error) {
	return sp.bufferedWriter, nil
}

type multiPrintWriter struct {
	treeNavigator  DataTreeNavigator
	nameExpression *ExpressionNode
	extension      string
	index          int
}

func NewMultiPrinterWriter(expression *ExpressionNode, format PrinterOutputFormat) PrinterWriter {
	extension := "yml"

	switch format {
	case JSONOutputFormat:
		extension = "json"
	case PropsOutputFormat:
		extension = "properties"
	}

	return &multiPrintWriter{
		nameExpression: expression,
		extension:      extension,
		treeNavigator:  NewDataTreeNavigator(),
		index:          0,
	}
}

func (sp *multiPrintWriter) GetWriter(node *CandidateNode) (*bufio.Writer, error) {
	name := ""

	indexVariableNode := yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: fmt.Sprintf("%v", sp.index)}
	indexVariableCandidate := CandidateNode{Node: &indexVariableNode}

	context := Context{MatchingNodes: node.AsList()}
	context.SetVariable("index", indexVariableCandidate.AsList())
	result, err := sp.treeNavigator.GetMatchingNodes(context, sp.nameExpression)
	if err != nil {
		return nil, err
	}
	if result.MatchingNodes.Len() > 0 {
		name = result.MatchingNodes.Front().Value.(*CandidateNode).Node.Value
	}
	var extensionRegexp = regexp.MustCompile(`\.[a-zA-Z0-9]+$`)
	if !extensionRegexp.MatchString(name) {
		name = fmt.Sprintf("%v.%v", name, sp.extension)
	}

	f, err := os.Create(name)

	if err != nil {
		return nil, err
	}
	sp.index = sp.index + 1

	return bufio.NewWriter(f), nil

}
