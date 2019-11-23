package yqlib

import (
	mergo "gopkg.in/imdario/mergo.v0"
	logging "gopkg.in/op/go-logging.v1"
)

var log = logging.MustGetLogger("yq")

func SetLogger(l *logging.Logger) {
	log = l
}

func ReadPath(dataBucket interface{}, path string) (interface{}, error) {
	var paths = ParsePath(path)
	return Recurse(dataBucket, paths[0], paths[1:])
}

func WritePath(dataBucket interface{}, path string, value interface{}) (interface{}) {
	var paths = ParsePath(path)
	return UpdatedChildValue(dataBucket, paths, value)
}

func PrefixPath(dataBucket interface{}, prefix string) (interface{}) {
	var paths = ParsePath(prefix)

	// Inverse order
	for i := len(paths)/2 - 1; i >= 0; i-- {
		opp := len(paths) - 1 - i
		paths[i], paths[opp] = paths[opp], paths[i]
	}

	var mapDataBucket = dataBucket
	for _, key := range paths {
		singlePath := []string{key}
		mapDataBucket = UpdatedChildValue(nil, singlePath, mapDataBucket)
	}

	return mapDataBucket
}

func DeletePath(dataBucket interface{}, path string) (interface{}, error) {
	var paths = ParsePath(path)
	return DeleteChildValue(dataBucket, paths)
}

func Merge(dst interface{}, src interface{}, overwrite bool, append bool) error {
	if overwrite {
		return mergo.Merge(dst, src, mergo.WithOverride)
	} else if append {
		return mergo.Merge(dst, src, mergo.WithAppendSlice)
	}
	return mergo.Merge(dst, src)
}