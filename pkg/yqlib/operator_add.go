package yqlib

import (
	"fmt"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

func createAddOp(lhs *ExpressionNode, rhs *ExpressionNode) *ExpressionNode {
	return &ExpressionNode{Operation: &Operation{OperationType: addOpType},
		Lhs: lhs,
		Rhs: rhs}
}

func addAssignOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	return compoundAssignFunction(d, context, expressionNode, createAddOp)
}

func toNodes(candidate *CandidateNode) []*yaml.Node {
	if candidate.Node.Tag == "!!null" {
		return []*yaml.Node{}
	}

	switch candidate.Node.Kind {
	case yaml.SequenceNode:
		return candidate.Node.Content
	default:
		return []*yaml.Node{candidate.Node}
	}

}

func addOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("Add operator")

	return crossFunction(d, context.ReadOnlyClone(), expressionNode, add, false)
}

func add(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	lhs.Node = unwrapDoc(lhs.Node)
	rhs.Node = unwrapDoc(rhs.Node)

	lhsNode := lhs.Node

	if lhsNode.Tag == "!!null" {
		return lhs.CreateReplacement(rhs.Node), nil
	}

	target := lhs.CreateReplacement(&yaml.Node{})

	switch lhsNode.Kind {
	case yaml.MappingNode:
		return nil, fmt.Errorf("maps not yet supported for addition")
	case yaml.SequenceNode:
		target.Node.Kind = yaml.SequenceNode
		target.Node.Style = lhsNode.Style
		target.Node.Tag = lhsNode.Tag
		target.Node.Content = append(lhsNode.Content, toNodes(rhs)...)
	case yaml.ScalarNode:
		if rhs.Node.Kind != yaml.ScalarNode {
			return nil, fmt.Errorf("%v (%v) cannot be added to a 2%v", rhs.Node.Tag, rhs.Path, lhsNode.Tag)
		}
		target.Node.Kind = yaml.ScalarNode
		target.Node.Style = lhsNode.Style
		return addScalars(target, lhsNode, rhs.Node)
	}

	return target, nil
}

func guessTagFromCustomType(node *yaml.Node) string {
	decoder := NewYamlDecoder()
	decoder.Init(strings.NewReader(node.Value))
	var dataBucket yaml.Node
	errorReading := decoder.Decode(&dataBucket)
	if errorReading != nil {
		log.Warning("could not guess underlying tag type %w", errorReading)
		return node.Tag
	}
	guessedTag := unwrapDoc(&dataBucket).Tag
	log.Info("im guessing the tag %v is a %v", node.Tag, guessedTag)
	return guessedTag
}

func addScalars(target *CandidateNode, lhs *yaml.Node, rhs *yaml.Node) (*CandidateNode, error) {
	lhsTag := lhs.Tag
	rhsTag := rhs.Tag
	lhsIsCustom := false
	if !strings.HasPrefix(lhsTag, "!!") {
		// custom tag - we have to have a guess
		lhsTag = guessTagFromCustomType(lhs)
		lhsIsCustom = true
	}

	if !strings.HasPrefix(rhsTag, "!!") {
		// custom tag - we have to have a guess
		rhsTag = guessTagFromCustomType(rhs)
	}

	if lhsTag == "!!str" {
		target.Node.Tag = lhs.Tag
		target.Node.Value = lhs.Value + rhs.Value
	} else if lhsTag == "!!int" && rhsTag == "!!int" {
		format, lhsNum, err := parseInt(lhs.Value)
		if err != nil {
			return nil, err
		}
		_, rhsNum, err := parseInt(rhs.Value)
		if err != nil {
			return nil, err
		}
		sum := lhsNum + rhsNum
		target.Node.Tag = lhs.Tag
		target.Node.Value = fmt.Sprintf(format, sum)
	} else if (lhsTag == "!!int" || lhsTag == "!!float") && (rhsTag == "!!int" || rhsTag == "!!float") {
		lhsNum, err := strconv.ParseFloat(lhs.Value, 64)
		if err != nil {
			return nil, err
		}
		rhsNum, err := strconv.ParseFloat(rhs.Value, 64)
		if err != nil {
			return nil, err
		}
		sum := lhsNum + rhsNum
		if lhsIsCustom {
			target.Node.Tag = lhs.Tag
		} else {
			target.Node.Tag = "!!float"
		}
		target.Node.Value = fmt.Sprintf("%v", sum)
	} else {
		return nil, fmt.Errorf("%v cannot be added to %v", lhsTag, rhsTag)
	}

	return target, nil
}
