package yqlib

import (
	"container/list"
	"fmt"
	"log/slog"
	"time"
)

type Context struct {
	MatchingNodes       *list.List
	Variables           map[string]*list.List
	DontAutoCreate      bool
	AutoCreateOnMissing bool
	datetimeLayout      string
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

// CollectingChildContext evaluates a collect expression without allowing a
// read-only parent context to suppress missing-key nulls. Missing values are
// materialised by traversal without being added to a read-only input node.
func (n *Context) CollectingChildContext(candidate *CandidateNode) Context {
	newContext := n.SingleChildContext(candidate)
	newContext.AutoCreateOnMissing = true
	return newContext
}

// CollectingReadonlyChildContext is CollectingChildContext with the read-only
// flag forced: missing keys materialise as detached nulls and the input node
// is never modified. Used when collecting over documents that must not gain
// auto-created keys (e.g. collect-together over multiple documents).
func (n *Context) CollectingReadonlyChildContext(candidate *CandidateNode) Context {
	newContext := n.SingleReadonlyChildContext(candidate)
	newContext.AutoCreateOnMissing = true
	return newContext
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
	clone := Context{
		DontAutoCreate:      n.DontAutoCreate,
		AutoCreateOnMissing: n.AutoCreateOnMissing,
		datetimeLayout:      n.datetimeLayout,
	}
	clone.Variables = make(map[string]*list.List)
	for variableKey, originalValueList := range n.Variables {

		variableCopyList := list.New()
		for el := originalValueList.Front(); el != nil; el = el.Next() {
			// note that we dont make a copy of the candidate node
			// this is so the 'ref' operator can work correctly.
			clonedNode := el.Value.(*CandidateNode)
			variableCopyList.PushBack(clonedNode)
		}

		clone.Variables[variableKey] = variableCopyList
	}

	clone.MatchingNodes = results
	return clone
}

func (n *Context) ToString() string {
	if !log.IsEnabledFor(slog.LevelDebug) {
		return ""
	}
	result := fmt.Sprintf("Context\nDontAutoCreate: %v\n", n.DontAutoCreate)
	return result + NodesToString(n.MatchingNodes)
}

func (n *Context) DeepClone() Context {

	clonedContent := list.New()
	for el := n.MatchingNodes.Front(); el != nil; el = el.Next() {
		clonedNode := el.Value.(*CandidateNode).Copy()
		clonedContent.PushBack(clonedNode)
	}

	return n.ChildContext(clonedContent)
}

func (n *Context) Clone() Context {
	return n.ChildContext(n.MatchingNodes)
}

func (n *Context) ReadOnlyClone() Context {
	clone := n.Clone()
	clone.DontAutoCreate = true
	// a read-only lookup wants existing nodes only: never materialise
	// missing-key nulls here, even inside a collect (see CollectingChildContext)
	clone.AutoCreateOnMissing = false
	return clone
}

func (n *Context) WritableClone() Context {
	clone := n.Clone()
	clone.DontAutoCreate = false
	return clone
}
