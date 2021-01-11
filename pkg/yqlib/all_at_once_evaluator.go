package yqlib

import "container/list"

// A yaml expression evaluator that runs the expression once against all files/nodes in memory.
type Evaluator interface {
	EvaluateFiles(expression string, filenames []string, printer Printer) error

	// Runs the expression once against the list of candidate nodes, returns the
	// resulting nodes.
	EvaluateNodes(expression string, inputCandidateNodes *list.List) (*list.List, error)
}

type allAtOnceEvaluator struct {
	treeNavigator DataTreeNavigator
	treeCreator   PathTreeCreator
}

func NewAllAtOnceEvaluator() Evaluator {
	return &allAtOnceEvaluator{treeNavigator: NewDataTreeNavigator(), treeCreator: NewPathTreeCreator()}
}

func (e *allAtOnceEvaluator) EvaluateNodes(expression string, inputCandidates *list.List) (*list.List, error) {
	node, err := treeCreator.ParsePath(expression)
	if err != nil {
		return nil, err
	}
	return treeNavigator.GetMatchingNodes(inputCandidates, node)
}

func (e *allAtOnceEvaluator) EvaluateFiles(expression string, filenames []string, printer Printer) error {
	fileIndex := 0

	var allDocuments *list.List = list.New()
	for _, filename := range filenames {
		reader, err := readStream(filename)
		if err != nil {
			return err
		}
		fileDocuments, err := readDocuments(reader, filename, fileIndex)
		if err != nil {
			return err
		}
		allDocuments.PushBackList(fileDocuments)
		fileIndex = fileIndex + 1
	}
	matches, err := e.EvaluateNodes(expression, allDocuments)
	if err != nil {
		return err
	}
	return printer.PrintResults(matches)
}
