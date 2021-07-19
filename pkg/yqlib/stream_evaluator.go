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
	Evaluate(filename string, reader io.Reader, node *ExpressionNode, printer Printer, leadingContent string) (uint, error)
	EvaluateFiles(expression string, filenames []string, printer Printer) error
	EvaluateNew(expression string, printer Printer, leadingContent string) error
}

type streamEvaluator struct {
	treeNavigator DataTreeNavigator
	treeCreator   ExpressionParser
	fileIndex     int
}

func NewStreamEvaluator() StreamEvaluator {
	return &streamEvaluator{treeNavigator: NewDataTreeNavigator(), treeCreator: NewExpressionParser()}
}

func (s *streamEvaluator) EvaluateNew(expression string, printer Printer, leadingContent string) error {
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
	return printer.PrintResults(result.MatchingNodes, leadingContent)
}

func (s *streamEvaluator) EvaluateFiles(expression string, filenames []string, printer Printer) error {
	var totalProcessDocs uint = 0
	node, err := s.treeCreator.ParseExpression(expression)
	if err != nil {
		return err
	}

	var firstFileLeadingContent string

	for index, filename := range filenames {
		reader, leadingContent, err := readStream(filename)

		if index == 0 {
			firstFileLeadingContent = leadingContent
		}

		if err != nil {
			return err
		}
		processedDocs, err := s.Evaluate(filename, reader, node, printer, leadingContent)
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
		return s.EvaluateNew(expression, printer, firstFileLeadingContent)
	}

	return nil
}

func (s *streamEvaluator) Evaluate(filename string, reader io.Reader, node *ExpressionNode, printer Printer, leadingContent string) (uint, error) {

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
		var err error
		if currentIndex == 0 {
			err = printer.PrintResults(result.MatchingNodes, leadingContent)
		} else {
			err = printer.PrintResults(result.MatchingNodes, "")
		}

		if err != nil {
			return currentIndex, err
		}
		currentIndex = currentIndex + 1
	}
}
