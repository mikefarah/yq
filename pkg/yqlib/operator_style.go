package yqlib

import (
	"container/list"
	"fmt"
)

func parseStyle(customStyle string) (Style, error) {
	if customStyle == "tagged" {
		return TaggedStyle, nil
	} else if customStyle == "double" {
		return DoubleQuotedStyle, nil
	} else if customStyle == "single" {
		return SingleQuotedStyle, nil
	} else if customStyle == "literal" {
		return LiteralStyle, nil
	} else if customStyle == "folded" {
		return FoldedStyle, nil
	} else if customStyle == "flow" {
		return FlowStyle, nil
	} else if customStyle != "" {
		return 0, fmt.Errorf("unknown style %v", customStyle)
	}
	return 0, nil
}

func assignStyleOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debugf("AssignStyleOperator: %v")
	var style Style
	if !expressionNode.Operation.UpdateAssign {
		rhs, err := d.GetMatchingNodes(context.ReadOnlyClone(), expressionNode.RHS)
		if err != nil {
			return Context{}, err
		}

		if rhs.MatchingNodes.Front() != nil {
			style, err = parseStyle(rhs.MatchingNodes.Front().Value.(*CandidateNode).Value)
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
		log.Debugf("Setting style of : %v", NodeToString(candidate))
		if expressionNode.Operation.UpdateAssign {
			rhs, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(candidate), expressionNode.RHS)
			if err != nil {
				return Context{}, err
			}

			if rhs.MatchingNodes.Front() != nil {
				style, err = parseStyle(rhs.MatchingNodes.Front().Value.(*CandidateNode).Value)
				if err != nil {
					return Context{}, err
				}
			}
		}

		candidate.Style = style
	}

	return context, nil
}

func getStyleOperator(_ *dataTreeNavigator, context Context, _ *ExpressionNode) (Context, error) {
	log.Debugf("GetStyleOperator")

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		var style string
		switch candidate.Style {
		case TaggedStyle:
			style = "tagged"
		case DoubleQuotedStyle:
			style = "double"
		case SingleQuotedStyle:
			style = "single"
		case LiteralStyle:
			style = "literal"
		case FoldedStyle:
			style = "folded"
		case FlowStyle:
			style = "flow"
		case 0:
			style = ""
		default:
			style = "<unknown>"
		}
		result := candidate.CreateReplacement(ScalarNode, "!!str", style)
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}
