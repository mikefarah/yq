package yqlib

import (
	"container/list"
	"fmt"
	"time"

	"github.com/jinzhu/copier"
	logging "gopkg.in/op/go-logging.v1"
)

type Context struct {
	MatchingNodes  *list.List
	Variables      map[string]*list.List
	DontAutoCreate bool
	datetimeLayout string
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

func (n *Context) SetDateTimeLayout(newDateTimeLayout string) {
	n.datetimeLayout = newDateTimeLayout
}

func (n *Context) GetDateTimeLayout() string {
	if n.datetimeLayout != "" {
		return n.datetimeLayout
	}
	return time.RFC3339
}

func (n *Context) GetVariable(name string) *list.List {
	if n.Variables == nil {
		return nil
	}
	return n.Variables[name]
}

func (n *Context) SetVariable(name string, value *list.List) {
	if n.Variables == nil {
		n.Variables = make(map[string]*list.List)
	}
	n.Variables[name] = value
}

func (n *Context) ChildContext(results *list.List) Context {
	clone := Context{DontAutoCreate: n.DontAutoCreate, datetimeLayout: n.datetimeLayout}
	clone.Variables = make(map[string]*list.List)
	if len(n.Variables) > 0 {
		err := copier.Copy(&clone.Variables, n.Variables)
		if err != nil {
			log.Error("Error cloning context :(")
			panic(err)
		}
	}
	clone.MatchingNodes = results
	return clone
}

func (n *Context) ToString() string {
	if !log.IsEnabledFor(logging.DEBUG) {
		return ""
	}
	result := fmt.Sprintf("Context\nDontAutoCreate: %v\n", n.DontAutoCreate)
	return result + NodesToString(n.MatchingNodes)
}

func (n *Context) DeepClone() Context {
	clone := Context{}
	err := copier.Copy(&clone, n)
	// copier doesn't do lists properly for some reason
	clone.MatchingNodes = list.New()
	for el := n.MatchingNodes.Front(); el != nil; el = el.Next() {
		clonedNode, err := el.Value.(*CandidateNode).Copy()
		if err != nil {
			log.Error("Error cloning context :(")
			panic(err)
		}
		clone.MatchingNodes.PushBack(clonedNode)
	}

	if err != nil {
		log.Error("Error cloning context :(")
		panic(err)
	}
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
