package yqlib

import (
	yaml "gopkg.in/yaml.v3"
)

func DeleteNavigationStrategy(pathElementToDelete interface{}) NavigationStrategy {
	parser := NewPathParser()
	return &NavigationStrategyImpl{
		visitedNodes: []*NodeContext{},
		pathParser:   parser,
		followAlias: func(nodeContext NodeContext) bool {
			return false
		},
		shouldOnlyDeeplyVisitLeaves: func(nodeContext NodeContext) bool {
			return false
		},
		visit: func(nodeContext NodeContext) error {
			node := nodeContext.Node
			log.Debug("need to find and delete %v in here (%v)", pathElementToDelete, pathStackToString(nodeContext.PathStack))
			DebugNode(node)
			if node.Kind == yaml.SequenceNode {
				newContent := deleteFromArray(parser, node.Content, nodeContext.PathStack, pathElementToDelete)
				node.Content = newContent
			} else if node.Kind == yaml.MappingNode {
				node.Content = deleteFromMap(parser, node.Content, nodeContext.PathStack, pathElementToDelete)
			}
			return nil
		},
	}
}
func deleteFromMap(pathParser PathParser, contents []*yaml.Node, pathStack []interface{}, pathElementToDelete interface{}) []*yaml.Node {
	newContents := make([]*yaml.Node, 0)
	for index := 0; index < len(contents); index = index + 2 {
		keyNode := contents[index]
		valueNode := contents[index+1]
		if !pathParser.MatchesNextPathElement(NewNodeContext(keyNode, pathElementToDelete, make([]interface{}, 0), pathStack), keyNode.Value) {
			log.Debug("adding node %v", keyNode.Value)
			newContents = append(newContents, keyNode, valueNode)
		} else {
			log.Debug("skipping node %v", keyNode.Value)
		}
	}
	return newContents
}

func deleteFromArray(pathParser PathParser, content []*yaml.Node, pathStack []interface{}, pathElementToDelete interface{}) []*yaml.Node {

	switch pathElementToDelete := pathElementToDelete.(type) {
	case int64:
		return deleteIndexInArray(content, pathElementToDelete)
	default:
		log.Debug("%v is not a numeric index, finding matching patterns", pathElementToDelete)
		var newArray = make([]*yaml.Node, 0)

		for _, childValue := range content {
			if !pathParser.MatchesNextPathElement(NewNodeContext(childValue, pathElementToDelete, make([]interface{}, 0), pathStack), childValue.Value) {
				newArray = append(newArray, childValue)
			}
		}
		return newArray
	}

}

func deleteIndexInArray(content []*yaml.Node, index int64) []*yaml.Node {
	log.Debug("deleting index %v in array", index)
	if index >= int64(len(content)) {
		log.Debug("index %v is greater than content length %v", index, len(content))
		return content
	}
	return append(content[:index], content[index+1:]...)
}
