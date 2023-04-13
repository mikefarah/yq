package yqlib

import (
	"bytes"

	"github.com/goccy/go-json"
)

func (o *CandidateNode) MarshalJSON() ([]byte, error) {
	log.Debugf("MarshalJSON %v", NodeToString(o))
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetIndent("", " ")
	enc.SetEscapeHTML(false) // do not escape html chars e.g. &, <, >

	switch o.Kind {
	case DocumentNode:
		log.Debugf("MarshalJSON DocumentNode")
		err := enc.Encode(o.Content[0])
		return buf.Bytes(), err
	case AliasNode:
		log.Debugf("MarshalJSON AliasNode")
		err := enc.Encode(o.Alias)
		return buf.Bytes(), err
	case ScalarNode:
		log.Debugf("MarshalJSON ScalarNode")
		value, err := o.GetValueRep()
		if err != nil {
			return buf.Bytes(), err
		}
		err = enc.Encode(value)
		return buf.Bytes(), err
	case MappingNode:
		log.Debugf("MarshalJSON MappingNode")
		buf.WriteByte('{')
		for i := 0; i < len(o.Content); i += 2 {
			if err := enc.Encode(o.Content[i].Value); err != nil {
				return nil, err
			}
			buf.WriteByte(':')
			if err := enc.Encode(o.Content[i+1]); err != nil {
				return nil, err
			}
			if i != len(o.Content)-2 {
				buf.WriteByte(',')
			}
		}
		buf.WriteByte('}')
	case SequenceNode:
		log.Debugf("MarshalJSON SequenceNode")
		err := enc.Encode(o.Content)
		return buf.Bytes(), err
	}
	log.Debug("none of those things?")
	return buf.Bytes(), nil
}
