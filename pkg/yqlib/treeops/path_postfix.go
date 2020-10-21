package treeops

import (
	"errors"

	"gopkg.in/op/go-logging.v1"
)

type PathPostFixer interface {
	ConvertToPostfix([]*Token) ([]*Operation, error)
}

type pathPostFixer struct {
}

func NewPathPostFixer() PathPostFixer {
	return &pathPostFixer{}
}

func popOpToResult(opStack []*Token, result []*Operation) ([]*Token, []*Operation) {
	var newOp *Token
	opStack, newOp = opStack[0:len(opStack)-1], opStack[len(opStack)-1]
	return opStack, append(result, newOp.Operation)
}

func (p *pathPostFixer) ConvertToPostfix(infixTokens []*Token) ([]*Operation, error) {
	var result []*Operation
	// surround the whole thing with quotes
	var opStack = []*Token{&Token{TokenType: OpenBracket}}
	var tokens = append(infixTokens, &Token{TokenType: CloseBracket})

	for _, token := range tokens {
		log.Debugf("postfix processing token %v, %v", token.toString(), token.Operation)
		switch token.TokenType {
		case OpenBracket, OpenCollect, OpenCollectObject:
			opStack = append(opStack, token)
		case CloseCollect, CloseCollectObject:
			var opener TokenType = OpenCollect
			var collectOperator *OperationType = Collect
			if token.TokenType == CloseCollectObject {
				opener = OpenCollectObject
				collectOperator = CollectObject
			}
			for len(opStack) > 0 && opStack[len(opStack)-1].TokenType != opener {
				opStack, result = popOpToResult(opStack, result)
			}
			if len(opStack) == 0 {
				return nil, errors.New("Bad path expression, got close collect brackets without matching opening bracket")
			}
			// now we should have [] as the last element on the opStack, get rid of it
			opStack = opStack[0 : len(opStack)-1]
			//and append a collect to the opStack
			opStack = append(opStack, &Token{TokenType: OperationToken, Operation: &Operation{OperationType: Pipe}})
			opStack = append(opStack, &Token{TokenType: OperationToken, Operation: &Operation{OperationType: collectOperator}})
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
			var currentPrecedence = token.Operation.OperationType.Precedence
			// pop off higher precedent operators onto the result
			for len(opStack) > 0 &&
				opStack[len(opStack)-1].TokenType == OperationToken &&
				opStack[len(opStack)-1].Operation.OperationType.Precedence >= currentPrecedence {
				opStack, result = popOpToResult(opStack, result)
			}
			// add this operator to the opStack
			opStack = append(opStack, token)
		}
	}

	if log.IsEnabledFor(logging.DEBUG) {
		log.Debugf("PostFix Result:")
		for _, token := range result {
			log.Debugf("> %v", token.toString())
		}
	}

	return result, nil
}
