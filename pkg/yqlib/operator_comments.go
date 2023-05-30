package yqlib

import (
	"bufio"
	"bytes"
	"container/list"
	"regexp"
)

type commentOpPreferences struct {
	LineComment bool
	HeadComment bool
	FootComment bool
}

func assignCommentsOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debugf("AssignComments operator!")

	lhs, err := d.GetMatchingNodes(context, expressionNode.LHS)

	if err != nil {
		return Context{}, err
	}

	preferences := expressionNode.Operation.Preferences.(commentOpPreferences)

	comment := ""
	if !expressionNode.Operation.UpdateAssign {
		rhs, err := d.GetMatchingNodes(context.ReadOnlyClone(), expressionNode.RHS)
		if err != nil {
			return Context{}, err
		}

		if rhs.MatchingNodes.Front() != nil {
			comment = rhs.MatchingNodes.Front().Value.(*CandidateNode).Value
		}
	}

	log.Debugf("AssignComments comment is %v", comment)

	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		log.Debugf("AssignComments lhs %v", NodeToString(candidate))

		if expressionNode.Operation.UpdateAssign {
			rhs, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(candidate), expressionNode.RHS)
			if err != nil {
				return Context{}, err
			}

			if rhs.MatchingNodes.Front() != nil {
				comment = rhs.MatchingNodes.Front().Value.(*CandidateNode).Value
			}
		}

		log.Debugf("Setting comment of : %v", candidate.GetKey())
		if preferences.LineComment {
			log.Debugf("Setting line comment of : %v to %v", candidate.GetKey(), comment)
			candidate.LineComment = comment
		}
		if preferences.HeadComment {
			candidate.HeadComment = comment
			candidate.LeadingContent = "" // clobber the leading content, if there was any.
		}
		// if preferences.FootComment && candidate.Kind == DocumentNode && comment != "" {
		// 	log.Debugf("AssignComments - setting line comment to %v", comment)
		// 	candidate.TrailingContent = "# " + comment
		// } else if preferences.FootComment && candidate.Kind == DocumentNode {
		// 	log.Debugf("AssignComments - setting line comment to %v", comment)
		// 	candidate.TrailingContent = comment
		if preferences.FootComment { //&& candidate.Kind != DocumentNode {
			candidate.FootComment = comment
			candidate.TrailingContent = ""
		}

	}
	return context, nil
}

func getCommentsOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	preferences := expressionNode.Operation.Preferences.(commentOpPreferences)
	var startCommentCharaterRegExp = regexp.MustCompile(`^# `)
	var subsequentCommentCharaterRegExp = regexp.MustCompile(`\n# `)

	log.Debugf("GetComments operator!")
	var results = list.New()

	yamlPrefs := NewDefaultYamlPreferences()
	yamlPrefs.PrintDocSeparators = false
	yamlPrefs.UnwrapScalar = false

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		comment := ""
		if preferences.LineComment {
			log.Debugf("Reading line comment of : %v to %v", candidate.GetKey(), candidate.LineComment)
			comment = candidate.LineComment
		} else if preferences.HeadComment && candidate.LeadingContent != "" {
			var chompRegexp = regexp.MustCompile(`\n$`)
			var output bytes.Buffer
			var writer = bufio.NewWriter(&output)
			var encoder = NewYamlEncoder(2, false, yamlPrefs)
			if err := encoder.PrintLeadingContent(writer, candidate.LeadingContent); err != nil {
				return Context{}, err
			}
			if err := writer.Flush(); err != nil {
				return Context{}, err
			}
			comment = output.String()
			comment = chompRegexp.ReplaceAllString(comment, "")
		} else if preferences.HeadComment {
			comment = candidate.HeadComment
		} else if preferences.FootComment && candidate.Kind == DocumentNode && candidate.TrailingContent != "" {
			comment = candidate.TrailingContent
		} else if preferences.FootComment {
			comment = candidate.FootComment
		}
		comment = startCommentCharaterRegExp.ReplaceAllString(comment, "")
		comment = subsequentCommentCharaterRegExp.ReplaceAllString(comment, "\n")

		result := candidate.CreateReplacement(ScalarNode, "!!str", comment)
		if candidate.IsMapKey {
			result.IsMapKey = false
			result.Key = candidate
		}
		results.PushBack(result)
	}
	return context.ChildContext(results), nil
}
