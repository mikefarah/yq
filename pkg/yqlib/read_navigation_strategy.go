package yqlib

func ReadNavigationStrategy() NavigationStrategy {
	return &NavigationStrategyImpl{
		visitedNodes: []*NodeContext{},
		followAlias: func(nodeContext NodeContext) bool {
			return true
		},
		autoCreateMap: func(nodeContext NodeContext) bool {
			return false
		},
		visit: func(nodeContext NodeContext) error {
			return nil
		},
	}
}
