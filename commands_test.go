package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"gopkg.in/spf13/cobra.v0"
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

func TestRootCmd_VersionShort(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "-V")
	if result.Error != nil {
		t.Error(result.Error)
	}
	if !strings.Contains(result.Output, "yq version") {
		t.Error("expected version message to be printed out, but the message was not found.")
	}
}

func TestRootCmd_VersionLong(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "--version")
	if result.Error != nil {
		t.Error(result.Error)
	}
	if !strings.Contains(result.Output, "yq version") {
		t.Error("expected version message to be printed out, but the message was not found.")
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

func TestReadInvalidDocumentIndexCmd(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "read -df examples/sample.yaml b.c")
	if result.Error == nil {
		t.Error("Expected command to fail due to invalid path")
	}
	expectedOutput := `Document index f is not a integer or *: strconv.ParseInt: parsing "f": invalid syntax`
	assertResult(t, expectedOutput, result.Error.Error())
}

func TestReadBadDocumentIndexCmd(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "read -d1 examples/sample.yaml b.c")
	if result.Error == nil {
		t.Error("Expected command to fail due to invalid path")
	}
	expectedOutput := `Asked to process document index 1 but there are only 1 document(s)`
	assertResult(t, expectedOutput, result.Error.Error())
}

func TestReadOrderCmd(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "read examples/order.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	assertResult(t,
		`version: 3
application: MyApp
`,
		result.Output)
}

func TestReadMultiCmd(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "read -d 1 examples/multiple_docs.yaml another.document")
	if result.Error != nil {
		t.Error(result.Error)
	}
	assertResult(t, "here\n", result.Output)
}

func TestReadMultiAllCmd(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "read -d* examples/multiple_docs.yaml commonKey")
	if result.Error != nil {
		t.Error(result.Error)
	}
	assertResult(t,
		`- first document
- second document
- third document
`, result.Output)
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
	expectedOutput := `Error reading path in document index 0: Error accessing array: strconv.ParseInt: parsing "x": invalid syntax`
	assertResult(t, expectedOutput, result.Error.Error())
}

func TestReadCmd_ArrayYaml_Splat_ErrorBadPath(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "read examples/array.yaml [*].roles[x]")
	if result.Error == nil {
		t.Error("Expected command to fail due to invalid path")
	}
	expectedOutput := `Error reading path in document index 0: Error accessing array: strconv.ParseInt: parsing "x": invalid syntax`
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
	expectedOutput := `Error reading path in document index 0: Error accessing array: strconv.ParseInt: parsing "x": invalid syntax`
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
	result := runCmd(cmd, "read -j examples/sample.yaml b.c")
	if result.Error != nil {
		t.Error(result.Error)
	}
	assertResult(t, "2\n", result.Output)
}

func TestReadCmd_ToJsonLong(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "read --tojson examples/sample.yaml b.c")
	if result.Error != nil {
		t.Error(result.Error)
	}
	assertResult(t, "2\n", result.Output)
}

