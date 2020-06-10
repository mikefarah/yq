package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v3/test"
)

func TestDeleteYamlCmd(t *testing.T) {
	content := `a: 2
b:
  c: things
  d: something else
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("delete %s b.c", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}

	expectedOutput := `a: 2
b:
  d: something else
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestDeleteDeepDoesNotExistCmd(t *testing.T) {
	content := `a: 2`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("delete %s b.c", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}

	expectedOutput := `a: 2
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestDeleteSplatYaml(t *testing.T) {
	content := `a: other
b: [3, 4]
c:
  toast: leave
  test: 1
  tell: 1
  tasty.taco: cool
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("delete %s c.te*", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}

	expectedOutput := `a: other
b: [3, 4]
c:
  toast: leave
  tasty.taco: cool
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestDeleteSplatArrayYaml(t *testing.T) {
	content := `a: 2
b:
 hi:
  - thing: item1
    name: fred
  - thing: item2
    name: sam
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("delete %s b.hi[*].thing", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}

	expectedOutput := `a: 2
b:
  hi:
    - name: fred
    - name: sam
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestDeleteSplatPrefixYaml(t *testing.T) {
	content := `a: 2
b:
 hi:
   c: things
   d: something else
 there:
   c: more things
   d: more something else
 there2:
   c: more things also
   d: more something else also
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("delete %s b.there*.c", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}

	expectedOutput := `a: 2
b:
  hi:
    c: things
    d: something else
  there:
    d: more something else
  there2:
    d: more something else also
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestDeleteYamlArrayCmd(t *testing.T) {
	content := `- 1
- 2
- 3
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("delete %s [1]", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}

	expectedOutput := `- 1
- 3
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestDeleteYamlArrayExpressionCmd(t *testing.T) {
	content := `- name: fred
- name: cat
- name: thing
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("delete %s (name==cat)", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}

	expectedOutput := `- name: fred
- name: thing
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestDeleteYamlMulti(t *testing.T) {
	content := `apples: great
---
- 1
- 2
- 3
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("delete -d 1 %s [1]", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}

	expectedOutput := `apples: great
---
- 1
- 3
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestDeleteYamlMultiAllCmd(t *testing.T) {
	content := `b:
  c: 3
apples: great
---
apples: great
something: else
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("delete %s -d * apples", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 3
---
something: else`
	test.AssertResult(t, expectedOutput, strings.Trim(result.Output, "\n "))
}
