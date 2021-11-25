package yqlib

import (
	"strings"
	"testing"
)

func TestDeepMatch(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		ok      bool
	}{
		{"", "", true},
		{"", "x", false},
		{"x", "", false},
		{"abc", "abc", true},
		{"abc", "*", true},
		{"abc", "*c", true},
		{"abc", "*b", false},
		{"abc", "a*", true},
		{"abc", "b*", false},
		{"a", "a*", true},
		{"a", "*a", true},
		{"axbxcxdxe", "a*b*c*d*e*", true},
		{"axbxcxdxexxx", "a*b*c*d*e*", true},
		{"abxbbxdbxebxczzx", "a*b?c*x", true},
		{"abxbbxdbxebxczzy", "a*b?c*x", false},
		{strings.Repeat("a", 100), "a*a*a*a*b", false},
		{"xxx", "*x", true},
	}

	for _, tt := range tests {
		t.Run(tt.name+" "+tt.pattern, func(t *testing.T) {
			if want, got := tt.ok, deepMatch(tt.name, tt.pattern); want != got {
				t.Errorf("Expected %v got %v", want, got)
			}
		})
	}
}
