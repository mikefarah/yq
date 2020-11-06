package yqlib

import (
	"container/list"
	"strings"

	"gopkg.in/yaml.v3"
)

type CommentOpPreferences struct {
	LineComment bool
	HeadComment bool
	FootComment bool
}

func AssignCommentsOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {

	log.Debugf("AssignComments operator!")

	rhs, err := d.GetMatchingNodes(matchingNodes, pathNode.Rhs)
	if err != nil {
		return nil, err
	}
	comment := ""
	if rhs.Front() != nil {
		comment = rhs.Front().Value.(*CandidateNode).Node.Value
	}

	lhs, err := d.GetMatchingNodes(matchingNodes, pathNode.Lhs)

	if err != nil {
		return nil, err
	}

	preferences := pathNode.Operation.Preferences.(*CommentOpPreferences)

	for el := lhs.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
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

func GetCommentsOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	preferences := pathNode.Operation.Preferences.(*CommentOpPreferences)
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
		lengthCand := &CandidateNode{Node: node, Document: candidate.Document, Path: candidate.Path}
		results.PushBack(lengthCand)
	}
	return results, nil
}
