package yqlib

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type envOpPreferences struct {
	StringValue bool
}

func envOperator(d *dataTreeNavigator, matchMap *list.List, pathNode *PathTreeNode) (*list.List, error) {
	envName := pathNode.Operation.CandidateNode.Node.Value
	log.Debug("EnvOperator, env name:", envName)

	rawValue := os.Getenv(envName)

	preferences := pathNode.Operation.Preferences.(*envOpPreferences)

	var node *yaml.Node
	if preferences.StringValue {
		node = &yaml.Node{
			Kind:  yaml.ScalarNode,
			Tag:   "!!str",
			Value: rawValue,
		}
	} else if rawValue == "" {
		return nil, fmt.Errorf("Value for env variable '%v' not provided in env()", envName)
	} else {
		var dataBucket yaml.Node
		decoder := yaml.NewDecoder(strings.NewReader(rawValue))
		errorReading := decoder.Decode(&dataBucket)
		if errorReading != nil {
			return nil, errorReading
		}
		//first node is a doc
		node = UnwrapDoc(&dataBucket)
	}
	log.Debug("ENV tag", node.Tag)
	log.Debug("ENV value", node.Value)
	log.Debug("ENV Kind", node.Kind)

	target := &CandidateNode{
		Path:     make([]interface{}, 0),
		Document: 0,
		Filename: "",
		Node:     node,
	}

	return nodeToMap(target), nil
}
