package yqlib

import (
	"bytes"
	"container/list"
	"fmt"
	"os/exec"
	"strings"
)

func systemOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	if !ConfiguredSecurityPreferences.EnableSystemOps {
		log.Warning("system operator is disabled, use --enable-system-operator flag to enable")
		results := list.New()
		for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
			candidate := el.Value.(*CandidateNode)
			results.PushBack(candidate.CreateReplacement(ScalarNode, "!!null", "null"))
		}
		return context.ChildContext(results), nil
	}

	var command string
	var argsExpression *ExpressionNode

	// check if it's a block operator (command; args) or just (command)
	if expressionNode.RHS.Operation.OperationType == blockOpType {
		block := expressionNode.RHS
		commandNodes, err := d.GetMatchingNodes(context.ReadOnlyClone(), block.LHS)
		if err != nil {
			return Context{}, err
		}
		if commandNodes.MatchingNodes.Front() == nil {
			return Context{}, fmt.Errorf("system operator: command expression returned no results")
		}
		command = commandNodes.MatchingNodes.Front().Value.(*CandidateNode).Value
		argsExpression = block.RHS
	} else {
		commandNodes, err := d.GetMatchingNodes(context.ReadOnlyClone(), expressionNode.RHS)
		if err != nil {
			return Context{}, err
		}
		if commandNodes.MatchingNodes.Front() == nil {
			return Context{}, fmt.Errorf("system operator: command expression returned no results")
		}
		command = commandNodes.MatchingNodes.Front().Value.(*CandidateNode).Value
	}

	// evaluate args if present
	var args []string
	if argsExpression != nil {
		argsNodes, err := d.GetMatchingNodes(context.ReadOnlyClone(), argsExpression)
		if err != nil {
			return Context{}, err
		}
		if argsNodes.MatchingNodes.Front() != nil {
			argsNode := argsNodes.MatchingNodes.Front().Value.(*CandidateNode)
			if argsNode.Kind == SequenceNode {
				for _, child := range argsNode.Content {
					args = append(args, child.Value)
				}
			} else if argsNode.Tag != "!!null" {
				args = []string{argsNode.Value}
			}
		}
	}

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		var stdin bytes.Buffer
		if candidate.Tag != "!!null" {
			encoded, err := encodeToYamlString(candidate)
			if err != nil {
				return Context{}, err
			}
			stdin.WriteString(encoded)
		}

		// #nosec G204 - intentional: user must explicitly enable this operator
		cmd := exec.Command(command, args...)
		cmd.Stdin = &stdin
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		output, err := cmd.Output()
		if err != nil {
			stderrStr := strings.TrimSpace(stderr.String())
			if stderrStr != "" {
				return Context{}, fmt.Errorf("system command '%v' failed: %w\nstderr: %v", command, err, stderrStr)
			}
			return Context{}, fmt.Errorf("system command '%v' failed: %w", command, err)
		}

		result := strings.TrimRight(string(output), "\n")
		newNode := candidate.CreateReplacement(ScalarNode, "!!str", result)
		results.PushBack(newNode)
	}

	return context.ChildContext(results), nil
}
