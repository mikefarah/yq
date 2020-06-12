package cmd

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/mikefarah/yq/v3/test"
)

func TestReadCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/sample.yaml b.c")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "2\n", result.Output)
}

func TestReadCmdWithExitStatus(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/sample.yaml b.c -e")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "2\n", result.Output)
}

func TestReadCmdWithExitStatusNotExist(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/sample.yaml caterpillar -e")
	test.AssertResult(t, "No matches found", result.Error.Error())
}

func TestReadCmdNotExist(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/sample.yaml caterpillar")
	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "", result.Output)
}

func TestReadUnwrapCmd(t *testing.T) {

	content := `b: 'frog' # my favourite`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read %s b --unwrapScalar=false", filename))

	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "'frog' # my favourite\n", result.Output)
}

func TestReadStripCommentsCmd(t *testing.T) {

	content := `# this is really cool
b: # my favourite
  c: 5 # cats
# blah
`

	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read %s --stripComments", filename))

	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := `b:
  c: 5
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestReadUnwrapJsonByDefaultCmd(t *testing.T) {

	content := `b: 'frog' # my favourite`
	filename := test.WriteTempYamlFile(content)
	defer test.RemoveTempYamlFile(filename)

	cmd := getRootCommand()
	result := test.RunCmd(cmd, fmt.Sprintf("read %s b -j", filename))

	if result.Error != nil {
		t.Error(result.Error)
	}
	test.AssertResult(t, "\"frog\"\n", result.Output)
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

func TestReadNotFoundWithExitStatus(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/sample.yaml adsf -e")
	if result.Error == nil {
		t.Error("Expected command to fail")
	}
	expectedOutput := `No matches found`
	test.AssertResult(t, expectedOutput, result.Error.Error())
}

func TestReadNotFoundWithoutExitStatus(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "read ../examples/sample.yaml adsf")
	if result.Error != nil {
		t.Error("Expected command to succeed!")
	}
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
