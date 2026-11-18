package yqlib

import (
	"container/list"
	"fmt"

	"github.com/elliotchance/orderedmap"
)

func unique(d *dataTreeNavigator, context Context, _ *ExpressionNode) (Context, error) {
	selfExpression := &ExpressionNode{Operation: &Operation{OperationType: selfReferenceOpType}}
	uniqueByExpression := &ExpressionNode{Operation: &Operation{OperationType: uniqueByOpType}, RHS: selfExpression}
	return uniqueBy(d, context, uniqueByExpression)

}

func uniqueBy(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debugf("uniqueBy Operator")
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		if candidate.Kind != SequenceNode {
			return Context{}, fmt.Errorf("only arrays are supported for unique")
		}

		var newMatches = orderedmap.NewOrderedMap()
		for _, child := range candidate.Content {
			rhs, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(child), expressionNode.RHS)

			if err != nil {
				return Context{}, err
			}

			keyValue, err := getUniqueKeyValue(rhs)
			if err != nil {
				return Context{}, err
			}

			_, exists := newMatches.Get(keyValue)

			if !exists {
				newMatches.Set(keyValue, child)
			}
		}
		resultNode := candidate.CreateReplacementWithComments(SequenceNode, "!!seq", candidate.Style)
		for el := newMatches.Front(); el != nil; el = el.Next() {
			resultNode.AddChild(el.Value.(*CandidateNode))
		}

		results.PushBack(resultNode)
	}

	return context.ChildContext(results), nil

}

func getUniqueKeyValue(rhs Context) (string, error) {
	keyValue := "null"
	var err error

	if rhs.MatchingNodes.Len() > 0 {
		first := rhs.MatchingNodes.Front()
		keyCandidate := first.Value.(*CandidateNode)
		keyValue = keyCandidate.Value
		if keyCandidate.Kind != ScalarNode {
			// Encode a comment-free copy so that uniqueness is decided by the
			// data only. Otherwise a comment attached to one of two otherwise
			// equal nodes changes its encoding and they are wrongly kept as
			// distinct (see #2491).
			keyValue, err = encodeToString(withoutComments(keyCandidate), encoderPreferences{YamlFormat, 0})
		}
	}
	return keyValue, err
}

// withoutComments returns a deep copy of the node with all head, line and foot
// comments removed throughout the tree, so comments do not affect equality.
func withoutComments(node *CandidateNode) *CandidateNode {
	clone := node.Copy()
	clearComments(clone)
	return clone
}

func clearComments(node *CandidateNode) {
	node.HeadComment = ""
	node.LineComment = ""
	node.FootComment = ""
	if node.Key != nil {
		clearComments(node.Key)
	}
	for _, child := range node.Content {
		clearComments(child)
	}
}
