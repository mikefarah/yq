package yqlib

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"strings"

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
	printerWriter      printerWriter
	firstTimePrinting  bool
	previousDocIndex   uint
	previousFileIndex  int
	printedMatches     bool
	treeNavigator      DataTreeNavigator
	appendixReader     io.Reader
}

func NewPrinter(printerWriter printerWriter, outputFormat PrinterOutputFormat, unwrapScalar bool, colorsEnabled bool, indent int, printDocSeparators bool) Printer {
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
		return p.writeString(writer, node.Value+"\n")
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

func (p *resultsPrinter) processLeadingContent(mappedDoc *CandidateNode, writer io.Writer) error {
	if strings.Contains(mappedDoc.Node.HeadComment, "$yqLeadingContent$") {
		log.Debug("headcommentwas %v", mappedDoc.Node.HeadComment)
		log.Debug("finished headcomment")
		reader := bufio.NewReader(strings.NewReader(mappedDoc.Node.HeadComment))
		mappedDoc.Node.HeadComment = ""

		for {

			readline, errReading := reader.ReadString('\n')
			if errReading != nil && errReading != io.EOF {
				return errReading
			}
			if strings.Contains(readline, "$yqLeadingContent$") {
				// skip this

			} else if strings.Contains(readline, "$yqDocSeperator$") {
				if p.printDocSeparators {
					if err := p.writeString(writer, "---\n"); err != nil {
						return err
					}
				}
			} else if p.outputFormat == YamlOutputFormat {
				if err := p.writeString(writer, readline); err != nil {
					return err
				}
			}

			if errReading == io.EOF {
				if readline != "" {
					// the last comment we read didn't have a new line, put one in
					if err := p.writeString(writer, "\n"); err != nil {
						return err
					}
				}
				break
			}
		}

	}
	return nil
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

	index := 0

	for el := matchingNodes.Front(); el != nil; el = el.Next() {

		mappedDoc := el.Value.(*CandidateNode)
		log.Debug("-- print sep logic: p.firstTimePrinting: %v, previousDocIndex: %v, mappedDoc.Document: %v, printDocSeparators: %v", p.firstTimePrinting, p.previousDocIndex, mappedDoc.Document, p.printDocSeparators)

		writer, errorWriting := p.printerWriter.GetWriter(mappedDoc, index)
		if errorWriting != nil {
			return errorWriting
		}

		commentStartsWithSeparator := strings.Contains(mappedDoc.Node.HeadComment, "$yqLeadingContent$\n$yqDocSeperator$")

		if (p.previousDocIndex != mappedDoc.Document || p.previousFileIndex != mappedDoc.FileIndex) && p.printDocSeparators && !commentStartsWithSeparator {
			log.Debug("-- writing doc sep")
			if err := p.writeString(writer, "---\n"); err != nil {
				return err
			}
		}

		if err := p.processLeadingContent(mappedDoc, writer); err != nil {
			return err
		}

		if err := p.printNode(mappedDoc.Node, writer); err != nil {
			return err
		}

		p.previousDocIndex = mappedDoc.Document
		if err := writer.Flush(); err != nil {
			return err
		}

		index++

	}

	if p.appendixReader != nil && p.outputFormat == YamlOutputFormat {
		writer := p.printerWriter.GetWriter(nil, index)
		log.Debug("Piping appendix reader...")
		betterReader := bufio.NewReader(p.appendixReader)
		_, err := io.Copy(writer, betterReader)
		if err != nil {
			return err
		}
	}

	return nil
}
