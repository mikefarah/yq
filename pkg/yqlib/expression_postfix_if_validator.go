package yqlib

import (
	"errors"
	"fmt"
)

// ifBlockValidator tracks then/elif/else/end ordering for each open `if` block
// encountered while converting to postfix, keyed by that block's opening token.
type ifBlockValidator struct {
	lastKeyword map[*token]string
}

func newIfBlockValidator() *ifBlockValidator {
	return &ifBlockValidator{lastKeyword: map[*token]string{}}
}

func findEnclosingBracket(opStack []*token) *token {
	for i := len(opStack) - 1; i >= 0; i-- {
		if opStack[i].TokenType != operationToken {
			return opStack[i]
		}
	}
	return nil
}

// onOpen records the start of a new if block, if currentToken opens one.
func (v *ifBlockValidator) onOpen(currentToken *token) {
	if currentToken.Match == "if" {
		v.lastKeyword[currentToken] = "if"
	}
}

// onKeyword validates a then/elif/else token; it is a no-op for anything else.
func (v *ifBlockValidator) onKeyword(opStack []*token, currentToken *token) error {
	if !tokenIsOpType(currentToken, ifThenOpType) && !tokenIsOpType(currentToken, ifElseOpType) {
		return nil
	}
	keyword := currentToken.Operation.StringValue
	ifToken := findEnclosingBracket(opStack)
	if ifToken == nil || ifToken.TokenType != openBracket || ifToken.Match != "if" {
		return fmt.Errorf("bad expression, `%v` without matching `if`", keyword)
	}
	previous := v.lastKeyword[ifToken]
	if keyword == "then" && previous != "if" && previous != "elif" {
		return fmt.Errorf("bad expression, `then` must follow `if` or `elif`")
	} else if keyword != "then" && previous != "then" {
		return fmt.Errorf("bad expression, `%v` must follow `then`", keyword)
	}
	v.lastKeyword[ifToken] = keyword
	return nil
}

// onClose validates a closing bracket against its opener, in case either side is part of an if block.
func (v *ifBlockValidator) onClose(opener *token, closer *token) error {
	if opener.Match != "if" && closer.Match != "end" {
		return nil
	} else if opener.Match != "if" {
		return errors.New("bad expression, got `end` without matching `if`")
	} else if closer.Match != "end" {
		return errors.New("bad expression, could not find matching `end`")
	}
	previous := v.lastKeyword[opener]
	if previous != "then" && previous != "else" {
		return fmt.Errorf("bad expression, `%v` must be followed by `then`", previous)
	}
	return nil
}
