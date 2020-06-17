package yqlib

import "gopkg.in/yaml.v3"

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

			if updateCommand.Overwrite || node.Value == "" {
				log.Debug("going to update")
				DebugNode(node)
				log.Debug("with")
				DebugNode(changesToApply)
				node.Value = changesToApply.Value
				node.Tag = changesToApply.Tag
				node.Kind = changesToApply.Kind
				node.Style = changesToApply.Style
				node.Anchor = changesToApply.Anchor
				node.Alias = changesToApply.Alias
				node.HeadComment = changesToApply.HeadComment
				node.LineComment = changesToApply.LineComment
				node.FootComment = changesToApply.FootComment

				if !updateCommand.DontUpdateNodeContent {
					node.Content = changesToApply.Content
				}

				// // TODO: mergeComments flag
				// if node.HeadComment != "" && changesToApply.HeadComment != "" {
				// 	node.HeadComment = node.HeadComment + "\n" + changesToApply.HeadComment
				// 	log.Debug("merged comments with a space, %v", node.HeadComment)
				// } else {
				// 	node.HeadComment = node.HeadComment + changesToApply.HeadComment
				// 	if node.HeadComment != "" {
				// 		log.Debug("merged comments with no space, %v", node.HeadComment)
				// 	}
				// }
				// node.LineComment = node.LineComment + changesToApply.LineComment
				// node.FootComment = node.FootComment + changesToApply.FootComment
			} else {
				log.Debug("skipping update as node already has value %v and overwriteFlag is ", node.Value, updateCommand.Overwrite)
			}
			return nil
		},
	}
}
