package yqlib

import (
	"bytes"
	"strconv"
	"strings"

	logging "gopkg.in/op/go-logging.v1"
	yaml "gopkg.in/yaml.v3"
)

type DataNavigator interface {
	DebugNode(node *yaml.Node)
	Get(rootNode *yaml.Node, path []string) ([]MatchingNode, error)
	Update(rootNode *yaml.Node, path []string, changesToApply *yaml.Node) error
	Delete(rootNode *yaml.Node, path []string) error
	GuessKind(tail []string, guess yaml.Kind) yaml.Kind
}

type navigator struct {
	log           *logging.Logger
	followAliases bool
}

type VisitorFn func(matchingNode *yaml.Node, pathStack []interface{}) error

func NewDataNavigator(l *logging.Logger, followAliases bool) DataNavigator {
	return &navigator{
		log:           l,
		followAliases: followAliases,
	}
}

type MatchingNode struct {
	Node      *yaml.Node
	PathStack []interface{}
}

func (n *navigator) Get(value *yaml.Node, path []string) ([]MatchingNode, error) {
	matchingNodes := make([]MatchingNode, 0)

	n.Visit(value, path, func(matchedNode *yaml.Node, pathStack []interface{}) error {
		matchingNodes = append(matchingNodes, MatchingNode{matchedNode, pathStack})
		n.log.Debug("Matched")
		for _, pathElement := range pathStack {
			n.log.Debug("%v", pathElement)
		}
		n.DebugNode(matchedNode)
		return nil
	})
	return matchingNodes, nil
}

func (n *navigator) Update(rootNode *yaml.Node, path []string, changesToApply *yaml.Node) error {
	errorVisiting := n.Visit(rootNode, path, func(nodeToUpdate *yaml.Node, pathStack []interface{}) error {
		n.log.Debug("going to update")
		n.DebugNode(nodeToUpdate)
		n.log.Debug("with")
		n.DebugNode(changesToApply)
		nodeToUpdate.Value = changesToApply.Value
		nodeToUpdate.Tag = changesToApply.Tag
		nodeToUpdate.Kind = changesToApply.Kind
		nodeToUpdate.Style = changesToApply.Style
		nodeToUpdate.Content = changesToApply.Content
		nodeToUpdate.HeadComment = changesToApply.HeadComment
		nodeToUpdate.LineComment = changesToApply.LineComment
		nodeToUpdate.FootComment = changesToApply.FootComment
		return nil
	})
	return errorVisiting
}

// TODO: refactor delete..
func (n *navigator) Delete(rootNode *yaml.Node, path []string) error {

	lastBit, newTail := path[len(path)-1], path[:len(path)-1]
	n.log.Debug("splitting path, %v", lastBit)
	n.log.Debug("new tail, %v", newTail)
	errorVisiting := n.Visit(rootNode, newTail, func(nodeToUpdate *yaml.Node, pathStack []interface{}) error {
		n.log.Debug("need to find %v in here", lastBit)
		n.DebugNode(nodeToUpdate)
		original := nodeToUpdate.Content
		if nodeToUpdate.Kind == yaml.SequenceNode {
			var index, err = strconv.ParseInt(lastBit, 10, 64) // nolint
			if err != nil {
				return err
			}
			if index >= int64(len(nodeToUpdate.Content)) {
				n.log.Debug("index %v is greater than content length %v", index, len(nodeToUpdate.Content))
				return nil
			}
			nodeToUpdate.Content = append(original[:index], original[index+1:]...)

		} else if nodeToUpdate.Kind == yaml.MappingNode {
			// need to delete in reverse - otherwise the matching indexes
			// become incorrect.
			matchingIndices := make([]int, 0)
			_, errorVisiting := n.visitMatchingEntries(nodeToUpdate.Content, lastBit, func(matchingNode []*yaml.Node, indexInMap int) error {
				matchingIndices = append(matchingIndices, indexInMap)
				n.log.Debug("matchingIndices %v", indexInMap)
				return nil
			})
			n.log.Debug("delete matching indices now")
			n.log.Debug("%v", matchingIndices)
			if errorVisiting != nil {
				return errorVisiting
			}
			for i := len(matchingIndices) - 1; i >= 0; i-- {
				indexToDelete := matchingIndices[i]
				n.log.Debug("deleting index %v, %v", indexToDelete, nodeToUpdate.Content[indexToDelete].Value)
				nodeToUpdate.Content = append(nodeToUpdate.Content[:indexToDelete], nodeToUpdate.Content[indexToDelete+2:]...)
			}

		}

		return nil
	})
	return errorVisiting
}

func (n *navigator) Visit(value *yaml.Node, path []string, visitor VisitorFn) error {
	realValue := value
	emptyArray := make([]interface{}, 0)
	if realValue.Kind == yaml.DocumentNode {
		n.log.Debugf("its a document! returning the first child")
		return n.doVisit(value.Content[0], path, visitor, emptyArray)
	}
	return n.doVisit(value, path, visitor, emptyArray)
}

