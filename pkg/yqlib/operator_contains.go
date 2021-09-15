package yqlib

import (
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

func containsOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	return crossFunction(d, context.ReadOnlyClone(), expressionNode, containsWithNodes, false)
}

func containsArrayElement(array *yaml.Node, item *yaml.Node) (bool, error) {
	for index := 0; index < len(array.Content); index = index + 1 {
		containedInArray, err := contains(array.Content[index], item)
		if err != nil {
			return false, err
		}
		if containedInArray {
			return true, nil
		}
	}
	return false, nil
}

func containsArray(lhs *yaml.Node, rhs *yaml.Node) (bool, error) {
	if rhs.Kind != yaml.SequenceNode {
		return containsArrayElement(lhs, rhs)
	}
	for index := 0; index < len(rhs.Content); index = index + 1 {
		itemInArray, err := containsArrayElement(lhs, rhs.Content[index])
		if err != nil {
			return false, err
		}
		if !itemInArray {
			return false, nil
		}
	}
	return true, nil
}

func containsObject(lhs *yaml.Node, rhs *yaml.Node) (bool, error) {
	if rhs.Kind != yaml.MappingNode {
		return false, nil
	}
	for index := 0; index < len(rhs.Content); index = index + 2 {
		rhsKey := rhs.Content[index]
		rhsValue := rhs.Content[index+1]
		log.Debugf("Looking for %v in the lhs", rhsKey.Value)
		lhsKeyIndex := findInArray(lhs, rhsKey)
		log.Debugf("index is %v", lhsKeyIndex)
		if lhsKeyIndex < 0 || lhsKeyIndex%2 != 0 {
			return false, nil
		}
		lhsValue := lhs.Content[lhsKeyIndex+1]
		log.Debugf("lhsValue is %v", lhsValue.Value)

		itemInArray, err := contains(lhsValue, rhsValue)
		log.Debugf("rhsValue is %v", rhsValue.Value)
		if err != nil {
			return false, err
		}
		if !itemInArray {
			return false, nil
		}
	}
	return true, nil
}

func containsScalars(lhs *yaml.Node, rhs *yaml.Node) (bool, error) {
	if lhs.Tag == "!!str" {
		return strings.Contains(lhs.Value, rhs.Value), nil
	}
	return lhs.Value == rhs.Value, nil
}

func contains(lhs *yaml.Node, rhs *yaml.Node) (bool, error) {
	switch lhs.Kind {
	case yaml.MappingNode:
		return containsObject(lhs, rhs)
	case yaml.SequenceNode:
		return containsArray(lhs, rhs)
	case yaml.ScalarNode:
		if rhs.Kind != yaml.ScalarNode || lhs.Tag != rhs.Tag {
			return false, nil
		}
		if lhs.Tag == "!!null" {
			return rhs.Tag == "!!null", nil
		}
		return containsScalars(lhs, rhs)
	}

	return false, fmt.Errorf("%v not yet supported for contains", lhs.Tag)
}

func containsWithNodes(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	lhs.Node = unwrapDoc(lhs.Node)
	rhs.Node = unwrapDoc(rhs.Node)

	if lhs.Node.Kind != rhs.Node.Kind {
		return nil, fmt.Errorf("%v cannot check contained in %v", rhs.Node.Tag, lhs.Node.Tag)
	}

	result, err := contains(lhs.Node, rhs.Node)
	if err != nil {
		return nil, err
	}

	return createBooleanCandidate(lhs, result), nil
}
