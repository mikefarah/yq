package treeops

import (
	"container/list"
	"fmt"

	"gopkg.in/yaml.v3"
)

func AssignStyleOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	customStyle := pathNode.Rhs.Operation.StringValue
	log.Debugf("AssignStyleOperator: %v", customStyle)

	var style yaml.Style
	if customStyle == "tagged" {
		style = yaml.TaggedStyle
	} else if customStyle == "double" {
		style = yaml.DoubleQuotedStyle
	} else if customStyle == "single" {
		style = yaml.SingleQuotedStyle
	} else if customStyle == "literal" {
		style = yaml.LiteralStyle
	} else if customStyle == "folded" {
		style = yaml.FoldedStyle
	} else if customStyle == "flow" {
		style = yaml.FlowStyle
	} else if customStyle != "" {
		return nil, fmt.Errorf("Unknown style %v", customStyle)
	}
	lhs, err := d.GetMatchingNodes(matchingNodes, pathNode.Lhs)

	if err != nil {
		return nil, err
	}

	for el := lhs.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		log.Debugf("Setting style of : %v", candidate.GetKey())
		candidate.Node.Style = style
	}

	return matchingNodes, nil
}

func GetStyleOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("GetStyleOperator")

	var results = list.New()

	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		var style = ""
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
		lengthCand := &CandidateNode{Node: node, Document: candidate.Document, Path: candidate.Path}
		results.PushBack(lengthCand)
	}

	return results, nil
}
