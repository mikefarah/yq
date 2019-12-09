package yqlib

import (
	"fmt"

	logging "gopkg.in/op/go-logging.v1"
	yaml "gopkg.in/yaml.v3"
)

type UpdateCommand struct {
	Command string
	Path    string
	Value   yaml.Node
}

type YqLib interface {
	DebugNode(node *yaml.Node)
	Get(rootNode *yaml.Node, path string) (*yaml.Node, error)
	Update(rootNode *yaml.Node, updateCommand UpdateCommand) error
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

func (l *lib) Update(rootNode *yaml.Node, updateCommand UpdateCommand) error {
	// later - support other command types
	l.log.Debugf("%v to %v", updateCommand.Command, updateCommand.Path)
	switch updateCommand.Command {
	case "update":
		var paths = l.parser.ParsePath(updateCommand.Path)
		return l.navigator.Update(rootNode, paths, updateCommand.Value)
	case "delete":
		l.log.Debugf("need to implement delete")
		return nil
	default:
		return fmt.Errorf("Unknown command %v", updateCommand.Command)
	}

}
