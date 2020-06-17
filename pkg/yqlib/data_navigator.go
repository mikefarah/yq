package yqlib

import (
	"fmt"
	"strconv"

	yaml "gopkg.in/yaml.v3"
)

type DataNavigator interface {
	Traverse(value *yaml.Node, path []interface{}) error
}

type navigator struct {
	navigationStrategy NavigationStrategy
}

func NewDataNavigator(NavigationStrategy NavigationStrategy) DataNavigator {
	return &navigator{
		navigationStrategy: NavigationStrategy,
	}
}

func (n *navigator) Traverse(value *yaml.Node, path []interface{}) error {
	emptyArray := make([]interface{}, 0)
	log.Debugf("Traversing path %v", pathStackToString(path))
	return n.doTraverse(value, "", path, emptyArray)
}

func (n *navigator) doTraverse(value *yaml.Node, head interface{}, tail []interface{}, pathStack []interface{}) error {

	log.Debug("head %v", head)
	DebugNode(value)
	var nodeContext = NewNodeContext(value, head, tail, pathStack)

	var errorDeepSplatting error
	// no need to deeply traverse the DocumentNode, as it's already covered by its first child.
	if head == "**" && value.Kind != yaml.DocumentNode && value.Kind != yaml.ScalarNode && n.navigationStrategy.ShouldDeeplyTraverse(nodeContext) {
		if len(pathStack) == 0 || pathStack[len(pathStack)-1] != "<<" {
			errorDeepSplatting = n.recurse(value, head, tail, pathStack)
		}
		// ignore errors here, we are deep splatting so we may accidently give a string key
		// to an array sequence
		if len(tail) > 0 {
			_ = n.recurse(value, tail[0], tail[1:], pathStack)
		}
		return errorDeepSplatting
	}

	if value.Kind == yaml.DocumentNode {
		log.Debugf("its a document, diving into %v", head)
		DebugNode(value)
		return n.recurse(value, head, tail, pathStack)
	} else if len(tail) > 0 && value.Kind != yaml.ScalarNode {
		log.Debugf("diving into %v", tail[0])
		DebugNode(value)
		return n.recurse(value, tail[0], tail[1:], pathStack)
	}
	return n.navigationStrategy.Visit(nodeContext)
}

func (n *navigator) getOrReplace(original *yaml.Node, expectedKind yaml.Kind) *yaml.Node {
	if original.Kind != expectedKind {
		log.Debug("wanted %v but it was %v, overriding", KindString(expectedKind), KindString(original.Kind))
		return &yaml.Node{Kind: expectedKind}
	}
	return original
}

func (n *navigator) recurse(value *yaml.Node, head interface{}, tail []interface{}, pathStack []interface{}) error {
	log.Debug("recursing, processing %v, pathStack %v", head, pathStackToString(pathStack))

	nodeContext := NewNodeContext(value, head, tail, pathStack)

	if head == "**" && !n.navigationStrategy.ShouldOnlyDeeplyVisitLeaves(nodeContext) {
		nodeContext.IsMiddleNode = true
		errorVisitingDeeply := n.navigationStrategy.Visit(nodeContext)
		if errorVisitingDeeply != nil {
			return errorVisitingDeeply
		}
	}

	switch value.Kind {
	case yaml.MappingNode:
		log.Debug("its a map with %v entries", len(value.Content)/2)
		headString := fmt.Sprintf("%v", head)
		return n.recurseMap(value, headString, tail, pathStack)

	case yaml.SequenceNode:
		log.Debug("its a sequence of %v things!", len(value.Content))

		switch head := head.(type) {
		case int64:
			return n.recurseArray(value, head, head, tail, pathStack)
		default:

			if head == "+" {
				return n.appendArray(value, head, tail, pathStack)
			} else if len(value.Content) == 0 && head == "**" {
				return n.navigationStrategy.Visit(nodeContext)
			}
			return n.splatArray(value, head, tail, pathStack)
		}
	case yaml.AliasNode:
		log.Debug("its an alias!")
		DebugNode(value.Alias)
		if n.navigationStrategy.FollowAlias(nodeContext) {
			log.Debug("following the alias")
			return n.recurse(value.Alias, head, tail, pathStack)
		}
		return nil
	case yaml.DocumentNode:
		return n.doTraverse(value.Content[0], head, tail, pathStack)
	default:
		return n.navigationStrategy.Visit(nodeContext)
	}
}

