package yqlib

import (
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

type NodeContext struct {
	Node      *yaml.Node
	Head      string
	Tail      []string
	PathStack []interface{}
}

func NewNodeContext(node *yaml.Node, head string, tail []string, pathStack []interface{}) NodeContext {
	newTail := make([]string, len(tail))
	copy(newTail, tail)

	newPathStack := make([]interface{}, len(pathStack))
	copy(newPathStack, pathStack)
	return NodeContext{
		Node:      node,
		Head:      head,
		Tail:      newTail,
		PathStack: newPathStack,
	}
}

type NavigationStrategy interface {
	FollowAlias(nodeContext NodeContext) bool
	AutoCreateMap(nodeContext NodeContext) bool
	Visit(nodeContext NodeContext) error
	// node key is the string value of the last element in the path stack
	// we use it to match against the pathExpression in head.
	ShouldTraverse(nodeContext NodeContext, nodeKey string) bool
	GetVisitedNodes() []*NodeContext
	DebugVisitedNodes()
}

type NavigationStrategyImpl struct {
	followAlias   func(nodeContext NodeContext) bool
	autoCreateMap func(nodeContext NodeContext) bool
	visit         func(nodeContext NodeContext) error
	visitedNodes  []*NodeContext
}

func (ns *NavigationStrategyImpl) GetVisitedNodes() []*NodeContext {
	return ns.visitedNodes
}

func (ns *NavigationStrategyImpl) FollowAlias(nodeContext NodeContext) bool {
	return ns.followAlias(nodeContext)
}

func (ns *NavigationStrategyImpl) AutoCreateMap(nodeContext NodeContext) bool {
	return ns.autoCreateMap(nodeContext)
}

func (ns *NavigationStrategyImpl) ShouldTraverse(nodeContext NodeContext, nodeKey string) bool {
	// we should traverse aliases (if enabled), but not visit them :/
	if len(nodeContext.PathStack) == 0 {
		return true
	}

	if ns.alreadyVisited(nodeContext.PathStack) {
		return false
	}

	parser := NewPathParser()

	return (nodeKey == "<<" && ns.FollowAlias(nodeContext)) || (nodeKey != "<<" &&
		parser.MatchesNextPathElement(nodeContext, nodeKey))
}

func (ns *NavigationStrategyImpl) shouldVisit(nodeContext NodeContext) bool {
	// we should traverse aliases (if enabled), but not visit them :/
	pathStack := nodeContext.PathStack
	if len(pathStack) == 0 {
		return true
	}

	if ns.alreadyVisited(pathStack) {
		return false
	}

	nodeKey := fmt.Sprintf("%v", pathStack[len(pathStack)-1])
	parser := NewPathParser()

	// only visit aliases if its an exact match
	return (nodeKey == "<<" && nodeContext.Head == "<<") || (nodeKey != "<<" &&
		parser.MatchesNextPathElement(nodeContext, nodeKey))
}

func (ns *NavigationStrategyImpl) Visit(nodeContext NodeContext) error {
	log.Debug("Visit?, %v, %v", nodeContext.Head, PathStackToString(nodeContext.PathStack))
	DebugNode(nodeContext.Node)
	if ns.shouldVisit(nodeContext) {
		log.Debug("yep, visiting")
		//  pathStack array must be
		// copied, as append() may sometimes reuse and modify the array
		ns.visitedNodes = append(ns.visitedNodes, &nodeContext)
		ns.DebugVisitedNodes()
		return ns.visit(nodeContext)
	}
	log.Debug("nope, skip it")
	return nil
}

func (ns *NavigationStrategyImpl) DebugVisitedNodes() {
	log.Debug("%v", ns.visitedNodes)
	for _, candidate := range ns.visitedNodes {
		log.Debug(" - %v", PathStackToString(candidate.PathStack))
	}
}

func (ns *NavigationStrategyImpl) alreadyVisited(pathStack []interface{}) bool {
	log.Debug("checking already visited pathStack: %v", PathStackToString(pathStack))
	for _, candidate := range ns.visitedNodes {
		candidatePathStack := candidate.PathStack
		if patchStacksMatch(candidatePathStack, pathStack) {
			log.Debug("paths match, already seen it")
			return true
		}

	}
	log.Debug("never seen it before!")
	return false
}

func patchStacksMatch(path1 []interface{}, path2 []interface{}) bool {
	log.Debug("checking against path: %v", PathStackToString(path1))

	if len(path1) != len(path2) {
		return false
	}
	for index, p1Value := range path1 {

		p2Value := path2[index]
		if p1Value != p2Value {
			return false
		}
	}
	return true

}
