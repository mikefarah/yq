package yqlib

import (
	"errors"
	"fmt"

	logging "gopkg.in/op/go-logging.v1"
)

type expressionPostFixer interface {
	ConvertToPostfix([]*token) ([]*Operation, error)
}

type expressionPostFixerImpl struct {
}

func newExpressionPostFixer() expressionPostFixer {
	return &expressionPostFixerImpl{}
}

func popOpToResult(opStack []*token, result []*Operation) ([]*token, []*Operation) {
	var newOp *token
	opStack, newOp = opStack[0:len(opStack)-1], opStack[len(opStack)-1]
	log.Debugf("popped %v from opstack to results", newOp.toString(true))
	return opStack, append(result, newOp.Operation)
}

func validateNoOpenTokens(token *token) error {
	switch token.TokenType {
	case openCollect:
		return fmt.Errorf(("bad expression, could not find matching `]`"))
	case openCollectObject:
		return fmt.Errorf(("bad expression, could not find matching `}`"))
	case openBracket:
		return fmt.Errorf(("bad expression, could not find matching `)`"))
	}
	return nil
}

func (p *expressionPostFixerImpl) ConvertToPostfix(infixTokens []*token) ([]*Operation, error) {
	var result []*Operation
	// surround the whole thing with brackets
	var opStack = []*token{{TokenType: openBracket}}
	var tokens = append(infixTokens, &token{TokenType: closeBracket})

	for _, currentToken := range tokens {
		log.Debugf("postfix processing currentToken %v", currentToken.toString(true))
		switch currentToken.TokenType {
		case openBracket, openCollect, openCollectObject:
			opStack = append(opStack, currentToken)
			log.Debugf("put %v onto the opstack", currentToken.toString(true))
		case closeCollect, closeCollectObject:
			var opener tokenType = openCollect
			var collectOperator = collectOpType
			if currentToken.TokenType == closeCollectObject {
				opener = openCollectObject
				collectOperator = collectObjectOpType
			}

			for len(opStack) > 0 && opStack[len(opStack)-1].TokenType != opener {
				missingClosingTokenErr := validateNoOpenTokens(opStack[len(opStack)-1])
				if missingClosingTokenErr != nil {
					return nil, missingClosingTokenErr
				}
				opStack, result = popOpToResult(opStack, result)
			}
			if len(opStack) == 0 {
				return nil, errors.New("bad path expression, got close collect brackets without matching opening bracket")
			}
			// now we should have [ as the last element on the opStack, get rid of it
			opStack = opStack[0 : len(opStack)-1]
			log.Debugf("deleting open bracket from opstack")

			//and append a collect to the result

			// hack - see if there's the optional traverse flag
			// on the close op - move it to the traverse array op
			// allows for .["cat"]?
			prefs := traversePreferences{}
			closeTokenMatch := currentToken.Match
			if closeTokenMatch[len(closeTokenMatch)-1:] == "?" {
				prefs.OptionalTraverse = true
			}
			result = append(result, &Operation{OperationType: collectOperator})
			log.Debugf("put collect onto the result")
			if opener != openCollect {
				result = append(result, &Operation{OperationType: shortPipeOpType})
				log.Debugf("put shortpipe onto the result")
			}

			//traverseArrayCollect is a sneaky op that needs to be included too
			//when closing a ]
			if len(opStack) > 0 && opStack[len(opStack)-1].Operation != nil && opStack[len(opStack)-1].Operation.OperationType == traverseArrayOpType {
				opStack[len(opStack)-1].Operation.Preferences = prefs
				opStack, result = popOpToResult(opStack, result)
			}

		case closeBracket:
			for len(opStack) > 0 && opStack[len(opStack)-1].TokenType != openBracket {
				missingClosingTokenErr := validateNoOpenTokens(opStack[len(opStack)-1])
				if missingClosingTokenErr != nil {
					return nil, missingClosingTokenErr
				}

				opStack, result = popOpToResult(opStack, result)
			}
			if len(opStack) == 0 {
				return nil, errors.New("bad expression, got close brackets without matching opening bracket")
			}
			// now we should have ( as the last element on the opStack, get rid of it
			opStack = opStack[0 : len(opStack)-1]

		default:
			var currentPrecedence = currentToken.Operation.OperationType.Precedence
			// pop off higher precedent operators onto the result
			for len(opStack) > 0 &&
				opStack[len(opStack)-1].TokenType == operationToken &&
				opStack[len(opStack)-1].Operation.OperationType.Precedence > currentPrecedence {
				opStack, result = popOpToResult(opStack, result)
			}
			// add this operator to the opStack
			opStack = append(opStack, currentToken)
			log.Debugf("put %v onto the opstack", currentToken.toString(true))
		}
	}

	log.Debugf("opstackLen: %v", len(opStack))
	if len(opStack) > 0 {
		log.Debugf("opstack:")
		for _, token := range opStack {
			log.Debugf("- %v", token.toString(true))
		}

		return nil, fmt.Errorf("bad expression - probably missing close bracket on %v", opStack[len(opStack)-1].toString(false))
	}

	if log.IsEnabledFor(logging.DEBUG) {
		log.Debugf("PostFix Result:")
		for _, currentToken := range result {
			log.Debugf("> %v", currentToken.toString())
		}
	}

	return result, nil
}
