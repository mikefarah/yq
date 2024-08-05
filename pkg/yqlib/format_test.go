package yqlib

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

type formatStringScenario struct {
	description    string
	input          string
	expectedFormat *Format
	expectedError  string
}

var formatStringScenarios = []formatStringScenario{
	{
		description:    "yaml",
		input:          "yaml",
		expectedFormat: YamlFormat,
	},
	{
		description:   "Unknown format type",
		input:         "doc",
		expectedError: "unknown format 'doc' please use",
	},
	{
		description:   "blank should error",
		input:         "",
		expectedError: "unknown format '' please use",
	},
}

func TestFormatFromString(t *testing.T) {
	for _, tt := range formatStringScenarios {
		actualFormat, actualError := FormatFromString(tt.input)

		if tt.expectedError != "" {
			if actualError == nil {
				t.Errorf("Expected [%v] error but found none", tt.expectedError)
			} else {
				test.AssertResultWithContext(t, true, strings.Contains(actualError.Error(), tt.expectedError),
					fmt.Sprintf("Expected [%v] to contain [%v]", actualError.Error(), tt.expectedError),
				)
			}
		} else {
			test.AssertResult(t, tt.expectedFormat, actualFormat)
		}
	}
}

func TestFormatStringFromFilename(t *testing.T) {
	test.AssertResult(t, "yaml", FormatStringFromFilename("test.Yaml"))
	test.AssertResult(t, "yaml", FormatStringFromFilename("test.index.Yaml"))
	test.AssertResult(t, "yaml", FormatStringFromFilename("test"))
	test.AssertResult(t, "json", FormatStringFromFilename("test.json"))
	test.AssertResult(t, "json", FormatStringFromFilename("TEST.JSON"))
	test.AssertResult(t, "yaml", FormatStringFromFilename("test.json/foo"))
	test.AssertResult(t, "yaml", FormatStringFromFilename(""))
}
