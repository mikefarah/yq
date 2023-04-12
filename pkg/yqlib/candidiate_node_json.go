package yqlib

import (
	"bytes"

	"github.com/goccy/go-json"
)

func (o *CandidateNode) MarshalJSON() ([]byte, error) {
	log.Debugf("going to encode %v - %v", o.GetNicePath(), o.Tag)
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetIndent("", " ")
	enc.SetEscapeHTML(false) // do not escape html chars e.g. &, <, >

	switch o.Kind {
	case DocumentNode:
		err := enc.Encode(o.Content[0])
		return buf.Bytes(), err
	case AliasNode:
		err := enc.Encode(o.Alias)
		return buf.Bytes(), err
	case ScalarNode:
		value, err := o.GetValueRep()
		if err != nil {
			return buf.Bytes(), err
		}
		err = enc.Encode(value)
		return buf.Bytes(), err
	case MappingNode:
		buf.WriteByte('{')
		for i := 0; i < len(o.Content); i += 2 {
			log.Debugf("writing key %v", NodeToString(o.Content[i]))
			if err := enc.Encode(o.Content[i].Value); err != nil {
				return nil, err
			}
			buf.WriteByte(':')
			log.Debugf("writing value %v", NodeToString(o.Content[i+1]))
			if err := enc.Encode(o.Content[i+1]); err != nil {
				return nil, err
			}
			if i != len(o.Content)-2 {
				buf.WriteByte(',')
			}
		}
		buf.WriteByte('}')
	case SequenceNode:
		err := enc.Encode(o.Content)
		return buf.Bytes(), err
	}
	return buf.Bytes(), nil
}
