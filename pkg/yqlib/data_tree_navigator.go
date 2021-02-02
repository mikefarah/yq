package yqlib

import (
	"fmt"

	logging "gopkg.in/op/go-logging.v1"
)

type DataTreeNavigator interface {
	// given the context and a expressionNode,
	// this will process the against the given expressionNode and return
	// a new context of matching candidates
	GetMatchingNodes(context Context, expressionNode *ExpressionNode) (Context, error)
}

type dataTreeNavigator struct {
}

func NewDataTreeNavigator() DataTreeNavigator {
	return &dataTreeNavigator{}
}

func (d *dataTreeNavigator) GetMatchingNodes(context Context, expressionNode *ExpressionNode) (Context, error) {
	if expressionNode == nil {
		log.Debugf("getMatchingNodes - nothing to do")
		return context, nil
	}
	log.Debugf("Processing Op: %v", expressionNode.Operation.toString())
	if log.IsEnabledFor(logging.DEBUG) {
		for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
			log.Debug(NodeToString(el.Value.(*CandidateNode)))
		}
	}
	log.Debug(">>")
	handler := expressionNode.Operation.OperationType.Handler
	if handler != nil {
		return handler(d, context, expressionNode)
	}
	return Context{}, fmt.Errorf("Unknown operator %v", expressionNode.Operation.OperationType)

}
