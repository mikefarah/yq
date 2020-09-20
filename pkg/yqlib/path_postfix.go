package yqlib

import (
	"errors"
	"fmt"

	lex "github.com/timtadh/lexmachine"
)

var precedenceMap map[int]int

type PathElementType uint32

const (
	PathKey PathElementType = 1 << iota
	ArrayIndex
	Operation
)

type OperationType uint32

const (
	None OperationType = 1 << iota
	Or
	And
	Equals
)

type PathElement struct {
	PathElementType PathElementType
	OperationType   OperationType
	Value           interface{}
	ChildElements   [][]*PathElement
	Finished        bool
}

// debugging purposes only
func (p *PathElement) toString() string {
	var result string = `Type: `
	switch p.PathElementType {
	case PathKey:
		result = result + fmt.Sprintf("PathKey - %v\n", p.Value)
	case ArrayIndex:
		result = result + fmt.Sprintf("ArrayIndex - %v\n", p.Value)
	case Operation:
		result = result + "Operation - "
		switch p.OperationType {
		case Or:
			result = result + "OR\n"
		case And:
			result = result + "AND\n"
		case Equals:
			result = result + "EQUALS\n"
		}
	}
	return result
}

var operationTypeMapper map[int]OperationType

func initMaps() {
	precedenceMap = make(map[int]int)
	operationTypeMapper = make(map[int]OperationType)

	precedenceMap[TokenIds["("]] = 0

	precedenceMap[TokenIds["OR_OPERATOR"]] = 10
	operationTypeMapper[TokenIds["OR_OPERATOR"]] = Or

	precedenceMap[TokenIds["AND_OPERATOR"]] = 20
	operationTypeMapper[TokenIds["AND_OPERATOR"]] = And

	precedenceMap[TokenIds["EQUALS_OPERATOR"]] = 30
	operationTypeMapper[TokenIds["EQUALS_OPERATOR"]] = Equals
}

func createOperationPathElement(opToken *lex.Token) PathElement {
	var childElements = make([][]*PathElement, 2)
	var pathElement = PathElement{PathElementType: Operation, OperationType: operationTypeMapper[opToken.Type], ChildElements: childElements}
	return pathElement
}

type PathPostFixer interface {
	ConvertToPostfix([]*lex.Token) ([]*PathElement, error)
}

type pathPostFixer struct {
}

func NewPathPostFixer() PathPostFixer {
	return &pathPostFixer{}
}

func popOpToResult(opStack []*lex.Token, result []*PathElement) ([]*lex.Token, []*PathElement) {
	var operatorToPushToPostFix *lex.Token
	opStack, operatorToPushToPostFix = opStack[0:len(opStack)-1], opStack[len(opStack)-1]
	var pathElement = createOperationPathElement(operatorToPushToPostFix)
	return opStack, append(result, &pathElement)
}

func finishPathKey(result []*PathElement) {
	if len(result) > 0 {
		//need to mark PathKey elements as finished so we
		//stop appending PathKeys as children
		result[len(result)-1].Finished = true
	}
}

func (p *pathPostFixer) ConvertToPostfix(infixTokens []*lex.Token) ([]*PathElement, error) {
	var result []*PathElement
	// surround the whole thing with quotes
	var opStack = []*lex.Token{&lex.Token{Type: TokenIds["("]}}
	var tokens = append(infixTokens, &lex.Token{Type: TokenIds[")"]})

	for _, token := range tokens {
		switch token.Type {
		case TokenIds["PATH_KEY"]: // handle splats and array appends here too
			var emptyArray = [][]*PathElement{make([]*PathElement, 0)}
			var pathElement = PathElement{PathElementType: PathKey, Value: token.Value, ChildElements: emptyArray}

			if len(result) > 0 && result[len(result)-1].PathElementType == PathKey && !result[len(result)-1].Finished {
				var lastElement = result[len(result)-1]
				lastElement.ChildElements[0] = append(lastElement.ChildElements[0], &pathElement)
			} else {
				result = append(result, &pathElement)
			}
		case TokenIds["("]:
			opStack = append(opStack, token)
			finishPathKey(result)
		case TokenIds["OR_OPERATOR"], TokenIds["AND_OPERATOR"], TokenIds["EQUALS_OPERATOR"]:
			var currentPrecedence = precedenceMap[token.Type]
			// pop off higher precedent operators onto the result
			for len(opStack) > 0 && precedenceMap[opStack[len(opStack)-1].Type] > currentPrecedence {
				opStack, result = popOpToResult(opStack, result)
			}
			// add this operator to the opStack
			opStack = append(opStack, token)
			finishPathKey(result)
		case TokenIds[")"]:
			for len(opStack) > 0 && opStack[len(opStack)-1].Type != TokenIds["("] {
				opStack, result = popOpToResult(opStack, result)
			}
			if len(opStack) == 0 {
				return nil, errors.New("Bad path expression, got close brackets without matching opening bracket")
			}
			// now we should have ( as the last element on the opStack, get rid of it
			opStack = opStack[0 : len(opStack)-1]
			finishPathKey(result)
		}
	}
	return result, nil
}
