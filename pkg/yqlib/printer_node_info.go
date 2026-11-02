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
		err := usePrinterWriter(p.printerWriter, mappedDoc, func(writer *bufio.Writer) error {
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
			return writer.Flush()
		})
		if err != nil {
			return err
		}
	}

	if p.appendixReader != nil {
		return usePrinterWriter(p.printerWriter, nil, func(writer *bufio.Writer) error {
			log.Debug("Piping appendix reader...")
			betterReader := bufio.NewReader(p.appendixReader)
			if _, err := io.Copy(writer, betterReader); err != nil {
				return err
			}
			return writer.Flush()
		})
	}

	return nil
}
