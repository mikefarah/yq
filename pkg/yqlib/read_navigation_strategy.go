package yqlib

func ReadNavigationStrategy() NavigationStrategy {
	return &NavigationStrategyImpl{
		visitedNodes: []*NodeContext{},
		pathParser:   NewPathParser(),
		visit: func(nodeContext NodeContext) error {
			return nil
		},
	}
}
