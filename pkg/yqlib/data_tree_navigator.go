package yqlib

import (
	"fmt"

	logging "gopkg.in/op/go-logging.v1"
)

type DataTreeNavigator interface {
	// given the context and an expressionNode,
	// this will process the against the given expressionNode and return
	// a new context of matching candidates
	GetMatchingNodes(context Context, expressionNode *ExpressionNode) (Context, error)

	DeeplyAssign(context Context, path []interface{}, rhsNode *CandidateNode) error
	DeeplyAssignKey(context Context, path []interface{}, rhsNode *CandidateNode) error
}

type dataTreeNavigator struct {
}

func NewDataTreeNavigator() DataTreeNavigator {
	return &dataTreeNavigator{}
}

func (d *dataTreeNavigator) DeeplyAssign(context Context, path []interface{}, rhsCandidateNode *CandidateNode) error {
	return d.deeplyAssign(context, path, rhsCandidateNode, false)
}

func (d *dataTreeNavigator) DeeplyAssignKey(context Context, path []interface{}, rhsCandidateNode *CandidateNode) error {
	return d.deeplyAssign(context, path, rhsCandidateNode, true)
}

func (d *dataTreeNavigator) deeplyAssign(context Context, path []interface{}, rhsCandidateNode *CandidateNode, targetKey bool) error {

	assignmentOp := &Operation{OperationType: assignOpType, Preferences: assignPreferences{}}

	if rhsCandidateNode.Kind == MappingNode {
		log.Debug("DeeplyAssign: deeply merging object")
		// if the rhs is a map, we need to deeply merge it in.
		// otherwise we'll clobber any existing fields
		assignmentOp = &Operation{OperationType: multiplyAssignOpType, Preferences: multiplyPreferences{
			AppendArrays:  true,
			TraversePrefs: traversePreferences{DontFollowAlias: true},
			AssignPrefs:   assignPreferences{},
		}}
	}

	rhsOp := &Operation{OperationType: valueOpType, CandidateNode: rhsCandidateNode}

	assignmentOpNode := &ExpressionNode{
		Operation: assignmentOp,
		LHS:       createTraversalTree(path, traversePreferences{}, targetKey),
		RHS:       &ExpressionNode{Operation: rhsOp},
	}

	_, err := d.GetMatchingNodes(context, assignmentOpNode)
	return err
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
	handler := expressionNode.Operation.OperationType.Handler
	if handler != nil {
		return handler(d, context, expressionNode)
	}
	return Context{}, fmt.Errorf("Unknown operator %v", expressionNode.Operation.OperationType)

}
