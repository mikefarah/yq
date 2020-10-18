package treeops

import "github.com/elliotchance/orderedmap"

func MultiplyOperator(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	lhs, err := d.getMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}
	for el := lhs.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		// TODO handle scalar mulitplication
		switch candidate.Node.Kind {
			case 
		}
		
	}
	return matchingNodes, nil
}
