package treeops

import (
	"container/list"

	"gopkg.in/yaml.v3"
)

func CreateMapOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("-- createMapOperation")
	var path []interface{} = nil
	var document uint = 0
	if matchingNodes.Front() != nil {
		sample := matchingNodes.Front().Value.(*CandidateNode)
		path = sample.Path
		document = sample.Document
	}

	mapPairs, err := crossFunction(d, matchingNodes, pathNode,
		func(d *dataTreeNavigator, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
			node := yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
			log.Debugf("LHS:", lhs.Node.Value)
			log.Debugf("RHS:", rhs.Node.Value)
			node.Content = []*yaml.Node{
				lhs.Node,
				rhs.Node,
			}

			return &CandidateNode{Node: &node, Document: document, Path: path}, nil
		})

	if err != nil {
		return nil, err
	}
	//wrap up all the pairs into an array
	node := yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	for mapPair := mapPairs.Front(); mapPair != nil; mapPair = mapPair.Next() {
		mapPairCandidate := mapPair.Value.(*CandidateNode)
		log.Debugf("Collecting %v into sequence", NodeToString(mapPairCandidate))
		node.Content = append(node.Content, mapPairCandidate.Node)
	}
	return nodeToMap(&CandidateNode{Node: &node, Document: document, Path: path}), nil
}
