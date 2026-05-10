//go:build !yq_nojson

package yqlib

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
)

func (o *CandidateNode) setScalarFromJson(value interface{}) error {
	o.Kind = ScalarNode
	switch rawData := value.(type) {
	case nil:
		o.Tag = "!!null"
		o.Value = "null"
	case float64, float32:
		o.Value = fmt.Sprintf("%v", value)
		o.Tag = "!!float"
		// json decoder returns ints as float.
		if value == float64(int64(rawData.(float64))) {
			// aha it's an int disguised as a float
			o.Tag = "!!int"
			o.Value = fmt.Sprintf("%v", int64(value.(float64)))
		}
	case int, int64, int32:
		o.Value = fmt.Sprintf("%v", value)
		o.Tag = "!!int"
	case bool:
		o.Value = fmt.Sprintf("%v", value)
		o.Tag = "!!bool"
	case string:
		o.Value = rawData
		o.Tag = "!!str"
	default:
		return fmt.Errorf("unrecognised type :( %v", rawData)
	}
	return nil
}

func (o *CandidateNode) UnmarshalJSON(data []byte) error {
	log.Debug("UnmarshalJSON")
	switch data[0] {
	case '{':
		log.Debug("UnmarshalJSON -  its a map!")
		// its a map
		o.Kind = MappingNode
		o.Tag = "!!map"

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

			childKey := o.CreateChild()
			childKey.IsMapKey = true
			childKey.Value = tok.(string)
			childKey.Kind = ScalarNode
			childKey.Tag = "!!str"

			childValue := o.CreateChild()
			childValue.Key = childKey

			if err := dec.Decode(childValue); err != nil {
				return err
			}
			o.Content = append(o.Content, childKey, childValue)
		}
		// unexpected error
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		return nil
	case '[':
		o.Kind = SequenceNode
		o.Tag = "!!seq"
		log.Debug("UnmarshalJSON -  its an array!")
		var children []*CandidateNode
		if err := json.Unmarshal(data, &children); err != nil {
			return err
		}
		// now we put the children into the content, and set a key value for them
		for i, child := range children {

			if child == nil {
				// need to represent it as a null scalar
				child = createScalarNode(nil, "null")
			}
			childKey := o.CreateChild()
			childKey.Kind = ScalarNode
			childKey.Tag = "!!int"
			childKey.Value = fmt.Sprintf("%v", i)
			childKey.IsMapKey = true

			child.Parent = o
			child.Key = childKey
			o.Content = append(o.Content, child)
		}
		return nil
	}
	log.Debug("UnmarshalJSON -  its a scalar!")
	// otherwise, must be a scalar
	var scalar interface{}
	err := json.Unmarshal(data, &scalar)

	if err != nil {
		return err
	}
	log.Debugf("UnmarshalJSON -  scalar is %v", scalar)

	return o.setScalarFromJson(scalar)

}

func (o *CandidateNode) MarshalJSON() ([]byte, error) {
	log.Debugf("MarshalJSON %v", NodeToString(o))
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetIndent("", " ")
	enc.SetEscapeHTML(false) // do not escape html chars e.g. &, <, >

	switch o.Kind {
	case AliasNode:
		log.Debugf("MarshalJSON AliasNode")
		err := enc.Encode(o.Alias)
		return buf.Bytes(), err
	case ScalarNode:
		log.Debugf("MarshalJSON ScalarNode")
		if o.guessTagFromCustomType() == "!!float" {
			if raw, ok := jsonFloatLiteral(o.Value); ok {
				buf.WriteString(raw)
				return buf.Bytes(), nil
			}
		}
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
		return buf.Bytes(), nil
	case SequenceNode:
		log.Debugf("MarshalJSON SequenceNode, %v, len: %v", o.Content, len(o.Content))
		var err error
		if len(o.Content) == 0 {
			buf.WriteString("[]")
		} else {
			err = enc.Encode(o.Content)
		}
		return buf.Bytes(), err
	default:
		err := enc.Encode(nil)
		return buf.Bytes(), err
	}
}

// jsonFloatLiteral returns a JSON-shaped representation of a YAML !!float scalar
// value, preserving the original textual form (e.g. "50.0" stays "50.0") whenever
// possible. The second return value is false when the value cannot be safely
// rendered as a JSON number (e.g. ".inf", ".nan", or anything that parses to a
// non-finite float); callers should fall back to the normal encoding path in
// that case, which preserves the existing behaviour for those inputs.
func jsonFloatLiteral(raw string) (string, bool) {
	if raw == "" {
		return "", false
	}
	f, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return "", false
	}
	if math.IsInf(f, 0) || math.IsNaN(f) {
		return "", false
	}
	if isJSONNumberLiteral(raw) {
		return raw, true
	}
	formatted := strconv.FormatFloat(f, 'f', -1, 64)
	if !strings.ContainsAny(formatted, ".eE") {
		formatted += ".0"
	}
	return formatted, true
}

// isJSONNumberLiteral reports whether s is already a valid JSON number literal
// representing a fractional value (i.e. contains a "." or an exponent), so it
// can be emitted verbatim without round-tripping through a float64.
func isJSONNumberLiteral(s string) bool {
	if s == "" {
		return false
	}
	i := 0
	if s[i] == '-' {
		i++
		if i == len(s) {
			return false
		}
	}
	// integer part: 0 or [1-9][0-9]*
	if s[i] == '0' {
		i++
	} else if s[i] >= '1' && s[i] <= '9' {
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
		}
	} else {
		return false
	}
	hasFraction := false
	if i < len(s) && s[i] == '.' {
		hasFraction = true
		i++
		if i == len(s) || s[i] < '0' || s[i] > '9' {
			return false
		}
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
		}
	}
	hasExponent := false
	if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		hasExponent = true
		i++
		if i < len(s) && (s[i] == '+' || s[i] == '-') {
			i++
		}
		if i == len(s) || s[i] < '0' || s[i] > '9' {
			return false
		}
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
		}
	}
	if i != len(s) {
		return false
	}
	return hasFraction || hasExponent
}
