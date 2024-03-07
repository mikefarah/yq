package yqlib

import (
	"bufio"
	"bytes"
	"container/list"
	"strings"
)

type StringEvaluator interface {
	Evaluate(expression string, input string, encoder Encoder, decoder Decoder) (string, error)
	EvaluateAll(expression string, input string, encoder Encoder, decoder Decoder) (string, error)
}

type stringEvaluator struct {
	treeNavigator DataTreeNavigator
}

func NewStringEvaluator() StringEvaluator {
	return &stringEvaluator{
		treeNavigator: NewDataTreeNavigator(),
	}
}

func (s *stringEvaluator) EvaluateAll(expression string, input string, encoder Encoder, decoder Decoder) (string, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var documents *list.List
	var results *list.List
	var err error

	if documents, err = ReadDocuments(reader, decoder); err != nil {
		return "", err
	}

	evaluator := NewAllAtOnceEvaluator()
	if results, err = evaluator.EvaluateCandidateNodes(expression, documents); err != nil {
		return "", err
	}

	out := new(bytes.Buffer)
	printer := NewPrinter(encoder, NewSinglePrinterWriter(out))
	if err := printer.PrintResults(results); err != nil {
		return "", err
	}
	return out.String(), nil
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
	evaluator := NewStreamEvaluator()
	if _, err := evaluator.Evaluate("", reader, node, printer, decoder); err != nil {
		return "", err
	}
	return out.String(), nil
}
