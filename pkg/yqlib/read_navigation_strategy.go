package yqlib

func ReadNavigationStrategy(deeplyTraverseArrays bool) NavigationStrategy {
	return &NavigationStrategyImpl{
		visitedNodes: []*NodeContext{},
		pathParser:   NewPathParser(),
		visit: func(nodeContext NodeContext) error {
			return nil
		},
		shouldDeeplyTraverse: func(nodeContext NodeContext) bool {
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
