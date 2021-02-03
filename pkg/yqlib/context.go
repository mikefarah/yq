package yqlib

import (
	"container/list"

	"github.com/jinzhu/copier"
)

type Context struct {
	MatchingNodes  *list.List
	Variables      map[string]*list.List
	DontAutoCreate bool
}

func (n *Context) SingleChildContext(candidate *CandidateNode) Context {
	list := list.New()
	list.PushBack(candidate)
	return n.ChildContext(list)
}

func (n *Context) ChildContext(results *list.List) Context {
	clone := Context{}
	err := copier.Copy(&clone, n)
	if err != nil {
		log.Error("Error cloning context :(")
		panic(err)
	}
	clone.MatchingNodes = results
	return clone
}
