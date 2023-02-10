package yqlib

import (
	"container/list"
	"fmt"
	"math/rand"

	yaml "gopkg.in/yaml.v3"
)

func shuffleOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	// ignore CWE-338 gosec issue of not using crypto/rand
	// this is just to shuffle an array rather generating a
	// secret or something that needs proper rand.
	myRand := rand.New(rand.NewSource(Now().UnixNano())) // #nosec

	results := list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		candidateNode := unwrapDoc(candidate.Node)

		if candidateNode.Kind != yaml.SequenceNode {
			return context, fmt.Errorf("node at path [%v] is not an array (it's a %v)", candidate.GetNicePath(), candidate.GetNiceTag())
		}

		result := deepClone(candidateNode)

		a := result.Content

		myRand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
		results.PushBack(candidate.CreateReplacement(result))
	}
	return context.ChildContext(results), nil
}
