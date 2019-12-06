package yqlib

import (
	logging "gopkg.in/op/go-logging.v1"
	yaml "gopkg.in/yaml.v3"
)

type YqLib interface {
	Get(rootNode *yaml.Node, path string) (*yaml.Node, error)
	Update(rootNode *yaml.Node, path string, writeCommand WriteCommand) error
}

type lib struct {
	navigator DataNavigator
	parser    PathParser
}

func NewYqLib(l *logging.Logger) YqLib {
	return &lib{
		navigator: NewDataNavigator(l),
		parser:    NewPathParser(),
	}
}

func (l *lib) Get(rootNode *yaml.Node, path string) (*yaml.Node, error) {
	if path == "" {
		return rootNode, nil
	}
	var paths = l.parser.ParsePath(path)
	return l.navigator.Get(rootNode, paths)
}

func (l *lib) Update(rootNode *yaml.Node, path string, writeCommand WriteCommand) error {
	var paths = l.parser.ParsePath(path)
	return l.navigator.Update(rootNode, paths, writeCommand)
}
