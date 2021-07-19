package yqlib

import (
	"bufio"
	"container/list"
	"io"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type Printer interface {
	PrintResults(matchingNodes *list.List, leadingContent string) error
	PrintedAnything() bool
	//e.g. when given a front-matter doc, like jekyll
	SetAppendix(reader io.Reader)
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
	previousFileIndex  int
	printedMatches     bool
	treeNavigator      DataTreeNavigator
	preambleReader     io.Reader
	appendixReader     io.Reader
}

func NewPrinter(writer io.Writer, outputToJSON bool, unwrapScalar bool, colorsEnabled bool, indent int, printDocSeparators bool) Printer {
	return &resultsPrinter{
		writer:             writer,
		outputToJSON:       outputToJSON,
		unwrapScalar:       unwrapScalar,
		colorsEnabled:      colorsEnabled,
		indent:             indent,
		printDocSeparators: !outputToJSON && printDocSeparators,
		firstTimePrinting:  true,
		treeNavigator:      NewDataTreeNavigator(),
	}
}

func (p *resultsPrinter) SetPreamble(reader io.Reader) {
	p.preambleReader = reader
}

func (p *resultsPrinter) SetAppendix(reader io.Reader) {
	p.appendixReader = reader
}

func (p *resultsPrinter) PrintedAnything() bool {
	return p.printedMatches
}

func (p *resultsPrinter) printNode(node *yaml.Node, writer io.Writer) error {
	p.printedMatches = p.printedMatches || (node.Tag != "!!null" &&
		(node.Tag != "!!bool" || node.Value != "false"))

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

func (p *resultsPrinter) PrintResults(matchingNodes *list.List, leadingContent string) error {
	log.Debug("PrintResults for %v matches", matchingNodes.Len())
	if p.outputToJSON {
		explodeOp := Operation{OperationType: explodeOpType}
		explodeNode := ExpressionNode{Operation: &explodeOp}
		context, err := p.treeNavigator.GetMatchingNodes(Context{MatchingNodes: matchingNodes}, &explodeNode)
		if err != nil {
			return err
		}
		matchingNodes = context.MatchingNodes
	}

	bufferedWriter := bufio.NewWriter(p.writer)
	defer p.safelyFlush(bufferedWriter)

	if matchingNodes.Len() == 0 {
		p.writeString(bufferedWriter, leadingContent)
		log.Debug("no matching results, nothing to print")
		return nil
	}
	if p.firstTimePrinting {
		node := matchingNodes.Front().Value.(*CandidateNode)
		p.previousDocIndex = node.Document
		p.previousFileIndex = node.FileIndex
		p.firstTimePrinting = false
	}

	printedLead := false

	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		mappedDoc := el.Value.(*CandidateNode)
		log.Debug("-- print sep logic: p.firstTimePrinting: %v, previousDocIndex: %v, mappedDoc.Document: %v, printDocSeparators: %v", p.firstTimePrinting, p.previousDocIndex, mappedDoc.Document, p.printDocSeparators)
		if (p.previousDocIndex != mappedDoc.Document || p.previousFileIndex != mappedDoc.FileIndex) && p.printDocSeparators &&
			(printedLead || !strings.HasPrefix(leadingContent, "---")) {
			log.Debug("-- writing doc sep")
			if err := p.writeString(bufferedWriter, "---\n"); err != nil {
				return err
			}
		}

		if !printedLead {
			// we want to print this after the seperator logic
			p.writeString(bufferedWriter, leadingContent)
			printedLead = true
		}

		if err := p.printNode(mappedDoc.Node, bufferedWriter); err != nil {
			return err
		}

		p.previousDocIndex = mappedDoc.Document
	}

	if p.appendixReader != nil && !p.outputToJSON {
		log.Debug("Piping appendix reader...")
		betterReader := bufio.NewReader(p.appendixReader)
		_, err := io.Copy(bufferedWriter, betterReader)
		if err != nil {
			return err
		}
	}

	return nil
}
