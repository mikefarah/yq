package yqlib

import (
	"fmt"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type envOpPreferences struct {
	StringValue bool
}

func envOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	envName := expressionNode.Operation.CandidateNode.Node.Value
	log.Debug("EnvOperator, env name:", envName)

	rawValue := os.Getenv(envName)

	preferences := expressionNode.Operation.Preferences.(envOpPreferences)

	var node *yaml.Node
	if preferences.StringValue {
		node = &yaml.Node{
			Kind:  yaml.ScalarNode,
			Tag:   "!!str",
			Value: rawValue,
		}
	} else if rawValue == "" {
		return Context{}, fmt.Errorf("Value for env variable '%v' not provided in env()", envName)
	} else {
		var dataBucket yaml.Node
		decoder := yaml.NewDecoder(strings.NewReader(rawValue))
		errorReading := decoder.Decode(&dataBucket)
		if errorReading != nil {
			return Context{}, errorReading
		}
		// first node is a doc
		node = unwrapDoc(&dataBucket)
	}
	log.Debug("ENV tag", node.Tag)
	log.Debug("ENV value", node.Value)
	log.Debug("ENV Kind", node.Kind)

	target := &CandidateNode{Node: node}

	return context.SingleChildContext(target), nil
}
