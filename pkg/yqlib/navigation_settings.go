package yqlib

import (
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type VisitedNode struct {
	Node      *yaml.Node
	Head      string
	Tail      []string
	PathStack []interface{}
}

type NavigationSettings interface {
	FollowAlias(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool
	AutoCreateMap(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool
	Visit(node *yaml.Node, head string, tail []string, pathStack []interface{}) error
	ShouldTraverse(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool
	GetVisitedNodes() []*VisitedNode
}

type NavigationSettingsImpl struct {
	followAlias   func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool
	autoCreateMap func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool
	visit         func(node *yaml.Node, head string, tail []string, pathStack []interface{}) error
	visitedNodes  []*VisitedNode
}

func matches(node *yaml.Node, head string) bool {
	var prefixMatch = strings.TrimSuffix(head, "*")
	if prefixMatch != head {
		log.Debug("prefix match, %v", strings.HasPrefix(node.Value, prefixMatch))
		return strings.HasPrefix(node.Value, prefixMatch)
	}
	log.Debug("equals match, %v", node.Value == head)
	return node.Value == head
}

func (ns *NavigationSettingsImpl) GetVisitedNodes() []*VisitedNode {
	return ns.visitedNodes
}

func (ns *NavigationSettingsImpl) FollowAlias(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
	return ns.followAlias(node, head, tail, pathStack)
}

func (ns *NavigationSettingsImpl) AutoCreateMap(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
	return ns.autoCreateMap(node, head, tail, pathStack)
}

func (ns *NavigationSettingsImpl) matchesNextPath(path string, candidate string) bool {
	var prefixMatch = strings.TrimSuffix(path, "*")
	if prefixMatch != path {
		log.Debug("prefix match, %v", strings.HasPrefix(candidate, prefixMatch))
		return strings.HasPrefix(candidate, prefixMatch)
	}
	return candidate == path
}

func (ns *NavigationSettingsImpl) ShouldTraverse(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
	// we should traverse aliases (if enabled), but not visit them :/
	if len(pathStack) == 0 {
		return true
	}

	if ns.alreadyVisited(pathStack) {
		return false
	}

	lastBit := fmt.Sprintf("%v", pathStack[len(pathStack)-1])

	return (lastBit == "<<" && ns.FollowAlias(node, head, tail, pathStack)) || (lastBit != "<<" && ns.matchesNextPath(head, lastBit))
}

func (ns *NavigationSettingsImpl) shouldVisit(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
	// we should traverse aliases (if enabled), but not visit them :/
	if len(pathStack) == 0 {
		return true
	}

	if ns.alreadyVisited(pathStack) {
		return false
	}

	lastBit := fmt.Sprintf("%v", pathStack[len(pathStack)-1])
	// only visit aliases if its an exact match
	return (lastBit == "<<" && head == "<<") || (lastBit != "<<" && ns.matchesNextPath(head, lastBit))

}

func (ns *NavigationSettingsImpl) Visit(node *yaml.Node, head string, tail []string, pathStack []interface{}) error {
	if ns.shouldVisit(node, head, tail, pathStack) {
		ns.visitedNodes = append(ns.visitedNodes, &VisitedNode{node, head, tail, pathStack})
		log.Debug("adding to visited nodes, %v", head)
		return ns.visit(node, head, tail, pathStack)
	}
	return nil
}

func (ns *NavigationSettingsImpl) alreadyVisited(pathStack []interface{}) bool {
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
