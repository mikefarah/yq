package yqlib

import (
	"fmt"
)

type flattenPreferences struct {
	depth int
}

func flatten(node *CandidateNode, depth int) {
	if depth == 0 {
		return
	}
	if node.Kind != SequenceNode {
		return
	}
	content := node.Content
	newSeq := make([]*CandidateNode, 0)

	for i := 0; i < len(content); i++ {
		child := content[i]
		if child.Kind == AliasNode && child.Alias != nil {
			child = child.Alias.Copy()
		}
		if child.Kind == SequenceNode {
			flatten(child, depth-1)
			for j := 0; j < len(child.Content); j++ {
				newSeq = append(newSeq, child.Content[j])
			}
		} else {
			newSeq = append(newSeq, content[i])
		}
	}
	node.Content = make([]*CandidateNode, 0)
	node.AddChildren(newSeq)
}

func flattenOp(_ *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debugf("flatten Operator")
	depth := expressionNode.Operation.Preferences.(flattenPreferences).depth

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		if candidate.Kind != SequenceNode {
			return Context{}, fmt.Errorf("only arrays are supported for flatten")
		}

		flatten(candidate, depth)

	}

	return context, nil

}
