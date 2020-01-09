package yqlib

import (
	"strconv"

	errors "github.com/pkg/errors"
	yaml "gopkg.in/yaml.v3"
)

type DataNavigator interface {
	Traverse(value *yaml.Node, path []string) error
}

type navigator struct {
	navigationStrategy NavigationStrategy
}

func NewDataNavigator(NavigationStrategy NavigationStrategy) DataNavigator {
	return &navigator{
		navigationStrategy: NavigationStrategy,
	}
}

func (n *navigator) Traverse(value *yaml.Node, path []string) error {
	realValue := value
	emptyArray := make([]interface{}, 0)
	if realValue.Kind == yaml.DocumentNode {
		log.Debugf("its a document! returning the first child")
		return n.doTraverse(value.Content[0], "", path, emptyArray)
	}
	return n.doTraverse(value, "", path, emptyArray)
}

func (n *navigator) doTraverse(value *yaml.Node, head string, tail []string, pathStack []interface{}) error {
	log.Debug("head %v", head)
	DebugNode(value)
	var errorDeepSplatting error
	if head == "**" && value.Kind != yaml.ScalarNode {
		errorDeepSplatting = n.recurse(value, head, tail, pathStack)
		// ignore errors here, we are deep splatting so we may accidently give a string key
		// to an array sequence
		if len(tail) > 0 {
			n.recurse(value, tail[0], tail[1:], pathStack)
		}
		return errorDeepSplatting
	}

	if len(tail) > 0 {
		log.Debugf("diving into %v", tail[0])
		DebugNode(value)
		return n.recurse(value, tail[0], tail[1:], pathStack)
	}
	return n.navigationStrategy.Visit(NewNodeContext(value, head, tail, pathStack))
}

func (n *navigator) getOrReplace(original *yaml.Node, expectedKind yaml.Kind) *yaml.Node {
	if original.Kind != expectedKind {
		log.Debug("wanted %v but it was %v, overriding", expectedKind, original.Kind)
		return &yaml.Node{Kind: expectedKind}
	}
	return original
}

func (n *navigator) recurse(value *yaml.Node, head string, tail []string, pathStack []interface{}) error {
	log.Debug("recursing, processing %v", head)
	switch value.Kind {
	case yaml.MappingNode:
		log.Debug("its a map with %v entries", len(value.Content)/2)
		return n.recurseMap(value, head, tail, pathStack)
	case yaml.SequenceNode:
		log.Debug("its a sequence of %v things!", len(value.Content))
		if head == "*" || head == "**" {
			return n.splatArray(value, head, tail, pathStack)
		} else if head == "+" {
			return n.appendArray(value, head, tail, pathStack)
		}
		return n.recurseArray(value, head, tail, pathStack)
	case yaml.AliasNode:
		log.Debug("its an alias!")
		DebugNode(value.Alias)
		if n.navigationStrategy.FollowAlias(NewNodeContext(value, head, tail, pathStack)) == true {
			log.Debug("following the alias")
			return n.recurse(value.Alias, head, tail, pathStack)
		}
		return nil
	default:
		return nil
	}
}

