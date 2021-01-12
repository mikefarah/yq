package yqlib

import (
	"container/list"
	"io"
	"os"

	yaml "gopkg.in/yaml.v3"
)

// A yaml expression evaluator that runs the expression multiple times for each given yaml document.
// Uses less memory than loading all documents and running the expression once, but this cannot process
// cross document expressions.
type StreamEvaluator interface {
	Evaluate(filename string, reader io.Reader, node *ExpressionNode, printer Printer) error
	EvaluateFiles(expression string, filenames []string, printer Printer) error
	EvaluateNew(expression string, printer Printer) error
}

type streamEvaluator struct {
	treeNavigator DataTreeNavigator
	treeCreator   ExpressionParser
	fileIndex     int
}

func NewStreamEvaluator() StreamEvaluator {
	return &streamEvaluator{treeNavigator: NewDataTreeNavigator(), treeCreator: NewExpressionParser()}
}

func (s *streamEvaluator) EvaluateNew(expression string, printer Printer) error {
	node, err := s.treeCreator.ParseExpression(expression)
	if err != nil {
		return err
	}
	candidateNode := &CandidateNode{
		Document:  0,
		Filename:  "",
		Node:      &yaml.Node{Tag: "!!null"},
		FileIndex: 0,
	}
	inputList := list.New()
	inputList.PushBack(candidateNode)

	matches, errorParsing := s.treeNavigator.GetMatchingNodes(inputList, node)
	if errorParsing != nil {
		return errorParsing
	}
	return printer.PrintResults(matches)
}

func (s *streamEvaluator) EvaluateFiles(expression string, filenames []string, printer Printer) error {

	node, err := s.treeCreator.ParseExpression(expression)
	if err != nil {
		return err
	}

	for _, filename := range filenames {
		reader, err := readStream(filename)
		if err != nil {
			return err
		}
		err = s.Evaluate(filename, reader, node, printer)
		if err != nil {
			return err
		}

		switch reader := reader.(type) {
		case *os.File:
			safelyCloseFile(reader)
		}
	}
	return nil
}

func (s *streamEvaluator) Evaluate(filename string, reader io.Reader, node *ExpressionNode, printer Printer) error {

	var currentIndex uint

	decoder := yaml.NewDecoder(reader)
	for {
		var dataBucket yaml.Node
		errorReading := decoder.Decode(&dataBucket)

		if errorReading == io.EOF {
			s.fileIndex = s.fileIndex + 1
			return nil
		} else if errorReading != nil {
			return errorReading
		}
		candidateNode := &CandidateNode{
			Document:  currentIndex,
			Filename:  filename,
			Node:      &dataBucket,
			FileIndex: s.fileIndex,
		}
		inputList := list.New()
		inputList.PushBack(candidateNode)

		matches, errorParsing := s.treeNavigator.GetMatchingNodes(inputList, node)
		if errorParsing != nil {
			return errorParsing
		}
		err := printer.PrintResults(matches)
		if err != nil {
			return err
		}
		currentIndex = currentIndex + 1
	}
}
