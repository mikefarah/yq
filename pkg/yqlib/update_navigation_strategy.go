package yqlib

import (
	"strings"

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
		shouldVisit: func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
			var prefixMatch = strings.TrimSuffix(head, "*")
			if prefixMatch != head {
				log.Debug("prefix match, %v", strings.HasPrefix(node.Value, prefixMatch))
				return strings.HasPrefix(node.Value, prefixMatch)
			}
			log.Debug("equals match, %v", node.Value == head)
			return node.Value == head
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
