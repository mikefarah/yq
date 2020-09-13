package yqlib

func UpdateNavigationStrategy(updateCommand UpdateCommand, autoCreate bool) NavigationStrategy {
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
			if updateCommand.Overwrite || node.Value == "" {
				log.Debug("going to update")
				DebugNode(node)
				log.Debug("with")
				DebugNode(changesToApply)
				if !updateCommand.DontUpdateNodeValue {
					node.Value = changesToApply.Value
				}
				node.Tag = changesToApply.Tag
				node.Kind = changesToApply.Kind
				node.Style = changesToApply.Style
				if !updateCommand.DontUpdateNodeContent {
					node.Content = changesToApply.Content
				}
				node.Anchor = changesToApply.Anchor
				node.Alias = changesToApply.Alias
				if updateCommand.CommentsMergeStrategy != IgnoreCommentsMergeStrategy {
					node.HeadComment = changesToApply.HeadComment
					node.LineComment = changesToApply.LineComment
					node.FootComment = changesToApply.FootComment
				}
			} else {
				log.Debug("skipping update as node already has value %v and overwriteFlag is ", node.Value, updateCommand.Overwrite)
			}
			return nil
		},
	}
}
