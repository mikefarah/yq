package yqlib

import (
	"strconv"

	yaml "gopkg.in/yaml.v3"
)

func DeleteNavigationStrategy(pathElementToDelete string) NavigationStrategy {
	parser := NewPathParser()
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
			log.Debug("need to find and delete %v in here", pathElementToDelete)
			DebugNode(node)
			if node.Kind == yaml.SequenceNode {
				newContent, errorDeleting := deleteFromArray(node.Content, pathElementToDelete)
				if errorDeleting != nil {
					return errorDeleting
				}
				node.Content = newContent
			} else if node.Kind == yaml.MappingNode {
				node.Content = deleteFromMap(parser, node.Content, nodeContext.PathStack, pathElementToDelete)
			}
			return nil
		},
	}
}
func deleteFromMap(pathParser PathParser, contents []*yaml.Node, pathStack []interface{}, pathElementToDelete string) []*yaml.Node {
	newContents := make([]*yaml.Node, 0)
	for index := 0; index < len(contents); index = index + 2 {
		keyNode := contents[index]
		valueNode := contents[index+1]
		if !pathParser.MatchesNextPathElement(NewNodeContext(keyNode, pathElementToDelete, []string{}, pathStack), keyNode.Value) {
			log.Debug("adding node %v", keyNode.Value)
			newContents = append(newContents, keyNode, valueNode)
		} else {
			log.Debug("skipping node %v", keyNode.Value)
		}
	}
	return newContents
}

func deleteFromArray(content []*yaml.Node, pathElementToDelete string) ([]*yaml.Node, error) {

	if pathElementToDelete == "*" {
		return make([]*yaml.Node, 0), nil
	}

	var index, err = strconv.ParseInt(pathElementToDelete, 10, 64) // nolint
	if err != nil {
		return content, err
	}
	if index >= int64(len(content)) {
		log.Debug("index %v is greater than content length %v", index, len(content))
		return content, nil
	}
	return append(content[:index], content[index+1:]...), nil
}
