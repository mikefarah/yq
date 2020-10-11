package treeops

import (
	"errors"
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
	var operatorToPushToPostFix *Token
	opStack, operatorToPushToPostFix = opStack[0:len(opStack)-1], opStack[len(opStack)-1]
	var pathElement = PathElement{PathElementType: Operation, OperationType: operatorToPushToPostFix.OperationType}
	return opStack, append(result, &pathElement)
}

func (p *pathPostFixer) ConvertToPostfix(infixTokens []*Token) ([]*PathElement, error) {
	var result []*PathElement
	// surround the whole thing with quotes
	var opStack = []*Token{&Token{PathElementType: OpenBracket, OperationType: None}}
	var tokens = append(infixTokens, &Token{PathElementType: CloseBracket, OperationType: None})

	for _, token := range tokens {
		switch token.PathElementType {
		case PathKey, ArrayIndex, SelfReference:
			var pathElement = PathElement{PathElementType: token.PathElementType, Value: token.Value, StringValue: token.StringValue}
			result = append(result, &pathElement)
		case OpenBracket:
			opStack = append(opStack, token)

		case CloseBracket:
			for len(opStack) > 0 && opStack[len(opStack)-1].PathElementType != OpenBracket {
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
	return result, nil
}