func (n *navigator) doVisit(value *yaml.Node, path []string, visitor VisitorFn, pathStack []interface{}) error {
	if len(path) > 0 {
		n.log.Debugf("diving into %v", path[0])
		n.DebugNode(value)
		return n.recurse(value, path[0], path[1:], visitor, pathStack)
	}
	return visitor(value, pathStack)
}

func (n *navigator) GuessKind(tail []string, guess yaml.Kind) yaml.Kind {
	n.log.Debug("tail %v", tail)
	if len(tail) == 0 && guess == 0 {
		n.log.Debug("end of path, must be a scalar")
		return yaml.ScalarNode
	} else if len(tail) == 0 {
		return guess
	}

	var _, errorParsingInt = strconv.ParseInt(tail[0], 10, 64)
	if tail[0] == "+" || errorParsingInt == nil {
		return yaml.SequenceNode
	}
	if tail[0] == "*" && (guess == yaml.SequenceNode || guess == yaml.MappingNode) {
		return guess
	}
	if guess == yaml.AliasNode {
		n.log.Debug("guess was an alias, okey doke.")
		return guess
	}
	n.log.Debug("forcing a mapping node")
	n.log.Debug("yaml.SequenceNode ?", guess == yaml.SequenceNode)
	n.log.Debug("yaml.ScalarNode ?", guess == yaml.ScalarNode)
	return yaml.MappingNode
}

func (n *navigator) getOrReplace(original *yaml.Node, expectedKind yaml.Kind) *yaml.Node {
	if original.Kind != expectedKind {
		n.log.Debug("wanted %v but it was %v, overriding", expectedKind, original.Kind)
		return &yaml.Node{Kind: expectedKind}
	}
	return original
}

func (n *navigator) DebugNode(value *yaml.Node) {
	if value == nil {
		n.log.Debug("-- node is nil --")
	} else if n.log.IsEnabledFor(logging.DEBUG) {
		buf := new(bytes.Buffer)
		encoder := yaml.NewEncoder(buf)
		encoder.Encode(value)
		encoder.Close()
		n.log.Debug("Tag: %v", value.Tag)
		n.log.Debug("%v", buf.String())
	}
}

func (n *navigator) recurse(value *yaml.Node, head string, tail []string, visitor VisitorFn, pathStack []interface{}) error {
	n.log.Debug("recursing, processing %v", head)
	switch value.Kind {
	case yaml.MappingNode:
		n.log.Debug("its a map with %v entries", len(value.Content)/2)
		if head == "*" {
			return n.splatMap(value, tail, visitor, pathStack)
		}
		return n.recurseMap(value, head, tail, visitor, pathStack)
	case yaml.SequenceNode:
		n.log.Debug("its a sequence of %v things!, %v", len(value.Content))
		if head == "*" {
			return n.splatArray(value, tail, visitor, pathStack)
		} else if head == "+" {
			return n.appendArray(value, tail, visitor, pathStack)
		}
		return n.recurseArray(value, head, tail, visitor, pathStack)
	case yaml.AliasNode:
		n.log.Debug("its an alias, followAliases: %v", n.followAliases)
		n.DebugNode(value.Alias)
		if n.followAliases == true {
			return n.recurse(value.Alias, head, tail, visitor, pathStack)
		}
		return nil
	default:
		return nil
	}
}

