package yqlib

import (
	"container/list"
	"errors"
	"fmt"
	"io"
	"os"

	yaml "gopkg.in/yaml.v3"
)

// A yaml expression evaluator that runs the expression multiple times for each given yaml document.
// Uses less memory than loading all documents and running the expression once, but this cannot process
// cross document expressions.
type StreamEvaluator interface {
	Evaluate(filename string, reader io.Reader, node *ExpressionNode, printer Printer, leadingContent string, decoder Decoder) (uint, error)
	EvaluateFiles(expression string, filenames []string, printer Printer, leadingContentPreProcessing bool, decoder Decoder) error
	EvaluateNew(expression string, printer Printer, leadingContent string) error
}

type streamEvaluator struct {
	treeNavigator DataTreeNavigator
	fileIndex     int
}

func NewStreamEvaluator() StreamEvaluator {
	return &streamEvaluator{treeNavigator: NewDataTreeNavigator()}
}

func (s *streamEvaluator) EvaluateNew(expression string, printer Printer, leadingContent string) error {
	node, err := ExpressionParser.ParseExpression(expression)
	if err != nil {
		return err
	}
	candidateNode := &CandidateNode{
		Document:       0,
		Filename:       "",
		Node:           &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{{Tag: "!!null", Kind: yaml.ScalarNode}}},
		FileIndex:      0,
		LeadingContent: leadingContent,
	}
	inputList := list.New()
	inputList.PushBack(candidateNode)

	result, errorParsing := s.treeNavigator.GetMatchingNodes(Context{MatchingNodes: inputList}, node)
	if errorParsing != nil {
		return errorParsing
	}
	return printer.PrintResults(result.MatchingNodes)
}

func (s *streamEvaluator) EvaluateFiles(expression string, filenames []string, printer Printer, leadingContentPreProcessing bool, decoder Decoder) error {
	var totalProcessDocs uint
	node, err := ExpressionParser.ParseExpression(expression)
	if err != nil {
		return err
	}

	var firstFileLeadingContent string

	for index, filename := range filenames {
		reader, leadingContent, err := readStream(filename, leadingContentPreProcessing)

		if index == 0 {
			firstFileLeadingContent = leadingContent
		}

		if err != nil {
			return err
		}
		processedDocs, err := s.Evaluate(filename, reader, node, printer, leadingContent, decoder)
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

func (s *streamEvaluator) Evaluate(filename string, reader io.Reader, node *ExpressionNode, printer Printer, leadingContent string, decoder Decoder) (uint, error) {

	var currentIndex uint
	decoder.Init(reader)
	for {
		var dataBucket yaml.Node
		errorReading := decoder.Decode(&dataBucket)

		if errors.Is(errorReading, io.EOF) {
			s.fileIndex = s.fileIndex + 1
			return currentIndex, nil
		} else if errorReading != nil {
			return currentIndex, fmt.Errorf("bad file '%v': %w", filename, errorReading)
		}

		candidateNode := &CandidateNode{
			Document:  currentIndex,
			Filename:  filename,
			Node:      &dataBucket,
			FileIndex: s.fileIndex,
		}
		if currentIndex == 0 {
			candidateNode.LeadingContent = leadingContent
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
