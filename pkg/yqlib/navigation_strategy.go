package yqlib

import (
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

type NodeContext struct {
	Node      *yaml.Node
	Head      interface{}
	Tail      []interface{}
	PathStack []interface{}
	// middle nodes are nodes that match along the original path, but not a
	// target match of the path. This is only relevant when ShouldOnlyDeeplyVisitLeaves is false.
	IsMiddleNode bool
}

func NewNodeContext(node *yaml.Node, head interface{}, tail []interface{}, pathStack []interface{}) NodeContext {
	newTail := make([]interface{}, len(tail))
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
	ShouldDeeplyTraverse(nodeContext NodeContext) bool
	// when deeply traversing, should we visit all matching nodes, or just leaves?
	ShouldOnlyDeeplyVisitLeaves(NodeContext) bool
	GetVisitedNodes() []*NodeContext
	DebugVisitedNodes()
	GetPathParser() PathParser
}

type NavigationStrategyImpl struct {
	followAlias                 func(nodeContext NodeContext) bool
	autoCreateMap               func(nodeContext NodeContext) bool
	visit                       func(nodeContext NodeContext) error
	shouldVisitExtraFn          func(nodeContext NodeContext) bool
	shouldDeeplyTraverse        func(nodeContext NodeContext) bool
	shouldOnlyDeeplyVisitLeaves func(nodeContext NodeContext) bool
	visitedNodes                []*NodeContext
	pathParser                  PathParser
}

func (ns *NavigationStrategyImpl) GetPathParser() PathParser {
	return ns.pathParser
}

func (ns *NavigationStrategyImpl) GetVisitedNodes() []*NodeContext {
	return ns.visitedNodes
}

func (ns *NavigationStrategyImpl) FollowAlias(nodeContext NodeContext) bool {
	if ns.followAlias != nil {
		return ns.followAlias(nodeContext)
	}
	return true
}

func (ns *NavigationStrategyImpl) AutoCreateMap(nodeContext NodeContext) bool {
	if ns.autoCreateMap != nil {
		return ns.autoCreateMap(nodeContext)
	}
	return false
}

func (ns *NavigationStrategyImpl) ShouldDeeplyTraverse(nodeContext NodeContext) bool {
	if ns.shouldDeeplyTraverse != nil {
		return ns.shouldDeeplyTraverse(nodeContext)
	}
	return true
}

func (ns *NavigationStrategyImpl) ShouldOnlyDeeplyVisitLeaves(nodeContext NodeContext) bool {
	if ns.shouldOnlyDeeplyVisitLeaves != nil {
		return ns.shouldOnlyDeeplyVisitLeaves(nodeContext)
	}
	return true

}

func (ns *NavigationStrategyImpl) ShouldTraverse(nodeContext NodeContext, nodeKey string) bool {
	// we should traverse aliases (if enabled), but not visit them :/
	if len(nodeContext.PathStack) == 0 {
		return true
	}

	if ns.alreadyVisited(nodeContext.PathStack) {
		return false
	}

	return (nodeKey == "<<" && ns.FollowAlias(nodeContext)) || (nodeKey != "<<" &&
		ns.pathParser.MatchesNextPathElement(nodeContext, nodeKey))
}

func (ns *NavigationStrategyImpl) shouldVisit(nodeContext NodeContext) bool {
	pathStack := nodeContext.PathStack
	if len(pathStack) == 0 {
		return true
	}
	log.Debug("tail len %v", len(nodeContext.Tail))

	if ns.alreadyVisited(pathStack) || len(nodeContext.Tail) != 0 {
		return false
	}

	nodeKey := fmt.Sprintf("%v", pathStack[len(pathStack)-1])
	log.Debug("nodeKey: %v, nodeContext.Head: %v", nodeKey, nodeContext.Head)

	// only visit aliases if its an exact match
	return ((nodeKey == "<<" && nodeContext.Head == "<<") || (nodeKey != "<<" &&
		ns.pathParser.MatchesNextPathElement(nodeContext, nodeKey))) && (ns.shouldVisitExtraFn == nil || ns.shouldVisitExtraFn(nodeContext))
}

func (ns *NavigationStrategyImpl) Visit(nodeContext NodeContext) error {
	log.Debug("Visit?, %v, %v", nodeContext.Head, pathStackToString(nodeContext.PathStack))
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
	log.Debug("Visited Nodes:")
	for _, candidate := range ns.visitedNodes {
		log.Debug(" - %v", pathStackToString(candidate.PathStack))
	}
}

func (ns *NavigationStrategyImpl) alreadyVisited(pathStack []interface{}) bool {
	log.Debug("checking already visited pathStack: %v", pathStackToString(pathStack))
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
	log.Debug("checking against path: %v", pathStackToString(path1))

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
