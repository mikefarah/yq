package yqlib

import (
	"fmt"
	"strings"
)

func containsOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	return crossFunction(d, context.ReadOnlyClone(), expressionNode, containsWithNodes, false)
}

func containsArrayElement(array *CandidateNode, item *CandidateNode) (bool, error) {
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

func containsArray(lhs *CandidateNode, rhs *CandidateNode) (bool, error) {
	if rhs.Kind != SequenceNode {
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

func containsObject(lhs *CandidateNode, rhs *CandidateNode) (bool, error) {
	if rhs.Kind != MappingNode {
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

func containsScalars(lhs *CandidateNode, rhs *CandidateNode) (bool, error) {
	if lhs.Tag == "!!str" {
		return strings.Contains(lhs.Value, rhs.Value), nil
	}
	return lhs.Value == rhs.Value, nil
}

func contains(lhs *CandidateNode, rhs *CandidateNode) (bool, error) {
	switch lhs.Kind {
	case MappingNode:
		return containsObject(lhs, rhs)
	case SequenceNode:
		return containsArray(lhs, rhs)
	case ScalarNode:
		if rhs.Kind != ScalarNode || lhs.Tag != rhs.Tag {
			return false, nil
		}
		if lhs.Tag == "!!null" {
			return rhs.Tag == "!!null", nil
		}
		return containsScalars(lhs, rhs)
	}

	return false, fmt.Errorf("%v not yet supported for contains", lhs.Tag)
}

func containsWithNodes(_ *dataTreeNavigator, _ Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	if lhs.Kind != rhs.Kind {
		return nil, fmt.Errorf("%v cannot check contained in %v", rhs.Tag, lhs.Tag)
	}

	result, err := contains(lhs, rhs)
	if err != nil {
		return nil, err
	}

	return createBooleanCandidate(lhs, result), nil
}
