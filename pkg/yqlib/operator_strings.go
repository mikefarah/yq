package yqlib

import (
	"container/list"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func joinStringOperator(d *dataTreeNavigator, matchMap *list.List, expressionNode *ExpressionNode) (*list.List, error) {
	log.Debugf("-- joinStringOperator")
	joinStr := ""

	rhs, err := d.GetMatchingNodes(matchMap, expressionNode.Rhs)
	if err != nil {
		return nil, err
	}
	if rhs.Front() != nil {
		joinStr = rhs.Front().Value.(*CandidateNode).Node.Value
	}

	var results = list.New()

	for el := matchMap.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := unwrapDoc(candidate.Node)
		if node.Kind != yaml.SequenceNode {
			return nil, fmt.Errorf("Cannot join with %v, can only join arrays of scalars", node.Tag)
		}
		targetNode := join(node.Content, joinStr)
		result := candidate.CreateChild(nil, targetNode)
		results.PushBack(result)
	}

	return results, nil
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
