package yqlib

import (
	"container/list"
	"os"

	yaml "gopkg.in/yaml.v3"
)

// type EnvOpPreferences struct {
// 	StringValue bool
// }

func EnvOperator(d *dataTreeNavigator, matchMap *list.List, pathNode *PathTreeNode) (*list.List, error) {
	envName := pathNode.Operation.CandidateNode.Node.Value
	log.Debug("EnvOperator, env name:", envName)

	rawValue := os.Getenv(envName)

	target := &CandidateNode{
		Path:     make([]interface{}, 0),
		Document: 0,
		Filename: "",
		Node: &yaml.Node{
			Kind:  yaml.ScalarNode,
			Tag:   "!!str",
			Value: rawValue,
		},
	}

	return nodeToMap(target), nil
}
