package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v3/test"
	"github.com/spf13/cobra"
)

func getRootCommand() *cobra.Command {
	return New()
}

func TestRootCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "")
	if result.Error != nil {
		t.Error(result.Error)
	}

	if !strings.Contains(result.Output, "Usage:") {
		t.Error("Expected usage message to be printed out, but the usage message was not found.")
	}
}

func TestRootCmd_Help(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "--help")
	if result.Error != nil {
		t.Error(result.Error)
	}

	if !strings.Contains(result.Output, "yq is a lightweight and portable command-line YAML processor. It aims to be the jq or sed of yaml files.") {
		t.Error("Expected usage message to be printed out, but the usage message was not found.")
	}
}

func TestRootCmd_VerboseLong(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "--verbose")
	if result.Error != nil {
		t.Error(result.Error)
	}

	if !verbose {
		t.Error("Expected verbose to be true")
	}
}

func TestRootCmd_VerboseShort(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "-v")
	if result.Error != nil {
		t.Error(result.Error)
	}

	if !verbose {
		t.Error("Expected verbose to be true")
	}
}

func TestRootCmd_VersionShort(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "-V")
	if result.Error != nil {
		t.Error(result.Error)
	}
	if !strings.Contains(result.Output, "yq version") {
		t.Error("expected version message to be printed out, but the message was not found.")
	}
}

func TestRootCmd_VersionLong(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "--version")
	if result.Error != nil {
		t.Error(result.Error)
	}
	if !strings.Contains(result.Output, "yq version") {
		t.Error("expected version message to be printed out, but the message was not found.")
	}
}

func TestReadCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/sample.yaml b.c")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "2\n", result.Output)
}

func TestCompareSameCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "compare ../examples/data1.yaml ../examples/data1.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := ``
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestCompareDifferentCmd(t *testing.T) {
	forceOsExit = false
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "compare ../examples/data1.yaml ../examples/data3.yaml")

	expectedOutput := `-a: simple # just the best
-b: [1, 2]
+a: "simple" # just the best
+b: [1, 3]
 c:
   test: 1
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestComparePrettyCmd(t *testing.T) {
	forceOsExit = false
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "compare -P ../examples/data1.yaml ../examples/data3.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := ` a: simple # just the best
 b:
 - 1
-- 2
+- 3
 c:
   test: 1
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestComparePathsCmd(t *testing.T) {
	forceOsExit = false
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "compare -P -ppv ../examples/data1.yaml ../examples/data3.yaml **")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := ` a: simple # just the best
 b.[0]: 1
-b.[1]: 2
+b.[1]: 3
 c.test: 1
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestValidateCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "validate ../examples/sample.yaml b.c")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "", result.Output)
}

func TestReadWithAdvancedFilterCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/sample.yaml b.e(name==sam).value")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "4\n", result.Output)
}

func TestReadWithAdvancedFilterMapCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/sample.yaml b.e[name==fr*]")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `name: fred
value: 3
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadWithKeyAndValueCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -p pv ../examples/sample.yaml b.c")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "b.c: 2\n", result.Output)
}

func TestReadArrayCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -p pv ../examples/sample.yaml b.e[1].name")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "b.e.[1].name: sam\n", result.Output)
}

func TestReadArrayBackwardsCmd(t *testing.T) {
	content := `- one
- two
- three`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read -p pv %s [-1]", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "'[-1]': three\n", result.Output)
}

func TestReadArrayBackwardsNegative0Cmd(t *testing.T) {
	content := `- one
- two
- three`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read -p pv %s [-0]", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "'[0]': one\n", result.Output)
}

func TestReadArrayBackwardsPastLimitCmd(t *testing.T) {
	content := `- one
- two
- three`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read -p pv %s [-4]", filename))
	expectedOutput := "Error reading path in document index 0: Index [-4] out of range, array size is 3"
	test.AssertResult(t, expectedOutput, result.Error.Error())
}

func TestReadArrayLengthCmd(t *testing.T) {
	content := `- things
- whatever
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read -l %s", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "2\n", result.Output)
}

func TestReadArrayLengthDeepCmd(t *testing.T) {
	content := `holder: 
- things
- whatever
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read -l %s holder", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "2\n", result.Output)
}

