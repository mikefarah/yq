package yqlib

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
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

func (sp *singlePrinterWriter) GetWriter(_ *CandidateNode) (*bufio.Writer, error) {
	return sp.bufferedWriter, nil
}

type multiPrintWriter struct {
	treeNavigator  DataTreeNavigator
	nameExpression *ExpressionNode
	extension      string
	index          int
}

func NewMultiPrinterWriter(expression *ExpressionNode, format *Format) PrinterWriter {
	extension := "yml"

	switch format {
	case JSONFormat:
		extension = "json"
	case PropertiesFormat:
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

	indexVariableNode := CandidateNode{Kind: ScalarNode, Tag: "!!int", Value: fmt.Sprintf("%v", sp.index)}

	context := Context{MatchingNodes: node.AsList()}
	context.SetVariable("index", indexVariableNode.AsList())
	result, err := sp.treeNavigator.GetMatchingNodes(context, sp.nameExpression)
	if err != nil {
		return nil, err
	}
	if result.MatchingNodes.Len() > 0 {
		name = result.MatchingNodes.Front().Value.(*CandidateNode).Value
	}
	var extensionRegexp = regexp.MustCompile(`\.[a-zA-Z0-9]+$`)
	if !extensionRegexp.MatchString(name) {
		name = fmt.Sprintf("%v.%v", name, sp.extension)
	}

	err = os.MkdirAll(filepath.Dir(name), 0750)
	if err != nil {
		return nil, err
	}
	f, err := os.Create(name)

	if err != nil {
		return nil, err
	}
	sp.index = sp.index + 1

	return bufio.NewWriter(f), nil

}
