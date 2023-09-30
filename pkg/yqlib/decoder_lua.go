package yqlib

import (
	"fmt"
	"io"
	"math"

	lua "github.com/yuin/gopher-lua"
	yaml "gopkg.in/yaml.v3"
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

func (dec *luaDecoder) convertToYamlNode(ls *lua.LState, lv lua.LValue) *yaml.Node {
	switch lv.Type() {
	case lua.LTNil:
		return &yaml.Node{
			Kind:  yaml.ScalarNode,
			Tag:   "!!null",
			Value: "",
		}
	case lua.LTBool:
		return &yaml.Node{
			Kind:  yaml.ScalarNode,
			Tag:   "!!bool",
			Value: lv.String(),
		}
	case lua.LTNumber:
		n := float64(lua.LVAsNumber(lv))
		// various special case floats
		if math.IsNaN(n) {
			return &yaml.Node{
				Kind:  yaml.ScalarNode,
				Tag:   "!!float",
				Value: ".nan",
			}
		}
		if math.IsInf(n, 1) {
			return &yaml.Node{
				Kind:  yaml.ScalarNode,
				Tag:   "!!float",
				Value: ".inf",
			}
		}
		if math.IsInf(n, -1) {
			return &yaml.Node{
				Kind:  yaml.ScalarNode,
				Tag:   "!!float",
				Value: "-.inf",
			}
		}

		// does it look like an integer?
		if n == float64(int(n)) {
			return &yaml.Node{
				Kind:  yaml.ScalarNode,
				Tag:   "!!int",
				Value: lv.String(),
			}
		}

		return &yaml.Node{
			Kind:  yaml.ScalarNode,
			Tag:   "!!float",
			Value: lv.String(),
		}
	case lua.LTString:
		return &yaml.Node{
			Kind:  yaml.ScalarNode,
			Tag:   "!!str",
			Value: lv.String(),
		}
	case lua.LTFunction:
		return &yaml.Node{
			Kind:  yaml.ScalarNode,
			Tag:   "tag:lua.org,2006,function",
			Value: lv.String(),
		}
	case lua.LTTable:
		// Simultaneously create a sequence and a map, pick which one to return
		// based on whether all keys were consecutive integers
		i := 1
		yaml_sequence := &yaml.Node{
			Kind: yaml.SequenceNode,
			Tag:  "!!seq",
		}
		yaml_map := &yaml.Node{
			Kind: yaml.MappingNode,
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
			yaml_map.Content = append(yaml_map.Content, dec.convertToYamlNode(ls, k))
			yv := dec.convertToYamlNode(ls, v)
			yaml_map.Content = append(yaml_map.Content, yv)
			if i != 0 {
				yaml_sequence.Content = append(yaml_sequence.Content, yv)
			}
			k, v = ls.Next(t, k)
		}
		if i != 0 {
			return yaml_sequence
		}
		return yaml_map
	default:
		return &yaml.Node{
			Kind:        yaml.ScalarNode,
			LineComment: fmt.Sprintf("Unhandled Lua type: %s", lv.Type().String()),
			Tag:         "!!null",
			Value:       lv.String(),
		}
	}
}

func (dec *luaDecoder) decideTopLevelNode(ls *lua.LState) *yaml.Node {
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
	return &CandidateNode{
		Node: &yaml.Node{
			Kind:    yaml.DocumentNode,
			Content: []*yaml.Node{firstNode},
		},
	}, nil
}
