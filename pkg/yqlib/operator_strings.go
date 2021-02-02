package yqlib

import (
	"container/list"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func joinStringOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- joinStringOperator")
	joinStr := ""

	rhs, err := d.GetMatchingNodes(context, expressionNode.Rhs)
	if err != nil {
		return Context{}, err
	}
	if rhs.MatchingNodes.Front() != nil {
		joinStr = rhs.MatchingNodes.Front().Value.(*CandidateNode).Node.Value
	}

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := unwrapDoc(candidate.Node)
		if node.Kind != yaml.SequenceNode {
			return Context{}, fmt.Errorf("Cannot join with %v, can only join arrays of scalars", node.Tag)
		}
		targetNode := join(node.Content, joinStr)
		result := candidate.CreateChild(nil, targetNode)
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}

func join(content []*yaml.Node, joinStr string) *yaml.Node {
	var stringsToJoin []string
	for _, node := range content {
		str := node.Value
		if node.Tag == "!!null" {
			str = ""
		}
		stringsToJoin = append(stringsToJoin, str)
	}

	return &yaml.Node{Kind: yaml.ScalarNode, Value: strings.Join(stringsToJoin, joinStr), Tag: "!!str"}
}

func splitStringOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- splitStringOperator")
	splitStr := ""

	rhs, err := d.GetMatchingNodes(context, expressionNode.Rhs)
	if err != nil {
		return Context{}, err
	}
	if rhs.MatchingNodes.Front() != nil {
		splitStr = rhs.MatchingNodes.Front().Value.(*CandidateNode).Node.Value
	}

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := unwrapDoc(candidate.Node)
		if node.Tag == "!!null" {
			continue
		}
		if node.Tag != "!!str" {
			return Context{}, fmt.Errorf("Cannot split %v, can only split strings", node.Tag)
		}
		targetNode := split(node.Value, splitStr)
		result := candidate.CreateChild(nil, targetNode)
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}

func split(value string, spltStr string) *yaml.Node {
	var contents []*yaml.Node

	if value != "" {
		var newStrings = strings.Split(value, spltStr)
		contents = make([]*yaml.Node, len(newStrings))

		for index, str := range newStrings {
			contents[index] = &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: str}
		}
	}

	return &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq", Content: contents}
}
