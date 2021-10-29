package yqlib

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type PrinterWriter interface {
	GetWriter(node *CandidateNode, index int) (*bufio.Writer, error)
}

type singlePrinterWriter struct {
	bufferedWriter *bufio.Writer
}

func NewSinglePrinterWriter(writer io.Writer) PrinterWriter {
	return &singlePrinterWriter{
		bufferedWriter: bufio.NewWriter(writer),
	}
}

func (sp *singlePrinterWriter) GetWriter(node *CandidateNode, i int) (*bufio.Writer, error) {
	return sp.bufferedWriter, nil
}

type multiPrintWriter struct {
	treeNavigator  DataTreeNavigator
	nameExpression *ExpressionNode
	extension      string
}

func NewMultiPrinterWriter(expression *ExpressionNode, format PrinterOutputFormat) PrinterWriter {
	extension := "yml"

	switch format {
	case JsonOutputFormat:
		extension = "json"
	case PropsOutputFormat:
		extension = "properties"
	}

	return &multiPrintWriter{
		nameExpression: expression,
		extension:      extension,
		treeNavigator:  NewDataTreeNavigator(),
	}
}

func (sp *multiPrintWriter) GetWriter(node *CandidateNode, index int) (*bufio.Writer, error) {
	name := ""

	if sp.nameExpression != nil {
		context := Context{MatchingNodes: node.AsList()}
		result, err := sp.treeNavigator.GetMatchingNodes(context, sp.nameExpression)
		if err != nil {
			return nil, err
		}
		if result.MatchingNodes.Len() > 0 {
			name = result.MatchingNodes.Front().Value.(*CandidateNode).Node.Value
		}
	}
	if name == "" {
		name = fmt.Sprintf("%v.%v", index, sp.extension)
	} else {
		name = fmt.Sprintf("%v.%v", name, sp.extension)
	}

	f, err := os.Create(name)

	if err != nil {
		return nil, err
	}

	return bufio.NewWriter(f), nil

}
