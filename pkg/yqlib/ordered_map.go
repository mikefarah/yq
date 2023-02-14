package yqlib

// orderedMap allows to marshal and unmarshal JSON and YAML values keeping the
// order of keys and values in a map or an object.
type orderedMap struct {
	// if this is an object, kv != nil. If this is not an object, kv == nil.
	kv     []orderedMapKV
	altVal interface{}
}

type orderedMapKV struct {
	K string
	V orderedMap
}
