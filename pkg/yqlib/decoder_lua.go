//go:build !yq_nolua

package yqlib

import (
	"fmt"
	"io"
	"math"

	lua "github.com/yuin/gopher-lua"
)

type luaDecoder struct {
	reader   io.Reader
	finished bool
	prefs    LuaPreferences
}

func NewLuaDecoder(prefs LuaPreferences) Decoder {
	return &luaDecoder{
		prefs: prefs,
	}
}

func (dec *luaDecoder) Init(reader io.Reader) error {
	dec.reader = reader
	return nil
}

func (dec *luaDecoder) convertToYamlNode(ls *lua.LState, lv lua.LValue) *CandidateNode {
	switch lv.Type() {
	case lua.LTNil:
		return &CandidateNode{
			Kind:  ScalarNode,
			Tag:   "!!null",
			Value: "",
		}
	case lua.LTBool:
		return &CandidateNode{
			Kind:  ScalarNode,
			Tag:   "!!bool",
			Value: lv.String(),
		}
	case lua.LTNumber:
		n := float64(lua.LVAsNumber(lv))
		// various special case floats
		if math.IsNaN(n) {
			return &CandidateNode{
				Kind:  ScalarNode,
				Tag:   "!!float",
				Value: ".nan",
			}
		}
		if math.IsInf(n, 1) {
			return &CandidateNode{
				Kind:  ScalarNode,
				Tag:   "!!float",
				Value: ".inf",
			}
		}
		if math.IsInf(n, -1) {
			return &CandidateNode{
				Kind:  ScalarNode,
				Tag:   "!!float",
				Value: "-.inf",
			}
		}

		// does it look like an integer?
		if n == float64(int(n)) {
			return &CandidateNode{
				Kind:  ScalarNode,
				Tag:   "!!int",
				Value: lv.String(),
			}
		}

		return &CandidateNode{
			Kind:  ScalarNode,
			Tag:   "!!float",
			Value: lv.String(),
		}
	case lua.LTString:
		return &CandidateNode{
			Kind:  ScalarNode,
			Tag:   "!!str",
			Value: lv.String(),
		}
	case lua.LTFunction:
		return &CandidateNode{
			Kind:  ScalarNode,
			Tag:   "tag:lua.org,2006,function",
			Value: lv.String(),
		}
	case lua.LTTable:
		// Simultaneously create a sequence and a map, pick which one to return
		// based on whether all keys were consecutive integers
		i := 1
		yaml_sequence := &CandidateNode{
			Kind: SequenceNode,
			Tag:  "!!seq",
		}
		yaml_map := &CandidateNode{
			Kind: MappingNode,
			Tag:  "!!map",
		}
		t := lv.(*lua.LTable)
		k, v := ls.Next(t, lua.LNil)
		for k != lua.LNil {
			if ki, ok := k.(lua.LNumber); i != 0 && ok && math.Mod(float64(ki), 1) == 0 && int(ki) == i {
				i++
			} else {
				i = 0
			}
			newKey := dec.convertToYamlNode(ls, k)

			yv := dec.convertToYamlNode(ls, v)
			yaml_map.AddKeyValueChild(newKey, yv)

			if i != 0 {
				yaml_sequence.AddChild(yv)
			}
			k, v = ls.Next(t, k)
		}
		if i != 0 {
			return yaml_sequence
		}
		return yaml_map
	default:
		return &CandidateNode{
			Kind:        ScalarNode,
			LineComment: fmt.Sprintf("Unhandled Lua type: %s", lv.Type().String()),
			Tag:         "!!null",
			Value:       lv.String(),
		}
	}
}

func (dec *luaDecoder) decideTopLevelNode(ls *lua.LState) *CandidateNode {
	if ls.GetTop() == 0 {
		// no items were explicitly returned, encode the globals table instead
		return dec.convertToYamlNode(ls, ls.Get(lua.GlobalsIndex))
	}
	return dec.convertToYamlNode(ls, ls.Get(1))
}

func (dec *luaDecoder) Decode() (*CandidateNode, error) {
	if dec.finished {
		return nil, io.EOF
	}
	ls := lua.NewState(lua.Options{SkipOpenLibs: true})
	defer ls.Close()
	fn, err := ls.Load(dec.reader, "@input")
	if err != nil {
		return nil, err
	}
	ls.Push(fn)
	err = ls.PCall(0, lua.MultRet, nil)
	if err != nil {
		return nil, err
	}
	firstNode := dec.decideTopLevelNode(ls)
	dec.finished = true
	return firstNode, nil
}
