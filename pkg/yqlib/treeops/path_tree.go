package treeops

import "fmt"

var myPathTokeniser = NewPathTokeniser()
var myPathPostfixer = NewPathPostFixer()

type PathTreeNode struct {
	PathElement *PathElement
	Lhs         *PathTreeNode
	Rhs         *PathTreeNode
}

type PathTreeCreator interface {
	ParsePath(path string) (*PathTreeNode, error)
	CreatePathTree(postFixPath []*PathElement) (*PathTreeNode, error)
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
	var pathElements []*PathElement
	pathElements, err = myPathPostfixer.ConvertToPostfix(tokens)
	if err != nil {
		return nil, err
	}
	return p.CreatePathTree(pathElements)
}

func (p *pathTreeCreator) CreatePathTree(postFixPath []*PathElement) (*PathTreeNode, error) {
	var stack = make([]*PathTreeNode, 0)

	for _, pathElement := range postFixPath {
		var newNode = PathTreeNode{PathElement: pathElement}
		if pathElement.PathElementType == Operation {
			remaining, lhs, rhs := stack[:len(stack)-2], stack[len(stack)-2], stack[len(stack)-1]
			newNode.Lhs = lhs
			newNode.Rhs = rhs
			stack = remaining
		}
		stack = append(stack, &newNode)
	}
	if len(stack) != 1 {
		return nil, fmt.Errorf("expected stack to have 1 thing but its %v", stack)
	}
	return stack[0], nil
}
