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

func assignCommentsOperator(d *dataTreeNavigator, matchingNodes *list.List, expressionNode *ExpressionNode) (*list.List, error) {

	log.Debugf("AssignComments operator!")

	lhs, err := d.GetMatchingNodes(matchingNodes, expressionNode.Lhs)

	if err != nil {
		return nil, err
	}

	preferences := expressionNode.Operation.Preferences.(commentOpPreferences)

	comment := ""
	if !expressionNode.Operation.UpdateAssign {
		rhs, err := d.GetMatchingNodes(matchingNodes, expressionNode.Rhs)
		if err != nil {
			return nil, err
		}

		if rhs.Front() != nil {
			comment = rhs.Front().Value.(*CandidateNode).Node.Value
		}
	}

	for el := lhs.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		if expressionNode.Operation.UpdateAssign {
			rhs, err := d.GetMatchingNodes(nodeToMap(candidate), expressionNode.Rhs)
			if err != nil {
				return nil, err
			}

			if rhs.Front() != nil {
				comment = rhs.Front().Value.(*CandidateNode).Node.Value
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
	return matchingNodes, nil
}

func getCommentsOperator(d *dataTreeNavigator, matchingNodes *list.List, expressionNode *ExpressionNode) (*list.List, error) {
	preferences := expressionNode.Operation.Preferences.(commentOpPreferences)
	log.Debugf("GetComments operator!")
	var results = list.New()

	for el := matchingNodes.Front(); el != nil; el = el.Next() {
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
	return results, nil
}
