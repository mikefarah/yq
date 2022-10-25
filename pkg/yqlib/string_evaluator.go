package yqlib

import (
	"bufio"
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"io"
	"strings"
)

type StringEvaluator interface {
	Evaluate(expression string, input string, encoder Encoder, decoder Decoder) (string, error)
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

func (s *stringEvaluator) Evaluate(expression string, input string, encoder Encoder, decoder Decoder) (string, error) {

	// Use bytes.Buffer for output of string
	out := new(bytes.Buffer)
	printer := NewPrinter(encoder, NewSinglePrinterWriter(out))

	InitExpressionParser()
	node, err := ExpressionParser.ParseExpression(expression)
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(strings.NewReader(input))

	var currentIndex uint
	err = decoder.Init(reader)
	if err != nil {
		return "", err
	}
	for {
		candidateNode, errorReading := decoder.Decode()

		if errors.Is(errorReading, io.EOF) {
			s.fileIndex = s.fileIndex + 1
			return out.String(), nil
		} else if errorReading != nil {
			return "", fmt.Errorf("bad input '%v': %w", input, errorReading)
		}
		candidateNode.Document = currentIndex
		candidateNode.FileIndex = s.fileIndex

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
