package treeops

import (
	"github.com/elliotchance/orderedmap"
)

func RecursiveDescentOperator(d *dataTreeNavigator, matchMap *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	var results = orderedmap.NewOrderedMap()

	err := recursiveDecent(d, results, matchMap)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func recursiveDecent(d *dataTreeNavigator, results *orderedmap.OrderedMap, matchMap *orderedmap.OrderedMap) error {
	splatPathElement := &PathElement{OperationType: TraversePath, Value: "[]"}
	splatTreeNode := &PathTreeNode{PathElement: splatPathElement}

	for el := matchMap.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		results.Set(candidate.GetKey(), candidate)

		children, err := TraversePathOperator(d, nodeToMap(candidate), splatTreeNode)

		if err != nil {
			return err
		}
		err = recursiveDecent(d, results, children)
		if err != nil {
			return err
		}
	}
	return nil
}
