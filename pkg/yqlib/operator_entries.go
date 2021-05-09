package yqlib

import (
	"container/list"
	"fmt"
	yaml "gopkg.in/yaml.v3"
)

func entrySeqFor(key *yaml.Node, value *yaml.Node) *yaml.Node {
	var keyKey = &yaml.Node{Kind:  yaml.ScalarNode, Tag: "!!str", Value: "key"}
	var valueKey = &yaml.Node{Kind:  yaml.ScalarNode, Tag: "!!str", Value: "value"}

	return &yaml.Node{
		Kind:  yaml.MappingNode, 
		Tag: "!!map", 
		Content: []*yaml.Node{keyKey, key, valueKey, value},
	}
}

func toEntriesFromMap(candidateNode *CandidateNode) *CandidateNode {
	var sequence = &yaml.Node{Kind:  yaml.SequenceNode, Tag: "!!seq"}
	var entriesNode = candidateNode.CreateChild(nil, sequence)

	var contents = unwrapDoc(candidateNode.Node).Content
	for index := 0; index < len(contents); index = index + 2 {
		key := contents[index]
		value := contents[index+1]

		sequence.Content = append(sequence.Content, entrySeqFor(key, value))
	}
	return entriesNode
}

func toEntriesfromSeq(candidateNode *CandidateNode) *CandidateNode {
	var sequence = &yaml.Node{Kind:  yaml.SequenceNode, Tag: "!!seq"}
	var entriesNode = candidateNode.CreateChild(nil, sequence)

	var contents = unwrapDoc(candidateNode.Node).Content
	for index := 0; index < len(contents); index = index + 1 {
		key := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: fmt.Sprintf("%v", index)}
		value := contents[index]

		sequence.Content = append(sequence.Content, entrySeqFor(key, value))
	}
	return entriesNode
}

func toEntriesOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	var results = list.New()
	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		candidateNode := unwrapDoc(candidate.Node)

		switch candidateNode.Kind {
		case yaml.MappingNode:
			results.PushBack(toEntriesFromMap(candidate))

		case yaml.SequenceNode:
			results.PushBack(toEntriesfromSeq(candidate))
		default:
			return Context{}, fmt.Errorf("%v has no keys", candidate.Node.Tag)
		}
	}

	return context.ChildContext(results), nil

}