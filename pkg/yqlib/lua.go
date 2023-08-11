package yqlib

type LuaPreferences struct {
	DocPrefix    string
	DocSuffix    string
	UnquotedKeys bool
	Globals      bool
}

func NewDefaultLuaPreferences() LuaPreferences {
	return LuaPreferences{
		DocPrefix:    "return ",
		DocSuffix:    ";\n",
		UnquotedKeys: false,
		Globals:      false,
	}
}

var ConfiguredLuaPreferences = NewDefaultLuaPreferences()
