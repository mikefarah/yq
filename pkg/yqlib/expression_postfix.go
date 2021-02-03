package yqlib

import (
	"errors"

	logging "gopkg.in/op/go-logging.v1"
)

type expressionPostFixer interface {
	ConvertToPostfix([]*token) ([]*Operation, error)
}

type expressionPostFixerImpl struct {
}

func newExpressionPostFixer() expressionPostFixer {
	return &expressionPostFixerImpl{}
}

func popOpToResult(opStack []*token, result []*Operation) ([]*token, []*Operation) {
	var newOp *token
	opStack, newOp = opStack[0:len(opStack)-1], opStack[len(opStack)-1]
	log.Debugf("popped %v from opstack to results", newOp.toString(true))
	return opStack, append(result, newOp.Operation)
}

func (p *expressionPostFixerImpl) ConvertToPostfix(infixTokens []*token) ([]*Operation, error) {
	var result []*Operation
	// surround the whole thing with brackets
	var opStack = []*token{{TokenType: openBracket}}
	var tokens = append(infixTokens, &token{TokenType: closeBracket})

	for _, currentToken := range tokens {
		log.Debugf("postfix processing currentToken %v", currentToken.toString(true))
		switch currentToken.TokenType {
		case openBracket, openCollect, openCollectObject:
			opStack = append(opStack, currentToken)
			log.Debugf("put %v onto the opstack", currentToken.toString(true))
		case closeCollect, closeCollectObject:
			var opener tokenType = openCollect
			var collectOperator *operationType = collectOpType
			if currentToken.TokenType == closeCollectObject {
				opener = openCollectObject
				collectOperator = collectObjectOpType
			}
			itemsInMiddle := false
			for len(opStack) > 0 && opStack[len(opStack)-1].TokenType != opener {
				opStack, result = popOpToResult(opStack, result)
				itemsInMiddle = true
			}
			if !itemsInMiddle {
				// must be an empty collection, add the empty object as a LHS parameter
				result = append(result, &Operation{OperationType: emptyOpType})
				log.Debugf("put EMPTY onto the result")
			}
			if len(opStack) == 0 {
				return nil, errors.New("Bad path expression, got close collect brackets without matching opening bracket")
			}
			// now we should have [ as the last element on the opStack, get rid of it
			opStack = opStack[0 : len(opStack)-1]
			log.Debugf("deleteing open bracket from opstack")

			//and append a collect to the opStack
			//WHY AM I NOT PUTTING THIS STRAIGHT ONTO THE RESULT LIKE EMPTY?
			result = append(result, &Operation{OperationType: collectOperator})
			log.Debugf("put collect onto the result")
			result = append(result, &Operation{OperationType: shortPipeOpType})
			log.Debugf("put shortpipe onto the result")

			// opStack = append(opStack, &token{TokenType: operationToken, Operation: &Operation{OperationType: shortPipeOpType}})
			// log.Debugf("put shortPipe onto the opstack")
			// opStack = append(opStack, &token{TokenType: operationToken, Operation: &Operation{OperationType: collectOperator}})
			// log.Debugf("put collect onto the opstack")
		case closeBracket:
			for len(opStack) > 0 && opStack[len(opStack)-1].TokenType != openBracket {
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
				opStack[len(opStack)-1].TokenType == operationToken &&
				opStack[len(opStack)-1].Operation.OperationType.Precedence > currentPrecedence {
				opStack, result = popOpToResult(opStack, result)
			}
			// add this operator to the opStack
			opStack = append(opStack, currentToken)
			log.Debugf("put %v onto the opstack", currentToken.toString(true))
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
