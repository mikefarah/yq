package yqlib

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	parse "github.com/a8m/envsubst/parse"
)

type envOpPreferences struct {
	StringValue bool
	NoUnset     bool
	NoEmpty     bool
	FailFast    bool
}

func envOperator(_ *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	if ConfiguredSecurityPreferences.DisableEnvOps {
		return Context{}, fmt.Errorf("env operations have been disabled")
	}
	envName := expressionNode.Operation.CandidateNode.Value
	log.Debug("EnvOperator, env name:", envName)

	rawValue := os.Getenv(envName)

	preferences := expressionNode.Operation.Preferences.(envOpPreferences)

	var node *CandidateNode
	if preferences.StringValue {
		node = &CandidateNode{
			Kind:  ScalarNode,
			Tag:   "!!str",
			Value: rawValue,
		}
	} else if rawValue == "" {
		return Context{}, fmt.Errorf("value for env variable '%v' not provided in env()", envName)
	} else {
		decoder := NewYamlDecoder(ConfiguredYamlPreferences)
		if err := decoder.Init(strings.NewReader(rawValue)); err != nil {
			return Context{}, err
		}
		var err error
		node, err = decoder.Decode()

		if err != nil {
			return Context{}, err
		}

	}
	log.Debug("ENV tag", node.Tag)
	log.Debug("ENV value", node.Value)
	log.Debug("ENV Kind", node.Kind)

	return context.SingleChildContext(node), nil
}

func envsubstOperator(_ *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	if ConfiguredSecurityPreferences.DisableEnvOps {
		return Context{}, fmt.Errorf("env operations have been disabled")
	}
	var results = list.New()
	preferences := envOpPreferences{}
	if expressionNode.Operation.Preferences != nil {
		preferences = expressionNode.Operation.Preferences.(envOpPreferences)
	}

	parser := parse.New("string", os.Environ(),
		&parse.Restrictions{NoUnset: preferences.NoUnset, NoEmpty: preferences.NoEmpty})

	if preferences.FailFast {
		parser.Mode = parse.Quick
	} else {
		parser.Mode = parse.AllErrors
	}

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		node := el.Value.(*CandidateNode)
		if node.Tag != "!!str" {
			log.Warning("EnvSubstOperator, env name:", node.Tag, node.Value)
			return Context{}, fmt.Errorf("cannot substitute with %v, can only substitute strings. Hint: Most often you'll want to use '|=' over '=' for this operation", node.Tag)
		}

		value, err := parser.Parse(node.Value)
		if err != nil {
			return Context{}, err
		}
		result := node.CreateReplacement(ScalarNode, "!!str", value)
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}
