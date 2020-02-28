package yqlib

import (
	"fmt"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type PathParser interface {
	ParsePath(path string) []interface{}
	MatchesNextPathElement(nodeContext NodeContext, nodeKey string) bool
	IsPathExpression(pathElement string) bool
}

type pathParser struct{}

func NewPathParser() PathParser {
	return &pathParser{}
}

func matchesString(expression string, value string) bool {
	var prefixMatch = strings.TrimSuffix(expression, "*")
	if prefixMatch != expression {
		log.Debug("prefix match, %v", strings.HasPrefix(value, prefixMatch))
		return strings.HasPrefix(value, prefixMatch)
	}
	return value == expression
}

func (p *pathParser) IsPathExpression(pathElement string) bool {
	return pathElement == "*" || pathElement == "**" || strings.Contains(pathElement, "==")
}

/**
 * node: node that we may traverse/visit
 * head: path element expression to match against
 * tail: remaining path element expressions
 * pathStack: stack of actual paths we've matched to get to node
 * nodeKey: actual value of this nodes 'key' or index.
 */
func (p *pathParser) MatchesNextPathElement(nodeContext NodeContext, nodeKey string) bool {
	head := nodeContext.Head
	if head == "**" || head == "*" {
		return true
	}
	var headString = fmt.Sprintf("%v", head)

	if strings.Contains(headString, "==") && nodeContext.Node.Kind != yaml.ScalarNode {
		log.Debug("ooh deep recursion time")
		result := strings.SplitN(headString, "==", 2)
		path := strings.TrimSpace(result[0])
		value := strings.TrimSpace(result[1])
		log.Debug("path %v", path)
		log.Debug("value %v", value)
		DebugNode(nodeContext.Node)
		navigationStrategy := FilterMatchingNodesNavigationStrategy(value)

		navigator := NewDataNavigator(navigationStrategy)
		err := navigator.Traverse(nodeContext.Node, p.ParsePath(path))
		if err != nil {
			log.Error("Error deep recursing - ignoring")
			log.Error(err.Error())
		}
		log.Debug("done deep recursing, found %v matches", len(navigationStrategy.GetVisitedNodes()))
		return len(navigationStrategy.GetVisitedNodes()) > 0
	} else if strings.Contains(headString, "==") && nodeContext.Node.Kind == yaml.ScalarNode {
		result := strings.SplitN(headString, "==", 2)
		path := strings.TrimSpace(result[0])
		value := strings.TrimSpace(result[1])
		if path == "." {
			log.Debug("need to match scalar")
			return matchesString(value, nodeContext.Node.Value)
		}
	}

	if head == "+" {
		log.Debug("head is +, nodeKey is %v", nodeKey)
		var _, err = strconv.ParseInt(nodeKey, 10, 64) // nolint
		if err == nil {
			return true
		}
	}

	return matchesString(headString, nodeKey)
}

func (p *pathParser) ParsePath(path string) []interface{} {
	var paths = make([]interface{}, 0)
	if path == "" {
		return paths
	}
	return p.parsePathAccum(paths, path)
}

func (p *pathParser) parsePathAccum(paths []interface{}, remaining string) []interface{} {
	head, tail := p.nextYamlPath(remaining)
	if tail == "" {
		return append(paths, head)
	}
	return p.parsePathAccum(append(paths, head), tail)
}

func (p *pathParser) nextYamlPath(path string) (pathElement interface{}, remaining string) {
	switch path[0] {
	case '[':
		// e.g [0].blah.cat -> we need to return "0" and "blah.cat"
		var value, remainingBit = p.search(path[1:], []uint8{']'}, true)
		var number, errParsingInt = strconv.ParseInt(value, 10, 64) // nolint
		if errParsingInt == nil {
			return number, remainingBit
		}
		return value, remainingBit
	case '"':
		// e.g "a.b".blah.cat -> we need to return "a.b" and "blah.cat"
		return p.search(path[1:], []uint8{'"'}, true)
	case '(':
		// e.g "a.b".blah.cat -> we need to return "a.b" and "blah.cat"
		return p.search(path[1:], []uint8{')'}, true)
	default:
		// e.g "a.blah.cat" -> return "a" and "blah.cat"
		return p.search(path[0:], []uint8{'.', '[', '"', '('}, false)
	}
}

func (p *pathParser) search(path string, matchingChars []uint8, skipNext bool) (pathElement string, remaining string) {
	for i := 0; i < len(path); i++ {
		var char = path[i]
		if p.contains(matchingChars, char) {
			var remainingStart = i + 1
			if skipNext {
				remainingStart = remainingStart + 1
			} else if !skipNext && char != '.' {
				remainingStart = i
			}
			if remainingStart > len(path) {
				remainingStart = len(path)
			}
			return path[0:i], path[remainingStart:]
		}
	}
	return path, ""
}

func (p *pathParser) contains(matchingChars []uint8, candidate uint8) bool {
	for _, a := range matchingChars {
		if a == candidate {
			return true
		}
	}
	return false
}
