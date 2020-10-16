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

func nodeToMap(candidate *CandidateNode) *orderedmap.OrderedMap {
	elMap := orderedmap.NewOrderedMap()
	elMap.Set(candidate.GetKey(), candidate)
	return elMap
}

func AssignOperator(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	lhs, err := d.getMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}
	for el := lhs.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		rhs, err := d.getMatchingNodes(nodeToMap(candidate), pathNode.Rhs)

		if err != nil {
			return nil, err
		}

		// grab the first value
		first := rhs.Front()

		if first != nil {
			candidate.UpdateFrom(first.Value.(*CandidateNode))
		}
	}
	return lhs, nil
}

func UnionOperator(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	lhs, err := d.getMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}
	rhs, err := d.getMatchingNodes(matchingNodes, pathNode.Rhs)
	if err != nil {
		return nil, err
	}
	for el := rhs.Front(); el != nil; el = el.Next() {
		node := el.Value.(*CandidateNode)
		lhs.Set(node.GetKey(), node)
	}
	return lhs, nil
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

func splatNode(d *dataTreeNavigator, candidate *CandidateNode) (*orderedmap.OrderedMap, error) {
	//need to splat matching nodes, then search through them
	splatter := &PathTreeNode{PathElement: &PathElement{
		PathElementType: PathKey,
		Value:           "*",
		StringValue:     "*",
	}}
	return d.getMatchingNodes(nodeToMap(candidate), splatter)
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

func CollectOperator(d *dataTreeNavigator, matchMap *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	log.Debugf("-- collectOperation")

	var results = orderedmap.NewOrderedMap()

	node := &yaml.Node{Kind: yaml.SequenceNode}

	var document uint = 0
	var path []interface{}

	for el := matchMap.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		if path == nil && candidate.Path != nil {
			path = candidate.Path
			document = candidate.Document
		}
		node.Content = append(node.Content, candidate.Node)
	}

	collectC := &CandidateNode{Node: node, Document: document, Path: path}
	results.Set(collectC.GetKey(), collectC)

	return results, nil

}
