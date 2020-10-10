package treeops

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
	Traverse
	Or
	And
	Equals
	EqualsSelf
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
	case Operation:
		result = result + "Operation - "
		switch p.OperationType {
		case Or:
			result = result + "OR\n"
		case And:
			result = result + "AND\n"
		case Equals:
			result = result + "EQUALS\n"
		case EqualsSelf:
			result = result + "EQUALS SELF\n"
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

	precedenceMap[TokenIds["EQUALS_SELF_OPERATOR"]] = 30
	operationTypeMapper[TokenIds["EQUALS_SELF_OPERATOR"]] = EqualsSelf

	precedenceMap[TokenIds["DELETE_CHILD_OPERATOR"]] = 30
	operationTypeMapper[TokenIds["DELETE_CHILD_OPERATOR"]] = DeleteChild

	precedenceMap[TokenIds["ASSIGN_OPERATOR"]] = 35
	operationTypeMapper[TokenIds["ASSIGN_OPERATOR"]] = Assign

	precedenceMap[TokenIds["TRAVERSE_OPERATOR"]] = 40
	operationTypeMapper[TokenIds["TRAVERSE_OPERATOR"]] = Traverse
}

func createOperationPathElement(opToken *lex.Token) PathElement {
	var pathElement = PathElement{PathElementType: Operation, OperationType: operationTypeMapper[opToken.Type]}
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

func (p *pathPostFixer) ConvertToPostfix(infixTokens []*lex.Token) ([]*PathElement, error) {
	var result []*PathElement
	// surround the whole thing with quotes
	var opStack = []*lex.Token{&lex.Token{Type: TokenIds["("]}}
	var tokens = append(infixTokens, &lex.Token{Type: TokenIds[")"]})

	for _, token := range tokens {
		switch token.Type {
		case TokenIds["PATH_KEY"], TokenIds["ARRAY_INDEX"], TokenIds["[+]"], TokenIds["[*]"], TokenIds["**"]:
			var pathElement = PathElement{PathElementType: PathKey, Value: token.Value, StringValue: fmt.Sprintf("%v", token.Value)}
			result = append(result, &pathElement)
		case TokenIds["("]:
			opStack = append(opStack, token)
		case TokenIds[")"]:
			for len(opStack) > 0 && opStack[len(opStack)-1].Type != TokenIds["("] {
				opStack, result = popOpToResult(opStack, result)
			}
			if len(opStack) == 0 {
				return nil, errors.New("Bad path expression, got close brackets without matching opening bracket")
			}
			// now we should have ( as the last element on the opStack, get rid of it
			opStack = opStack[0 : len(opStack)-1]
		default:
			var currentPrecedence = precedenceMap[token.Type]
			// pop off higher precedent operators onto the result
			for len(opStack) > 0 && precedenceMap[opStack[len(opStack)-1].Type] >= currentPrecedence {
				opStack, result = popOpToResult(opStack, result)
			}
			// add this operator to the opStack
			opStack = append(opStack, token)
		}
	}
	return result, nil
}
