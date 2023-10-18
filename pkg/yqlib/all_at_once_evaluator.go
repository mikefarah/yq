package yqlib

import (
	"container/list"
)

// A yaml expression evaluator that runs the expression once against all files/nodes in memory.
type Evaluator interface {
	EvaluateFiles(expression string, filenames []string, printer Printer, decoder Decoder) error

	// EvaluateNodes takes an expression and one or more yaml nodes, returning a list of matching candidate nodes
	EvaluateNodes(expression string, nodes ...*CandidateNode) (*list.List, error)

	// EvaluateCandidateNodes takes an expression and list of candidate nodes, returning a list of matching candidate nodes
	EvaluateCandidateNodes(expression string, inputCandidateNodes *list.List) (*list.List, error)
}

type allAtOnceEvaluator struct {
	treeNavigator DataTreeNavigator
}

func NewAllAtOnceEvaluator() Evaluator {
	InitExpressionParser()
	return &allAtOnceEvaluator{treeNavigator: NewDataTreeNavigator()}
}

func (e *allAtOnceEvaluator) EvaluateNodes(expression string, nodes ...*CandidateNode) (*list.List, error) {
	inputCandidates := list.New()
	for _, node := range nodes {
		inputCandidates.PushBack(node)
	}
	return e.EvaluateCandidateNodes(expression, inputCandidates)
}

func (e *allAtOnceEvaluator) EvaluateCandidateNodes(expression string, inputCandidates *list.List) (*list.List, error) {
	node, err := ExpressionParser.ParseExpression(expression)
	if err != nil {
		return nil, err
	}
	context, err := e.treeNavigator.GetMatchingNodes(Context{MatchingNodes: inputCandidates}, node)
	if err != nil {
		return nil, err
	}
	return context.MatchingNodes, nil
}

func (e *allAtOnceEvaluator) EvaluateFiles(expression string, filenames []string, printer Printer, decoder Decoder) error {
	fileIndex := 0

	var allDocuments = list.New()
	for _, filename := range filenames {
		reader, err := readStream(filename)
		if err != nil {
			return err
		}

		fileDocuments, err := readDocuments(reader, filename, fileIndex, decoder)
		if err != nil {
			return err
		}
		allDocuments.PushBackList(fileDocuments)
		fileIndex = fileIndex + 1
	}

	if allDocuments.Len() == 0 {
		candidateNode := createScalarNode(nil, "")
		allDocuments.PushBack(candidateNode)
	}

	matches, err := e.EvaluateCandidateNodes(expression, allDocuments)
	if err != nil {
		return err
	}
	return printer.PrintResults(matches)
}
