package yqlib

import lex "github.com/timtadh/lexmachine"

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
	ChildEquals
)

type PathElement struct {
	PathElementType PathElementType
	OperationType   OperationType
	Value           interface{}
	ChildElements   [][]*PathElement
}

func parseTree(tokens []*lex.Token, currentElement *PathElement, allElements []*PathElement) []*PathElement {
	currentToken, remainingTokens := tokens[0], tokens[1:]

	switch currentToken.Type {
	case TokenIds["PATH_KEY"]:
		currentElement.PathElementType = PathKey
		currentElement.OperationType = None
		currentElement.Value = currentToken.Value
	}

	if len(remainingTokens) == 0 {
		return append(allElements, currentElement)
	}
	return parseTree(remainingTokens, &PathElement{}, append(allElements, currentElement))

}

func ParseTree(tokens []*lex.Token) []*PathElement {
	if len(tokens) == 0 {
		return make([]*PathElement, 0)
	}
	return parseTree(tokens, &PathElement{}, make([]*PathElement, 0))
}
