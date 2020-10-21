package treeops

import (
	"container/list"
)

func RecursiveDescentOperator(d *dataTreeNavigator, matchMap *list.List, pathNode *PathTreeNode) (*list.List, error) {
	var results = list.New()

	err := recursiveDecent(d, results, matchMap)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func recursiveDecent(d *dataTreeNavigator, results *list.List, matchMap *list.List) error {
	splatOperation := &Operation{OperationType: TraversePath, Value: "[]"}
	splatTreeNode := &PathTreeNode{Operation: splatOperation}

	for el := matchMap.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		results.PushBack(candidate)

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
