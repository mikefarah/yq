package yqlib

import (
	"container/list"
	"fmt"

	"gopkg.in/yaml.v3"
)

type operatorHandler func(d *dataTreeNavigator, matchingNodes *list.List, expressionNode *ExpressionNode) (*list.List, error)

func unwrapDoc(node *yaml.Node) *yaml.Node {
	if node.Kind == yaml.DocumentNode {
		return node.Content[0]
	}
	return node
}

func emptyOperator(d *dataTreeNavigator, matchingNodes *list.List, expressionNode *ExpressionNode) (*list.List, error) {
	return list.New(), nil
}

func createBooleanCandidate(owner *CandidateNode, value bool) *CandidateNode {
	valString := "true"
	if !value {
		valString = "false"
	}
	node := &yaml.Node{Kind: yaml.ScalarNode, Value: valString, Tag: "!!bool"}
	return owner.CreateChild(nil, node)
}

func nodeToMap(candidate *CandidateNode) *list.List {
	elMap := list.New()
	elMap.PushBack(candidate)
	return elMap
}

func createTraversalTree(path []interface{}, traversePrefs traversePreferences) *ExpressionNode {
	if len(path) == 0 {
		return &ExpressionNode{Operation: &Operation{OperationType: selfReferenceOpType}}
	} else if len(path) == 1 {
		return &ExpressionNode{Operation: &Operation{OperationType: traversePathOpType, Preferences: traversePrefs, Value: path[0], StringValue: fmt.Sprintf("%v", path[0])}}
	}

	return &ExpressionNode{
		Operation: &Operation{OperationType: shortPipeOpType},
		Lhs:       createTraversalTree(path[0:1], traversePrefs),
		Rhs:       createTraversalTree(path[1:], traversePrefs),
	}
}
