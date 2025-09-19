package yqlib

import (
	"bufio"
	"container/list"
	"io"

	"go.yaml.in/yaml/v4"
)

type nodeInfoPrinter struct {
	printerWriter  PrinterWriter
	appendixReader io.Reader
	printedMatches bool
}

func NewNodeInfoPrinter(printerWriter PrinterWriter) Printer {
	return &nodeInfoPrinter{
		printerWriter: printerWriter,
	}
}

func (p *nodeInfoPrinter) SetNulSepOutput(_ bool) {
}

func (p *nodeInfoPrinter) SetAppendix(reader io.Reader) {
	p.appendixReader = reader
}

func (p *nodeInfoPrinter) PrintedAnything() bool {
	return p.printedMatches
}

func (p *nodeInfoPrinter) PrintResults(matchingNodes *list.List) error {

	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		mappedDoc := el.Value.(*CandidateNode)
		writer, errorWriting := p.printerWriter.GetWriter(mappedDoc)
		if errorWriting != nil {
			return errorWriting
		}
		bytes, err := yaml.Marshal(mappedDoc.ConvertToNodeInfo())
		if err != nil {
			return err
		}
		if _, err := writer.Write(bytes); err != nil {
			return err
		}
		if _, err := writer.Write([]byte("\n")); err != nil {
			return err
		}
		p.printedMatches = true
		if err := writer.Flush(); err != nil {
			return err
		}
	}

	if p.appendixReader != nil {
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
