package yqlib

func ReadNavigationStrategy() NavigationStrategy {
	return &NavigationStrategyImpl{
		visitedNodes: []*NodeContext{},
		pathParser:   NewPathParser(),
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
