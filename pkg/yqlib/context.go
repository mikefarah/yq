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

func (n *Context) SingleReadonlyChildContext(candidate *CandidateNode) Context {
	list := list.New()
	list.PushBack(candidate)
	newContext := n.ChildContext(list)
	newContext.DontAutoCreate = true
	return newContext
}

func (n *Context) SingleChildContext(candidate *CandidateNode) Context {
	list := list.New()
	list.PushBack(candidate)
	return n.ChildContext(list)
}

func (n *Context) GetVariable(name string) *list.List {
	if n.Variables == nil {
		return nil
	}
	log.Debug("GetVariable - %v to %v", name, NodesToString(n.Variables[name]))
	return n.Variables[name]
}

func (n *Context) SetVariable(name string, value *list.List) {
	if n.Variables == nil {
		n.Variables = make(map[string]*list.List)
	}
	log.Debug("SetVariable - %v to %v", name, NodesToString(value))
	n.Variables[name] = value
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

func (n *Context) Clone() Context {
	clone := Context{}
	err := copier.Copy(&clone, n)
	if err != nil {
		log.Error("Error cloning context :(")
		panic(err)
	}
	return clone
}

func (n *Context) ReadOnlyClone() Context {
	clone := n.Clone()
	clone.DontAutoCreate = true
	return clone
}

func (n *Context) WritableClone() Context {
	clone := n.Clone()
	clone.DontAutoCreate = false
	return clone
}
