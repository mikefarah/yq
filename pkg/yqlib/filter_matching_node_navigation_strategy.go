package yqlib

func FilterMatchingNodesNavigationStrategy(value string) NavigationStrategy {
	return &NavigationStrategyImpl{
		visitedNodes: []*NodeContext{},
		pathParser:   NewPathParser(),
		visit: func(nodeContext NodeContext) error {
			return nil
		},
		shouldVisitExtraFn: func(nodeContext NodeContext) bool {
			log.Debug("does %v match %v ? %v", nodeContext.Node.Value, value, nodeContext.Node.Value == value)
			return matchesString(value, nodeContext.Node.Value)
		},
	}
}
