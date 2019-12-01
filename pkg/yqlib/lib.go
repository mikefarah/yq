package yqlib

import (
	mergo "gopkg.in/imdario/mergo.v0"
	logging "gopkg.in/op/go-logging.v1"
)

type YqLib interface {
	ReadPath(dataBucket interface{}, path string) (interface{}, error)
	WritePath(dataBucket interface{}, path string, value interface{}) interface{}
	PrefixPath(dataBucket interface{}, prefix string) interface{}
	DeletePath(dataBucket interface{}, path string) (interface{}, error)
	Merge(dst interface{}, src interface{}, overwrite bool, append bool) error
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

func (l *lib) ReadPath(dataBucket interface{}, path string) (interface{}, error) {
	var paths = l.parser.ParsePath(path)
	return l.navigator.ReadChildValue(dataBucket, paths)
}

func (l *lib) WritePath(dataBucket interface{}, path string, value interface{}) interface{} {
	var paths = l.parser.ParsePath(path)
	return l.navigator.UpdatedChildValue(dataBucket, paths, value)
}

func (l *lib) PrefixPath(dataBucket interface{}, prefix string) interface{} {
	var paths = l.parser.ParsePath(prefix)

	// Inverse order
	for i := len(paths)/2 - 1; i >= 0; i-- {
		opp := len(paths) - 1 - i
		paths[i], paths[opp] = paths[opp], paths[i]
	}

	var mapDataBucket = dataBucket
	for _, key := range paths {
		singlePath := []string{key}
		mapDataBucket = l.navigator.UpdatedChildValue(nil, singlePath, mapDataBucket)
	}

	return mapDataBucket
}

func (l *lib) DeletePath(dataBucket interface{}, path string) (interface{}, error) {
	var paths = l.parser.ParsePath(path)
	return l.navigator.DeleteChildValue(dataBucket, paths)
}

func (l *lib) Merge(dst interface{}, src interface{}, overwrite bool, append bool) error {
	if overwrite {
		return mergo.Merge(dst, src, mergo.WithOverride)
	} else if append {
		return mergo.Merge(dst, src, mergo.WithAppendSlice)
	}
	return mergo.Merge(dst, src)
}
