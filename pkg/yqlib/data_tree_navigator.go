package yqlib

// import yaml "gopkg.in/yaml.v3"

// type NodeLeafContext struct {
// 	Node      *yaml.Node
// 	Head      interface{}
// 	PathStack []interface{}
// }

// func newNodeLeafContext(node *yaml.Node, head interface{}, tailpathStack []interface{}) NodeLeafContext {
// 	newPathStack := make([]interface{}, len(pathStack))
// 	copy(newPathStack, pathStack)
// 	return NodeContext{
// 		Node:      node,
// 		Head:      head,
// 		PathStack: newPathStack,
// 	}
// }

// type DataTreeNavigator interface {
// 	Traverse(value *NodeLeafContext)
// }
