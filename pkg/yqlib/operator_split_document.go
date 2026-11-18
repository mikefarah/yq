package yqlib

func splitDocumentOperator(_ *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("splitDocumentOperator")

	index := context.splitDocumentIndices[expressionNode]
	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		candidate.SetDocument(index)
		candidate.SetParent(nil)
		index = index + 1
	}
	context.splitDocumentIndices[expressionNode] = index

	return context, nil
}
