package yqlib

import (
	yaml "gopkg.in/yaml.v3"
)

func UpdateNavigationSettings(changesToApply *yaml.Node) NavigationSettings {
	return &NavigationSettingsImpl{
		visitedNodes: []*VisitedNode{},
		followAlias: func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
			return false
		},
		autoCreateMap: func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
			return true
		},
		visit: func(node *yaml.Node, head string, tail []string, pathStack []interface{}) error {
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
