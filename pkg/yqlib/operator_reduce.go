package yqlib

import (
	"container/list"
	"fmt"
)

func reduceOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- reduceOp")
	//.a as $var reduce (0; . + $var)
	//lhs is the assignment operator
	//rhs is the reduce block
	// '.' refers to the current accumulator, initialised to 0
	// $var references a single element from the .a

	// ensure lhs is actually an assignment
	// and rhs is a block (empty)
	if expressionNode.Lhs.Operation.OperationType != assignVariableOpType {
		return Context{}, fmt.Errorf("reduce must be given a variables assignment, got %v instead", expressionNode.Lhs.Operation.OperationType.Type)
	} else if expressionNode.Rhs.Operation.OperationType != blockOpType {
		return Context{}, fmt.Errorf("reduce must be given a block, got %v instead", expressionNode.Rhs.Operation.OperationType.Type)
	}

	arrayExpNode := expressionNode.Lhs.Lhs
	array, err := d.GetMatchingNodes(context, arrayExpNode)

	log.Debugf("array of %v things", array.MatchingNodes.Len())

	if err != nil {
		return Context{}, err
	}

	variableName := expressionNode.Lhs.Rhs.Operation.StringValue

	initExp := expressionNode.Rhs.Lhs

	accum, err := d.GetMatchingNodes(context, initExp)
	if err != nil {
		return Context{}, err
	}

	log.Debugf("with variable %v", variableName)

	blockExp := expressionNode.Rhs.Rhs
	for el := array.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		log.Debugf("REDUCING WITH %v", NodeToString(candidate))
		l := list.New()
		l.PushBack(candidate)
		accum.SetVariable(variableName, l)

		accum, err = d.GetMatchingNodes(accum, blockExp)
		if err != nil {
			return Context{}, err
		}
	}

	return accum, nil
}
