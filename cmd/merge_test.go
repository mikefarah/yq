package cmd

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/mikefarah/yq/v3/test"
)

func TestMergeCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "merge ../examples/data1.yaml ../examples/data2.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `a: simple # just the best
b: [1, 2]
c:
  test: 1
  toast: leave
  tell: 1
  tasty.taco: cool
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestMergeOneFileCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "merge ../examples/data1.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `a: simple # just the best
b: [1, 2]
c:
  test: 1
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestMergeNoAutoCreateCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "merge -c=false ../examples/data1.yaml ../examples/data2.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `a: simple # just the best
b: [1, 2]
c:
  test: 1
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestMergeOverwriteCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "merge -c=false --overwrite ../examples/data1.yaml ../examples/data2.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `a: other # better than the original
b: [3, 4]
c:
  test: 1
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestMergeAppendCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "merge --autocreate=false --append ../examples/data1.yaml ../examples/data2.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `a: simple # just the best
b: [1, 2, 3, 4]
c:
  test: 1
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestMergeAppendArraysCmd(t *testing.T) {
	content := `people:
  - name: Barry
    age: 21`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	mergeContent := `people:
  - name: Roger
    age: 44`
	mergeFilename := test.WriteTempYamlFile(mergeContent)
	defer test.RemoveTempYamlFile(mergeFilename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("merge --append -d* %s %s", filename, mergeFilename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `people:
  - name: Barry
    age: 21
  - name: Roger
    age: 44
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestMergeOverwriteAndAppendCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "merge --autocreate=false --append --overwrite ../examples/data1.yaml ../examples/data2.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `a: other # better than the original
b: [1, 2, 3, 4]
c:
  test: 1
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestMergeArraysCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "merge --append ../examples/sample_array.yaml ../examples/sample_array_2.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `[1, 2, 3, 4, 5]
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestMergeCmd_Multi(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "merge -d1 ../examples/multiple_docs_small.yaml ../examples/data1.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `a: Easy! as one two three
---
another:
  document: here
a: simple # just the best
b:
  - 1
  - 2
c:
  test: 1
---
- 1
- 2
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestMergeYamlMultiAllCmd(t *testing.T) {
	content := `b:
  c: 3
apples: green
---
something: else`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	mergeContent := `apples: red
something: good`
	mergeFilename := test.WriteTempYamlFile(mergeContent)
	defer test.RemoveTempYamlFile(mergeFilename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("merge -d* %s %s", filename, mergeFilename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 3
apples: green
something: good
---
something: else
apples: red
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestMergeSpecialCharacterKeysCmd(t *testing.T) {
	content := ``
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	mergeContent := `key[bracket]: value
key.bracket: value
key"value": value
key'value': value
`
	mergeFilename := test.WriteTempYamlFile(mergeContent)
	defer test.RemoveTempYamlFile(mergeFilename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("merge %s %s", filename, mergeFilename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, mergeContent, result.Output)
}

func TestMergeYamlMultiAllOverwriteCmd(t *testing.T) {
	content := `b:
  c: 3
apples: green
---
something: else`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	mergeContent := `apples: red
something: good`
	mergeFilename := test.WriteTempYamlFile(mergeContent)
	defer test.RemoveTempYamlFile(mergeFilename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("merge --overwrite -d* %s %s", filename, mergeFilename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 3
apples: red
something: good
---
something: good
apples: red
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestMergeYamlNullMapCmd(t *testing.T) {
	content := `b:`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	mergeContent := `b:
  thing: a frog
`
	mergeFilename := test.WriteTempYamlFile(mergeContent)
	defer test.RemoveTempYamlFile(mergeFilename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("merge %s %s", filename, mergeFilename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, mergeContent, result.Output)
}

func TestMergeCmd_Error(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "merge")
	if result.Error == nil {
		t.Error("Expected command to fail due to missing arg")
	}
	expectedOutput := `Must provide at least 1 yaml file`
	test.AssertResult(t, expectedOutput, result.Error.Error())
}

func TestMergeCmd_ErrorUnreadableFile(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "merge ../examples/data1.yaml fake-unknown")
	if result.Error == nil {
		t.Error("Expected command to fail due to unknown file")
	}
	var expectedOutput string
	if runtime.GOOS == "windows" {
		expectedOutput = `open fake-unknown: The system cannot find the file specified.`
	} else {
		expectedOutput = `open fake-unknown: no such file or directory`
	}
	test.AssertResult(t, expectedOutput, result.Error.Error())
}

func TestMergeCmd_Inplace(t *testing.T) {
	filename := test.WriteTempYamlFile(test.ReadTempYamlFile("../examples/data1.yaml"))
	err := os.Chmod(filename, os.FileMode(int(0666)))
	if err != nil {
		t.Error(err)
	}
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("merge -i %s ../examples/data2.yaml", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	info, _ := os.Stat(filename)
	gotOutput := test.ReadTempYamlFile(filename)
	expectedOutput := `a: simple # just the best
b: [1, 2]
c:
  test: 1
  toast: leave
  tell: 1
  tasty.taco: cool
`
	test.AssertResult(t, expectedOutput, gotOutput)
	test.AssertResult(t, os.FileMode(int(0666)), info.Mode())
}

func TestMergeAllowEmptyTargetCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "merge ../examples/empty.yaml ../examples/data1.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `a: simple # just the best
b:
  - 1
  - 2
c:
  test: 1
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestMergeAllowEmptyMergeCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "merge ../examples/data1.yaml ../examples/empty.yaml")
	expectedOutput := `a: simple # just the best
b: [1, 2]
c:
  test: 1
`
	test.AssertResult(t, expectedOutput, result.Output)
}
