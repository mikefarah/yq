package treeops

import (
	"fmt"

	"github.com/elliotchance/orderedmap"
	"gopkg.in/yaml.v3"
)

type OperatorHandler func(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error)

func TraverseOperator(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	lhs, err := d.getMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}
	return d.getMatchingNodes(lhs, pathNode.Rhs)
}

func AssignOperator(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	lhs, err := d.getMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}
	for el := lhs.Front(); el != nil; el = el.Next() {
		node := el.Value.(*CandidateNode)
		log.Debugf("Assiging %v to %v", node.GetKey(), pathNode.Rhs.PathElement.StringValue)
		node.Node.Value = pathNode.Rhs.PathElement.StringValue
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
	elMap := orderedmap.NewOrderedMap()
	elMap.Set(candidate.GetKey(), candidate)
	//need to splat matching nodes, then search through them
	splatter := &PathTreeNode{PathElement: &PathElement{
		PathElementType: PathKey,
		Value:           "*",
		StringValue:     "*",
	}}
	return d.getMatchingNodes(elMap, splatter)
}

func EqualsOperator(d *dataTreeNavigator, matchMap *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	log.Debugf("-- equalsOperation")
	var results = orderedmap.NewOrderedMap()

	for el := matchMap.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		valuePattern := pathNode.Rhs.PathElement.StringValue
		log.Debug("checking %v", candidate)

		errInChild := findMatchingChildren(d, results, candidate, pathNode.Lhs, valuePattern)
		if errInChild != nil {
			return nil, errInChild
		}
	}

	return results, nil
}

func CountOperator(d *dataTreeNavigator, matchMap *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	log.Debugf("-- countOperation")
	var results = orderedmap.NewOrderedMap()

	for el := matchMap.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		elMap := orderedmap.NewOrderedMap()
		elMap.Set(el.Key, el.Value)
		childMatches, errChild := d.getMatchingNodes(elMap, pathNode.Rhs)

		if errChild != nil {
			return nil, errChild
		}

		length := childMatches.Len()
		node := &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%v", length), Tag: "!!int"}
		lengthCand := &CandidateNode{Node: node, Document: candidate.Document, Path: candidate.Path}
		results.Set(candidate.GetKey(), lengthCand)

	}

	return results, nil
}

func findMatchingChildren(d *dataTreeNavigator, results *orderedmap.OrderedMap, candidate *CandidateNode, lhs *PathTreeNode, valuePattern string) error {
	var children *orderedmap.OrderedMap
	var err error
	// don't splat scalars.
	if candidate.Node.Kind != yaml.ScalarNode {
		children, err = splatNode(d, candidate)
		log.Debugf("-- splatted matches, ")
		if err != nil {
			return err
		}
	} else {
		children = orderedmap.NewOrderedMap()
		children.Set(candidate.GetKey(), candidate)
	}

	for childEl := children.Front(); childEl != nil; childEl = childEl.Next() {
		childMap := orderedmap.NewOrderedMap()
		childMap.Set(childEl.Key, childEl.Value)
		childMatches, errChild := d.getMatchingNodes(childMap, lhs)
		log.Debug("got the LHS")
		if errChild != nil {
			return errChild
		}

		if containsMatchingValue(childMatches, valuePattern) {
			results.Set(childEl.Key, childEl.Value)
		}
	}
	return nil
}

func containsMatchingValue(matchMap *orderedmap.OrderedMap, valuePattern string) bool {
	log.Debugf("-- findMatchingValues")

	for el := matchMap.Front(); el != nil; el = el.Next() {
		node := el.Value.(*CandidateNode)
		log.Debugf("-- comparing %v to %v", node.Node.Value, valuePattern)
		if Match(node.Node.Value, valuePattern) {
			return true
		}
	}
	log.Debugf("-- done findMatchingValues")

	return false
}
