package main

import (
	"github.com/imdario/mergo"
)

func merge(dst, src interface{}, overwrite bool) error {
	if overwrite {
		return mergo.MergeWithOverwrite(dst, src)
	}
	return mergo.Merge(dst, src)
}
