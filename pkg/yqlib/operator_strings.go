package yqlib

import (
	"container/list"
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type changeCasePrefs struct {
	ToUpperCase bool
}

func changeCaseOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	results := list.New()
	prefs := expressionNode.Operation.Preferences.(changeCasePrefs)

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		node := unwrapDoc(candidate.Node)

		if guessTagFromCustomType(node) != "!!str" {
			return Context{}, fmt.Errorf("cannot change case with %v, can only operate on strings. ", node.Tag)
		}

		newStringNode := &yaml.Node{Kind: yaml.ScalarNode, Tag: node.Tag, Style: node.Style}
		if prefs.ToUpperCase {
			newStringNode.Value = strings.ToUpper(node.Value)
		} else {
			newStringNode.Value = strings.ToLower(node.Value)
		}
		results.PushBack(candidate.CreateReplacement(newStringNode))

	}

	return context.ChildContext(results), nil

}

func getSubstituteParameters(d *dataTreeNavigator, block *ExpressionNode, context Context) (string, string, error) {
	regEx := ""
	replacementText := ""

	regExNodes, err := d.GetMatchingNodes(context.ReadOnlyClone(), block.LHS)
	if err != nil {
		return "", "", err
	}
	if regExNodes.MatchingNodes.Front() != nil {
		regEx = regExNodes.MatchingNodes.Front().Value.(*CandidateNode).Node.Value
	}

	log.Debug("regEx %v", regEx)

	replacementNodes, err := d.GetMatchingNodes(context, block.RHS)
	if err != nil {
		return "", "", err
	}
	if replacementNodes.MatchingNodes.Front() != nil {
		replacementText = replacementNodes.MatchingNodes.Front().Value.(*CandidateNode).Node.Value
	}

	return regEx, replacementText, nil
}

func substitute(original string, regex *regexp.Regexp, replacement string) *yaml.Node {
	replacedString := regex.ReplaceAllString(original, replacement)
	return &yaml.Node{Kind: yaml.ScalarNode, Value: replacedString, Tag: "!!str"}
}

func substituteStringOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	//rhs  block operator
	//lhs of block = regex
	//rhs of block = replacement expression
	block := expressionNode.RHS

	regExStr, replacementText, err := getSubstituteParameters(d, block, context)

	if err != nil {
		return Context{}, err
	}

	regEx, err := regexp.Compile(regExStr)
	if err != nil {
		return Context{}, err
	}

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := unwrapDoc(candidate.Node)

		if guessTagFromCustomType(node) != "!!str" {
			return Context{}, fmt.Errorf("cannot substitute with %v, can only substitute strings. Hint: Most often you'll want to use '|=' over '=' for this operation", node.Tag)
		}

		targetNode := substitute(node.Value, regEx, replacementText)
		result := candidate.CreateReplacement(targetNode)
		results.PushBack(result)
	}

	return context.ChildContext(results), nil

}

func addMatch(original []*yaml.Node, match string, offset int, name string) []*yaml.Node {

	newContent := append(original,
		createScalarNode("string", "string"))

	if offset < 0 {
		// offset of -1 means there was no match, force a null value like jq
		newContent = append(newContent,
			createScalarNode(nil, "null"),
		)
	} else {
		newContent = append(newContent,
			createScalarNode(match, match),
		)
	}

	newContent = append(newContent,
		createScalarNode("offset", "offset"),
		createScalarNode(offset, fmt.Sprintf("%v", offset)),
		createScalarNode("length", "length"),
		createScalarNode(len(match), fmt.Sprintf("%v", len(match))))

	if name != "" {
		newContent = append(newContent,
			createScalarNode("name", "name"),
			createScalarNode(name, name),
		)
	}
	return newContent
}

type matchPreferences struct {
	Global bool
}

func getMatches(matchPrefs matchPreferences, regEx *regexp.Regexp, value string) ([][]string, [][]int) {
	var allMatches [][]string
	var allIndices [][]int

	if matchPrefs.Global {
		allMatches = regEx.FindAllStringSubmatch(value, -1)
		allIndices = regEx.FindAllStringSubmatchIndex(value, -1)
	} else {
		allMatches = [][]string{regEx.FindStringSubmatch(value)}
		allIndices = [][]int{regEx.FindStringSubmatchIndex(value)}
	}

	log.Debug("allMatches, %v", allMatches)
	return allMatches, allIndices
}

