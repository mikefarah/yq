package test

import (
	"testing"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
)

// only test format detection based on filename extension
func TestFormatStringFromFilename(t *testing.T) {
	cases := []struct {
		filename string
		expected string
	}{
		// filenames that have extensions
		{"file.yaml", "yaml"},
		{"FILE.JSON", "json"},
		{"file.properties", "properties"},
		{"file.xml", "xml"},
		{"file.unknown", "unknown"},

		// filenames without extensions
		{"file", "yaml"},
		{"a.dir/file", "yaml"},
		{"file.", "yaml"},
		{".", "yaml"},
		{"", "yaml"},
	}

	for _, c := range cases {
		result := yqlib.FormatStringFromFilename(c.filename)
		if result != c.expected {
			t.Errorf("FormatStringFromFilename(%q) = %q, wanted: %q", c.filename, result, c.expected)
		}
	}
}
