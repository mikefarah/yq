package yqlib

import (
	"container/list"
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

func parseStyle(customStyle string) (yaml.Style, error) {
	if customStyle == "tagged" {
		return yaml.TaggedStyle, nil
	} else if customStyle == "double" {
		return yaml.DoubleQuotedStyle, nil
	} else if customStyle == "single" {
		return yaml.SingleQuotedStyle, nil
	} else if customStyle == "literal" {
		return yaml.LiteralStyle, nil
	} else if customStyle == "folded" {
		return yaml.FoldedStyle, nil
	} else if customStyle == "flow" {
		return yaml.FlowStyle, nil
	} else if customStyle != "" {
		return 0, fmt.Errorf("Unknown style %v", customStyle)
	}
	return 0, nil
}

func assignStyleOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debugf("AssignStyleOperator: %v")
	var style yaml.Style
	if !expressionNode.Operation.UpdateAssign {
		rhs, err := d.GetMatchingNodes(context.ReadOnlyClone(), expressionNode.RHS)
		if err != nil {
			return Context{}, err
		}

		if rhs.MatchingNodes.Front() != nil {
			style, err = parseStyle(rhs.MatchingNodes.Front().Value.(*CandidateNode).Node.Value)
			if err != nil {
				return Context{}, err
			}
		}
	}

	lhs, err := d.GetMatchingNodes(context, expressionNode.LHS)

	if err != nil {
		return Context{}, err
	}

	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		log.Debugf("Setting style of : %v", candidate.GetKey())
		if expressionNode.Operation.UpdateAssign {
			rhs, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(candidate), expressionNode.RHS)
			if err != nil {
				return Context{}, err
			}

			if rhs.MatchingNodes.Front() != nil {
				style, err = parseStyle(rhs.MatchingNodes.Front().Value.(*CandidateNode).Node.Value)
				if err != nil {
					return Context{}, err
				}
			}
		}

		candidate.Node.Style = style
	}

	return context, nil
}

func getStyleOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("GetStyleOperator")

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		var style string
		switch candidate.Node.Style {
		case yaml.TaggedStyle:
			style = "tagged"
		case yaml.DoubleQuotedStyle:
			style = "double"
		case yaml.SingleQuotedStyle:
			style = "single"
		case yaml.LiteralStyle:
			style = "literal"
		case yaml.FoldedStyle:
			style = "folded"
		case yaml.FlowStyle:
			style = "flow"
		case 0:
			style = ""
		default:
			style = "<unknown>"
		}
		node := &yaml.Node{Kind: yaml.ScalarNode, Value: style, Tag: "!!str"}
		result := candidate.CreateReplacement(node)
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}
