package yqlib

import "gopkg.in/yaml.v3"

type ArrayMergeStrategy uint32

const (
	UpdateArrayMergeStrategy ArrayMergeStrategy = 1 << iota
	OverwriteArrayMergeStrategy
	AppendArrayMergeStrategy
)

type CommentsMergeStrategy uint32

const (
	SetWhenBlankCommentsMergeStrategy CommentsMergeStrategy = 1 << iota
	IgnoreCommentsMergeStrategy
	OverwriteCommentsMergeStrategy
	AppendCommentsMergeStrategy
)

func MergeNavigationStrategy(updateCommand UpdateCommand, autoCreate bool) NavigationStrategy {
	return &NavigationStrategyImpl{
		visitedNodes: []*NodeContext{},
		pathParser:   NewPathParser(),
		followAlias: func(nodeContext NodeContext) bool {
			return false
		},
		autoCreateMap: func(nodeContext NodeContext) bool {
			return autoCreate
		},
		visit: func(nodeContext NodeContext) error {
			node := nodeContext.Node
			changesToApply := updateCommand.Value

			if node.Kind == yaml.DocumentNode && changesToApply.Kind != yaml.DocumentNode {
				// when the path is empty, it matches both the top level pseudo document node
				// and the actual top level node (e.g. map/sequence/whatever)
				// so when we are updating with no path, make sure we update the right node.
				node = node.Content[0]
			}

			log.Debug("going to update")
			DebugNode(node)
			log.Debug("with")
			DebugNode(changesToApply)

			if updateCommand.Overwrite || node.Value == "" {
				node.Value = changesToApply.Value
				node.Tag = changesToApply.Tag
				node.Kind = changesToApply.Kind
				node.Style = changesToApply.Style
				node.Anchor = changesToApply.Anchor
				node.Alias = changesToApply.Alias

				if !updateCommand.DontUpdateNodeContent {
					node.Content = changesToApply.Content
				}
			} else {
				log.Debug("skipping update as node already has value %v and overwriteFlag is ", node.Value, updateCommand.Overwrite)
			}

			switch updateCommand.CommentsMergeStrategy {
			case OverwriteCommentsMergeStrategy:
				node.HeadComment = changesToApply.HeadComment
				node.LineComment = changesToApply.LineComment
				node.FootComment = changesToApply.FootComment
			case SetWhenBlankCommentsMergeStrategy:
				if node.HeadComment == "" {
					node.HeadComment = changesToApply.HeadComment
				}
				if node.LineComment == "" {
					node.LineComment = changesToApply.LineComment
				}
				if node.FootComment == "" {
					node.FootComment = changesToApply.FootComment
				}
			case AppendCommentsMergeStrategy:
				if node.HeadComment == "" {
					node.HeadComment = changesToApply.HeadComment
				} else {
					node.HeadComment = node.HeadComment + "\n" + changesToApply.HeadComment
				}
				if node.LineComment == "" {
					node.LineComment = changesToApply.LineComment
				} else {
					node.LineComment = node.LineComment + " " + changesToApply.LineComment
				}
				if node.FootComment == "" {
					node.FootComment = changesToApply.FootComment
				} else {
					node.FootComment = node.FootComment + "\n" + changesToApply.FootComment
				}
			default:
			}

			log.Debug("result")
			DebugNode(node)

			return nil
		},
	}
}
