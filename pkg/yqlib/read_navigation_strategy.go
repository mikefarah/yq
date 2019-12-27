package yqlib

import (
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
		visit: func(node *yaml.Node, head string, tail []string, pathStack []interface{}) error {
			return nil
		},
	}
}
