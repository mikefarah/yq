package yqlib

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"io"

	yaml "gopkg.in/yaml.v3"
)

type StringEvaluator interface {
	Evaluate(expression string, input string, encoder Encoder, leadingContentPreProcessing bool, decoder Decoder) (string, error)
}

type stringEvaluator struct {
	treeNavigator DataTreeNavigator
	fileIndex     int
}

func NewStringEvaluator() StringEvaluator {
	return &stringEvaluator{
		treeNavigator: NewDataTreeNavigator(),
	}
}

func (s *stringEvaluator) Evaluate(expression string, input string, encoder Encoder, leadingContentPreProcessing bool, decoder Decoder) (string, error) {

	// Use bytes.Buffer for output of string
	out := new(bytes.Buffer)
	printer := NewPrinter(encoder, NewSinglePrinterWriter(out))

	InitExpressionParser()
	node, err := ExpressionParser.ParseExpression(expression)
	if err != nil {
		return "", err
	}

	reader, leadingContent, err := readString(input, leadingContentPreProcessing)
	if err != nil {
		return "", err
	}

	var currentIndex uint
	decoder.Init(reader)
	for {
		var dataBucket yaml.Node
		errorReading := decoder.Decode(&dataBucket)

		if errors.Is(errorReading, io.EOF) {
			s.fileIndex = s.fileIndex + 1
			return out.String(), nil
		} else if errorReading != nil {
			return "", fmt.Errorf("bad input '%v': %w", input, errorReading)
		}

		candidateNode := &CandidateNode{
			Document:  currentIndex,
			Node:      &dataBucket,
			FileIndex: s.fileIndex,
		}
		// move document comments into candidate node
		// otherwise unwrap drops them.
		candidateNode.TrailingContent = dataBucket.FootComment
		dataBucket.FootComment = ""

		if currentIndex == 0 {
			candidateNode.LeadingContent = leadingContent
		}
		inputList := list.New()
		inputList.PushBack(candidateNode)

		result, errorParsing := s.treeNavigator.GetMatchingNodes(Context{MatchingNodes: inputList}, node)
		if errorParsing != nil {
			return "", errorParsing
		}
		err = printer.PrintResults(result.MatchingNodes)

		if err != nil {
			return "", err
		}
		currentIndex = currentIndex + 1
	}
}