func (n *navigator) recurseMap(value *yaml.Node, head string, tail []string, pathStack []interface{}) error {
	traversedEntry := false
	errorVisiting := n.visitMatchingEntries(value, head, tail, pathStack, func(contents []*yaml.Node, indexInMap int) error {
		log.Debug("recurseMap: visitMatchingEntries")
		n.navigationStrategy.DebugVisitedNodes()
		newPathStack := append(pathStack, contents[indexInMap].Value)
		log.Debug("appended %v", contents[indexInMap].Value)
		n.navigationStrategy.DebugVisitedNodes()
		log.Debug("should I traverse? %v, %v", head, pathStackToString(newPathStack))
		DebugNode(value)
		if n.navigationStrategy.ShouldTraverse(NewNodeContext(contents[indexInMap+1], head, tail, newPathStack), contents[indexInMap].Value) == true {
			log.Debug("recurseMap: Going to traverse")
			traversedEntry = true
			// contents[indexInMap+1] = n.getOrReplace(contents[indexInMap+1], guessKind(head, tail, contents[indexInMap+1].Kind))
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

	if traversedEntry == true || head == "*" || head == "**" || n.navigationStrategy.AutoCreateMap(NewNodeContext(value, head, tail, pathStack)) == false {
		return nil
	}

	mapEntryKey := yaml.Node{Value: head, Kind: yaml.ScalarNode}
	value.Content = append(value.Content, &mapEntryKey)
	mapEntryValue := yaml.Node{Kind: guessKind(head, tail, 0)}
	value.Content = append(value.Content, &mapEntryValue)
	log.Debug("adding new node %v", head)
	return n.doTraverse(&mapEntryValue, head, tail, append(pathStack, head))
}

// need to pass the node in, as it may be aliased
type mapVisitorFn func(contents []*yaml.Node, index int) error

func (n *navigator) visitDirectMatchingEntries(node *yaml.Node, head string, tail []string, pathStack []interface{}, visit mapVisitorFn) error {
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

func (n *navigator) visitMatchingEntries(node *yaml.Node, head string, tail []string, pathStack []interface{}, visit mapVisitorFn) error {
	var contents = node.Content
	log.Debug("visitMatchingEntries %v", head)
	DebugNode(node)
	// value.Content is a concatenated array of key, value,
	// so keys are in the even indexes, values in odd.
	// merge aliases are defined first, but we only want to traverse them
	// if we don't find a match directly on this node first.
	errorVisitedDirectEntries := n.visitDirectMatchingEntries(node, head, tail, pathStack, visit)

	if errorVisitedDirectEntries != nil || n.navigationStrategy.FollowAlias(NewNodeContext(node, head, tail, pathStack)) == false {
		return errorVisitedDirectEntries
	}
	return n.visitAliases(contents, head, tail, pathStack, visit)
}

func (n *navigator) visitAliases(contents []*yaml.Node, head string, tail []string, pathStack []interface{}, visit mapVisitorFn) error {
	// merge aliases are defined first, but we only want to traverse them
	// if we don't find a match on this node first.
	// traverse them backwards so that the last alias overrides the preceding.
	// a node can either be
	// an alias to one other node (e.g. <<: *blah)
	// or a sequence of aliases   (e.g. <<: [*blah, *foo])
	log.Debug("checking for aliases")
	for index := len(contents) - 2; index >= 0; index = index - 2 {

		if contents[index+1].Kind == yaml.AliasNode {
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

func (n *navigator) visitAliasSequence(possibleAliasArray []*yaml.Node, head string, tail []string, pathStack []interface{}, visit mapVisitorFn) error {
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

func (n *navigator) splatArray(value *yaml.Node, head string, tail []string, pathStack []interface{}) error {
	for index, childValue := range value.Content {
		log.Debug("processing")
		DebugNode(childValue)
		childValue = n.getOrReplace(childValue, guessKind(head, tail, childValue.Kind))
		var err = n.doTraverse(childValue, head, tail, append(pathStack, index))
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *navigator) appendArray(value *yaml.Node, head string, tail []string, pathStack []interface{}) error {
	var newNode = yaml.Node{Kind: guessKind(head, tail, 0)}
	value.Content = append(value.Content, &newNode)
	log.Debug("appending a new node, %v", value.Content)
	return n.doTraverse(&newNode, head, tail, append(pathStack, len(value.Content)-1))
}

func (n *navigator) recurseArray(value *yaml.Node, head string, tail []string, pathStack []interface{}) error {
	var index, err = strconv.ParseInt(head, 10, 64) // nolint
	if err != nil {
		return errors.Wrapf(err, "Error parsing array index '%v' for '%v'", head, pathStackToString(pathStack))
	}

	for int64(len(value.Content)) <= index {
		value.Content = append(value.Content, &yaml.Node{Kind: guessKind(head, tail, 0)})
	}

	value.Content[index] = n.getOrReplace(value.Content[index], guessKind(head, tail, value.Content[index].Kind))

	return n.doTraverse(value.Content[index], head, tail, append(pathStack, index))
}
