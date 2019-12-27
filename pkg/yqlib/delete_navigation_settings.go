package yqlib

import (
	"strconv"

	yaml "gopkg.in/yaml.v3"
)

func DeleteNavigationSettings(lastBit string) NavigationSettings {
	parser := NewPathParser()
	return &NavigationSettingsImpl{
		visitedNodes: []*VisitedNode{},
		followAlias: func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
			return false
		},
		autoCreateMap: func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
			return true
		},
		visit: func(node *yaml.Node, head string, tail []string, pathStack []interface{}) error {
			log.Debug("need to find and delete %v in here", lastBit)
			DebugNode(node)
			if node.Kind == yaml.SequenceNode {
				newContent, errorDeleting := deleteFromArray(node.Content, lastBit)
				if errorDeleting != nil {
					return errorDeleting
				}
				node.Content = newContent
			} else if node.Kind == yaml.MappingNode {
				node.Content = deleteFromMap(parser, node.Content, pathStack, lastBit)
			}
			return nil
		},
	}
}
func deleteFromMap(pathParser PathParser, contents []*yaml.Node, pathStack []interface{}, lastBit string) []*yaml.Node {
	newContents := make([]*yaml.Node, 0)
	for index := 0; index < len(contents); index = index + 2 {
		keyNode := contents[index]
		valueNode := contents[index+1]
		if pathParser.MatchesNextPathElement(keyNode, lastBit, []string{}, pathStack, keyNode.Value) == false {
			log.Debug("adding node %v", keyNode.Value)
			newContents = append(newContents, keyNode, valueNode)
		} else {
			log.Debug("skipping node %v", keyNode.Value)
		}
	}
	return newContents
}

func deleteFromArray(content []*yaml.Node, lastBit string) ([]*yaml.Node, error) {

	if lastBit == "*" {
		return make([]*yaml.Node, 0), nil
	}

	var index, err = strconv.ParseInt(lastBit, 10, 64) // nolint
	if err != nil {
		return content, err
	}
	if index >= int64(len(content)) {
		log.Debug("index %v is greater than content length %v", index, len(content))
		return content, nil
	}
	return append(content[:index], content[index+1:]...), nil
}
