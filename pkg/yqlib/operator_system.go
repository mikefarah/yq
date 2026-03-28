package yqlib

import (
	"bytes"
	"container/list"
	"fmt"
	"os/exec"
	"strings"
)

func resolveSystemArgs(argsNode *CandidateNode) []string {
	if argsNode.Kind == SequenceNode {
		args := make([]string, 0, len(argsNode.Content))
		for _, child := range argsNode.Content {
			args = append(args, child.Value)
		}
		return args
	}
	if argsNode.Tag != "!!null" {
		return []string{argsNode.Value}
	}
	return nil
}

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

	// determine at parse time whether we have (command; args) or just (command)
	hasArgs := expressionNode.RHS.Operation.OperationType == blockOpType

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		nodeContext := context.SingleReadonlyChildContext(candidate)

		var command string
		var args []string

		if hasArgs {
			block := expressionNode.RHS
			commandNodes, err := d.GetMatchingNodes(nodeContext, block.LHS)
			if err != nil {
				return Context{}, err
			}
			if commandNodes.MatchingNodes.Front() == nil {
				return Context{}, fmt.Errorf("system operator: command expression returned no results")
			}
			command = commandNodes.MatchingNodes.Front().Value.(*CandidateNode).Value

			argsNodes, err := d.GetMatchingNodes(nodeContext, block.RHS)
			if err != nil {
				return Context{}, err
			}
			if argsNodes.MatchingNodes.Front() != nil {
				args = resolveSystemArgs(argsNodes.MatchingNodes.Front().Value.(*CandidateNode))
			}
		} else {
			commandNodes, err := d.GetMatchingNodes(nodeContext, expressionNode.RHS)
			if err != nil {
				return Context{}, err
			}
			if commandNodes.MatchingNodes.Front() == nil {
				return Context{}, fmt.Errorf("system operator: command expression returned no results")
			}
			command = commandNodes.MatchingNodes.Front().Value.(*CandidateNode).Value
		}

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

		result := string(output)
		if strings.HasSuffix(result, "\r\n") {
			result = result[:len(result)-2]
		} else if strings.HasSuffix(result, "\n") {
			result = result[:len(result)-1]
		}
		newNode := candidate.CreateReplacement(ScalarNode, "!!str", result)
		results.PushBack(newNode)
	}

	return context.ChildContext(results), nil
}
