package yqlib

import (
	"container/list"
	"sort"

	yaml "gopkg.in/yaml.v3"
)

func sortKeysOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {

	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		rhs, err := d.GetMatchingNodes(nodeToMap(candidate), pathNode.Rhs)
		if err != nil {
			return nil, err
		}

		for childEl := rhs.Front(); childEl != nil; childEl = childEl.Next() {
			node := unwrapDoc(childEl.Value.(*CandidateNode).Node)
			if node.Kind == yaml.MappingNode {
				sortKeys(node)
			}
			if err != nil {
				return nil, err
			}
		}

	}
	return matchingNodes, nil
}

func sortKeys(node *yaml.Node) {
	keys := make([]string, len(node.Content)/2)
	keyBucket := map[string]*yaml.Node{}
	valueBucket := map[string]*yaml.Node{}
	var contents = node.Content
	for index := 0; index < len(contents); index = index + 2 {
		key := contents[index]
		value := contents[index+1]
		keys[index/2] = key.Value
		keyBucket[key.Value] = key
		valueBucket[key.Value] = value
	}
	sort.Strings(keys)
	sortedContent := make([]*yaml.Node, len(node.Content))
	for index := 0; index < len(keys); index = index + 1 {
		keyString := keys[index]
		sortedContent[index*2] = keyBucket[keyString]
		sortedContent[1+(index*2)] = valueBucket[keyString]
	}
	node.Content = sortedContent
}