func match(matchPrefs matchPreferences, regEx *regexp.Regexp, candidate *CandidateNode, value string, results *list.List) {
	subNames := regEx.SubexpNames()
	allMatches, allIndices := getMatches(matchPrefs, regEx, value)

	// if all matches just has an empty array in it,
	// then nothing matched
	if len(allMatches) > 0 && len(allMatches[0]) == 0 {
		return
	}

	for i, matches := range allMatches {
		capturesListNode := &yaml.Node{Kind: yaml.SequenceNode}
		match, submatches := matches[0], matches[1:]
		for j, submatch := range submatches {
			captureNode := &yaml.Node{Kind: yaml.MappingNode}
			captureNode.Content = addMatch(captureNode.Content, submatch, allIndices[i][2+j*2], subNames[j+1])
			capturesListNode.Content = append(capturesListNode.Content, captureNode)
		}

		node := &yaml.Node{Kind: yaml.MappingNode}
		node.Content = addMatch(node.Content, match, allIndices[i][0], "")
		node.Content = append(node.Content,
			createScalarNode("captures", "captures"),
			capturesListNode,
		)
		results.PushBack(candidate.CreateReplacement(node))

	}

}

func capture(matchPrefs matchPreferences, regEx *regexp.Regexp, candidate *CandidateNode, value string, results *list.List) {
	subNames := regEx.SubexpNames()
	allMatches, allIndices := getMatches(matchPrefs, regEx, value)

	// if all matches just has an empty array in it,
	// then nothing matched
	if len(allMatches) > 0 && len(allMatches[0]) == 0 {
		return
	}

	for i, matches := range allMatches {
		capturesNode := &yaml.Node{Kind: yaml.MappingNode}

		_, submatches := matches[0], matches[1:]
		for j, submatch := range submatches {
			capturesNode.Content = append(capturesNode.Content,
				createScalarNode(subNames[j+1], subNames[j+1]))

			offset := allIndices[i][2+j*2]
			// offset of -1 means there was no match, force a null value like jq
			if offset < 0 {
				capturesNode.Content = append(capturesNode.Content,
					createScalarNode(nil, "null"),
				)
			} else {
				capturesNode.Content = append(capturesNode.Content,
					createScalarNode(submatch, submatch),
				)
			}
		}

		results.PushBack(candidate.CreateReplacement(capturesNode))

	}

}

func extractMatchArguments(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (*regexp.Regexp, matchPreferences, error) {
	regExExpNode := expressionNode.RHS

	matchPrefs := matchPreferences{}

	// we got given parameters e.g. match(exp; params)
	if expressionNode.RHS.Operation.OperationType == blockOpType {
		block := expressionNode.RHS
		regExExpNode = block.LHS
		replacementNodes, err := d.GetMatchingNodes(context, block.RHS)
		if err != nil {
			return nil, matchPrefs, err
		}
		paramText := ""
		if replacementNodes.MatchingNodes.Front() != nil {
			paramText = replacementNodes.MatchingNodes.Front().Value.(*CandidateNode).Node.Value
		}
		if strings.Contains(paramText, "g") {
			paramText = strings.ReplaceAll(paramText, "g", "")
			matchPrefs.Global = true
		}
		if strings.Contains(paramText, "i") {
			return nil, matchPrefs, fmt.Errorf(`'i' is not a valid option for match. To ignore case, use an expression like match("(?i)cat")`)
		}
		if len(paramText) > 0 {
			return nil, matchPrefs, fmt.Errorf(`Unrecognised match params '%v', please see docs at https://mikefarah.gitbook.io/yq/operators/string-operators`, paramText)
		}
	}

	regExNodes, err := d.GetMatchingNodes(context.ReadOnlyClone(), regExExpNode)
	if err != nil {
		return nil, matchPrefs, err
	}
	log.Debug(NodesToString(regExNodes.MatchingNodes))
	regExStr := ""
	if regExNodes.MatchingNodes.Front() != nil {
		regExStr = regExNodes.MatchingNodes.Front().Value.(*CandidateNode).Node.Value
	}
	log.Debug("regEx %v", regExStr)
	regEx, err := regexp.Compile(regExStr)
	return regEx, matchPrefs, err
}

func matchOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	regEx, matchPrefs, err := extractMatchArguments(d, context, expressionNode)
	if err != nil {
		return Context{}, err
	}

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := unwrapDoc(candidate.Node)

		if guessTagFromCustomType(node) != "!!str" {
			return Context{}, fmt.Errorf("cannot match with %v, can only match strings. Hint: Most often you'll want to use '|=' over '=' for this operation", node.Tag)
		}

		match(matchPrefs, regEx, candidate, node.Value, results)
	}

	return context.ChildContext(results), nil
}

func captureOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	regEx, matchPrefs, err := extractMatchArguments(d, context, expressionNode)
	if err != nil {
		return Context{}, err
	}

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := unwrapDoc(candidate.Node)

		if guessTagFromCustomType(node) != "!!str" {
			return Context{}, fmt.Errorf("cannot match with %v, can only match strings. Hint: Most often you'll want to use '|=' over '=' for this operation", node.Tag)
		}
		capture(matchPrefs, regEx, candidate, node.Value, results)

	}

	return context.ChildContext(results), nil
}

func testOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	regEx, _, err := extractMatchArguments(d, context, expressionNode)
	if err != nil {
		return Context{}, err
	}

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := unwrapDoc(candidate.Node)

		if guessTagFromCustomType(node) != "!!str" {
			return Context{}, fmt.Errorf("cannot match with %v, can only match strings. Hint: Most often you'll want to use '|=' over '=' for this operation", node.Tag)
		}
		matches := regEx.FindStringSubmatch(node.Value)
		results.PushBack(createBooleanCandidate(candidate, len(matches) > 0))

	}

	return context.ChildContext(results), nil
}

func joinStringOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- joinStringOperator")
	joinStr := ""

	rhs, err := d.GetMatchingNodes(context.ReadOnlyClone(), expressionNode.RHS)
	if err != nil {
		return Context{}, err
	}
	if rhs.MatchingNodes.Front() != nil {
		joinStr = rhs.MatchingNodes.Front().Value.(*CandidateNode).Node.Value
	}

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := unwrapDoc(candidate.Node)
		if node.Kind != yaml.SequenceNode {
			return Context{}, fmt.Errorf("cannot join with %v, can only join arrays of scalars", node.Tag)
		}
		targetNode := join(node.Content, joinStr)
		result := candidate.CreateReplacement(targetNode)
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}

func join(content []*yaml.Node, joinStr string) *yaml.Node {
	var stringsToJoin []string
	for _, node := range content {
		str := node.Value
		if node.Tag == "!!null" {
			str = ""
		}
		stringsToJoin = append(stringsToJoin, str)
	}

	return &yaml.Node{Kind: yaml.ScalarNode, Value: strings.Join(stringsToJoin, joinStr), Tag: "!!str"}
}

func splitStringOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- splitStringOperator")
	splitStr := ""

	rhs, err := d.GetMatchingNodes(context.ReadOnlyClone(), expressionNode.RHS)
	if err != nil {
		return Context{}, err
	}
	if rhs.MatchingNodes.Front() != nil {
		splitStr = rhs.MatchingNodes.Front().Value.(*CandidateNode).Node.Value
	}

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := unwrapDoc(candidate.Node)
		if node.Tag == "!!null" {
			continue
		}

		if guessTagFromCustomType(node) != "!!str" {
			return Context{}, fmt.Errorf("Cannot split %v, can only split strings", node.Tag)
		}
		targetNode := split(node.Value, splitStr)
		result := candidate.CreateReplacement(targetNode)
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}

func split(value string, spltStr string) *yaml.Node {
	var contents []*yaml.Node

	if value != "" {
		var newStrings = strings.Split(value, spltStr)
		contents = make([]*yaml.Node, len(newStrings))

		for index, str := range newStrings {
			contents[index] = &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: str}
		}
	}

	return &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq", Content: contents}
}
