package yqlib

import "container/list"

/**
	Loads all yaml documents of all files given into memory, then runs the given expression once.
**/
type Evaluator interface {
	EvaluateFiles(expression string, filenames []string, printer Printer) error
}

type allAtOnceEvaluator struct {
	treeNavigator DataTreeNavigator
	treeCreator   PathTreeCreator
}

func NewAllAtOnceEvaluator() Evaluator {
	return &allAtOnceEvaluator{treeNavigator: NewDataTreeNavigator(), treeCreator: NewPathTreeCreator()}
}

func (e *allAtOnceEvaluator) EvaluateFiles(expression string, filenames []string, printer Printer) error {
	fileIndex := 0
	node, err := treeCreator.ParsePath(expression)
	if err != nil {
		return err
	}
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
	matches, err := treeNavigator.GetMatchingNodes(allDocuments, node)
	if err != nil {
		return err
	}
	return printer.PrintResults(matches)
}
