package yqlib

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

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
		for tok, err = dec.Token(); err == nil; tok, err = dec.Token() {
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
		var res []*orderedMap
		if err := json.Unmarshal(data, &res); err != nil {
			return err
		}
		o.altVal = res
		o.kv = nil
		return nil
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
