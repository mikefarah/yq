package cmd

import (
	"fmt"
	"testing"

	"github.com/mikefarah/yq/v3/test"
)

func TestValidateCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "validate ../examples/sample.yaml b.c")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "", result.Output)
}

func TestValidateBadDataCmd(t *testing.T) {
	content := `[!Whatever]`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("validate %s", filename))
	if result.Error == nil {
		t.Error("Expected command to fail")
	}
	expectedOutput := `yaml: line 1: did not find expected ',' or ']'`
	test.AssertResult(t, expectedOutput, result.Error.Error())
}