func (n *navigator) recurseMap(value *yaml.Node, head string, tail []interface{}, pathStack []interface{}) error {
	traversedEntry := false
	errorVisiting := n.visitMatchingEntries(value, head, tail, pathStack, func(contents []*yaml.Node, indexInMap int) error {
		log.Debug("recurseMap: visitMatchingEntries for %v", contents[indexInMap].Value)
		n.navigationStrategy.DebugVisitedNodes()
		newPathStack := append(pathStack, contents[indexInMap].Value)
		log.Debug("should I traverse? head: %v, path: %v", head, pathStackToString(newPathStack))
		DebugNode(value)
		if n.navigationStrategy.ShouldTraverse(NewNodeContext(contents[indexInMap+1], head, tail, newPathStack), contents[indexInMap].Value) {
			log.Debug("recurseMap: Going to traverse")
			traversedEntry = true
			contents[indexInMap+1] = n.getOrReplace(contents[indexInMap+1], guessKind(head, tail, contents[indexInMap+1].Kind))
			errorTraversing := n.doTraverse(contents[indexInMap+1], head, tail, newPathStack)
			log.Debug("recurseMap: Finished traversing")
			n.navigationStrategy.DebugVisitedNodes()
			return errorTraversing
		} else {
			log.Debug("nope not traversing")
		}
		return nil
	})

	if errorVisiting != nil {
		return errorVisiting
	}

	if len(value.Content) == 0 && head == "**" {
		return n.navigationStrategy.Visit(NewNodeContext(value, head, tail, pathStack))
	} else if traversedEntry || n.navigationStrategy.GetPathParser().IsPathExpression(head) || !n.navigationStrategy.AutoCreateMap(NewNodeContext(value, head, tail, pathStack)) {
		return nil
	}

	_, errorParsingInt := strconv.ParseInt(head, 10, 64)

	mapEntryKey := yaml.Node{Value: head, Kind: yaml.ScalarNode}

	if errorParsingInt == nil {
		// fixes a json encoding problem where keys that look like numbers
		// get treated as numbers and cannot be used in a json map
		mapEntryKey.Style = yaml.LiteralStyle
	}

	value.Content = append(value.Content, &mapEntryKey)
	mapEntryValue := yaml.Node{Kind: guessKind(head, tail, 0)}
	value.Content = append(value.Content, &mapEntryValue)
	log.Debug("adding a new node %v - def a string", head)
	return n.doTraverse(&mapEntryValue, head, tail, append(pathStack, head))
}

// need to pass the node in, as it may be aliased
type mapVisitorFn func(contents []*yaml.Node, index int) error

func (n *navigator) visitDirectMatchingEntries(node *yaml.Node, head string, tail []interface{}, pathStack []interface{}, visit mapVisitorFn) error {
	var contents = node.Content
	for index := 0; index < len(contents); index = index + 2 {
		content := contents[index]

		log.Debug("index %v, checking %v, %v", index, content.Value, content.Tag)
		n.navigationStrategy.DebugVisitedNodes()
		errorVisiting := visit(contents, index)
		if errorVisiting != nil {
			return errorVisiting
		}
	}
	return nil
}

func (n *navigator) visitMatchingEntries(node *yaml.Node, head string, tail []interface{}, pathStack []interface{}, visit mapVisitorFn) error {
	var contents = node.Content
	log.Debug("visitMatchingEntries %v", head)
	DebugNode(node)
	// value.Content is a concatenated array of key, value,
	// so keys are in the even indexes, values in odd.
	// merge aliases are defined first, but we only want to traverse them
	// if we don't find a match directly on this node first.
	errorVisitedDirectEntries := n.visitDirectMatchingEntries(node, head, tail, pathStack, visit)

	if errorVisitedDirectEntries != nil || !n.navigationStrategy.FollowAlias(NewNodeContext(node, head, tail, pathStack)) {
		return errorVisitedDirectEntries
	}
	return n.visitAliases(contents, head, tail, pathStack, visit)
}

