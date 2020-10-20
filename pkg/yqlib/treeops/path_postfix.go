package treeops

import (
	"errors"

	"gopkg.in/op/go-logging.v1"
)

type PathPostFixer interface {
	ConvertToPostfix([]*Token) ([]*PathElement, error)
}

type pathPostFixer struct {
}

func NewPathPostFixer() PathPostFixer {
	return &pathPostFixer{}
}

func popOpToResult(opStack []*Token, result []*PathElement) ([]*Token, []*PathElement) {
	var newOp *Token
	opStack, newOp = opStack[0:len(opStack)-1], opStack[len(opStack)-1]
	var pathElement = PathElement{OperationType: newOp.OperationType, Value: newOp.Value, StringValue: newOp.StringValue}

	if newOp.OperationType == ValueOp {
		var candidateNode = BuildCandidateNodeFrom(newOp)
		pathElement.CandidateNode = candidateNode
	}

	return opStack, append(result, &pathElement)
}

func (p *pathPostFixer) ConvertToPostfix(infixTokens []*Token) ([]*PathElement, error) {
	var result []*PathElement
	// surround the whole thing with quotes
	var opStack = []*Token{&Token{TokenType: OpenBracket, Value: "("}}
	var tokens = append(infixTokens, &Token{TokenType: CloseBracket, Value: ")"})

	for _, token := range tokens {
		log.Debugf("postfix processing token %v", token.Value)
		switch token.TokenType {
		case OpenBracket, OpenCollect:
			opStack = append(opStack, token)
		case CloseCollect:
			for len(opStack) > 0 && opStack[len(opStack)-1].TokenType != OpenCollect {
				opStack, result = popOpToResult(opStack, result)
			}
			if len(opStack) == 0 {
				return nil, errors.New("Bad path expression, got close collect brackets without matching opening bracket")
			}
			// now we should have [] as the last element on the opStack, get rid of it
			opStack = opStack[0 : len(opStack)-1]
			//and append a collect to the opStack
			opStack = append(opStack, &Token{TokenType: Operation, OperationType: Pipe})
			opStack = append(opStack, &Token{TokenType: Operation, OperationType: Collect})
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
			var currentPrecedence = token.OperationType.Precedence
			// pop off higher precedent operators onto the result
			for len(opStack) > 0 && opStack[len(opStack)-1].OperationType.Precedence >= currentPrecedence {
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
