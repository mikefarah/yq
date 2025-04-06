package yqlib

import (
	"fmt"
	"strings"
)

type ExpressionNode struct {
	Operation *Operation
	LHS       *ExpressionNode
	RHS       *ExpressionNode
}

type ExpressionParserInterface interface {
	ParseExpression(expression string) (*ExpressionNode, error)
}

type expressionParserImpl struct {
	pathTokeniser expressionTokeniser
	pathPostFixer expressionPostFixer
}

func newExpressionParser() ExpressionParserInterface {
	return &expressionParserImpl{newParticipleLexer(), newExpressionPostFixer()}
}

func (p *expressionParserImpl) ParseExpression(expression string) (*ExpressionNode, error) {
	log.Debug("Parsing expression: [%v]", expression)
	tokens, err := p.pathTokeniser.Tokenise(expression)
	if err != nil {
		return nil, err
	}
	var Operations []*Operation
	Operations, err = p.pathPostFixer.ConvertToPostfix(tokens)
	if err != nil {
		return nil, err
	}
	return p.createExpressionTree(Operations)
}

func (p *expressionParserImpl) createExpressionTree(postFixPath []*Operation) (*ExpressionNode, error) {
	var stack = make([]*ExpressionNode, 0)

	if len(postFixPath) == 0 {
		return nil, nil
	}

	for _, Operation := range postFixPath {
		var newNode = ExpressionNode{Operation: Operation}
		log.Debugf("pathTree %v ", Operation.toString())
		if Operation.OperationType.NumArgs > 0 {
			numArgs := Operation.OperationType.NumArgs
			switch numArgs {
			case 1:
				if len(stack) < 1 {
					return nil, fmt.Errorf("'%v' expects 1 arg but received none", strings.TrimSpace(Operation.StringValue))
				}
				remaining, rhs := stack[:len(stack)-1], stack[len(stack)-1]
				newNode.RHS = rhs
				stack = remaining
			case 2:
				if len(stack) < 2 {
					return nil, fmt.Errorf("'%v' expects 2 args but there is %v", strings.TrimSpace(Operation.StringValue), len(stack))
				}
				remaining, lhs, rhs := stack[:len(stack)-2], stack[len(stack)-2], stack[len(stack)-1]
				newNode.LHS = lhs
				newNode.RHS = rhs
				stack = remaining
			}
		}
		stack = append(stack, &newNode)
	}
	if len(stack) != 1 {
		return nil, fmt.Errorf("bad expression, please check expression syntax")
	}
	return stack[0], nil
}
