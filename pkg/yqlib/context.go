package yqlib

import "container/list"

type Context struct {
	MatchingNodes *list.List
	Variables     map[string]*list.List
}

func (n *Context) SingleChildContext(candidate *CandidateNode) Context {
	elMap := list.New()
	elMap.PushBack(candidate)
	return Context{MatchingNodes: elMap, Variables: n.Variables}
}

func (n *Context) ChildContext(results *list.List) Context {
	return Context{MatchingNodes: results, Variables: n.Variables}
}
