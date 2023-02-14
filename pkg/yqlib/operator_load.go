package yqlib

import (
	"bufio"
	"container/list"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var LoadYamlPreferences = YamlPreferences{
	LeadingContentPreProcessing: false,
	PrintDocSeparators:          true,
	UnwrapScalar:                true,
	EvaluateTogether:            false,
}

type loadPrefs struct {
	loadAsString bool
	decoder      Decoder
}

func loadString(filename string) (*CandidateNode, error) {
	// ignore CWE-22 gosec issue - that's more targeted for http based apps that run in a public directory,
	// and ensuring that it's not possible to give a path to a file outside that directory.

	filebytes, err := os.ReadFile(filename) // #nosec
	if err != nil {
		return nil, err
	}

	return &CandidateNode{Node: &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: string(filebytes)}}, nil
}

func loadYaml(filename string, decoder Decoder) (*CandidateNode, error) {
	if decoder == nil {
		return nil, fmt.Errorf("could not load %s", filename)
	}

	file, err := os.Open(filename) // #nosec
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)

	documents, err := readDocuments(reader, filename, 0, decoder)
	if err != nil {
		return nil, err
	}

	if documents.Len() == 0 {
		// return null candidate
		return &CandidateNode{Node: &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null"}}, nil
	} else if documents.Len() == 1 {
		candidate := documents.Front().Value.(*CandidateNode)
		return candidate, nil

	} else {
		sequenceNode := &CandidateNode{Node: &yaml.Node{Kind: yaml.SequenceNode}}
		for doc := documents.Front(); doc != nil; doc = doc.Next() {
			sequenceNode.Node.Content = append(sequenceNode.Node.Content, unwrapDoc(doc.Value.(*CandidateNode).Node))
		}
		return sequenceNode, nil
	}
}

func loadYamlOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("loadYamlOperator")

	loadPrefs := expressionNode.Operation.Preferences.(loadPrefs)

	// need to evaluate the 1st parameter against the context
	// and return the data accordingly.

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		rhs, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(candidate), expressionNode.RHS)
		if err != nil {
			return Context{}, err
		}
		if rhs.MatchingNodes.Front() == nil {
			return Context{}, fmt.Errorf("Filename expression returned nil")
		}
		nameCandidateNode := rhs.MatchingNodes.Front().Value.(*CandidateNode)

		filename := nameCandidateNode.Node.Value

		var contentsCandidate *CandidateNode

		if loadPrefs.loadAsString {
			contentsCandidate, err = loadString(filename)
		} else {
			contentsCandidate, err = loadYaml(filename, loadPrefs.decoder)
		}
		if err != nil {
			return Context{}, fmt.Errorf("Failed to load %v: %w", filename, err)
		}

		results.PushBack(contentsCandidate)

	}

	return context.ChildContext(results), nil
}
