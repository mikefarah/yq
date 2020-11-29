package yqlib

import (
	"bufio"
	"container/list"
	"io"

	yaml "gopkg.in/yaml.v3"
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
	previousDocIndex   uint
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

func (p *resultsPrinter) safelyFlush(writer *bufio.Writer) {
	if err := writer.Flush(); err != nil {
		log.Error("Error flushing writer!")
		log.Error(err.Error())
	}
}

func (p *resultsPrinter) PrintResults(matchingNodes *list.List) error {
	log.Debug("PrintResults for %v matches", matchingNodes.Len())
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
	defer p.safelyFlush(bufferedWriter)

	if matchingNodes.Len() == 0 {
		log.Debug("no matching results, nothing to print")
		return nil
	}
	if p.firstTimePrinting {
		p.previousDocIndex = matchingNodes.Front().Value.(*CandidateNode).Document
		p.firstTimePrinting = false
	}

	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		mappedDoc := el.Value.(*CandidateNode)
		log.Debug("-- print sep logic: p.firstTimePrinting: %v, previousDocIndex: %v, mappedDoc.Document: %v, printDocSeparators: %v", p.firstTimePrinting, p.previousDocIndex, mappedDoc.Document, p.printDocSeparators)
		if (p.previousDocIndex != mappedDoc.Document) && p.printDocSeparators {
			log.Debug("-- writing doc sep")
			if err := p.writeString(bufferedWriter, "---\n"); err != nil {
				return err
			}

		}

		if err := p.printNode(mappedDoc.Node, bufferedWriter); err != nil {
			return err
		}

		p.previousDocIndex = mappedDoc.Document
	}

	return nil
}
