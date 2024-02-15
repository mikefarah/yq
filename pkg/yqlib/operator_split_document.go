package yqlib

func splitDocumentOperator(_ *dataTreeNavigator, context Context, _ *ExpressionNode) (Context, error) {
	log.Debugf("splitDocumentOperator")

	var index uint
	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		candidate.SetDocument(index)
		candidate.SetParent(nil)
		index = index + 1
	}

	return context, nil
}
