package yqlib

import (
	"container/list"
	"strconv"

	yaml "gopkg.in/yaml.v3"
)

func hasOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {

	log.Debugf("-- hasOperation")
	var results = list.New()

	rhs, err := d.GetMatchingNodes(matchingNodes, pathNode.Rhs)
	wanted := rhs.Front().Value.(*CandidateNode).Node
	wantedKey := wanted.Value

	if err != nil {
		return nil, err
	}

	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		// grab the first value
		var contents = candidate.Node.Content
		switch candidate.Node.Kind {
		case yaml.MappingNode:
			candidateHasKey := false
			for index := 0; index < len(contents) && !candidateHasKey; index = index + 2 {
				key := contents[index]
				if key.Value == wantedKey {
					candidateHasKey = true
				}
			}
			results.PushBack(createBooleanCandidate(candidate, candidateHasKey))
		case yaml.SequenceNode:
			candidateHasKey := false
			if wanted.Tag == "!!int" {
				var number, errParsingInt = strconv.ParseInt(wantedKey, 10, 64) // nolint
				if errParsingInt != nil {
					return nil, errParsingInt
				}
				candidateHasKey = int64(len(contents)) > number
			}
			results.PushBack(createBooleanCandidate(candidate, candidateHasKey))
		default:
			results.PushBack(createBooleanCandidate(candidate, false))
		}
	}
	return results, nil
}
