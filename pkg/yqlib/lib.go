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

// TODO: enumerate
type UpdateCommand struct {
	Command   string
	Path      string
	Value     *yaml.Node
	Overwrite bool
}

func DebugNode(value *yaml.Node) {
	if value == nil {
		log.Debug("-- node is nil --")
	} else if log.IsEnabledFor(logging.DEBUG) {
		buf := new(bytes.Buffer)
		encoder := yaml.NewEncoder(buf)
		encoder.Encode(value)
		encoder.Close()
		log.Debug("Tag: %v", value.Tag)
		log.Debug("%v", buf.String())
	}
}

func PathStackToString(pathStack []interface{}) string {
	return MergePathStackToString(pathStack, false)
}

func MergePathStackToString(pathStack []interface{}, appendArrays bool) string {
	var sb strings.Builder
	for index, path := range pathStack {
		switch path.(type) {
		case int:
			if appendArrays {
				sb.WriteString("[+]")
			} else {
				sb.WriteString(fmt.Sprintf("[%v]", path))
			}

		default:
			sb.WriteString(fmt.Sprintf("%v", path))
		}

		if index < len(pathStack)-1 {
			sb.WriteString(".")
		}
	}
	return sb.String()
}

func guessKind(head string, tail []string, guess yaml.Kind) yaml.Kind {
	log.Debug("tail %v", tail)
	if len(tail) == 0 && guess == 0 {
		log.Debug("end of path, must be a scalar")
		return yaml.ScalarNode
	} else if len(tail) == 0 {
		return guess
	}

	var _, errorParsingInt = strconv.ParseInt(tail[0], 10, 64)
	if tail[0] == "+" || errorParsingInt == nil {
		return yaml.SequenceNode
	}
	if (tail[0] == "*" || tail[0] == "**" || head == "**") && (guess == yaml.SequenceNode || guess == yaml.MappingNode) {
		return guess
	}
	if guess == yaml.AliasNode {
		log.Debug("guess was an alias, okey doke.")
		return guess
	}
	log.Debug("forcing a mapping node")
	log.Debug("yaml.SequenceNode %v", guess == yaml.SequenceNode)
	log.Debug("yaml.ScalarNode %v", guess == yaml.ScalarNode)
	return yaml.MappingNode
}

type YqLib interface {
	Get(rootNode *yaml.Node, path string) ([]*NodeContext, error)
	Update(rootNode *yaml.Node, updateCommand UpdateCommand, autoCreate bool) error
	New(path string) yaml.Node
}

type lib struct {
	navigator DataNavigator
	parser    PathParser
}

func NewYqLib() YqLib {
	return &lib{
		parser: NewPathParser(),
	}
}

func (l *lib) Get(rootNode *yaml.Node, path string) ([]*NodeContext, error) {
	var paths = l.parser.ParsePath(path)
	NavigationStrategy := ReadNavigationStrategy()
	navigator := NewDataNavigator(NavigationStrategy)
	error := navigator.Traverse(rootNode, paths)
	return NavigationStrategy.GetVisitedNodes(), error

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
	case "delete":
		var paths = l.parser.ParsePath(updateCommand.Path)
		lastBit, newTail := paths[len(paths)-1], paths[:len(paths)-1]
		navigator := NewDataNavigator(DeleteNavigationStrategy(lastBit))
		return navigator.Traverse(rootNode, newTail)
	default:
		return fmt.Errorf("Unknown command %v", updateCommand.Command)
	}

}
