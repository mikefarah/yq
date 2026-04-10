package yqlib

import (
	"bytes"
	"container/list"
	"fmt"
	"os/exec"
	"strings"
)

func resolveSystemArgs(argsNode *CandidateNode) ([]string, error) {
	if argsNode == nil {
		return nil, nil
	}

	if argsNode.Kind == SequenceNode {
		args := make([]string, 0, len(argsNode.Content))
		for _, child := range argsNode.Content {
			// Only non-null scalar children are valid arguments.
			if child == nil {
				continue
			}
			if child.Kind != ScalarNode || child.Tag == "!!null" {
				return nil, fmt.Errorf("system operator: argument must be a non-null scalar; got kind=%v tag=%v", child.Kind, child.Tag)
			}
			args = append(args, child.Value)
		}
		if len(args) == 0 {
			return nil, nil
		}
		return args, nil
	}

	// Single-argument case: only accept a non-null scalar node.
	if argsNode.Tag == "!!null" {
		return nil, nil
	}
	if argsNode.Kind != ScalarNode {
		return nil, fmt.Errorf("system operator: args must be a non-null scalar or sequence of non-null scalars; got kind=%v tag=%v", argsNode.Kind, argsNode.Tag)
	}
	return []string{argsNode.Value}, nil
}

func resolveCommandNode(commandNodes Context) (string, error) {
	if commandNodes.MatchingNodes.Front() == nil {
		return "", fmt.Errorf("system operator: command expression returned no results")
	}
	if commandNodes.MatchingNodes.Len() > 1 {
		log.Debugf("system operator: command expression returned %d results, using first", commandNodes.MatchingNodes.Len())
	}
	cmdNode := commandNodes.MatchingNodes.Front().Value.(*CandidateNode)
	if cmdNode.Kind != ScalarNode || cmdNode.guessTagFromCustomType() != "!!str" {
		return "", fmt.Errorf("system operator: command must be a string scalar")
	}
	if cmdNode.Value == "" {
		return "", fmt.Errorf("system operator: command must be a non-empty string")
	}
	return cmdNode.Value, nil
}

func systemOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	if !ConfiguredSecurityPreferences.EnableSystemOps {
		return Context{}, fmt.Errorf("system operations are disabled, use --security-enable-system-operator to enable")
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
			command, err = resolveCommandNode(commandNodes)
			if err != nil {
				return Context{}, err
			}

			argsNodes, err := d.GetMatchingNodes(nodeContext, block.RHS)
			if err != nil {
				return Context{}, err
			}
			if argsNodes.MatchingNodes.Len() > 1 {
				log.Debugf("system operator: args expression returned %d results, using first", argsNodes.MatchingNodes.Len())
			}
			if argsNodes.MatchingNodes.Front() != nil {
				args, err = resolveSystemArgs(argsNodes.MatchingNodes.Front().Value.(*CandidateNode))
				if err != nil {
					return Context{}, err
				}
			}
		} else {
			commandNodes, err := d.GetMatchingNodes(nodeContext, expressionNode.RHS)
			if err != nil {
				return Context{}, err
			}
			command, err = resolveCommandNode(commandNodes)
			if err != nil {
				return Context{}, err
			}
		}

		var stdin bytes.Buffer
		encoded, err := encodeToYamlString(candidate)
		if err != nil {
			return Context{}, err
		}
		stdin.WriteString(encoded)

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