func (n *navigator) visitAliases(contents []*yaml.Node, head string, tail []interface{}, pathStack []interface{}, visit mapVisitorFn) error {
	// merge aliases are defined first, but we only want to traverse them
	// if we don't find a match on this node first.
	// traverse them backwards so that the last alias overrides the preceding.
	// a node can either be
	// an alias to one other node (e.g. <<: *blah)
	// or a sequence of aliases   (e.g. <<: [*blah, *foo])
	log.Debug("checking for aliases, head: %v, pathstack: %v", head, pathStackToString(pathStack))
	for index := len(contents) - 2; index >= 0; index = index - 2 {

		if contents[index+1].Kind == yaml.AliasNode && contents[index].Value == "<<" {
			valueNode := contents[index+1]
			log.Debug("found an alias")
			DebugNode(contents[index])
			DebugNode(valueNode)

			errorInAlias := n.visitMatchingEntries(valueNode.Alias, head, tail, pathStack, visit)
			if errorInAlias != nil {
				return errorInAlias
			}
		} else if contents[index+1].Kind == yaml.SequenceNode {
			// could be an array of aliases...
			errorVisitingAliasSeq := n.visitAliasSequence(contents[index+1].Content, head, tail, pathStack, visit)
			if errorVisitingAliasSeq != nil {
				return errorVisitingAliasSeq
			}
		}
	}
	return nil
}

func (n *navigator) visitAliasSequence(possibleAliasArray []*yaml.Node, head string, tail []interface{}, pathStack []interface{}, visit mapVisitorFn) error {
	// need to search this backwards too, so that aliases defined last override the preceding.
	for aliasIndex := len(possibleAliasArray) - 1; aliasIndex >= 0; aliasIndex = aliasIndex - 1 {
		child := possibleAliasArray[aliasIndex]
		if child.Kind == yaml.AliasNode {
			log.Debug("found an alias")
			DebugNode(child)
			errorInAlias := n.visitMatchingEntries(child.Alias, head, tail, pathStack, visit)
			if errorInAlias != nil {
				return errorInAlias
			}
		}
	}
	return nil
}

func (n *navigator) splatArray(value *yaml.Node, head interface{}, tail []interface{}, pathStack []interface{}) error {
	for index, childValue := range value.Content {
		log.Debug("processing")
		DebugNode(childValue)
		childValue = n.getOrReplace(childValue, guessKind(head, tail, childValue.Kind))

		newPathStack := append(pathStack, index)
		if n.navigationStrategy.ShouldTraverse(NewNodeContext(childValue, head, tail, newPathStack), childValue.Value) {
			// here we should not deeply traverse the array if we are appending..not sure how to do that.
			// need to visit instead...
			// easiest way is to pop off the head and pass the rest of the tail in.
			var err = n.doTraverse(childValue, head, tail, newPathStack)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (n *navigator) appendArray(value *yaml.Node, head interface{}, tail []interface{}, pathStack []interface{}) error {
	var newNode = yaml.Node{Kind: guessKind(head, tail, 0)}
	value.Content = append(value.Content, &newNode)
	log.Debug("appending a new node, %v", value.Content)
	return n.doTraverse(&newNode, head, tail, append(pathStack, len(value.Content)-1))
}

func (n *navigator) recurseArray(value *yaml.Node, index int64, head interface{}, tail []interface{}, pathStack []interface{}) error {
	var contentLength = int64(len(value.Content))
	for contentLength <= index {
		value.Content = append(value.Content, &yaml.Node{Kind: guessKind(head, tail, 0)})
		contentLength = int64(len(value.Content))
	}
	var indexToUse = index

	if indexToUse < 0 {
		indexToUse = contentLength + indexToUse
	}

	if indexToUse < 0 {
		return fmt.Errorf("Index [%v] out of range, array size is %v", index, contentLength)
	}

	value.Content[indexToUse] = n.getOrReplace(value.Content[indexToUse], guessKind(head, tail, value.Content[indexToUse].Kind))

	return n.doTraverse(value.Content[indexToUse], head, tail, append(pathStack, index))
}
