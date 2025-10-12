package yqlib

import (
	"bytes"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name     string
		attr     color.Attribute
		expected string
	}{
		{
			name:     "reset color",
			attr:     color.Reset,
			expected: "\x1b[0m",
		},
		{
			name:     "red color",
			attr:     color.FgRed,
			expected: "\x1b[31m",
		},
		{
			name:     "green color",
			attr:     color.FgGreen,
			expected: "\x1b[32m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := format(tt.attr)
			if result != tt.expected {
				t.Errorf("format(%d) = %q, want %q", tt.attr, result, tt.expected)
			}
		})
	}
}

func TestColorizeAndPrint(t *testing.T) {
	tests := []struct {
		name      string
		yamlBytes []byte
		expectErr bool
	}{
		{
			name:      "simple yaml",
			yamlBytes: []byte("name: test\nage: 25\n"),
			expectErr: false,
		},
		{
			name:      "yaml with strings",
			yamlBytes: []byte("name: \"hello world\"\nactive: true\ncount: 42\n"),
			expectErr: false,
		},
		{
			name:      "yaml with anchors and aliases",
			yamlBytes: []byte("default: &default\n  name: test\nuser: *default\n"),
			expectErr: false,
		},
		{
			name:      "yaml with comments",
			yamlBytes: []byte("# This is a comment\nname: test\n"),
			expectErr: false,
		},
		{
			name:      "empty yaml",
			yamlBytes: []byte(""),
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := colorizeAndPrint(tt.yamlBytes, &buf)

			if tt.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check that output contains escape sequences (color codes)
			if !tt.expectErr && len(tt.yamlBytes) > 0 {
				output := buf.String()
				if !strings.Contains(output, "\x1b[") {
					t.Error("Expected output to contain color escape sequences")
				}
			}
		})
	}
}

func TestColorizeAndPrintWithDifferentYamlTypes(t *testing.T) {
	testCases := []struct {
		name      string
		yaml      string
		expectErr bool
	}{
		{
			name: "boolean values",
			yaml: "active: true\ninactive: false\n",
		},
		{
			name: "numeric values",
			yaml: "integer: 42\nfloat: 3.14\nnegative: -10\n",
		},
		{
			name: "map keys",
			yaml: "user:\n  name: john\n  age: 30\n",
		},
		{
			name: "string values",
			yaml: "message: \"hello world\"\ndescription: 'single quotes'\n",
		},
		{
			name: "mixed types",
			yaml: "config:\n  debug: true\n  port: 8080\n  host: \"localhost\"\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := colorizeAndPrint([]byte(tc.yaml), &buf)

			if tc.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Verify output contains color codes
			if !tc.expectErr {
				output := buf.String()
				if !strings.Contains(output, "\x1b[") {
					t.Error("Expected output to contain color escape sequences")
				}
				// Should end with newline
				if !strings.HasSuffix(output, "\n") {
					t.Error("Expected output to end with newline")
				}
			}
		})
	}
}
