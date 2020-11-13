package yqlib

import (
	"bufio"
	"container/list"
	"io"

	"gopkg.in/yaml.v3"
)

type Printer interface {
	PrintResults(matchingNodes *list.List) error
}

type resultsPrinter struct {
	outputToJSON       bool
	unwrapScalar       bool
	colorsEnabled      bool
	indent             int
	printDocSeparators bool
	writer             io.Writer
	firstTimePrinting  bool
}

func NewPrinter(writer io.Writer, outputToJSON bool, unwrapScalar bool, colorsEnabled bool, indent int, printDocSeparators bool) Printer {
	return &resultsPrinter{
		writer:             writer,
		outputToJSON:       outputToJSON,
		unwrapScalar:       unwrapScalar,
		colorsEnabled:      colorsEnabled,
		indent:             indent,
		printDocSeparators: printDocSeparators,
		firstTimePrinting:  true,
	}
}

func (p *resultsPrinter) printNode(node *yaml.Node, writer io.Writer) error {
	var encoder Encoder
	if node.Kind == yaml.ScalarNode && p.unwrapScalar && !p.outputToJSON {
		return p.writeString(writer, node.Value+"\n")
	}
	if p.outputToJSON {
		encoder = NewJsonEncoder(writer, p.indent)
	} else {
		encoder = NewYamlEncoder(writer, p.indent, p.colorsEnabled)
	}
	return encoder.Encode(node)
}

func (p *resultsPrinter) writeString(writer io.Writer, txt string) error {
	_, errorWriting := writer.Write([]byte(txt))
	return errorWriting
}

func (p *resultsPrinter) PrintResults(matchingNodes *list.List) error {
	var err error
	if p.outputToJSON {
		explodeOp := Operation{OperationType: Explode}
		explodeNode := PathTreeNode{Operation: &explodeOp}
		matchingNodes, err = treeNavigator.GetMatchingNodes(matchingNodes, &explodeNode)
		if err != nil {
			return err
		}
	}

	bufferedWriter := bufio.NewWriter(p.writer)
	defer safelyFlush(bufferedWriter)

	if matchingNodes.Len() == 0 {
		log.Debug("no matching results, nothing to print")
		return nil
	}

	previousDocIndex := matchingNodes.Front().Value.(*CandidateNode).Document

	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		mappedDoc := el.Value.(*CandidateNode)

		if (!p.firstTimePrinting || (previousDocIndex != mappedDoc.Document)) && p.printDocSeparators {
			p.writeString(bufferedWriter, "---\n")
		}

		if err := p.printNode(mappedDoc.Node, bufferedWriter); err != nil {
			return err
		}

		previousDocIndex = mappedDoc.Document
	}
	p.firstTimePrinting = false

	return nil
}