func (n *navigator) splatMap(value *yaml.Node, tail []string, visitor VisitorFn, pathStack []interface{}) error {
	for index, content := range value.Content {
		if index%2 == 0 {
			continue
		}
		content = n.getOrReplace(content, n.GuessKind(tail, content.Kind))
		var err = n.doVisit(content, tail, visitor, append(pathStack, value.Content[index-1].Value))
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *navigator) recurseMap(value *yaml.Node, head string, tail []string, visitor VisitorFn, pathStack []interface{}) error {
	visited, errorVisiting := n.visitMatchingEntries(value.Content, head, func(contents []*yaml.Node, indexInMap int) error {
		contents[indexInMap+1] = n.getOrReplace(contents[indexInMap+1], n.GuessKind(tail, contents[indexInMap+1].Kind))
		return n.doVisit(contents[indexInMap+1], tail, visitor, append(pathStack, contents[indexInMap].Value))
	})

	if errorVisiting != nil {
		return errorVisiting
	}

	if visited {
		return nil
	}

	//didn't find it, lets add it.
	mapEntryKey := yaml.Node{Value: head, Kind: yaml.ScalarNode}
	value.Content = append(value.Content, &mapEntryKey)
	mapEntryValue := yaml.Node{Kind: n.GuessKind(tail, 0)}
	value.Content = append(value.Content, &mapEntryValue)
	n.log.Debug("adding new node %v", value.Content)
	return n.doVisit(&mapEntryValue, tail, visitor, append(pathStack, head))
}

// need to pass the node in, as it may be aliased
type mapVisitorFn func(contents []*yaml.Node, index int) error

func (n *navigator) visitDirectMatchingEntries(contents []*yaml.Node, key string, visit mapVisitorFn) (bool, error) {
	visited := false
	for index := 0; index < len(contents); index = index + 2 {
		content := contents[index]
		n.log.Debug("index %v, checking %v, %v", index, content.Value, content.Tag)

		if n.matchesKey(key, content.Value) {
			n.log.Debug("found a match! %v", content.Value)
			errorVisiting := visit(contents, index)
			if errorVisiting != nil {
				return visited, errorVisiting
			}
			visited = true
		}
	}
	return visited, nil
}

func (n *navigator) visitMatchingEntries(contents []*yaml.Node, key string, visit mapVisitorFn) (bool, error) {

	n.log.Debug("visitMatchingEntries %v in %v", key, contents)
	// value.Content is a concatenated array of key, value,
	// so keys are in the even indexes, values in odd.
	// merge aliases are defined first, but we only want to traverse them
	// if we don't find a match directly on this node first.
	visited, errorVisitedDirectEntries := n.visitDirectMatchingEntries(contents, key, visit)
	if errorVisitedDirectEntries != nil || visited == true || n.followAliases == false {
		return visited, errorVisitedDirectEntries
	}
	// didnt find a match, lets check the aliases.

	return n.visitAliases(contents, key, visit)
}

func (n *navigator) visitAliases(contents []*yaml.Node, key string, visit mapVisitorFn) (bool, error) {
	// merge aliases are defined first, but we only want to traverse them
	// if we don't find a match on this node first.
	// traverse them backwards so that the last alias overrides the preceding.
	// a node can either be
	// an alias to one other node (e.g. <<: *blah)
	// or a sequence of aliases   (e.g. <<: [*blah, *foo])
	n.log.Debug("checking for aliases")
	for index := len(contents) - 2; index >= 0; index = index - 2 {

		if contents[index+1].Kind == yaml.AliasNode {
			valueNode := contents[index+1]
			n.log.Debug("found an alias")
			n.DebugNode(contents[index])
			n.DebugNode(valueNode)

			visitedAlias, errorInAlias := n.visitMatchingEntries(valueNode.Alias.Content, key, visit)
			if visitedAlias == true || errorInAlias != nil {
				return visitedAlias, errorInAlias
			}
		} else if contents[index+1].Kind == yaml.SequenceNode {
			// could be an array of aliases...
			visitedAliasSeq, errorVisitingAliasSeq := n.visitAliasSequence(contents[index+1].Content, key, visit)
			if visitedAliasSeq == true || errorVisitingAliasSeq != nil {
				return visitedAliasSeq, errorVisitingAliasSeq
			}
		}
	}
	n.log.Debug("nope no matching aliases found")
	return false, nil
}

func (n *navigator) visitAliasSequence(possibleAliasArray []*yaml.Node, key string, visit mapVisitorFn) (bool, error) {
	// need to search this backwards too, so that aliases defined last override the preceding.
	for aliasIndex := len(possibleAliasArray) - 1; aliasIndex >= 0; aliasIndex = aliasIndex - 1 {
		child := possibleAliasArray[aliasIndex]
		if child.Kind == yaml.AliasNode {
			n.log.Debug("found an alias")
			n.DebugNode(child)
			visitedAlias, errorInAlias := n.visitMatchingEntries(child.Alias.Content, key, visit)
			if visitedAlias == true || errorInAlias != nil {
				return visitedAlias, errorInAlias
			}
		}
	}
	return false, nil
}

func (n *navigator) matchesKey(key string, actual string) bool {
	var prefixMatch = strings.TrimSuffix(key, "*")
	if prefixMatch != key {
		return strings.HasPrefix(actual, prefixMatch)
	}
	return actual == key
}

func (n *navigator) splatArray(value *yaml.Node, tail []string, visitor VisitorFn, pathStack []interface{}) error {
	for index, childValue := range value.Content {
		n.log.Debug("processing")
		n.DebugNode(childValue)
		childValue = n.getOrReplace(childValue, n.GuessKind(tail, childValue.Kind))
		var err = n.doVisit(childValue, tail, visitor, append(pathStack, index))
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *navigator) appendArray(value *yaml.Node, tail []string, visitor VisitorFn, pathStack []interface{}) error {
	var newNode = yaml.Node{Kind: n.GuessKind(tail, 0)}
	value.Content = append(value.Content, &newNode)
	n.log.Debug("appending a new node, %v", value.Content)
	return n.doVisit(&newNode, tail, visitor, append(pathStack, len(value.Content)-1))
}

func (n *navigator) recurseArray(value *yaml.Node, head string, tail []string, visitor VisitorFn, pathStack []interface{}) error {
	var index, err = strconv.ParseInt(head, 10, 64) // nolint
	if err != nil {
		return err
	}
	if index >= int64(len(value.Content)) {
		return nil
	}
	value.Content[index] = n.getOrReplace(value.Content[index], n.GuessKind(tail, value.Content[index].Kind))
	return n.doVisit(value.Content[index], tail, visitor, append(pathStack, index))
}
