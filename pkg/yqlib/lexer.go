package yqlib

import (
	"fmt"
	"regexp"
	"strconv"
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
	if t.TokenType == operationToken {
		if detail {
			return fmt.Sprintf("%v (%v)", t.Operation.toString(), t.Operation.OperationType.Precedence)
		}
		return t.Operation.toString()
	} else if t.TokenType == openBracket {
		return "("
	} else if t.TokenType == closeBracket {
		return ")"
	} else if t.TokenType == openCollect {
		return "["
	} else if t.TokenType == closeCollect {
		return "]"
	} else if t.TokenType == openCollectObject {
		return "{"
	} else if t.TokenType == closeCollectObject {
		return "}"
	} else if t.TokenType == traverseArrayCollect {
		return ".["

	} else {
		return "NFI"
	}
}

func unwrap(value string) string {
	return value[1 : len(value)-1]
}

func extractNumberParameter(value string) (int, error) {
	parameterParser := regexp.MustCompile(`.*\(([0-9]+)\)`)
	matches := parameterParser.FindStringSubmatch(value)
	var indent, errParsingInt = strconv.ParseInt(matches[1], 10, 32)
	if errParsingInt != nil {
		return 0, errParsingInt
	}
	return int(indent), nil
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

func handleToken(tokens []*token, index int, postProcessedTokens []*token) (tokensAccum []*token, skipNextToken bool) {
	skipNextToken = false
	currentToken := tokens[index]

	log.Debug("processing %v", currentToken.toString(true))

	if currentToken.TokenType == traverseArrayCollect {
		// `.[exp]`` works by creating a traversal array of [self, exp] and piping that into the traverse array operator
		//need to put a traverse array then a collect currentToken
		// do this by adding traverse then converting currentToken to collect

		log.Debug("  adding self")
		op := &Operation{OperationType: selfReferenceOpType, StringValue: "SELF"}
		postProcessedTokens = append(postProcessedTokens, &token{TokenType: operationToken, Operation: op})

		log.Debug("  adding traverse array")
		op = &Operation{OperationType: traverseArrayOpType, StringValue: "TRAVERSE_ARRAY"}
		postProcessedTokens = append(postProcessedTokens, &token{TokenType: operationToken, Operation: op})

		currentToken = &token{TokenType: openCollect}

	}

	if index != len(tokens)-1 && currentToken.AssignOperation != nil &&
		tokens[index+1].TokenType == operationToken &&
		tokens[index+1].Operation.OperationType == assignOpType {
		log.Debug("  its an update assign")
		currentToken.Operation = currentToken.AssignOperation
		currentToken.Operation.UpdateAssign = tokens[index+1].Operation.UpdateAssign
		skipNextToken = true
	}

	log.Debug("  adding token to the fixed list")
	postProcessedTokens = append(postProcessedTokens, currentToken)

	if index != len(tokens)-1 &&
		((currentToken.TokenType == openCollect && tokens[index+1].TokenType == closeCollect) ||
			(currentToken.TokenType == openCollectObject && tokens[index+1].TokenType == closeCollectObject)) {
		log.Debug("  adding empty")
		op := &Operation{OperationType: emptyOpType, StringValue: "EMPTY"}
		postProcessedTokens = append(postProcessedTokens, &token{TokenType: operationToken, Operation: op})
	}

	if index != len(tokens)-1 && currentToken.CheckForPostTraverse &&
		((tokens[index+1].TokenType == operationToken && (tokens[index+1].Operation.OperationType == traversePathOpType)) ||
			(tokens[index+1].TokenType == traverseArrayCollect)) {
		log.Debug("  adding pipe because the next thing is traverse")
		op := &Operation{OperationType: shortPipeOpType, Value: "PIPE", StringValue: "."}
		postProcessedTokens = append(postProcessedTokens, &token{TokenType: operationToken, Operation: op})
	}
	if index != len(tokens)-1 && currentToken.CheckForPostTraverse &&
		tokens[index+1].TokenType == openCollect {

		log.Debug("  adding traverArray because next is opencollect")
		op := &Operation{OperationType: traverseArrayOpType}
		postProcessedTokens = append(postProcessedTokens, &token{TokenType: operationToken, Operation: op})
	}
	return postProcessedTokens, skipNextToken
}
