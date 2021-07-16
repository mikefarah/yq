package yqlib

import (
	"container/list"

	yaml "gopkg.in/yaml.v3"
)

// A yaml expression evaluator that runs the expression once against all files/nodes in memory.
type Evaluator interface {
	EvaluateFiles(expression string, filenames []string, printer Printer) error

	// EvaluateNodes takes an expression and one or more yaml nodes, returning a list of matching candidate nodes
	EvaluateNodes(expression string, nodes ...*yaml.Node) (*list.List, error)

	// EvaluateCandidateNodes takes an expression and list of candidate nodes, returning a list of matching candidate nodes
	EvaluateCandidateNodes(expression string, inputCandidateNodes *list.List) (*list.List, error)
}

type allAtOnceEvaluator struct {
	treeNavigator DataTreeNavigator
	treeCreator   ExpressionParser
}

func NewAllAtOnceEvaluator() Evaluator {
	return &allAtOnceEvaluator{treeNavigator: NewDataTreeNavigator(), treeCreator: NewExpressionParser()}
}

func (e *allAtOnceEvaluator) EvaluateNodes(expression string, nodes ...*yaml.Node) (*list.List, error) {
	inputCandidates := list.New()
	for _, node := range nodes {
		inputCandidates.PushBack(&CandidateNode{Node: node})
	}
	return e.EvaluateCandidateNodes(expression, inputCandidates)
}

func (e *allAtOnceEvaluator) EvaluateCandidateNodes(expression string, inputCandidates *list.List) (*list.List, error) {
	node, err := e.treeCreator.ParseExpression(expression)
	if err != nil {
		return nil, err
	}
	context, err := e.treeNavigator.GetMatchingNodes(Context{MatchingNodes: inputCandidates}, node)
	if err != nil {
		return nil, err
	}
	return context.MatchingNodes, nil
}

func (e *allAtOnceEvaluator) EvaluateFiles(expression string, filenames []string, printer Printer) error {
	fileIndex := 0
	firstFileLeadingSeperator := false

	var allDocuments *list.List = list.New()
	for _, filename := range filenames {
		reader, leadingSeperator, err := readStream(filename)
		if err != nil {
			return err
		}

		if fileIndex == 0 && leadingSeperator {
			firstFileLeadingSeperator = leadingSeperator
		}

		fileDocuments, err := readDocuments(reader, filename, fileIndex)
		if err != nil {
			return err
		}
		allDocuments.PushBackList(fileDocuments)
		fileIndex = fileIndex + 1
	}

	if allDocuments.Len() == 0 {
		candidateNode := &CandidateNode{
			Document:  0,
			Filename:  "",
			Node:      &yaml.Node{Tag: "!!null", Kind: yaml.ScalarNode},
			FileIndex: 0,
		}
		allDocuments.PushBack(candidateNode)
	}

	matches, err := e.EvaluateCandidateNodes(expression, allDocuments)
	if err != nil {
		return err
	}
	printer.SetPrintLeadingSeperator(firstFileLeadingSeperator)
	return printer.PrintResults(matches)
}
