package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func getRootCommand() *cobra.Command {
	return newCommandCLI()
}

func TestRootCmd(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "")
	if result.Error != nil {
		t.Error(result.Error)
	}

	if !strings.Contains(result.Output, "Usage:") {
		t.Error("Expected usage message to be printed out, but the usage message was not found.")
	}

}

func TestRootCmd_VerboseLong(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "--verbose")
	if result.Error != nil {
		t.Error(result.Error)
	}

	if !verbose {
		t.Error("Expected verbose to be true")
	}
}

func TestRootCmd_VerboseShort(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "-v")
	if result.Error != nil {
		t.Error(result.Error)
	}

	if !verbose {
		t.Error("Expected verbose to be true")
	}
}

func TestRootCmd_TrimLong(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "--trim")
	if result.Error != nil {
		t.Error(result.Error)
	}

	if !trimOutput {
		t.Error("Expected trimOutput to be true")
	}
}

func TestRootCmd_TrimShort(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "-t")
	if result.Error != nil {
		t.Error(result.Error)
	}

	if !trimOutput {
		t.Error("Expected trimOutput to be true")
	}
}

func TestRootCmd_ToJsonLong(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "--tojson")
	if result.Error != nil {
		t.Error(result.Error)
	}

	if !outputToJSON {
		t.Error("Expected outputToJSON to be true")
	}
}

func TestRootCmd_ToJsonShort(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "-j")
	if result.Error != nil {
		t.Error(result.Error)
	}

	if !outputToJSON {
		t.Error("Expected outputToJSON to be true")
	}
}

func TestReadCmd(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "read examples/sample.yaml b.c")
	if result.Error != nil {
		t.Error(result.Error)
	}
	assertResult(t, "2\n", result.Output)
}

func TestReadCmd_ArrayYaml(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "read examples/array.yaml [0].gather_facts")
	if result.Error != nil {
		t.Error(result.Error)
	}
	assertResult(t, "false\n", result.Output)
}

func TestReadCmd_ArrayYaml_NoPath(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "read examples/array.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `- become: true
  gather_facts: false
  hosts: lalaland
  name: Apply smth
  roles:
  - lala
  - land
  serial: 1
`
	assertResult(t, expectedOutput, result.Output)
}

func TestReadCmd_ArrayYaml_OneElement(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "read examples/array.yaml [0]")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `become: true
gather_facts: false
hosts: lalaland
name: Apply smth
roles:
- lala
- land
serial: 1
`
	assertResult(t, expectedOutput, result.Output)
}

func TestReadCmd_ArrayYaml_Splat(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "read examples/array.yaml [*]")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `- become: true
  gather_facts: false
  hosts: lalaland
  name: Apply smth
  roles:
  - lala
  - land
  serial: 1
`
	assertResult(t, expectedOutput, result.Output)
}

func TestReadCmd_ArrayYaml_SplatKey(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "read examples/array.yaml [*].gather_facts")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := "- false\n"
	assertResult(t, expectedOutput, result.Output)
}

func TestReadCmd_ArrayYaml_ErrorBadPath(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "read examples/array.yaml [x].gather_facts")
	if result.Error == nil {
		t.Error("Expected command to fail due to invalid path")
	}
	expectedOutput := `Error accessing array: strconv.ParseInt: parsing "x": invalid syntax`
	assertResult(t, expectedOutput, result.Error.Error())
}

func TestReadCmd_ArrayYaml_Splat_ErrorBadPath(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "read examples/array.yaml [*].roles[x]")
	if result.Error == nil {
		t.Error("Expected command to fail due to invalid path")
	}
	expectedOutput := `Error accessing array: strconv.ParseInt: parsing "x": invalid syntax`
	assertResult(t, expectedOutput, result.Error.Error())
}

func TestReadCmd_Error(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "read")
	if result.Error == nil {
		t.Error("Expected command to fail due to missing arg")
	}
	expectedOutput := `Must provide filename`
	assertResult(t, expectedOutput, result.Error.Error())
}

func TestReadCmd_ErrorEmptyFilename(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "read  ")
	if result.Error == nil {
		t.Error("Expected command to fail due to missing arg")
	}
	expectedOutput := `Must provide filename`
	assertResult(t, expectedOutput, result.Error.Error())
}

func TestReadCmd_ErrorUnreadableFile(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "read fake-unknown")
	if result.Error == nil {
		t.Error("Expected command to fail due to unknown file")
	}
	expectedOutput := `open fake-unknown: no such file or directory`
	assertResult(t, expectedOutput, result.Error.Error())
}

func TestReadCmd_ErrorBadPath(t *testing.T) {
	content := `b:
  d:
    e:
      - 3
      - 4
    f:
      - 1
      - 2
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("read %s b.d.*.[x]", filename))
	if result.Error == nil {
		t.Fatal("Expected command to fail due to invalid path")
	}
	expectedOutput := `Error accessing array: strconv.ParseInt: parsing "x": invalid syntax`
	assertResult(t, expectedOutput, result.Error.Error())
}

