package yqlib

import "container/list"

type AssignCommentPreferences struct {
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

	preferences := pathNode.Operation.Preferences.(*AssignCommentPreferences)

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
