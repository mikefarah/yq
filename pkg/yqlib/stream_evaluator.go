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
	Evaluate(filename string, reader io.Reader, node *ExpressionNode, printer Printer) (uint, error)
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
		Node:      &yaml.Node{Tag: "!!null", Kind: yaml.ScalarNode},
		FileIndex: 0,
	}
	inputList := list.New()
	inputList.PushBack(candidateNode)

	result, errorParsing := s.treeNavigator.GetMatchingNodes(Context{MatchingNodes: inputList}, node)
	if errorParsing != nil {
		return errorParsing
	}
	return printer.PrintResults(result.MatchingNodes)
}

func (s *streamEvaluator) EvaluateFiles(expression string, filenames []string, printer Printer) error {
	var totalProcessDocs uint = 0
	node, err := s.treeCreator.ParseExpression(expression)
	if err != nil {
		return err
	}

	for index, filename := range filenames {
		reader, leadingSeperator, err := readStream(filename)
		if index == 0 && leadingSeperator {
			printer.SetPrintLeadingSeperator(leadingSeperator)
		}
		if err != nil {
			return err
		}
		processedDocs, err := s.Evaluate(filename, reader, node, printer)
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

		if len(filenames) > 0 {
			reader, _, err := readStream(filenames[0])
			if err != nil {
				return err
			}
			switch reader := reader.(type) {
			case *os.File:
				defer safelyCloseFile(reader)
			}
			printer.SetPreamble(reader)
		}

		return s.EvaluateNew(expression, printer)
	}

	return nil
}

func (s *streamEvaluator) Evaluate(filename string, reader io.Reader, node *ExpressionNode, printer Printer) (uint, error) {

	var currentIndex uint
	decoder := yaml.NewDecoder(reader)
	for {
		var dataBucket yaml.Node
		errorReading := decoder.Decode(&dataBucket)

		if errorReading == io.EOF {
			s.fileIndex = s.fileIndex + 1
			return currentIndex, nil
		} else if errorReading != nil {
			return currentIndex, errorReading
		}
		candidateNode := &CandidateNode{
			Document:  currentIndex,
			Filename:  filename,
			Node:      &dataBucket,
			FileIndex: s.fileIndex,
		}
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
