package yqlib

import (
	"bytes"
	"fmt"
	"strconv"

	logging "gopkg.in/op/go-logging.v1"
	yaml "gopkg.in/yaml.v3"
)

var log = logging.MustGetLogger("yq")

type UpdateCommand struct {
	Command string
	Path    string
	Value   *yaml.Node
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

func guessKind(tail []string, guess yaml.Kind) yaml.Kind {
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
	if tail[0] == "*" && (guess == yaml.SequenceNode || guess == yaml.MappingNode) {
		return guess
	}
	if guess == yaml.AliasNode {
		log.Debug("guess was an alias, okey doke.")
		return guess
	}
	log.Debug("forcing a mapping node")
	log.Debug("yaml.SequenceNode ?", guess == yaml.SequenceNode)
	log.Debug("yaml.ScalarNode ?", guess == yaml.ScalarNode)
	return yaml.MappingNode
}

type YqLib interface {
	Get(rootNode *yaml.Node, path string) ([]*VisitedNode, error)
	Update(rootNode *yaml.Node, updateCommand UpdateCommand) error
	New(path string) yaml.Node
}

type lib struct {
	navigator DataNavigator
	parser    PathParser
}

func NewYqLib(l *logging.Logger) YqLib {
	return &lib{
		parser: NewPathParser(),
	}
}

func (l *lib) Get(rootNode *yaml.Node, path string) ([]*VisitedNode, error) {
	var paths = l.parser.ParsePath(path)
	navigationSettings := ReadNavigationSettings()
	navigator := NewDataNavigator(navigationSettings)
	error := navigator.Traverse(rootNode, paths)
	return navigationSettings.GetVisitedNodes(), error

}

func (l *lib) New(path string) yaml.Node {
	var paths = l.parser.ParsePath(path)
	newNode := yaml.Node{Kind: guessKind(paths, 0)}
	return newNode
}

func (l *lib) Update(rootNode *yaml.Node, updateCommand UpdateCommand) error {
	log.Debugf("%v to %v", updateCommand.Command, updateCommand.Path)
	switch updateCommand.Command {
	case "update":
		var paths = l.parser.ParsePath(updateCommand.Path)
		navigator := NewDataNavigator(UpdateNavigationSettings(updateCommand.Value))
		return navigator.Traverse(rootNode, paths)
	case "delete":
		var paths = l.parser.ParsePath(updateCommand.Path)
		lastBit, newTail := paths[len(paths)-1], paths[:len(paths)-1]
		navigator := NewDataNavigator(DeleteNavigationSettings(lastBit))
		return navigator.Traverse(rootNode, newTail)
	default:
		return fmt.Errorf("Unknown command %v", updateCommand.Command)
	}

}