func TestReadCmd_Verbose(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "-v read examples/sample.yaml b.c")
	if result.Error != nil {
		t.Error(result.Error)
	}
	assertResult(t, "2\n", result.Output)
}

func TestReadCmd_NoTrim(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "--trim=false read examples/sample.yaml b.c")
	if result.Error != nil {
		t.Error(result.Error)
	}
	assertResult(t, "2\n\n", result.Output)
}

func TestReadCmd_ToJson(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "-j read examples/sample.yaml b.c")
	if result.Error != nil {
		t.Error(result.Error)
	}
	assertResult(t, "2\n", result.Output)
}

func TestNewCmd(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "new b.c 3")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 3
`
	assertResult(t, expectedOutput, result.Output)
}

func TestNewCmd_Error(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "new b.c")
	if result.Error == nil {
		t.Error("Expected command to fail due to missing arg")
	}
	expectedOutput := `Must provide <path_to_update> <value>`
	assertResult(t, expectedOutput, result.Error.Error())
}

func TestNewCmd_Verbose(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "-v new b.c 3")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 3
`
	assertResult(t, expectedOutput, result.Output)
}

func TestNewCmd_ToJson(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "-j new b.c 3")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `{"b":{"c":3}}
`
	assertResult(t, expectedOutput, result.Output)
}

func TestWriteCmd(t *testing.T) {
	content := `b:
  c: 3
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("write %s b.c 7", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 7
`
	assertResult(t, expectedOutput, result.Output)
}

func TestWriteCmd_EmptyArray(t *testing.T) {
	content := `b: 3`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("write %s a []", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b: 3
a: []
`
	assertResult(t, expectedOutput, result.Output)
}

func TestWriteCmd_Error(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "write")
	if result.Error == nil {
		t.Error("Expected command to fail due to missing arg")
	}
	expectedOutput := `Must provide <filename> <path_to_update> <value>`
	assertResult(t, expectedOutput, result.Error.Error())
}

func TestWriteCmd_ErrorUnreadableFile(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "write fake-unknown a.b 3")
	if result.Error == nil {
		t.Error("Expected command to fail due to unknown file")
	}
	expectedOutput := `open fake-unknown: no such file or directory`
	assertResult(t, expectedOutput, result.Error.Error())
}

func TestWriteCmd_Verbose(t *testing.T) {
	content := `b:
  c: 3
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("-v write %s b.c 7", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 7
`
	assertResult(t, expectedOutput, result.Output)
}

func TestWriteCmd_Inplace(t *testing.T) {
	content := `b:
  c: 3
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("write -i %s b.c 7", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	gotOutput := readTempYamlFile(filename)
	expectedOutput := `b:
  c: 7`
	assertResult(t, expectedOutput, gotOutput)
}

func TestMergeCmd(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "merge examples/data1.yaml examples/data2.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `a: simple
b:
- 1
- 2
c:
  test: 1
`
	assertResult(t, expectedOutput, result.Output)
}

func TestMergeCmd_Error(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "merge examples/data1.yaml")
	if result.Error == nil {
		t.Error("Expected command to fail due to missing arg")
	}
	expectedOutput := `Must provide at least 2 yaml files`
	assertResult(t, expectedOutput, result.Error.Error())
}

func TestMergeCmd_ErrorUnreadableFile(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "merge examples/data1.yaml fake-unknown")
	if result.Error == nil {
		t.Error("Expected command to fail due to unknown file")
	}
	expectedOutput := `open fake-unknown: no such file or directory`
	assertResult(t, expectedOutput, result.Error.Error())
}

func TestMergeCmd_Verbose(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "-v merge examples/data1.yaml examples/data2.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `a: simple
b:
- 1
- 2
c:
  test: 1
`
	assertResult(t, expectedOutput, result.Output)
}

func TestMergeCmd_Inplace(t *testing.T) {
	filename := writeTempYamlFile(readTempYamlFile("examples/data1.yaml"))
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("merge -i %s examples/data2.yaml", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	gotOutput := readTempYamlFile(filename)
	expectedOutput := `a: simple
b:
- 1
- 2
c:
  test: 1`
	assertResult(t, expectedOutput, gotOutput)
}
