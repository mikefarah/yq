package treeops

import (
	"errors"
	"fmt"
)

var precedenceMap map[int]int

type PathElementType uint32

const (
	PathKey PathElementType = 1 << iota
	ArrayIndex
	Operation
	SelfReference
	OpenBracket
	CloseBracket
)

type OperationType uint32

const (
	None OperationType = 1 << iota
	Traverse
	Or
	And
	Equals
	Assign
	DeleteChild
)

type PathElement struct {
	PathElementType PathElementType
	OperationType   OperationType
	Value           interface{}
	StringValue     string
}

// debugging purposes only
func (p *PathElement) toString() string {
	var result string = ``
	switch p.PathElementType {
	case PathKey:
		result = result + fmt.Sprintf("PathKey - '%v'\n", p.Value)
	case ArrayIndex:
		result = result + fmt.Sprintf("ArrayIndex - '%v'\n", p.Value)
	case SelfReference:
		result = result + fmt.Sprintf("SELF\n")
	case Operation:
		result = result + "Operation - "
		switch p.OperationType {
		case Or:
			result = result + "OR\n"
		case And:
			result = result + "AND\n"
		case Equals:
			result = result + "EQUALS\n"
		case Assign:
			result = result + "ASSIGN\n"
		case Traverse:
			result = result + "TRAVERSE\n"
		case DeleteChild:
			result = result + "DELETE CHILD\n"

		}

	}
	return result
}

func createOperationPathElement(opToken *Token) PathElement {
	var pathElement = PathElement{PathElementType: Operation, OperationType: opToken.OperationType}
	return pathElement
}

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
	var pathElement = createOperationPathElement(operatorToPushToPostFix)
	return opStack, append(result, &pathElement)
}

func (p *pathPostFixer) ConvertToPostfix(infixTokens []*Token) ([]*PathElement, error) {
	var result []*PathElement
	// surround the whole thing with quotes
	var opStack = []*Token{&Token{PathElementType: OpenBracket}}
	var tokens = append(infixTokens, &Token{PathElementType: CloseBracket})

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
			var currentPrecedence = p.precendenceOf(token)
			// pop off higher precedent operators onto the result
			for len(opStack) > 0 && p.precendenceOf(opStack[len(opStack)-1]) >= currentPrecedence {
				opStack, result = popOpToResult(opStack, result)
			}
			// add this operator to the opStack
			opStack = append(opStack, token)
		}
	}
	return result, nil
}

func (p *pathPostFixer) precendenceOf(token *Token) int {
	switch token.OperationType {
	case Or:
		return 10
	case And:
		return 20
	case Equals, DeleteChild:
		return 30
	case Assign:
		return 35
	case Traverse:
		return 40
	}
	return 0
}
