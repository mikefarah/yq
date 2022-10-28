package yqlib

import (
	"bufio"
	"bytes"
	"container/list"
	"regexp"

	yaml "gopkg.in/yaml.v3"
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
			comment = rhs.MatchingNodes.Front().Value.(*CandidateNode).Node.Value
		}
	}

	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		if expressionNode.Operation.UpdateAssign {
			rhs, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(candidate), expressionNode.RHS)
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
			candidate.LeadingContent = "" // clobber the leading content, if there was any.
		}
		if preferences.FootComment && candidate.Node.Kind == yaml.DocumentNode && comment != "" {
			candidate.TrailingContent = "# " + comment
		} else if preferences.FootComment && candidate.Node.Kind == yaml.DocumentNode {
			candidate.TrailingContent = comment

		} else if preferences.FootComment && candidate.Node.Kind != yaml.DocumentNode {
			candidate.Node.FootComment = comment
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
			comment = candidate.Node.LineComment
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
			comment = candidate.Node.HeadComment
		} else if preferences.FootComment && candidate.Node.Kind == yaml.DocumentNode && candidate.TrailingContent != "" {
			comment = candidate.TrailingContent
		} else if preferences.FootComment {
			comment = candidate.Node.FootComment
		}
		comment = startCommentCharaterRegExp.ReplaceAllString(comment, "")
		comment = subsequentCommentCharaterRegExp.ReplaceAllString(comment, "\n")

		node := &yaml.Node{Kind: yaml.ScalarNode, Value: comment, Tag: "!!str"}
		result := candidate.CreateReplacement(node)
		result.LeadingContent = "" // don't include the leading yaml content when retrieving a comment
		results.PushBack(result)
	}
	return context.ChildContext(results), nil
}
