package yqlib

import (
	"container/list"
	"errors"
	"fmt"
	"io"
	"os"
)

// A yaml expression evaluator that runs the expression multiple times for each given yaml document.
// Uses less memory than loading all documents and running the expression once, but this cannot process
// cross document expressions.
type StreamEvaluator interface {
	Evaluate(filename string, reader io.Reader, node *ExpressionNode, printer Printer, decoder Decoder) (uint, error)
	EvaluateFiles(expression string, filenames []string, printer Printer, decoder Decoder) error
	EvaluateNew(expression string, printer Printer) error
}

type streamEvaluator struct {
	treeNavigator DataTreeNavigator
	fileIndex     int
}

func NewStreamEvaluator() StreamEvaluator {
	return &streamEvaluator{treeNavigator: NewDataTreeNavigator()}
}

func (s *streamEvaluator) EvaluateNew(expression string, printer Printer) error {
	node, err := ExpressionParser.ParseExpression(expression)
	if err != nil {
		return err
	}
	candidateNode := createScalarNode(nil, "")
	inputList := list.New()
	inputList.PushBack(candidateNode)

	result, errorParsing := s.treeNavigator.GetMatchingNodes(Context{MatchingNodes: inputList}, node)
	if errorParsing != nil {
		return errorParsing
	}
	return printer.PrintResults(result.MatchingNodes)
}

func (s *streamEvaluator) EvaluateFiles(expression string, filenames []string, printer Printer, decoder Decoder) error {
	var totalProcessDocs uint
	node, err := ExpressionParser.ParseExpression(expression)
	if err != nil {
		return err
	}

	for _, filename := range filenames {
		reader, err := readStream(filename)

		if err != nil {
			return err
		}
		processedDocs, err := s.Evaluate(filename, reader, node, printer, decoder)
		if err != nil {
			return err
		}
		totalProcessDocs = totalProcessDocs + processedDocs

		switch reader := reader.(type) {
		case *os.File:
			safelyCloseFile(reader)
		}
	}

	if totalProcessDocs == 0 {
		// problem is I've already slurped the leading content sadface
		return s.EvaluateNew(expression, printer)
	}

	return nil
}

func (s *streamEvaluator) Evaluate(filename string, reader io.Reader, node *ExpressionNode, printer Printer, decoder Decoder) (uint, error) {

	var currentIndex uint
	err := decoder.Init(reader)
	if err != nil {
		return 0, err
	}
	for {
		candidateNode, errorReading := decoder.Decode()

		if errors.Is(errorReading, io.EOF) {
			s.fileIndex = s.fileIndex + 1
			return currentIndex, nil
		} else if errorReading != nil {
			return currentIndex, fmt.Errorf("bad file '%v': %w", filename, errorReading)
		}
		candidateNode.document = currentIndex
		candidateNode.filename = filename
		candidateNode.fileIndex = s.fileIndex

		inputList := list.New()
		inputList.PushBack(candidateNode)

		result, errorParsing := s.treeNavigator.GetMatchingNodes(Context{MatchingNodes: inputList}, node)
		if errorParsing != nil {
			return currentIndex, errorParsing
		}
		err := printer.PrintResults(result.MatchingNodes)

		if err != nil {
			return currentIndex, err
		}
		currentIndex = currentIndex + 1
	}
}
