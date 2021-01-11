package yqlib

import (
	"errors"

	logging "gopkg.in/op/go-logging.v1"
)

type pathPostFixerInterface interface {
	ConvertToPostfix([]*token) ([]*Operation, error)
}

type pathPostFixer struct {
}

func newPathPostFixer() pathPostFixerInterface {
	return &pathPostFixer{}
}

func popOpToResult(opStack []*token, result []*Operation) ([]*token, []*Operation) {
	var newOp *token
	opStack, newOp = opStack[0:len(opStack)-1], opStack[len(opStack)-1]
	return opStack, append(result, newOp.Operation)
}

func (p *pathPostFixer) ConvertToPostfix(infixTokens []*token) ([]*Operation, error) {
	var result []*Operation
	// surround the whole thing with quotes
	var opStack = []*token{&token{TokenType: OpenBracket}}
	var tokens = append(infixTokens, &token{TokenType: CloseBracket})

	for _, currentToken := range tokens {
		log.Debugf("postfix processing currentToken %v, %v", currentToken.toString(), currentToken.Operation)
		switch currentToken.TokenType {
		case OpenBracket, OpenCollect, OpenCollectObject:
			opStack = append(opStack, currentToken)
		case CloseCollect, CloseCollectObject:
			var opener tokenType = OpenCollect
			var collectOperator *OperationType = Collect
			if currentToken.TokenType == CloseCollectObject {
				opener = OpenCollectObject
				collectOperator = CollectObject
			}
			itemsInMiddle := false
			for len(opStack) > 0 && opStack[len(opStack)-1].TokenType != opener {
				opStack, result = popOpToResult(opStack, result)
				itemsInMiddle = true
			}
			if !itemsInMiddle {
				// must be an empty collection, add the empty object as a LHS parameter
				result = append(result, &Operation{OperationType: Empty})
			}
			if len(opStack) == 0 {
				return nil, errors.New("Bad path expression, got close collect brackets without matching opening bracket")
			}
			// now we should have [] as the last element on the opStack, get rid of it
			opStack = opStack[0 : len(opStack)-1]
			//and append a collect to the opStack
			opStack = append(opStack, &token{TokenType: OperationToken, Operation: &Operation{OperationType: ShortPipe}})
			opStack = append(opStack, &token{TokenType: OperationToken, Operation: &Operation{OperationType: collectOperator}})
		case CloseBracket:
			for len(opStack) > 0 && opStack[len(opStack)-1].TokenType != OpenBracket {
				opStack, result = popOpToResult(opStack, result)
			}
			if len(opStack) == 0 {
				return nil, errors.New("Bad path expression, got close brackets without matching opening bracket")
			}
			// now we should have ( as the last element on the opStack, get rid of it
			opStack = opStack[0 : len(opStack)-1]

		default:
			var currentPrecedence = currentToken.Operation.OperationType.Precedence
			// pop off higher precedent operators onto the result
			for len(opStack) > 0 &&
				opStack[len(opStack)-1].TokenType == OperationToken &&
				opStack[len(opStack)-1].Operation.OperationType.Precedence >= currentPrecedence {
				opStack, result = popOpToResult(opStack, result)
			}
			// add this operator to the opStack
			opStack = append(opStack, currentToken)
		}
	}

	if log.IsEnabledFor(logging.DEBUG) {
		log.Debugf("PostFix Result:")
		for _, currentToken := range result {
			log.Debugf("> %v", currentToken.toString())
		}
	}

	return result, nil
}
