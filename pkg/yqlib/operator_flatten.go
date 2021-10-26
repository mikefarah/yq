package yqlib

import (
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

type flattenPreferences struct {
	depth int
}

func flatten(node *yaml.Node, depth int) {
	if depth == 0 {
		return
	}
	if node.Kind != yaml.SequenceNode {
		return
	}
	content := node.Content
	newSeq := make([]*yaml.Node, 0)

	for i := 0; i < len(content); i++ {
		if content[i].Kind == yaml.SequenceNode {
			flatten(content[i], depth-1)
			for j := 0; j < len(content[i].Content); j++ {
				newSeq = append(newSeq, content[i].Content[j])
			}
		} else {
			newSeq = append(newSeq, content[i])
		}
	}
	node.Content = newSeq
}

func flattenOp(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debugf("-- flatten Operator")
	depth := expressionNode.Operation.Preferences.(flattenPreferences).depth

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		candidateNode := unwrapDoc(candidate.Node)
		if candidateNode.Kind != yaml.SequenceNode {
			return Context{}, fmt.Errorf("Only arrays are supported for flatten")
		}

		flatten(candidateNode, depth)

	}

	return context, nil

}
