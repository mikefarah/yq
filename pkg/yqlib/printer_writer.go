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
	noOverwrite    bool
}

func NewMultiPrinterWriter(expression *ExpressionNode, format *Format) PrinterWriter {
	return NewMultiPrinterWriterWithOptions(expression, format, false)
}

// NewMultiPrinterWriterWithOptions creates a multi-file printer writer.
// When noOverwrite is true, attempting to write to a file that already
// exists will fail with an error instead of silently overwriting it.
func NewMultiPrinterWriterWithOptions(expression *ExpressionNode, format *Format, noOverwrite bool) PrinterWriter {
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
		noOverwrite:    noOverwrite,
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
	var f *os.File
	if sp.noOverwrite {
		f, err = os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			if os.IsExist(err) {
				return nil, fmt.Errorf("refusing to overwrite existing file %q (--no-overwrite is set)", name)
			}
			return nil, err
		}
	} else {
		f, err = os.Create(name)
		if err != nil {
			return nil, err
		}
	}
	sp.index = sp.index + 1

	return bufio.NewWriter(f), nil

}
