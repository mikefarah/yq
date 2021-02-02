package yqlib

func unionOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	lhs, err := d.GetMatchingNodes(context, expressionNode.Lhs)
	if err != nil {
		return Context{}, err
	}
	rhs, err := d.GetMatchingNodes(context, expressionNode.Rhs)
	if err != nil {
		return Context{}, err
	}
	for el := rhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		node := el.Value.(*CandidateNode)
		lhs.MatchingNodes.PushBack(node)
	}
	return lhs, nil
}
