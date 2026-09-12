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
	noClobber      bool
}

func NewMultiPrinterWriter(expression *ExpressionNode, format *Format, noClobber bool) PrinterWriter {
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
		noClobber:      noClobber,
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

	openFlags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	if sp.noClobber {
		openFlags = os.O_WRONLY | os.O_CREATE | os.O_EXCL
	}
	f, err := os.OpenFile(name, openFlags, 0666)

	if err != nil {
		if sp.noClobber && os.IsExist(err) {
			return nil, fmt.Errorf("split file %q already exists", name)
		}
		return nil, err
	}
	sp.index = sp.index + 1

	return bufio.NewWriter(f), nil

}
