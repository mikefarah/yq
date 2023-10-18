package yqlib

import (
	"container/list"
)

func evalOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("Eval")
	pathExpStrResults, err := d.GetMatchingNodes(context.ReadOnlyClone(), expressionNode.RHS)
	if err != nil {
		return Context{}, err
	}

	expressions := make([]*ExpressionNode, pathExpStrResults.MatchingNodes.Len())
	expIndex := 0
	//parse every expression
	for pathExpStrEntry := pathExpStrResults.MatchingNodes.Front(); pathExpStrEntry != nil; pathExpStrEntry = pathExpStrEntry.Next() {
		expressionStrCandidate := pathExpStrEntry.Value.(*CandidateNode)

		expressions[expIndex], err = ExpressionParser.ParseExpression(expressionStrCandidate.Value)
		if err != nil {
			return Context{}, err
		}

		expIndex++
	}

	results := list.New()

	for matchEl := context.MatchingNodes.Front(); matchEl != nil; matchEl = matchEl.Next() {
		for expIndex = 0; expIndex < len(expressions); expIndex++ {
			result, err := d.GetMatchingNodes(context, expressions[expIndex])
			if err != nil {
				return Context{}, err
			}
			results.PushBackList(result.MatchingNodes)
		}
	}

	return context.ChildContext(results), nil

}