func TestPrefixCmd(t *testing.T) {
	content := `b:
  c: 3
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("prefix %s d", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `d:
  b:
    c: 3
`
	assertResult(t, expectedOutput, result.Output)
}

func TestPrefixCmdArray(t *testing.T) {
	content := `b:
  c: 3
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("prefix %s [0].d.[1]", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `- d:
  - null
  - b:
      c: 3
`
	assertResult(t, expectedOutput, result.Output)
}

func TestPrefixCmd_MultiLayer(t *testing.T) {
	content := `b:
  c: 3
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("prefix %s d.e.f", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `d:
  e:
    f:
      b:
        c: 3
`
	assertResult(t, expectedOutput, result.Output)
}

func TestPrefixMultiCmd(t *testing.T) {
	content := `b:
  c: 3
---
apples: great
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("prefix %s -d 1 d", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 3
---
d:
  apples: great
`
	assertResult(t, expectedOutput, result.Output)
}
func TestPrefixInvalidDocumentIndexCmd(t *testing.T) {
	content := `b:
  c: 3
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("prefix %s -df d", filename))
	if result.Error == nil {
		t.Error("Expected command to fail due to invalid path")
	}
	expectedOutput := `Document index f is not a integer or *: strconv.ParseInt: parsing "f": invalid syntax`
	assertResult(t, expectedOutput, result.Error.Error())
}

func TestPrefixBadDocumentIndexCmd(t *testing.T) {
	content := `b:
  c: 3
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("prefix %s -d 1 d", filename))
	if result.Error == nil {
		t.Error("Expected command to fail due to invalid path")
	}
	expectedOutput := `Asked to process document index 1 but there are only 1 document(s)`
	assertResult(t, expectedOutput, result.Error.Error())
}
func TestPrefixMultiAllCmd(t *testing.T) {
	content := `b:
  c: 3
---
apples: great
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("prefix %s -d * d", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `d:
  b:
    c: 3
---
d:
  apples: great`
	assertResult(t, expectedOutput, strings.Trim(result.Output, "\n "))
}

func TestPrefixCmd_Error(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "prefix")
	if result.Error == nil {
		t.Error("Expected command to fail due to missing arg")
	}
	expectedOutput := `Must provide <filename> <prefixed_path>`
	assertResult(t, expectedOutput, result.Error.Error())
}

func TestPrefixCmd_ErrorUnreadableFile(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "prefix fake-unknown a.b")
	if result.Error == nil {
		t.Error("Expected command to fail due to unknown file")
	}
	expectedOutput := `open fake-unknown: no such file or directory`
	assertResult(t, expectedOutput, result.Error.Error())
}

