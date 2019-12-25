package yqlib

import (
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
	ShouldVisit(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool
	Visit(node *yaml.Node, head string, tail []string, pathStack []interface{}) error
	GetVisitedNodes() []*VisitedNode
}

type NavigationSettingsImpl struct {
	followAlias   func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool
	autoCreateMap func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool
	shouldVisit   func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool
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

func (ns *NavigationSettingsImpl) Visit(node *yaml.Node, head string, tail []string, pathStack []interface{}) error {
	ns.visitedNodes = append(ns.visitedNodes, &VisitedNode{node, head, tail, pathStack})
	log.Debug("adding to visited nodes")
	return ns.visit(node, head, tail, pathStack)
}

func (ns *NavigationSettingsImpl) ShouldVisit(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
	if !ns.alreadyVisited(node) {
		return ns.shouldVisit(node, head, tail, pathStack)
	} else {
		log.Debug("Skipping over %v as we have seen it already", node.Value)
	}
	return false
}

func (ns *NavigationSettingsImpl) alreadyVisited(node *yaml.Node) bool {
	for _, candidate := range ns.visitedNodes {
		if candidate.Node.Value == node.Value {
			return true
		}
	}
	return false
}
