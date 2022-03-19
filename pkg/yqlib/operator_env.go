package yqlib

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	envsubst "github.com/a8m/envsubst"
	yaml "gopkg.in/yaml.v3"
)

type envOpPreferences struct {
	StringValue bool
	NoUnset     bool
	NoEmpty     bool
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
		//first node is a doc
		node = unwrapDoc(&dataBucket)
	}
	log.Debug("ENV tag", node.Tag)
	log.Debug("ENV value", node.Value)
	log.Debug("ENV Kind", node.Kind)

	target := &CandidateNode{Node: node}

	return context.SingleChildContext(target), nil
}

func envsubstOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	var results = list.New()
	preferences := expressionNode.Operation.Preferences.(envOpPreferences)

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := unwrapDoc(candidate.Node)
		if node.Tag != "!!str" {
			log.Warning("EnvSubstOperator, env name:", node.Tag, node.Value)
			return Context{}, fmt.Errorf("cannot substitute with %v, can only substitute strings. Hint: Most often you'll want to use '|=' over '=' for this operation", node.Tag)
		}

		value, err := envsubst.StringRestricted(node.Value, preferences.NoUnset, preferences.NoEmpty)
		if err != nil {
			return Context{}, err
		}
		targetNode := &yaml.Node{Kind: yaml.ScalarNode, Value: value, Tag: "!!str"}
		result := candidate.CreateReplacement(targetNode)
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}
