package cmd

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v3/test"
)

func TestWriteCmd(t *testing.T) {
	content := `b:
  c: 3
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s b.c 7", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 7
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteKeepCommentsCmd(t *testing.T) {
	content := `b:
  c: 3 # comment
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s b.c 7", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 7 # comment
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteWithTaggedStyleCmd(t *testing.T) {
	content := `b:
  c: dog
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s b.c cat --tag=!!str --style=tagged", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: !!str cat
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteWithDoubleQuotedStyleCmd(t *testing.T) {
	content := `b:
  c: dog
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s b.c cat --style=double", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: "cat"
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteUpdateStyleOnlyCmd(t *testing.T) {
	content := `b:
  c: dog
  d: things
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s b.* --style=single", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 'dog'
  d: 'things'
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteUpdateTagOnlyCmd(t *testing.T) {
	content := `b:
  c: true
  d: false
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s b.* --tag=!!str", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: "true"
  d: "false"
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteWithSingleQuotedStyleCmd(t *testing.T) {
	content := `b:
  c: dog
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s b.c cat --style=single", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 'cat'
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteWithLiteralStyleCmd(t *testing.T) {
	content := `b:
  c: dog
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s b.c cat --style=literal", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: |-
    cat
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteWithFoldedStyleCmd(t *testing.T) {
	content := `b:
  c: dog
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s b.c cat --style=folded", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: >-
    cat
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteEmptyMultiDocCmd(t *testing.T) {
	content := `# this is empty
---
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s c 7", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `c: 7

# this is empty
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteSurroundingEmptyMultiDocCmd(t *testing.T) {
	content := `---
# empty
---
cat: frog
---

# empty
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s -d1 c 7", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `

# empty
---
cat: frog
c: 7
---


# empty
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteFromFileCmd(t *testing.T) {
	content := `b:
  c: 3
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	source := `kittens: are cute # sure are!`
	fromFilename := test.WriteTempYamlFile(source)
	defer test.RemoveTempYamlFile(fromFilename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s b.c -f %s", filename, fromFilename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c:
    kittens: are cute # sure are!
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteEmptyCmd(t *testing.T) {
	content := ``
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s b.c 7", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 7
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteAutoCreateCmd(t *testing.T) {
	content := `applications:
  - name: app
    env:`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s applications[0].env.hello world", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `applications:
  - name: app
    env:
      hello: world
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteCmdScript(t *testing.T) {
	content := `b:
  c: 3
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	updateScript := `- command: update
  path: b.c
  value: 7`
	scriptFilename := test.WriteTempYamlFile(updateScript)
	defer test.RemoveTempYamlFile(scriptFilename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write --script %s %s", scriptFilename, filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 7
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteCmdEmptyScript(t *testing.T) {
	content := `b:
  c: 3
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	updateScript := ``
	scriptFilename := test.WriteTempYamlFile(updateScript)
	defer test.RemoveTempYamlFile(scriptFilename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write --script %s %s", scriptFilename, filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 3
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteMultiCmd(t *testing.T) {
	content := `b:
  c: 3
---
apples: great
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s -d 1 apples ok", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 3
---
apples: ok
`
	test.AssertResult(t, expectedOutput, result.Output)
}
func TestWriteInvalidDocumentIndexCmd(t *testing.T) {
	content := `b:
  c: 3
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s -df apples ok", filename))
	if result.Error == nil {
		t.Error("Expected command to fail due to invalid path")
	}
	expectedOutput := `Document index f is not a integer or *: strconv.ParseInt: parsing "f": invalid syntax`
	test.AssertResult(t, expectedOutput, result.Error.Error())
}

func TestWriteBadDocumentIndexCmd(t *testing.T) {
	content := `b:
  c: 3
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s -d 1 apples ok", filename))
	if result.Error == nil {
		t.Error("Expected command to fail due to invalid path")
	}
	expectedOutput := `asked to process document index 1 but there are only 1 document(s)`
	test.AssertResult(t, expectedOutput, result.Error.Error())
}
func TestWriteMultiAllCmd(t *testing.T) {
	content := `b:
  c: 3
---
apples: great
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s -d * apples ok", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 3
apples: ok
---
apples: ok`
	test.AssertResult(t, expectedOutput, strings.Trim(result.Output, "\n "))
}

func TestWriteCmd_EmptyArray(t *testing.T) {
	content := `b: 3`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s a []", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b: 3
a: []
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteCmd_Error(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "write")
	if result.Error == nil {
		t.Error("Expected command to fail due to missing arg")
	}
	expectedOutput := `Must provide <filename> <path_to_update> <value>`
	test.AssertResult(t, expectedOutput, result.Error.Error())
}

func TestWriteCmd_ErrorUnreadableFile(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "write fake-unknown a.b 3")
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

func TestWriteCmd_Inplace(t *testing.T) {
	content := `b:
  c: 3
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write -i %s b.c 7", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	gotOutput := test.ReadTempYamlFile(filename)
	expectedOutput := `b:
  c: 7`
	test.AssertResult(t, expectedOutput, strings.Trim(gotOutput, "\n "))
}

func TestWriteCmd_InplaceError(t *testing.T) {
	content := `b: cat
  c: 3
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write -i %s b.c 7", filename))
	if result.Error == nil {
		t.Error("Expected Error to occur!")
	}
	gotOutput := test.ReadTempYamlFile(filename)
	test.AssertResult(t, content, gotOutput)
}

func TestWriteCmd_Append(t *testing.T) {
	content := `b:
  - foo
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s b[+] 7", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  - foo
  - 7
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteCmd_AppendInline(t *testing.T) {
	content := `b: [foo]`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s b[+] 7", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b: [foo, 7]
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteCmd_AppendInlinePretty(t *testing.T) {
	content := `b: [foo]`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s -P b[+] 7", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  - foo
  - 7
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteCmd_AppendEmptyArray(t *testing.T) {
	content := `a: 2
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s b[+] v", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `a: 2
b:
  - v
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteCmd_SplatArray(t *testing.T) {
	content := `b:
- c: thing
- c: another thing
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s b[*].c new", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  - c: new
  - c: new
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteCmd_SplatMap(t *testing.T) {
	content := `b:
  c: thing
  d: another thing
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s b.* new", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: new
  d: new
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestWriteCmd_SplatMapEmpty(t *testing.T) {
	content := `b:
  c: thing
  d: another thing
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("write %s b.c.* new", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: {}
  d: another thing
`
	test.AssertResult(t, expectedOutput, result.Output)
}
