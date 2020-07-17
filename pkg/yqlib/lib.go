package yqlib

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	logging "gopkg.in/op/go-logging.v1"
	yaml "gopkg.in/yaml.v3"
)

var log = logging.MustGetLogger("yq")

type UpdateCommand struct {
	Command               string
	Path                  string
	Value                 *yaml.Node
	Overwrite             bool
	DontUpdateNodeValue   bool
	DontUpdateNodeContent bool
	CommentsMergeStrategy CommentsMergeStrategy
}

func KindString(kind yaml.Kind) string {
	switch kind {
	case yaml.ScalarNode:
		return "ScalarNode"
	case yaml.SequenceNode:
		return "SequenceNode"
	case yaml.MappingNode:
		return "MappingNode"
	case yaml.DocumentNode:
		return "DocumentNode"
	case yaml.AliasNode:
		return "AliasNode"
	default:
		return "unknown!"
	}
}

func DebugNode(value *yaml.Node) {
	if value == nil {
		log.Debug("-- node is nil --")
	} else if log.IsEnabledFor(logging.DEBUG) {
		buf := new(bytes.Buffer)
		encoder := yaml.NewEncoder(buf)
		errorEncoding := encoder.Encode(value)
		if errorEncoding != nil {
			log.Error("Error debugging node, %v", errorEncoding.Error())
		}
		encoder.Close()
		log.Debug("Tag: %v, Kind: %v, Anchor: %v", value.Tag, KindString(value.Kind), value.Anchor)
		log.Debug("Head Comment: %v", value.HeadComment)
		log.Debug("Line Comment: %v", value.LineComment)
		log.Debug("FootComment Comment: %v", value.FootComment)
		log.Debug("\n%v", buf.String())
	}
}

func pathStackToString(pathStack []interface{}) string {
	return mergePathStackToString(pathStack, UpdateArrayMergeStrategy)
}

func mergePathStackToString(pathStack []interface{}, arrayMergeStrategy ArrayMergeStrategy) string {
	var sb strings.Builder
	for index, path := range pathStack {
		switch path.(type) {
		case int, int64:
			if arrayMergeStrategy == AppendArrayMergeStrategy {
				sb.WriteString("[+]")
			} else {
				sb.WriteString(fmt.Sprintf("[%v]", path))
			}

		default:
			s := fmt.Sprintf("%v", path)
			var _, errParsingInt = strconv.ParseInt(s, 10, 64) // nolint

			hasSpecial := strings.Contains(s, ".") || strings.Contains(s, "[") || strings.Contains(s, "]") || strings.Contains(s, "\"")
			hasDoubleQuotes := strings.Contains(s, "\"")
			wrappingCharacterStart := "\""
			wrappingCharacterEnd := "\""
			if hasDoubleQuotes {
				wrappingCharacterStart = "("
				wrappingCharacterEnd = ")"
			}
			if hasSpecial || errParsingInt == nil {
				sb.WriteString(wrappingCharacterStart)
			}
			sb.WriteString(s)
			if hasSpecial || errParsingInt == nil {
				sb.WriteString(wrappingCharacterEnd)
			}
		}

		if index < len(pathStack)-1 {
			sb.WriteString(".")
		}
	}
	return sb.String()
}

func guessKind(head interface{}, tail []interface{}, guess yaml.Kind) yaml.Kind {
	log.Debug("guessKind: tail %v", tail)
	if len(tail) == 0 && guess == 0 {
		log.Debug("end of path, must be a scalar")
		return yaml.ScalarNode
	} else if len(tail) == 0 {
		return guess
	}
	var next = tail[0]
	switch next.(type) {
	case int64:
		return yaml.SequenceNode
	default:
		var nextString = fmt.Sprintf("%v", next)
		if nextString == "+" {
			return yaml.SequenceNode
		}
		pathParser := NewPathParser()
		if pathParser.IsPathExpression(nextString) && (guess == yaml.SequenceNode || guess == yaml.MappingNode) {
			return guess
		} else if guess == yaml.AliasNode {
			log.Debug("guess was an alias, okey doke.")
			return guess
		} else if head == "**" {
			log.Debug("deep wildcard, go with the guess")
			return guess
		}
		log.Debug("forcing a mapping node")
		log.Debug("yaml.SequenceNode %v", guess == yaml.SequenceNode)
		log.Debug("yaml.ScalarNode %v", guess == yaml.ScalarNode)
		return yaml.MappingNode
	}
}

type YqLib interface {
	Get(rootNode *yaml.Node, path string) ([]*NodeContext, error)
	GetForMerge(rootNode *yaml.Node, path string, arrayMergeStrategy ArrayMergeStrategy) ([]*NodeContext, error)
	Update(rootNode *yaml.Node, updateCommand UpdateCommand, autoCreate bool) error
	New(path string) yaml.Node

	PathStackToString(pathStack []interface{}) string
	MergePathStackToString(pathStack []interface{}, arrayMergeStrategy ArrayMergeStrategy) string
}

type lib struct {
	parser PathParser
}

func NewYqLib() YqLib {
	return &lib{
		parser: NewPathParser(),
	}
}

func (l *lib) Get(rootNode *yaml.Node, path string) ([]*NodeContext, error) {
	var paths = l.parser.ParsePath(path)
	navigationStrategy := ReadNavigationStrategy()
	navigator := NewDataNavigator(navigationStrategy)
	error := navigator.Traverse(rootNode, paths)
	return navigationStrategy.GetVisitedNodes(), error
}

func (l *lib) GetForMerge(rootNode *yaml.Node, path string, arrayMergeStrategy ArrayMergeStrategy) ([]*NodeContext, error) {
	var paths = l.parser.ParsePath(path)
	navigationStrategy := ReadForMergeNavigationStrategy(arrayMergeStrategy)
	navigator := NewDataNavigator(navigationStrategy)
	error := navigator.Traverse(rootNode, paths)
	return navigationStrategy.GetVisitedNodes(), error
}

func (l *lib) PathStackToString(pathStack []interface{}) string {
	return pathStackToString(pathStack)
}

func (l *lib) MergePathStackToString(pathStack []interface{}, arrayMergeStrategy ArrayMergeStrategy) string {
	return mergePathStackToString(pathStack, arrayMergeStrategy)
}

func (l *lib) New(path string) yaml.Node {
	var paths = l.parser.ParsePath(path)
	newNode := yaml.Node{Kind: guessKind("", paths, 0)}
	return newNode
}

func (l *lib) Update(rootNode *yaml.Node, updateCommand UpdateCommand, autoCreate bool) error {
	log.Debugf("%v to %v", updateCommand.Command, updateCommand.Path)
	switch updateCommand.Command {
	case "update":
		var paths = l.parser.ParsePath(updateCommand.Path)
		navigator := NewDataNavigator(UpdateNavigationStrategy(updateCommand, autoCreate))
		return navigator.Traverse(rootNode, paths)
	case "merge":
		var paths = l.parser.ParsePath(updateCommand.Path)
		navigator := NewDataNavigator(MergeNavigationStrategy(updateCommand, autoCreate))
		return navigator.Traverse(rootNode, paths)
	case "delete":
		var paths = l.parser.ParsePath(updateCommand.Path)
		lastBit, newTail := paths[len(paths)-1], paths[:len(paths)-1]
		navigator := NewDataNavigator(DeleteNavigationStrategy(lastBit))
		return navigator.Traverse(rootNode, newTail)
	default:
		return fmt.Errorf("Unknown command %v", updateCommand.Command)
	}

}
