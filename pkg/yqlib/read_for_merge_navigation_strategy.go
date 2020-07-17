package yqlib

import "gopkg.in/yaml.v3"

func ReadForMergeNavigationStrategy(deeplyTraverseArrays bool, overwriteArray bool) NavigationStrategy {
	return &NavigationStrategyImpl{
		visitedNodes: []*NodeContext{},
		pathParser:   NewPathParser(),
		followAlias: func(nodeContext NodeContext) bool {
			return false
		},
		shouldOnlyDeeplyVisitLeaves: func(nodeContext NodeContext) bool {
			return false
		},
		visit: func(nodeContext NodeContext) error {
			return nil
		},
		shouldDeeplyTraverse: func(nodeContext NodeContext) bool {
			if nodeContext.Node.Kind == yaml.SequenceNode && overwriteArray {
				nodeContext.IsMiddleNode = false
				return false
			}

			var isInArray = false
			if len(nodeContext.PathStack) > 0 {
				var lastElement = nodeContext.PathStack[len(nodeContext.PathStack)-1]
				switch lastElement.(type) {
				case int:
					isInArray = true
				default:
					isInArray = false
				}
			}
			return deeplyTraverseArrays || !isInArray
		},
	}
}
