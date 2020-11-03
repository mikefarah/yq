package yqlib

import "fmt"

var myPathTokeniser = NewPathTokeniser()
var myPathPostfixer = NewPathPostFixer()

type PathTreeNode struct {
	Operation *Operation
	Lhs         *PathTreeNode
	Rhs         *PathTreeNode
}

type PathTreeCreator interface {
	ParsePath(path string) (*PathTreeNode, error)
	CreatePathTree(postFixPath []*Operation) (*PathTreeNode, error)
}

type pathTreeCreator struct {
}

func NewPathTreeCreator() PathTreeCreator {
	return &pathTreeCreator{}
}

func (p *pathTreeCreator) ParsePath(path string) (*PathTreeNode, error) {
	tokens, err := myPathTokeniser.Tokenise(path)
	if err != nil {
		return nil, err
	}
	var Operations []*Operation
	Operations, err = myPathPostfixer.ConvertToPostfix(tokens)
	if err != nil {
		return nil, err
	}
	return p.CreatePathTree(Operations)
}

func (p *pathTreeCreator) CreatePathTree(postFixPath []*Operation) (*PathTreeNode, error) {
	var stack = make([]*PathTreeNode, 0)

	if len(postFixPath) == 0 {
		return nil, nil
	}

	for _, Operation := range postFixPath {
		var newNode = PathTreeNode{Operation: Operation}
		log.Debugf("pathTree %v ", Operation.toString())
		if Operation.OperationType.NumArgs > 0 {
			numArgs := Operation.OperationType.NumArgs
			if numArgs == 1 {
				remaining, rhs := stack[:len(stack)-1], stack[len(stack)-1]
				newNode.Rhs = rhs
				stack = remaining
			} else if numArgs == 2 {
				remaining, lhs, rhs := stack[:len(stack)-2], stack[len(stack)-2], stack[len(stack)-1]
				newNode.Lhs = lhs
				newNode.Rhs = rhs
				stack = remaining
			}
		}
		stack = append(stack, &newNode)
	}
	if len(stack) != 1 {
		return nil, fmt.Errorf("expected stack to have 1 thing but its %v", stack)
	}
	return stack[0], nil
}
