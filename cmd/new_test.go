package cmd

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/mikefarah/yq/v3/test"
// )

// func TestNewCmd(t *testing.T) {
// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, "new b.c 3")
// 	if result.Error != nil {
// 		t.Error(result.Error)
// 	}
// 	expectedOutput := `b:
//   c: 3
// `
// 	test.AssertResult(t, expectedOutput, result.Output)
// }

// func TestNewCmdScript(t *testing.T) {
// 	updateScript := `- command: update
//   path: b.c
//   value: 7`
// 	scriptFilename := test.WriteTempYamlFile(updateScript)
// 	defer test.RemoveTempYamlFile(scriptFilename)

// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, fmt.Sprintf("new --script %s", scriptFilename))
// 	if result.Error != nil {
// 		t.Error(result.Error)
// 	}
// 	expectedOutput := `b:
//   c: 7
// `
// 	test.AssertResult(t, expectedOutput, result.Output)
// }

// func TestNewAnchorCmd(t *testing.T) {
// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, "new b.c 3 --anchorName=fred")
// 	if result.Error != nil {
// 		t.Error(result.Error)
// 	}
// 	expectedOutput := `b:
//   c: &fred 3
// `
// 	test.AssertResult(t, expectedOutput, result.Output)
// }

// func TestNewAliasCmd(t *testing.T) {
// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, "new b.c foo --makeAlias")
// 	if result.Error != nil {
// 		t.Error(result.Error)
// 	}
// 	expectedOutput := `b:
//   c: *foo
// `
// 	test.AssertResult(t, expectedOutput, result.Output)
// }

// func TestNewArrayCmd(t *testing.T) {
// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, "new b[0] 3")
// 	if result.Error != nil {
// 		t.Error(result.Error)
// 	}
// 	expectedOutput := `b:
//   - 3
// `
// 	test.AssertResult(t, expectedOutput, result.Output)
// }

// func TestNewCmd_Error(t *testing.T) {
// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, "new b.c")
// 	if result.Error == nil {
// 		t.Error("Expected command to fail due to missing arg")
// 	}
// 	expectedOutput := `Must provide <path_to_update> <value>`
// 	test.AssertResult(t, expectedOutput, result.Error.Error())
// }

// func TestNewWithTaggedStyleCmd(t *testing.T) {
// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, "new b.c cat --tag=!!str --style=tagged")
// 	if result.Error != nil {
// 		t.Error(result.Error)
// 	}
// 	expectedOutput := `b:
//   c: !!str cat
// `
// 	test.AssertResult(t, expectedOutput, result.Output)
// }

// func TestNewWithDoubleQuotedStyleCmd(t *testing.T) {
// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, "new b.c cat --style=double")
// 	if result.Error != nil {
// 		t.Error(result.Error)
// 	}
// 	expectedOutput := `b:
//   c: "cat"
// `
// 	test.AssertResult(t, expectedOutput, result.Output)
// }

// func TestNewWithSingleQuotedStyleCmd(t *testing.T) {
// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, "new b.c cat --style=single")
// 	if result.Error != nil {
// 		t.Error(result.Error)
// 	}
// 	expectedOutput := `b:
//   c: 'cat'
// `
// 	test.AssertResult(t, expectedOutput, result.Output)
// }
