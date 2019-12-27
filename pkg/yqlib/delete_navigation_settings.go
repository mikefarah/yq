package yqlib

import (
	"strconv"

	yaml "gopkg.in/yaml.v3"
)

func DeleteNavigationSettings(lastBit string) NavigationSettings {
	return &NavigationSettingsImpl{
		visitedNodes: []*VisitedNode{},
		followAlias: func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
			return false
		},
		autoCreateMap: func(node *yaml.Node, head string, tail []string, pathStack []interface{}) bool {
			return true
		},
		visit: func(node *yaml.Node, head string, tail []string, pathStack []interface{}) error {
			log.Debug("need to find %v in here", lastBit)
			DebugNode(node)
			if node.Kind == yaml.SequenceNode {
				newContent, errorDeleting := deleteFromArray(node.Content, lastBit)
				if errorDeleting != nil {
					return errorDeleting
				}
				node.Content = newContent
			} else if node.Kind == yaml.MappingNode {
				// need to delete in reverse - otherwise the matching indexes
				// become incorrect.
				// matchingIndices := make([]int, 0)
				// _, errorVisiting := n.visitMatchingEntries(node, lastBit, []string{}, pathStack, func(matchingNode []*yaml.Node, indexInMap int) error {
				// 	matchingIndices = append(matchingIndices, indexInMap)
				// 	log.Debug("matchingIndices %v", indexInMap)
				// 	return nil
				// })
				// log.Debug("delete matching indices now")
				// log.Debug("%v", matchingIndices)
				// if errorVisiting != nil {
				// 	return errorVisiting
				// }
				// for i := len(matchingIndices) - 1; i >= 0; i-- {
				// 	indexToDelete := matchingIndices[i]
				// 	log.Debug("deleting index %v, %v", indexToDelete, node.Content[indexToDelete].Value)
				// 	node.Content = append(node.Content[:indexToDelete], node.Content[indexToDelete+2:]...)
				// }
			}
			return nil
		},
	}
}

func deleteFromArray(content []*yaml.Node, lastBit string) ([]*yaml.Node, error) {
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
