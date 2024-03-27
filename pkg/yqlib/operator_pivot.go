package yqlib

import (
	"container/list"
	"fmt"
)

func getUniqueElementTag(seq *CandidateNode) (string, error) {
	switch l := len(seq.Content); l {
	case 0:
		return "", nil
	default:
		result := seq.Content[0].Tag
		for i := 1; i < l; i++ {
			t := seq.Content[i].Tag
			if t != result {
				return "", fmt.Errorf("sequence contains elements of %v and %v types", result, t)
			}
		}
		return result, nil
	}
}

var nullNodeFactory = func() *CandidateNode { return createScalarNode(nil, "") }

func pad[E any](array []E, length int, factory func() E) []E {
	sz := len(array)
	if sz >= length {
		return array
	}
	pad := make([]E, length-sz)
	for i := 0; i < len(pad); i++ {
		pad[i] = factory()
	}
	return append(array, pad...)
}

func pivotSequences(seq *CandidateNode) *CandidateNode {
	sz := len(seq.Content)
	if sz == 0 {
		return seq
	}
	m := make(map[int][]*CandidateNode)

	for i := 0; i < sz; i++ {
		row := seq.Content[i]
		for j := 0; j < len(row.Content); j++ {
			e := m[j]
			if e == nil {
				e = make([]*CandidateNode, 0, sz)
			}
			m[j] = append(pad(e, i, nullNodeFactory), row.Content[j])
		}
	}
	result := CandidateNode{Kind: SequenceNode}

	for i := 0; i < len(m); i++ {
		e := CandidateNode{Kind: SequenceNode}
		e.AddChildren(pad(m[i], sz, nullNodeFactory))
		result.AddChild(&e)
	}
	return &result
}

func pivotMaps(seq *CandidateNode) *CandidateNode {
	sz := len(seq.Content)
	if sz == 0 {
		return &CandidateNode{Kind: MappingNode}
	}
	m := make(map[string][]*CandidateNode)
	keys := make([]string, 0)

	for i := 0; i < sz; i++ {
		row := seq.Content[i]
		for j := 0; j < len(row.Content); j += 2 {
			k := row.Content[j].Value
			v := row.Content[j+1]
			e := m[k]
			if e == nil {
				keys = append(keys, k)
				e = make([]*CandidateNode, 0, sz)
			}
			m[k] = append(pad(e, i, nullNodeFactory), v)
		}
	}
	result := CandidateNode{Kind: MappingNode}
	for _, k := range keys {
		pivotRow := CandidateNode{Kind: SequenceNode}
		pivotRow.AddChildren(
			pad(m[k], sz, nullNodeFactory))
		result.AddKeyValueChild(createScalarNode(k, k), &pivotRow)
	}
	return &result
}

func pivotOperator(_ *dataTreeNavigator, context Context, _ *ExpressionNode) (Context, error) {
	log.Debug("Pivot")
	results := list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		if candidate.Tag != "!!seq" {
			return Context{}, fmt.Errorf("cannot pivot node of type %v", candidate.Tag)
		}
		tag, err := getUniqueElementTag(candidate)
		if err != nil {
			return Context{}, err
		}
		var pivot *CandidateNode
		switch tag {
		case "!!seq":
			pivot = pivotSequences(candidate)
		case "!!map":
			pivot = pivotMaps(candidate)
		default:
			return Context{}, fmt.Errorf("can only pivot elements of !!seq or !!map types, received %v", tag)
		}
		results.PushBack(pivot)
	}
	return context.ChildContext(results), nil
}
