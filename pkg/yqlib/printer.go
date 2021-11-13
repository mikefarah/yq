package yqlib

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"regexp"

	yaml "gopkg.in/yaml.v3"
)

type Printer interface {
	PrintResults(matchingNodes *list.List) error
	PrintedAnything() bool
	//e.g. when given a front-matter doc, like jekyll
	SetAppendix(reader io.Reader)
}

type PrinterOutputFormat uint32

const (
	YamlOutputFormat = 1 << iota
	JsonOutputFormat
	PropsOutputFormat
)

func OutputFormatFromString(format string) (PrinterOutputFormat, error) {
	switch format {
	case "yaml", "y":
		return YamlOutputFormat, nil
	case "json", "j":
		return JsonOutputFormat, nil
	case "props", "p":
		return PropsOutputFormat, nil
	default:
		return 0, fmt.Errorf("Unknown fromat '%v' please use [yaml|json|props]", format)
	}
}

type resultsPrinter struct {
	outputFormat       PrinterOutputFormat
	unwrapScalar       bool
	colorsEnabled      bool
	indent             int
	printDocSeparators bool
	printerWriter      PrinterWriter
	firstTimePrinting  bool
	previousDocIndex   uint
	previousFileIndex  int
	printedMatches     bool
	treeNavigator      DataTreeNavigator
	appendixReader     io.Reader
}

func NewPrinterWithSingleWriter(writer io.Writer, outputFormat PrinterOutputFormat, unwrapScalar bool, colorsEnabled bool, indent int, printDocSeparators bool) Printer {
	return NewPrinter(NewSinglePrinterWriter(writer), outputFormat, unwrapScalar, colorsEnabled, indent, printDocSeparators)
}

func NewPrinter(printerWriter PrinterWriter, outputFormat PrinterOutputFormat, unwrapScalar bool, colorsEnabled bool, indent int, printDocSeparators bool) Printer {
	return &resultsPrinter{
		printerWriter:      printerWriter,
		outputFormat:       outputFormat,
		unwrapScalar:       unwrapScalar,
		colorsEnabled:      colorsEnabled,
		indent:             indent,
		printDocSeparators: outputFormat == YamlOutputFormat && printDocSeparators,
		firstTimePrinting:  true,
		treeNavigator:      NewDataTreeNavigator(),
	}
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
	if node.Kind == yaml.ScalarNode && p.unwrapScalar && p.outputFormat == YamlOutputFormat {
		return writeString(writer, node.Value+"\n")
	}

	if p.outputFormat == JsonOutputFormat {
		encoder = NewJsonEncoder(writer, p.indent)
	} else if p.outputFormat == PropsOutputFormat {
		encoder = NewPropertiesEncoder(writer)
	} else {
		encoder = NewYamlEncoder(writer, p.indent, p.colorsEnabled)
	}
	return encoder.Encode(node)
}

func (p *resultsPrinter) PrintResults(matchingNodes *list.List) error {
	log.Debug("PrintResults for %v matches", matchingNodes.Len())

	if matchingNodes.Len() == 0 {
		log.Debug("no matching results, nothing to print")
		return nil
	}

	if p.outputFormat != YamlOutputFormat {
		explodeOp := Operation{OperationType: explodeOpType}
		explodeNode := ExpressionNode{Operation: &explodeOp}
		context, err := p.treeNavigator.GetMatchingNodes(Context{MatchingNodes: matchingNodes}, &explodeNode)
		if err != nil {
			return err
		}
		matchingNodes = context.MatchingNodes
	}

	if p.firstTimePrinting {
		node := matchingNodes.Front().Value.(*CandidateNode)
		p.previousDocIndex = node.Document
		p.previousFileIndex = node.FileIndex
		p.firstTimePrinting = false
	}

	for el := matchingNodes.Front(); el != nil; el = el.Next() {

		mappedDoc := el.Value.(*CandidateNode)
		log.Debug("-- print sep logic: p.firstTimePrinting: %v, previousDocIndex: %v, mappedDoc.Document: %v, printDocSeparators: %v", p.firstTimePrinting, p.previousDocIndex, mappedDoc.Document, p.printDocSeparators)

		writer, errorWriting := p.printerWriter.GetWriter(mappedDoc)
		if errorWriting != nil {
			return errorWriting
		}

		commentsStartWithSepExp := regexp.MustCompile(`^\$yqDocSeperator\$`)
		commentStartsWithSeparator := commentsStartWithSepExp.MatchString(mappedDoc.LeadingContent)

		if (p.previousDocIndex != mappedDoc.Document || p.previousFileIndex != mappedDoc.FileIndex) && p.printDocSeparators && !commentStartsWithSeparator {
			log.Debug("-- writing doc sep")
			if err := writeString(writer, "---\n"); err != nil {
				return err
			}
		}

		if err := processLeadingContent(mappedDoc, writer, p.printDocSeparators, p.outputFormat); err != nil {
			return err
		}

		if err := p.printNode(mappedDoc.Node, writer); err != nil {
			return err
		}

		p.previousDocIndex = mappedDoc.Document
		if err := writer.Flush(); err != nil {
			return err
		}
	}

	if p.appendixReader != nil && p.outputFormat == YamlOutputFormat {
		writer, err := p.printerWriter.GetWriter(nil)
		if err != nil {
			return err
		}

		log.Debug("Piping appendix reader...")
		betterReader := bufio.NewReader(p.appendixReader)
		_, err = io.Copy(writer, betterReader)
		if err != nil {
			return err
		}
		if err := writer.Flush(); err != nil {
			return err
		}
	}

	return nil
}
