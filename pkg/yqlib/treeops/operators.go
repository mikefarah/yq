package treeops

import "github.com/elliotchance/orderedmap"

type OperatorHandler func(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error)

func TraverseOperator(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	lhs, err := d.getMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}
	return d.getMatchingNodes(lhs, pathNode.Rhs)
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
		lhs.Set(node.getKey(), node)
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

func EqualsOperator(d *dataTreeNavigator, matchMap *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	log.Debugf("-- equalsOperation")
	var results = orderedmap.NewOrderedMap()

	for el := matchMap.Front(); el != nil; el = el.Next() {
		elMap := orderedmap.NewOrderedMap()
		elMap.Set(el.Key, el.Value)
		//need to splat matching nodes, then search through them
		splatter := &PathTreeNode{PathElement: &PathElement{
			PathElementType: PathKey,
			Value:           "*",
			StringValue:     "*",
		}}
		children, err := d.getMatchingNodes(elMap, splatter)
		log.Debugf("-- splatted matches, ")
		if err != nil {
			return nil, err
		}
		for childEl := children.Front(); childEl != nil; childEl = childEl.Next() {
			childMap := orderedmap.NewOrderedMap()
			childMap.Set(childEl.Key, childEl.Value)
			childMatches, errChild := d.getMatchingNodes(childMap, pathNode.Lhs)
			if errChild != nil {
				return nil, errChild
			}

			if containsMatchingValue(childMatches, pathNode.Rhs.PathElement.StringValue) {
				results.Set(childEl.Key, childEl.Value)
			}
		}
	}

	return results, nil
}

func containsMatchingValue(matchMap *orderedmap.OrderedMap, valuePattern string) bool {
	log.Debugf("-- findMatchingValues")

	for el := matchMap.Front(); el != nil; el = el.Next() {
		node := el.Value.(*CandidateNode)
		if Match(node.Node.Value, valuePattern) {
			return true
		}
	}

	return false
}
