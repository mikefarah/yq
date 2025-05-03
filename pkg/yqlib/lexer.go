package yqlib

import (
	"fmt"
	"regexp"
)

type expressionTokeniser interface {
	Tokenise(expression string) ([]*token, error)
}

type tokenType uint32

const (
	operationToken = 1 << iota
	openBracket
	closeBracket
	openCollect
	closeCollect
	openCollectObject
	closeCollectObject
	traverseArrayCollect
)

type token struct {
	TokenType            tokenType
	Operation            *Operation
	AssignOperation      *Operation // e.g. tag (GetTag) op becomes AssignTag if '=' follows it
	CheckForPostTraverse bool       // e.g. [1]cat should really be [1].cat
	Match                string
}

func (t *token) toString(detail bool) string {
	switch t.TokenType {
	case operationToken:
		if detail {
			return fmt.Sprintf("%v (%v)", t.Operation.toString(), t.Operation.OperationType.Precedence)
		}
		return t.Operation.toString()
	case openBracket:
		return "("
	case closeBracket:
		return ")"
	case openCollect:
		return "["
	case closeCollect:
		return "]"
	case openCollectObject:
		return "{"
	case closeCollectObject:
		return "}"
	case traverseArrayCollect:
		return ".["

	}
	return "NFI"
}

func unwrap(value string) string {
	return value[1 : len(value)-1]
}

func extractNumberParameter(value string) (int, error) {
	parameterParser := regexp.MustCompile(`.*\(([0-9]+)\)`)
	matches := parameterParser.FindStringSubmatch(value)
	var indent, errParsingInt = parseInt(matches[1])
	if errParsingInt != nil {
		return 0, errParsingInt
	}
	return indent, nil
}

func hasOptionParameter(value string, option string) bool {
	parameterParser := regexp.MustCompile(`.*\([^\)]*\)`)
	matches := parameterParser.FindStringSubmatch(value)
	if len(matches) == 0 {
		return false
	}
	parameterString := matches[0]
	optionParser := regexp.MustCompile(fmt.Sprintf("\\b%v\\b", option))
	return len(optionParser.FindStringSubmatch(parameterString)) > 0
}

func postProcessTokens(tokens []*token) []*token {
	var postProcessedTokens = make([]*token, 0)

	skipNextToken := false

	for index := range tokens {
		if skipNextToken {
			skipNextToken = false
		} else {
			postProcessedTokens, skipNextToken = handleToken(tokens, index, postProcessedTokens)
		}
	}

	return postProcessedTokens
}

func tokenIsOpType(token *token, opType *operationType) bool {
	return token.TokenType == operationToken && token.Operation.OperationType == opType
}

func handleToken(tokens []*token, index int, postProcessedTokens []*token) (tokensAccum []*token, skipNextToken bool) {
	skipNextToken = false
	currentToken := tokens[index]

	log.Debug("processing %v", currentToken.toString(true))

	if currentToken.TokenType == traverseArrayCollect {
		// `.[exp]`` works by creating a traversal array of [self, exp] and piping that into the traverse array operator
		//need to put a traverse array then a collect currentToken
		// do this by adding traverse then converting currentToken to collect

		log.Debug("adding self")
		op := &Operation{OperationType: selfReferenceOpType, StringValue: "SELF"}
		postProcessedTokens = append(postProcessedTokens, &token{TokenType: operationToken, Operation: op})

		log.Debug("adding traverse array")
		op = &Operation{OperationType: traverseArrayOpType, StringValue: "TRAVERSE_ARRAY"}
		postProcessedTokens = append(postProcessedTokens, &token{TokenType: operationToken, Operation: op})

		currentToken = &token{TokenType: openCollect}

	}

	if tokenIsOpType(currentToken, createMapOpType) {
		log.Debugf("tokenIsOpType: createMapOpType")
		// check the previous token is '[', means we are slice, but dont have a first number
		if index > 0 && tokens[index-1].TokenType == traverseArrayCollect {
			log.Debugf("previous token is : traverseArrayOpType")
			// need to put the number 0 before this token, as that is implied
			postProcessedTokens = append(postProcessedTokens, &token{TokenType: operationToken, Operation: createValueOperation(0, "0")})
		}
	}

	if index != len(tokens)-1 && currentToken.AssignOperation != nil &&
		tokenIsOpType(tokens[index+1], assignOpType) {
		log.Debug("its an update assign")
		currentToken.Operation = currentToken.AssignOperation
		currentToken.Operation.UpdateAssign = tokens[index+1].Operation.UpdateAssign
		skipNextToken = true
	}

	log.Debug("adding token to the fixed list")
	postProcessedTokens = append(postProcessedTokens, currentToken)

	if tokenIsOpType(currentToken, createMapOpType) {
		log.Debugf("tokenIsOpType: createMapOpType")
		// check the next token is ']', means we are slice, but dont have a second number
		if index != len(tokens)-1 && tokens[index+1].TokenType == closeCollect {
			log.Debugf("next token is : closeCollect")
			// need to put the number 0 before this token, as that is implied
			lengthOp := &Operation{OperationType: lengthOpType}
			postProcessedTokens = append(postProcessedTokens, &token{TokenType: operationToken, Operation: lengthOp})
		}
	}

	if index != len(tokens)-1 &&
		((currentToken.TokenType == openCollect && tokens[index+1].TokenType == closeCollect) ||
			(currentToken.TokenType == openCollectObject && tokens[index+1].TokenType == closeCollectObject)) {
		log.Debug("adding empty")
		op := &Operation{OperationType: emptyOpType, StringValue: "EMPTY"}
		postProcessedTokens = append(postProcessedTokens, &token{TokenType: operationToken, Operation: op})
	}

	if index != len(tokens)-1 && currentToken.CheckForPostTraverse &&

		(tokenIsOpType(tokens[index+1], traversePathOpType) ||
			(tokens[index+1].TokenType == traverseArrayCollect)) {
		log.Debug("adding pipe because the next thing is traverse")
		op := &Operation{OperationType: shortPipeOpType, Value: "PIPE", StringValue: "."}
		postProcessedTokens = append(postProcessedTokens, &token{TokenType: operationToken, Operation: op})
	}
	if index != len(tokens)-1 && currentToken.CheckForPostTraverse &&
		tokens[index+1].TokenType == openCollect {

		log.Debug("adding traverseArray because next is opencollect")
		op := &Operation{OperationType: traverseArrayOpType}
		postProcessedTokens = append(postProcessedTokens, &token{TokenType: operationToken, Operation: op})
	}
	return postProcessedTokens, skipNextToken
}
