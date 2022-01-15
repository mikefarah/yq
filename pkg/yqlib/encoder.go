package yqlib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	yaml "gopkg.in/yaml.v3"
)

type Encoder interface {
	Encode(writer io.Writer, node *yaml.Node) error
	PrintDocumentSeparator(writer io.Writer) error
	PrintLeadingContent(writer io.Writer, content string) error
	CanHandleAliases() bool
}

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

func (o *orderedMap) UnmarshalJSON(data []byte) error {
	switch data[0] {
	case '{':
		// initialise so that even if the object is empty it is not nil
		o.kv = []orderedMapKV{}

		// create decoder
		dec := json.NewDecoder(bytes.NewReader(data))
		_, err := dec.Token() // open object
		if err != nil {
			return err
		}

		// cycle through k/v
		var tok json.Token
		for tok, err = dec.Token(); !errors.Is(err, io.EOF); tok, err = dec.Token() {
			// we can expect two types: string or Delim. Delim automatically means
			// that it is the closing bracket of the object, whereas string means
			// that there is another key.
			if _, ok := tok.(json.Delim); ok {
				break
			}
			kv := orderedMapKV{
				K: tok.(string),
			}
			if err := dec.Decode(&kv.V); err != nil {
				return err
			}
			o.kv = append(o.kv, kv)
		}
		// unexpected error
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		return nil
	case '[':
		var arr []orderedMap
		return json.Unmarshal(data, &arr)
	}

	return json.Unmarshal(data, &o.altVal)
}

func (o orderedMap) MarshalJSON() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false) // do not escape html chars e.g. &, <, >
	if o.kv == nil {
		if err := enc.Encode(o.altVal); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	buf.WriteByte('{')
	for idx, el := range o.kv {
		if err := enc.Encode(el.K); err != nil {
			return nil, err
		}
		buf.WriteByte(':')
		if err := enc.Encode(el.V); err != nil {
			return nil, err
		}
		if idx != len(o.kv)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (o *orderedMap) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.DocumentNode:
		if len(node.Content) == 0 {
			return nil
		}
		return o.UnmarshalYAML(node.Content[0])
	case yaml.AliasNode:
		return o.UnmarshalYAML(node.Alias)
	case yaml.ScalarNode:
		return node.Decode(&o.altVal)
	case yaml.MappingNode:
		// set kv to non-nil
		o.kv = []orderedMapKV{}
		for i := 0; i < len(node.Content); i += 2 {
			var key string
			var val orderedMap
			if err := node.Content[i].Decode(&key); err != nil {
				return err
			}
			if err := node.Content[i+1].Decode(&val); err != nil {
				return err
			}
			o.kv = append(o.kv, orderedMapKV{
				K: key,
				V: val,
			})
		}
		return nil
	case yaml.SequenceNode:
		// note that this has to be a pointer, so that nulls can be represented.
		var res []*orderedMap
		if err := node.Decode(&res); err != nil {
			return err
		}
		o.altVal = res
		o.kv = nil
		return nil
	case 0:
		// null
		o.kv = nil
		o.altVal = nil
		return nil
	default:
		return fmt.Errorf("orderedMap: invalid yaml node")
	}
}

func (o *orderedMap) MarshalYAML() (interface{}, error) {
	// fast path: kv is nil, use altVal
	if o.kv == nil {
		return o.altVal, nil
	}
	content := make([]*yaml.Node, 0, len(o.kv)*2)
	for _, val := range o.kv {
		n := new(yaml.Node)
		if err := n.Encode(val.V); err != nil {
			return nil, err
		}
		content = append(content, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Tag:   "!!str",
			Value: val.K,
		}, n)
	}
	return &yaml.Node{
		Kind:    yaml.MappingNode,
		Tag:     "!!map",
		Content: content,
	}, nil
}
