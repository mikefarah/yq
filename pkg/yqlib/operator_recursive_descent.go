package yqlib

import (
	"container/list"

	yaml "gopkg.in/yaml.v3"
)

func RecursiveDescentOperator(d *dataTreeNavigator, matchMap *list.List, pathNode *PathTreeNode) (*list.List, error) {
	var results = list.New()

	err := recursiveDecent(d, results, matchMap, true)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func recursiveDecent(d *dataTreeNavigator, results *list.List, matchMap *list.List, recurseArray bool) error {
	for el := matchMap.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		candidate.Node = UnwrapDoc(candidate.Node)

		log.Debugf("Recursive Decent, added %v", NodeToString(candidate))
		results.PushBack(candidate)

		if candidate.Node.Kind != yaml.AliasNode && len(candidate.Node.Content) > 0 &&
			(recurseArray || candidate.Node.Kind != yaml.SequenceNode) {

			children, err := Splat(d, nodeToMap(candidate))

			if err != nil {
				return err
			}
			err = recursiveDecent(d, results, children, recurseArray)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
