package yqlib

import (
	"fmt"

	logging "gopkg.in/op/go-logging.v1"
	yaml "gopkg.in/yaml.v3"
)

type UpdateCommand struct {
	Command string
	Path    string
	Value   *yaml.Node
}

type YqLib interface {
	DebugNode(node *yaml.Node)
	Get(rootNode *yaml.Node, path string) (*yaml.Node, error)
	Update(rootNode *yaml.Node, updateCommand UpdateCommand) error
	New(path string) yaml.Node
}

type lib struct {
	navigator DataNavigator
	parser    PathParser
	log       *logging.Logger
}

func NewYqLib(l *logging.Logger) YqLib {
	return &lib{
		navigator: NewDataNavigator(l),
		parser:    NewPathParser(),
		log:       l,
	}
}

func (l *lib) DebugNode(node *yaml.Node) {
	l.navigator.DebugNode(node)
}

func (l *lib) Get(rootNode *yaml.Node, path string) (*yaml.Node, error) {
	var paths = l.parser.ParsePath(path)
	return l.navigator.Get(rootNode, paths)
}

func (l *lib) New(path string) yaml.Node {
	var paths = l.parser.ParsePath(path)
	newNode := yaml.Node{Kind: l.navigator.GuessKind(paths, 0)}
	return newNode
}

func (l *lib) Update(rootNode *yaml.Node, updateCommand UpdateCommand) error {
	// later - support other command types
	l.log.Debugf("%v to %v", updateCommand.Command, updateCommand.Path)
	switch updateCommand.Command {
	case "update":
		var paths = l.parser.ParsePath(updateCommand.Path)
		return l.navigator.Update(rootNode, paths, updateCommand.Value)
	case "delete":
		var paths = l.parser.ParsePath(updateCommand.Path)
		return l.navigator.Delete(rootNode, paths)
	default:
		return fmt.Errorf("Unknown command %v", updateCommand.Command)
	}

}
