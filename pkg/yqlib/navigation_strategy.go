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

type NavigationStrategy interface {
	FollowAlias(nodeContext NodeContext) bool
	AutoCreateMap(nodeContext NodeContext) bool
	Visit(nodeContext NodeContext) error
	// node key is the string value of the last element in the path stack
	// we use it to match against the pathExpression in head.
	ShouldTraverse(nodeContext NodeContext, nodeKey string) bool
	GetVisitedNodes() []*NodeContext
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
	if ns.shouldVisit(nodeContext) {
		ns.visitedNodes = append(ns.visitedNodes, &nodeContext)
		log.Debug("adding to visited nodes, %v", nodeContext.Head)
		return ns.visit(nodeContext)
	}
	return nil
}

func (ns *NavigationStrategyImpl) alreadyVisited(pathStack []interface{}) bool {
	log.Debug("looking for pathStack")
	for _, val := range pathStack {
		log.Debug("\t %v", val)
	}
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
	log.Debug("checking against path")
	for _, val := range path1 {
		log.Debug("\t %v", val)
	}

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
