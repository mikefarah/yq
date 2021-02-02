package yqlib

import (
	"container/list"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type commentOpPreferences struct {
	LineComment bool
	HeadComment bool
	FootComment bool
}

func assignCommentsOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debugf("AssignComments operator!")

	lhs, err := d.GetMatchingNodes(context, expressionNode.Lhs)

	if err != nil {
		return Context{}, err
	}

	preferences := expressionNode.Operation.Preferences.(commentOpPreferences)

	comment := ""
	if !expressionNode.Operation.UpdateAssign {
		rhs, err := d.GetMatchingNodes(context, expressionNode.Rhs)
		if err != nil {
			return Context{}, err
		}

		if rhs.MatchingNodes.Front() != nil {
			comment = rhs.MatchingNodes.Front().Value.(*CandidateNode).Node.Value
		}
	}

	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		if expressionNode.Operation.UpdateAssign {
			rhs, err := d.GetMatchingNodes(context.SingleChildContext(candidate), expressionNode.Rhs)
			if err != nil {
				return Context{}, err
			}

			if rhs.MatchingNodes.Front() != nil {
				comment = rhs.MatchingNodes.Front().Value.(*CandidateNode).Node.Value
			}
		}

		log.Debugf("Setting comment of : %v", candidate.GetKey())
		if preferences.LineComment {
			candidate.Node.LineComment = comment
		}
		if preferences.HeadComment {
			candidate.Node.HeadComment = comment
		}
		if preferences.FootComment {
			candidate.Node.FootComment = comment
		}

	}
	return context, nil
}

func getCommentsOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	preferences := expressionNode.Operation.Preferences.(commentOpPreferences)
	log.Debugf("GetComments operator!")
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		comment := ""
		if preferences.LineComment {
			comment = candidate.Node.LineComment
		} else if preferences.HeadComment {
			comment = candidate.Node.HeadComment
		} else if preferences.FootComment {
			comment = candidate.Node.FootComment
		}
		comment = strings.Replace(comment, "# ", "", 1)

		node := &yaml.Node{Kind: yaml.ScalarNode, Value: comment, Tag: "!!str"}
		result := candidate.CreateChild(nil, node)
		results.PushBack(result)
	}
	return context.ChildContext(results), nil
}
