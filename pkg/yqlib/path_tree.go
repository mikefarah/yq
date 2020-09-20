package yqlib

import lex "github.com/timtadh/lexmachine"

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
