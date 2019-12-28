package yqlib

import (
	yaml "gopkg.in/yaml.v3"
)

func UpdateNavigationStrategy(changesToApply *yaml.Node) NavigationStrategy {
	return &NavigationStrategyImpl{
		visitedNodes: []*NodeContext{},
		followAlias: func(nodeContext NodeContext) bool {
			return false
		},
		autoCreateMap: func(nodeContext NodeContext) bool {
			return true
		},
		visit: func(nodeContext NodeContext) error {
			node := nodeContext.Node
			log.Debug("going to update")
			DebugNode(node)
			log.Debug("with")
			DebugNode(changesToApply)
			node.Value = changesToApply.Value
			node.Tag = changesToApply.Tag
			node.Kind = changesToApply.Kind
			node.Style = changesToApply.Style
			node.Content = changesToApply.Content
			node.HeadComment = changesToApply.HeadComment
			node.LineComment = changesToApply.LineComment
			node.FootComment = changesToApply.FootComment
			return nil
		},
	}
}
