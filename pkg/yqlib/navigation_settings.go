package yqlib

import (
	"strings"

	logging "gopkg.in/op/go-logging.v1"
	yaml "gopkg.in/yaml.v3"
)

type NavigationSettings interface {
	FollowAlias(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool
	AutoCreateMap(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool
	ShouldVisit(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool
}

type NavigationSettingsImpl struct {
	followAlias   func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool
	autoCreateMap func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool
	shouldVisit   func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool
}

func (ns NavigationSettingsImpl) FollowAlias(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
	return ns.followAlias(node, head, tail, pathStack)
}

func (ns NavigationSettingsImpl) AutoCreateMap(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
	return ns.autoCreateMap(node, head, tail, pathStack)
}

func (ns NavigationSettingsImpl) ShouldVisit(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
	return ns.shouldVisit(node, head, tail, pathStack)
}

func UpdateNavigationSettings(l *logging.Logger) NavigationSettings {
	return NavigationSettingsImpl{
		followAlias: func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
			return false
		},
		autoCreateMap: func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
			return true
		},
		shouldVisit: func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
			var prefixMatch = strings.TrimSuffix(head, "*")
			if prefixMatch != head {
				l.Debug("prefix match, %v", strings.HasPrefix(node.Value, prefixMatch))
				return strings.HasPrefix(node.Value, prefixMatch)
			}
			l.Debug("equals match, %v", node.Value == head)
			return node.Value == head
		},
	}
}

func ReadNavigationSettings(l *logging.Logger) NavigationSettings {
	return NavigationSettingsImpl{
		followAlias: func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
			return true
		},
		autoCreateMap: func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
			return false
		},
		shouldVisit: func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
			l.Debug("shouldVisit h: %v, actual: %v", head, node.Value)
			if node.Value == "<<" {
				l.Debug("its an alias, skip it")
				// dont match alias keys, as we'll follow them instead
				return false
			}
			var prefixMatch = strings.TrimSuffix(head, "*")
			if prefixMatch != head {
				l.Debug("prefix match, %v", strings.HasPrefix(node.Value, prefixMatch))
				return strings.HasPrefix(node.Value, prefixMatch)
			}
			l.Debug("equals match, %v", node.Value == head)
			return node.Value == head
		},
	}
}
