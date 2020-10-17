package treeops

import (
	"fmt"

	"github.com/elliotchance/orderedmap"
	"gopkg.in/yaml.v3"
)

type OperatorHandler func(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error)

func PipeOperator(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	lhs, err := d.getMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}
	return d.getMatchingNodes(lhs, pathNode.Rhs)
}

func createBooleanCandidate(owner *CandidateNode, value bool) *CandidateNode {
	valString := "true"
	if !value {
		valString = "false"
	}
	node := &yaml.Node{Kind: yaml.ScalarNode, Value: valString, Tag: "!!bool"}
	return &CandidateNode{Node: node, Document: owner.Document, Path: owner.Path}
}

func nodeToMap(candidate *CandidateNode) *orderedmap.OrderedMap {
	elMap := orderedmap.NewOrderedMap()
	elMap.Set(candidate.GetKey(), candidate)
	return elMap
}

func IntersectionOperator(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	lhs, err := d.getMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}
	rhs, err := d.getMatchingNodes(matchingNodes, pathNode.Rhs)
	if err != nil {
		return nil, err
	}
	var matchingNodeMap = orderedmap.NewOrderedMap()
	for el := lhs.Front(); el != nil; el = el.Next() {
		_, exists := rhs.Get(el.Key)
		if exists {
			matchingNodeMap.Set(el.Key, el.Value)
		}
	}
	return matchingNodeMap, nil
}

func LengthOperator(d *dataTreeNavigator, matchMap *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	log.Debugf("-- lengthOperation")
	var results = orderedmap.NewOrderedMap()

	for el := matchMap.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		var length int
		switch candidate.Node.Kind {
		case yaml.ScalarNode:
			length = len(candidate.Node.Value)
		case yaml.MappingNode:
			length = len(candidate.Node.Content) / 2
		case yaml.SequenceNode:
			length = len(candidate.Node.Content)
		default:
			length = 0
		}

		node := &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%v", length), Tag: "!!int"}
		lengthCand := &CandidateNode{Node: node, Document: candidate.Document, Path: candidate.Path}
		results.Set(candidate.GetKey(), lengthCand)
	}

	return results, nil
}