func TestPrefixCmd_Verbose(t *testing.T) {
	content := `b:
  c: 3
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("-v prefix %s x", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `x:
  b:
    c: 3
`
	assertResult(t, expectedOutput, result.Output)
}

func TestPrefixCmd_Inplace(t *testing.T) {
	content := `b:
  c: 3
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("prefix -i %s d", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	gotOutput := readTempYamlFile(filename)
	expectedOutput := `d:
  b:
    c: 3`
	assertResult(t, expectedOutput, strings.Trim(gotOutput, "\n "))
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

func TestWriteMultiCmd(t *testing.T) {
	content := `b:
  c: 3
---
apples: great
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("write %s -d 1 apples ok", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 3
---
apples: ok
`
	assertResult(t, expectedOutput, result.Output)
}
func TestWriteInvalidDocumentIndexCmd(t *testing.T) {
	content := `b:
  c: 3
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("write %s -df apples ok", filename))
	if result.Error == nil {
		t.Error("Expected command to fail due to invalid path")
	}
	expectedOutput := `Document index f is not a integer or *: strconv.ParseInt: parsing "f": invalid syntax`
	assertResult(t, expectedOutput, result.Error.Error())
}

func TestWriteBadDocumentIndexCmd(t *testing.T) {
	content := `b:
  c: 3
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("write %s -d 1 apples ok", filename))
	if result.Error == nil {
		t.Error("Expected command to fail due to invalid path")
	}
	expectedOutput := `Asked to process document index 1 but there are only 1 document(s)`
	assertResult(t, expectedOutput, result.Error.Error())
}
func TestWriteMultiAllCmd(t *testing.T) {
	content := `b:
  c: 3
---
apples: great
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("write %s -d * apples ok", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 3
apples: ok
---
apples: ok`
	assertResult(t, expectedOutput, strings.Trim(result.Output, "\n "))
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
	assertResult(t, expectedOutput, strings.Trim(gotOutput, "\n "))
}

func TestWriteCmd_Append(t *testing.T) {
	content := `b:
  - foo
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("write %s b[+] 7", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
- foo
- 7
`
	assertResult(t, expectedOutput, result.Output)
}

func TestWriteCmd_AppendEmptyArray(t *testing.T) {
	content := `a: 2
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("write %s b[+] v", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `a: 2
b:
- v
`
	assertResult(t, expectedOutput, result.Output)
}

func TestDeleteYaml(t *testing.T) {
	content := `a: 2
b:
  c: things
  d: something else
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("delete %s b.c", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}

	expectedOutput := `a: 2
b:
  d: something else
`
	assertResult(t, expectedOutput, result.Output)
}

func TestDeleteYamlArray(t *testing.T) {
	content := `- 1
- 2
- 3
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("delete %s [1]", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}

	expectedOutput := `- 1
- 3
`
	assertResult(t, expectedOutput, result.Output)
}

func TestDeleteYamlMulti(t *testing.T) {
	content := `apples: great
---
- 1
- 2
- 3
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("delete -d 1 %s [1]", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}

	expectedOutput := `apples: great
---
- 1
- 3
`
	assertResult(t, expectedOutput, result.Output)
}

func TestDeleteYamlMultiAllCmd(t *testing.T) {
	content := `b:
  c: 3
apples: great
---
apples: great
something: else
`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("delete %s -d * apples", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 3
---
something: else`
	assertResult(t, expectedOutput, strings.Trim(result.Output, "\n "))
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

func TestMergeOverwriteCmd(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "merge --overwrite examples/data1.yaml examples/data2.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `a: other
b:
- 3
- 4
c:
  test: 1
`
	assertResult(t, expectedOutput, result.Output)
}

func TestMergeAppendCmd(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "merge --append examples/data1.yaml examples/data2.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `a: simple
b:
- 1
- 2
- 3
- 4
c:
  test: 1
`
	assertResult(t, expectedOutput, result.Output)
}
func TestMergeArraysCmd(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "merge --append examples/sample_array.yaml examples/sample_array_2.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `- 1
- 2
- 3
- 4
- 5
`
	assertResult(t, expectedOutput, result.Output)
}

func TestMergeCmd_Multi(t *testing.T) {
	cmd := getRootCommand()
	result := runCmd(cmd, "merge -d1 examples/multiple_docs_small.yaml examples/data2.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `a: Easy! as one two three
---
a: other
another:
  document: here
b:
- 3
- 4
c:
  test: 1
---
- 1
- 2`
	assertResult(t, expectedOutput, strings.Trim(result.Output, "\n "))
}

func TestMergeYamlMultiAllCmd(t *testing.T) {
	content := `b:
  c: 3
apples: green
---
something: else`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	mergeContent := `apples: red
something: good`
	mergeFilename := writeTempYamlFile(mergeContent)
	defer removeTempYamlFile(mergeFilename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("merge -d* %s %s", filename, mergeFilename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `apples: green
b:
  c: 3
something: good
---
apples: red
something: else`
	assertResult(t, expectedOutput, strings.Trim(result.Output, "\n "))
}

func TestMergeYamlMultiAllOverwriteCmd(t *testing.T) {
	content := `b:
  c: 3
apples: green
---
something: else`
	filename := writeTempYamlFile(content)
	defer removeTempYamlFile(filename)

	mergeContent := `apples: red
something: good`
	mergeFilename := writeTempYamlFile(mergeContent)
	defer removeTempYamlFile(mergeFilename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("merge --overwrite -d* %s %s", filename, mergeFilename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `apples: red
b:
  c: 3
something: good
---
apples: red
something: good`
	assertResult(t, expectedOutput, strings.Trim(result.Output, "\n "))
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
	expectedOutput := `Error updating document at index 0: open fake-unknown: no such file or directory`
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
	err := os.Chmod(filename, os.FileMode(int(0666)))
	if err != nil {
		t.Error(err)
	}
	defer removeTempYamlFile(filename)

	cmd := getRootCommand()
	result := runCmd(cmd, fmt.Sprintf("merge -i %s examples/data2.yaml", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	info, _ := os.Stat(filename)
	gotOutput := readTempYamlFile(filename)
	expectedOutput := `a: simple
b:
- 1
- 2
c:
  test: 1`
	assertResult(t, expectedOutput, strings.Trim(gotOutput, "\n "))
	assertResult(t, os.FileMode(int(0666)), info.Mode())
}
