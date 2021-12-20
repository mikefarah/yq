package yqlib

func splitDocumentOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- splitDocumentOperator")

	var index uint
	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		candidate.Document = index
		index = index + 1
	}

	return context, nil
}
