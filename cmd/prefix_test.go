package cmd

// import (
// 	"fmt"
// 	"runtime"
// 	"strings"
// 	"testing"

// 	"github.com/mikefarah/yq/v3/test"
// )

// func TestPrefixCmd(t *testing.T) {
// 	content := `b:
//   c: 3
// `
// 	filename := test.WriteTempYamlFile(content)
// 	defer test.RemoveTempYamlFile(filename)

// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, fmt.Sprintf("prefix %s d", filename))
// 	if result.Error != nil {
// 		t.Error(result.Error)
// 	}
// 	expectedOutput := `d:
//   b:
//     c: 3
// `
// 	test.AssertResult(t, expectedOutput, result.Output)
// }

// func TestPrefixCmdArray(t *testing.T) {
// 	content := `b:
//   c: 3
// `
// 	filename := test.WriteTempYamlFile(content)
// 	defer test.RemoveTempYamlFile(filename)

// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, fmt.Sprintf("prefix %s [+].d.[+]", filename))
// 	if result.Error != nil {
// 		t.Error(result.Error)
// 	}
// 	expectedOutput := `- d:
//     - b:
//         c: 3
// `
// 	test.AssertResult(t, expectedOutput, result.Output)
// }

// func TestPrefixCmd_MultiLayer(t *testing.T) {
// 	content := `b:
//   c: 3
// `
// 	filename := test.WriteTempYamlFile(content)
// 	defer test.RemoveTempYamlFile(filename)

// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, fmt.Sprintf("prefix %s d.e.f", filename))
// 	if result.Error != nil {
// 		t.Error(result.Error)
// 	}
// 	expectedOutput := `d:
//   e:
//     f:
//       b:
//         c: 3
// `
// 	test.AssertResult(t, expectedOutput, result.Output)
// }

// func TestPrefixMultiCmd(t *testing.T) {
// 	content := `b:
//   c: 3
// ---
// apples: great
// `
// 	filename := test.WriteTempYamlFile(content)
// 	defer test.RemoveTempYamlFile(filename)

// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, fmt.Sprintf("prefix %s -d 1 d", filename))
// 	if result.Error != nil {
// 		t.Error(result.Error)
// 	}
// 	expectedOutput := `b:
//   c: 3
// ---
// d:
//   apples: great
// `
// 	test.AssertResult(t, expectedOutput, result.Output)
// }
// func TestPrefixInvalidDocumentIndexCmd(t *testing.T) {
// 	content := `b:
//   c: 3
// `
// 	filename := test.WriteTempYamlFile(content)
// 	defer test.RemoveTempYamlFile(filename)

// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, fmt.Sprintf("prefix %s -df d", filename))
// 	if result.Error == nil {
// 		t.Error("Expected command to fail due to invalid path")
// 	}
// 	expectedOutput := `Document index f is not a integer or *: strconv.ParseInt: parsing "f": invalid syntax`
// 	test.AssertResult(t, expectedOutput, result.Error.Error())
// }

// func TestPrefixBadDocumentIndexCmd(t *testing.T) {
// 	content := `b:
//   c: 3
// `
// 	filename := test.WriteTempYamlFile(content)
// 	defer test.RemoveTempYamlFile(filename)

// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, fmt.Sprintf("prefix %s -d 1 d", filename))
// 	if result.Error == nil {
// 		t.Error("Expected command to fail due to invalid path")
// 	}
// 	expectedOutput := `asked to process document index 1 but there are only 1 document(s)`
// 	test.AssertResult(t, expectedOutput, result.Error.Error())
// }
// func TestPrefixMultiAllCmd(t *testing.T) {
// 	content := `b:
//   c: 3
// ---
// apples: great
// `
// 	filename := test.WriteTempYamlFile(content)
// 	defer test.RemoveTempYamlFile(filename)

// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, fmt.Sprintf("prefix %s -d * d", filename))
// 	if result.Error != nil {
// 		t.Error(result.Error)
// 	}
// 	expectedOutput := `d:
//   b:
//     c: 3
// ---
// d:
//   apples: great`
// 	test.AssertResult(t, expectedOutput, strings.Trim(result.Output, "\n "))
// }

// func TestPrefixCmd_Error(t *testing.T) {
// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, "prefix")
// 	if result.Error == nil {
// 		t.Error("Expected command to fail due to missing arg")
// 	}
// 	expectedOutput := `Must provide <filename> <prefixed_path>`
// 	test.AssertResult(t, expectedOutput, result.Error.Error())
// }

// func TestPrefixCmd_ErrorUnreadableFile(t *testing.T) {
// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, "prefix fake-unknown a.b")
// 	if result.Error == nil {
// 		t.Error("Expected command to fail due to unknown file")
// 	}
// 	var expectedOutput string
// 	if runtime.GOOS == "windows" {
// 		expectedOutput = `open fake-unknown: The system cannot find the file specified.`
// 	} else {
// 		expectedOutput = `open fake-unknown: no such file or directory`
// 	}
// 	test.AssertResult(t, expectedOutput, result.Error.Error())
// }

// func TestPrefixCmd_Inplace(t *testing.T) {
// 	content := `b:
//   c: 3
// `
// 	filename := test.WriteTempYamlFile(content)
// 	defer test.RemoveTempYamlFile(filename)

// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, fmt.Sprintf("prefix -i %s d", filename))
// 	if result.Error != nil {
// 		t.Error(result.Error)
// 	}
// 	gotOutput := test.ReadTempYamlFile(filename)
// 	expectedOutput := `d:
//   b:
//     c: 3`
// 	test.AssertResult(t, expectedOutput, strings.Trim(gotOutput, "\n "))
// }
