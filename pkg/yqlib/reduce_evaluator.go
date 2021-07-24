package yqlib

import (
	"container/list"
	"fmt"
	"io"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type ReduceEvaluator interface {
	EvaluateFiles(reduceExpression string, filenames []string, printer Printer, leadingContentPreProcessing bool) error
}

type reduceEvaluator struct {
	treeNavigator DataTreeNavigator
	treeCreator   ExpressionParser
	reduceLhs     *ExpressionNode
	fileIndex     int
}

func NewReduceEvaluator() ReduceEvaluator {
	treeCreator := NewExpressionParser()
	reduceLhs, err := treeCreator.ParseExpression(". as $doc")
	if err != nil {
		panic(err)
	}
	return &reduceEvaluator{treeNavigator: NewDataTreeNavigator(), treeCreator: treeCreator, reduceLhs: reduceLhs}
}

func (r *reduceEvaluator) EvaluateFiles(expression string, filenames []string, printer Printer, leadingContentPreProcessing bool) error {

	node, err := r.treeCreator.ParseExpression(expression)
	if err != nil {
		return err
	}
	log.Debug("node %v", node.Operation.toString())

	if node.Operation.OperationType != blockOpType {
		return fmt.Errorf("Invalid reduce expression - expected '<initialValue>; <block that uses $doc>' got '%v'", expression)
	}

	currentValue := node.Lhs
	reduceExp := node.Rhs
	firstLeadingContent := ""

	log.Debug("initialValue %v", currentValue.Operation.toString())

	log.Debug("reduce Exp %v", reduceExp.Operation.toString())

	for index, filename := range filenames {
		reader, leadingContent, err := readStream(filename, leadingContentPreProcessing)

		if index == 0 {
			firstLeadingContent = leadingContent
		}

		if err != nil {
			return err
		}

		currentValue, err = r.ReduceFile(filename, leadingContent, reader, currentValue, reduceExp)
		if err != nil {
			return err
		}

		switch reader := reader.(type) {
		case *os.File:
			safelyCloseFile(reader)
		}
	}

	result := currentValue.Operation.ValueNodes

	if result.Len() > 0 {
		result.Front().Value.(*CandidateNode).Node.HeadComment = firstLeadingContent
	}

	printer.PrintResults(result)

	return nil
}

func (r *reduceEvaluator) createReduceOp(initialValue *ExpressionNode, reduceExp *ExpressionNode) *ExpressionNode {
	reduceBlock := &ExpressionNode{
		Operation: &Operation{OperationType: blockOpType},
		Lhs:       initialValue,
		Rhs:       reduceExp,
	}

	return &ExpressionNode{
		Operation: &Operation{OperationType: reduceOpType},
		Lhs:       r.reduceLhs,
		Rhs:       reduceBlock,
	}
}

func (r *reduceEvaluator) ReduceFile(filename string, leadingContent string, reader io.Reader, initialValue *ExpressionNode, reduceExp *ExpressionNode) (*ExpressionNode, error) {

	var currentIndex uint
	var currentValue = initialValue
	decoder := yaml.NewDecoder(reader)
	for {
		var dataBucket yaml.Node
		errorReading := decoder.Decode(&dataBucket)

		if errorReading == io.EOF {
			r.fileIndex = r.fileIndex + 1
			return currentValue, nil
		} else if errorReading != nil {
			return currentValue, errorReading
		}
		candidateNode := &CandidateNode{
			Document:  currentIndex,
			Filename:  filename,
			Node:      &dataBucket,
			FileIndex: r.fileIndex,
		}
		inputList := list.New()
		inputList.PushBack(candidateNode)

		reduceOp := r.createReduceOp(currentValue, reduceExp)
		// log.Debug("reduce - currentValueBefore: %v", NodesToString(currentValue.Operation.ValueNodes))

		result, errorParsing := r.treeNavigator.GetMatchingNodes(Context{MatchingNodes: inputList}, reduceOp)
		if errorParsing != nil {
			return currentValue, errorParsing
		}

		currentValue = &ExpressionNode{
			Operation: &Operation{
				OperationType: valueOpType,
				ValueNodes:    result.MatchingNodes,
			},
		}

		log.Debug("reduce - currentValueAfter: %v", NodesToString(currentValue.Operation.ValueNodes))

		currentIndex = currentIndex + 1
	}
}
