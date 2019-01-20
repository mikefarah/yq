package main

import mergo "gopkg.in/imdario/mergo.v0"

func merge(dst interface{}, src interface{}, overwrite bool, append bool) error {
	if overwrite {
		return mergo.Merge(dst, src, mergo.WithOverride)
	} else if append {
		return mergo.Merge(dst, src, mergo.WithAppendSlice)
	}
	return mergo.Merge(dst, src)
}
