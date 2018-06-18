package main

import "gopkg.in/imdario/mergo.v0"

func merge(dst, src interface{}, overwrite bool) error {
	if overwrite {
		return mergo.MergeWithOverwrite(dst, src)
	}
	return mergo.Merge(dst, src)
}
