package yqlib

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

func yamlToProps(sampleYaml string) string {
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)

	propsEncoder := NewPropertiesEncoder(writer)
	inputs, err := readDocuments(strings.NewReader(sampleYaml), "sample.yml", 0)
	if err != nil {
		panic(err)
	}
	node := inputs.Front().Value.(*CandidateNode).Node
	err = propsEncoder.Encode(node)
	if err != nil {
		panic(err)
	}
	writer.Flush()

	return strings.TrimSuffix(output.String(), "\n")
}

func TestPropertiesEncoderSimple(t *testing.T) {
	sampleYaml := `a: 'bob cool'`

	expectedJson := `a = bob cool`
	actualProps := yamlToProps(sampleYaml)
	test.AssertResult(t, expectedJson, actualProps)
}

func TestPropertiesEncoderSimpleWithComments(t *testing.T) {
	sampleYaml := `a: 'bob cool' # line`

	expectedJson := `# line
a = bob cool`
	actualProps := yamlToProps(sampleYaml)
	test.AssertResult(t, expectedJson, actualProps)
}

func TestPropertiesEncoderDeep(t *testing.T) {
	sampleYaml := `a: 
  b: "bob cool"
`

	expectedJson := `a.b = bob cool`
	actualProps := yamlToProps(sampleYaml)
	test.AssertResult(t, expectedJson, actualProps)
}

func TestPropertiesEncoderDeepWithComments(t *testing.T) {
	sampleYaml := `a:  # a thing
  b: "bob cool" # b thing
`

	expectedJson := `# b thing
a.b = bob cool`
	actualProps := yamlToProps(sampleYaml)
	test.AssertResult(t, expectedJson, actualProps)
}

func TestPropertiesEncoderArray(t *testing.T) {
	sampleYaml := `a: 
  b: [{c: dog}, {c: cat}]
`

	expectedJson := `a.b.0.c = dog
a.b.1.c = cat`
	actualProps := yamlToProps(sampleYaml)
	test.AssertResult(t, expectedJson, actualProps)
}
