package yqlib

import (
	"strings"

	yaml "gopkg.in/yaml.v3"
)

func ReadNavigationSettings() NavigationSettings {
	return &NavigationSettingsImpl{
		visitedNodes: []*VisitedNode{},
		followAlias: func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
			return true
		},
		autoCreateMap: func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
			return false
		},
		shouldVisit: func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
			log.Debug("shouldVisit h: %v, actual: %v", head, node.Value)
			if node.Value == "<<" {
				log.Debug("its an alias, skip it")
				// dont match alias keys, as we'll follow them instead
				return false
			}
			var prefixMatch = strings.TrimSuffix(head, "*")
			if prefixMatch != head {
				log.Debug("prefix match, %v", strings.HasPrefix(node.Value, prefixMatch))
				return strings.HasPrefix(node.Value, prefixMatch)
			}
			log.Debug("equals match, %v", node.Value == head)
			return node.Value == head
		},
		visit: func(node *yaml.Node, head string, tail []string, pathStack []interface{}) error {
			return nil
		},
	}
}