func TestReadArrayLengthDeepMultipleCmd(t *testing.T) {
	content := `holderA: 
- things
- whatever
skipMe:
- yep
holderB: 
- other things
- cool
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read -l -c %s holder*", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "2\n", result.Output)
}

func TestReadCollectCmd(t *testing.T) {
	content := `holderA: yep
skipMe: not me
holderB: me too
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read -c %s holder*", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `- yep
- me too
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadCollectArrayCmd(t *testing.T) {
	content := `- name: fred
  value: 32
- name: sam
  value: 67
- name: fernie
  value: 103
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read -c %s (name==f*)", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `- name: fred
  value: 32
- name: fernie
  value: 103
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadArrayLengthDeepMultipleWithPathCmd(t *testing.T) {
	content := `holderA: 
- things
- whatever
holderB: 
- other things
- cool
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read -l %s -ppv holder*", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "holderA: 2\nholderB: 2\n", result.Output)
}

func TestReadObjectLengthCmd(t *testing.T) {
	content := `cat: meow
dog: bark
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read -l %s", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "2\n", result.Output)
}

func TestReadObjectLengthDeepCmd(t *testing.T) {
	content := `holder: 
  cat: meow
  dog: bark
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read -l %s holder", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "2\n", result.Output)
}

func TestReadObjectLengthDeepMultipleCmd(t *testing.T) {
	content := `holderA: 
  cat: meow
  dog: bark
holderB: 
  elephant: meow
  zebra: bark
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read -l -c %s holder*", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "2\n", result.Output)
}

func TestReadObjectLengthDeepMultipleWithPathsCmd(t *testing.T) {
	content := `holderA: 
  cat: meow
  dog: bark
holderB: 
  elephant: meow
  zebra: bark
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read -l -ppv %s holder*", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "holderA: 2\nholderB: 2\n", result.Output)
}

func TestReadScalarLengthCmd(t *testing.T) {
	content := `meow`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read -l %s", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "4\n", result.Output)
}

func TestReadDoubleQuotedStringCmd(t *testing.T) {
	content := `name: "meow face"`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read %s name", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "meow face\n", result.Output)
}

func TestReadSingleQuotedStringCmd(t *testing.T) {
	content := `name: 'meow face'`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read %s name", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "meow face\n", result.Output)
}

func TestReadQuotedMultinlineStringCmd(t *testing.T) {
	content := `test: |
  abcdefg    
  hijklmno
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read %s test", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `abcdefg    
hijklmno

`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadQuotedMultinlineNoNewLineStringCmd(t *testing.T) {
	content := `test: |-
  abcdefg    
  hijklmno
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read %s test", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `abcdefg    
hijklmno
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadBooleanCmd(t *testing.T) {
	content := `name: true`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read %s name", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "true\n", result.Output)
}

func TestReadNumberCmd(t *testing.T) {
	content := `name: 32.13`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read %s name", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "32.13\n", result.Output)
}

func TestReadDeepSplatCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -p pv ../examples/sample.yaml b.**")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b.c: 2
b.d.[0]: 3
b.d.[1]: 4
b.d.[2]: 5
b.e.[0].name: fred
b.e.[0].value: 3
b.e.[1].name: sam
b.e.[1].value: 4
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadDeepSplatWithSuffixCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -p pv ../examples/sample.yaml b.**.name")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b.e.[0].name: fred
b.e.[1].name: sam
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadWithKeyCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -p p ../examples/sample.yaml b.c")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "b.c\n", result.Output)
}

func TestReadAnchorsCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/simple-anchor.yaml foobar.a")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "1\n", result.Output)
}

func TestReadAnchorsWithKeyAndValueCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -p pv ../examples/simple-anchor.yaml foobar.a")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "foobar.a: 1\n", result.Output)
}

func TestReadAllAnchorsWithKeyAndValueCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -p pv ../examples/merge-anchor.yaml **")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `foo.a: original
foo.thing: coolasdf
foo.thirsty: yep
bar.b: 2
bar.thing: coconut
bar.c: oldbar
foobarList.c: newbar
foobarList.b: 2
foobarList.thing: coconut
foobarList.a: original
foobarList.thirsty: yep
foobar.thirty: well beyond
foobar.thing: ice
foobar.c: 3
foobar.a: original
foobar.thirsty: yep
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadMergeAnchorsOriginalCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/merge-anchor.yaml foobar.a")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "original\n", result.Output)
}

func TestReadMergeAnchorsExplodeJsonCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -j ../examples/merge-anchor.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `{"bar":{"b":2,"c":"oldbar","thing":"coconut"},"foo":{"a":"original","thing":"coolasdf","thirsty":"yep"},"foobar":{"a":"original","c":3,"thing":"ice","thirsty":"yep","thirty":"well beyond"},"foobarList":{"a":"original","b":2,"c":"newbar","thing":"coconut","thirsty":"yep"}}
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadMergeAnchorsExplodeSimpleCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -X ../examples/simple-anchor.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `foo:
  a: 1
foobar:
  a: 1
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadMergeAnchorsExplodeSimpleValueCmd(t *testing.T) {
	content := `value: &value-pointer the value
pointer: *value-pointer`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read -X %s pointer", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := "the value\n"
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadMergeAnchorsExplodeSimpleArrayCmd(t *testing.T) {
	content := `- things`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read -X %s", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `- things
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadNumberKeyJsonCmd(t *testing.T) {
	content := `data: {"40433437326": 10.833332}`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read -j %s", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `{"data":{"40433437326":10.833332}}
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadMergeAnchorsExplodeSimpleArrayJsonCmd(t *testing.T) {
	content := `- things`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read -j %s", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `["things"]
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadMergeAnchorsExplodeSimpleValueForValueCmd(t *testing.T) {
	content := `value: &value-pointer the value
pointer: *value-pointer`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read -X %s value", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := "the value\n"
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadMergeAnchorsExplodeCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -X ../examples/merge-anchor.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `foo:
  a: original
  thing: coolasdf
  thirsty: yep
bar:
  b: 2
  thing: coconut
  c: oldbar
foobarList:
  c: newbar
  b: 2
  thing: coconut
  a: original
  thirsty: yep
foobar:
  thirty: well beyond
  thing: ice
  c: 3
  a: original
  thirsty: yep
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadMergeAnchorsExplodeDeepCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -X ../examples/merge-anchor.yaml foobar")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `thirty: well beyond
thing: ice
c: 3
a: original
thirsty: yep
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadMergeAnchorsOverrideCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/merge-anchor.yaml foobar.thing")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "ice\n", result.Output)
}

func TestReadMergeAnchorsPrefixMatchCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "r -p pv ../examples/merge-anchor.yaml foobar.th*")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `foobar.thirty: well beyond
foobar.thing: ice
foobar.thirsty: yep
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadMergeAnchorsListOriginalCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/merge-anchor.yaml foobarList.a")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "original\n", result.Output)
}

func TestReadMergeAnchorsListOverrideInListCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/merge-anchor.yaml foobarList.thing")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "coconut\n", result.Output)
}

func TestReadMergeAnchorsListOverrideCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/merge-anchor.yaml foobarList.c")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "newbar\n", result.Output)
}

func TestReadInvalidDocumentIndexCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -df ../examples/sample.yaml b.c")
	if result.Error == nil {
		t.Error("Expected command to fail due to invalid path")
	}
	expectedOutput := `Document index f is not a integer or *: strconv.ParseInt: parsing "f": invalid syntax`
	test.AssertResult(t, expectedOutput, result.Error.Error())
}

func TestReadBadDocumentIndexCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -d1 ../examples/sample.yaml b.c")
	if result.Error == nil {
		t.Error("Expected command to fail due to invalid path")
	}
	expectedOutput := `Could not process document index 1 as there are only 1 document(s)`
	test.AssertResult(t, expectedOutput, result.Error.Error())
}

func TestReadOrderCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/order.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t,
		`version: 3
application: MyApp
`,
		result.Output)
}

func TestReadMultiCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -d 1 ../examples/multiple_docs.yaml another.document")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "here\n", result.Output)
}

func TestReadMultiWithKeyAndValueCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -p vp -d 1 ../examples/multiple_docs.yaml another.document")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "another.document: here\n", result.Output)
}

func TestReadMultiAllCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -d* ../examples/multiple_docs.yaml commonKey")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t,
		`first document
second document
third document
`, result.Output)
}

func TestReadMultiAllWithKeyAndValueCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -p pv -d* ../examples/multiple_docs.yaml commonKey")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t,
		`commonKey: first document
commonKey: second document
commonKey: third document
`, result.Output)
}

func TestReadCmd_ArrayYaml(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/array.yaml [0].gather_facts")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "false\n", result.Output)
}

func TestReadEmptyContentCmd(t *testing.T) {
	content := ``
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read %s", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := ``
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadEmptyNodesPrintPathCmd(t *testing.T) {
	content := `map: 
  that: {}
array: 
  great: []
null:
  indeed: ~`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read %s -ppv **", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `map.that: {}
array.great: []
null.indeed: ~
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadEmptyContentWithDefaultValueCmd(t *testing.T) {
	content := ``
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read --defaultValue things %s", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `things`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadPrettyPrintCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -P ../examples/sample.json")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `a: Easy! as one two three
b:
  c: 2
  d:
  - 3
  - 4
  e:
  - name: fred
    value: 3
  - name: sam
    value: 4
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadPrettyPrintWithIndentCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -P -I4 ../examples/sample.json")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `a: Easy! as one two three
b:
    c: 2
    d:
      - 3
      - 4
    e:
      - name: fred
        value: 3
      - name: sam
        value: 4
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadCmd_ArrayYaml_NoPath(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/array.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `- become: true
  gather_facts: false
  hosts: lalaland
  name: "Apply smth"
  roles:
  - lala
  - land
  serial: 1
- become: false
  gather_facts: true
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadCmd_ArrayYaml_OneElement(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/array.yaml [0]")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `become: true
gather_facts: false
hosts: lalaland
name: "Apply smth"
roles:
- lala
- land
serial: 1
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadCmd_ArrayYaml_SplatCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/array.yaml [*]")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `become: true
gather_facts: false
hosts: lalaland
name: "Apply smth"
roles:
- lala
- land
serial: 1
become: false
gather_facts: true
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadCmd_ArrayYaml_SplatWithKeyAndValueCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -p pv ../examples/array.yaml [*]")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `'[0]':
  become: true
  gather_facts: false
  hosts: lalaland
  name: "Apply smth"
  roles:
  - lala
  - land
  serial: 1
'[1]':
  become: false
  gather_facts: true
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadCmd_ArrayYaml_SplatWithKeyCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -p p ../examples/array.yaml [*]")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `[0]
[1]
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadCmd_ArrayYaml_SplatKey(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/array.yaml [*].gather_facts")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `false
true
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadCmd_ArrayYaml_ErrorBadPath(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/array.yaml [x].gather_facts")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := ``
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadCmd_ArrayYaml_Splat_ErrorBadPath(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/array.yaml [*].roles[x]")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := ``
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadCmd_Error(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read")
	if result.Error == nil {
		t.Error("Expected command to fail due to missing arg")
	}
	expectedOutput := `Must provide filename`
	test.AssertResult(t, expectedOutput, result.Error.Error())
}

func TestReadCmd_ErrorEmptyFilename(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read  ")
	if result.Error == nil {
		t.Error("Expected command to fail due to missing arg")
	}
	expectedOutput := `Must provide filename`
	test.AssertResult(t, expectedOutput, result.Error.Error())
}

func TestReadCmd_ErrorUnreadableFile(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read fake-unknown")
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
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read %s b.d.*.[x]", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := ``
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadToJsonCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -j ../examples/sample.yaml b")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `{"c":2,"d":[3,4,5],"e":[{"name":"fred","value":3},{"name":"sam","value":4}]}
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadToJsonPrettyCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -j -P ../examples/sample.yaml b")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `{
  "c": 2,
  "d": [
    3,
    4,
    5
  ],
  "e": [
    {
      "name": "fred",
      "value": 3
    },
    {
      "name": "sam",
      "value": 4
    }
  ]
}
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadToJsonPrettyIndentCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read -j -I4 -P ../examples/sample.yaml b")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `{
    "c": 2,
    "d": [
        3,
        4,
        5
    ],
    "e": [
        {
            "name": "fred",
            "value": 3
        },
        {
            "name": "sam",
            "value": 4
        }
    ]
}
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadBadDataCmd(t *testing.T) {
	content := `[!Whatever]`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read %s", filename))
	if result.Error == nil {
		t.Error("Expected command to fail")
	}
	expectedOutput := `yaml: line 1: did not find expected ',' or ']'`
	test.AssertResult(t, expectedOutput, result.Error.Error())
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

func TestReadSplatPrefixCmd(t *testing.T) {
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
	result := test.RunCmd(cmd, fmt.Sprintf("read %s b.there*.c", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}

	expectedOutput := `more things
more things also
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadSplatPrefixWithKeyAndValueCmd(t *testing.T) {
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
	result := test.RunCmd(cmd, fmt.Sprintf("read -p pv %s b.there*.c", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}

	expectedOutput := `b.there.c: more things
b.there2.c: more things also
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadSplatPrefixWithKeyCmd(t *testing.T) {
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
	result := test.RunCmd(cmd, fmt.Sprintf("read -p p %s b.there*.c", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}

	expectedOutput := `b.there.c
b.there2.c
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestPrefixCmd(t *testing.T) {
	content := `b:
  c: 3
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("prefix %s d", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `d:
  b:
    c: 3
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestPrefixCmdArray(t *testing.T) {
	content := `b:
  c: 3
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("prefix %s [+].d.[+]", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `- d:
  - b:
      c: 3
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestPrefixCmd_MultiLayer(t *testing.T) {
	content := `b:
  c: 3
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("prefix %s d.e.f", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `d:
  e:
    f:
      b:
        c: 3
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestPrefixMultiCmd(t *testing.T) {
	content := `b:
  c: 3
---
apples: great
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("prefix %s -d 1 d", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 3
---
d:
  apples: great
`
	test.AssertResult(t, expectedOutput, result.Output)
}
func TestPrefixInvalidDocumentIndexCmd(t *testing.T) {
	content := `b:
  c: 3
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("prefix %s -df d", filename))
	if result.Error == nil {
		t.Error("Expected command to fail due to invalid path")
	}
	expectedOutput := `Document index f is not a integer or *: strconv.ParseInt: parsing "f": invalid syntax`
	test.AssertResult(t, expectedOutput, result.Error.Error())
}

func TestPrefixBadDocumentIndexCmd(t *testing.T) {
	content := `b:
  c: 3
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("prefix %s -d 1 d", filename))
	if result.Error == nil {
		t.Error("Expected command to fail due to invalid path")
	}
	expectedOutput := `asked to process document index 1 but there are only 1 document(s)`
	test.AssertResult(t, expectedOutput, result.Error.Error())
}
func TestPrefixMultiAllCmd(t *testing.T) {
	content := `b:
  c: 3
---
apples: great
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("prefix %s -d * d", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `d:
  b:
    c: 3
---
d:
  apples: great`
	test.AssertResult(t, expectedOutput, strings.Trim(result.Output, "\n "))
}

func TestPrefixCmd_Error(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "prefix")
	if result.Error == nil {
		t.Error("Expected command to fail due to missing arg")
	}
	expectedOutput := `Must provide <filename> <prefixed_path>`
	test.AssertResult(t, expectedOutput, result.Error.Error())
}

func TestPrefixCmd_ErrorUnreadableFile(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "prefix fake-unknown a.b")
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

func TestPrefixCmd_Inplace(t *testing.T) {
	content := `b:
  c: 3
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("prefix -i %s d", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}
	gotOutput := test.ReadTempYamlFile(filename)
	expectedOutput := `d:
  b:
    c: 3`
	test.AssertResult(t, expectedOutput, strings.Trim(gotOutput, "\n "))
}

func TestNewCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "new b.c 3")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 3
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestNewArrayCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "new b[0] 3")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
- 3
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestNewCmd_Error(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "new b.c")
	if result.Error == nil {
		t.Error("Expected command to fail due to missing arg")
	}
	expectedOutput := `Must provide <path_to_update> <value>`
	test.AssertResult(t, expectedOutput, result.Error.Error())
}

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

func TestReadExpression(t *testing.T) {
	content := `name: value`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)
	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("r %s (x==f)", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}

	expectedOutput := ``
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadFindValueArrayCmd(t *testing.T) {
	content := `- cat
- dog
- rat
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)
	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("r %s (.==dog)", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}

	expectedOutput := `dog
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadFindValueDeepArrayCmd(t *testing.T) {
	content := `animals:
  - cat
  - dog
  - rat
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)
	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("r %s animals(.==dog)", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}

	expectedOutput := `dog
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadFindValueDeepObjectCmd(t *testing.T) {
	content := `animals:
  great: yes
  small: sometimes
`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)
	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("r %s animals(.==yes) -ppv", filename))
	if result.Error != nil {
		t.Error(result.Error)
	}

	expectedOutput := `animals.great: yes
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
