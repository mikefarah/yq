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

	var propsEncoder = NewPropertiesEncoder()
	inputs, err := readDocuments(strings.NewReader(sampleYaml), "sample.yml", 0, NewYamlDecoder())
	if err != nil {
		panic(err)
	}
	node := inputs.Front().Value.(*CandidateNode).Node
	err = propsEncoder.Encode(writer, node)
	if err != nil {
		panic(err)
	}
	writer.Flush()

	return strings.TrimSuffix(output.String(), "\n")
}

func TestPropertiesEncoderSimple(t *testing.T) {
	var sampleYaml = `a: 'bob cool'`

	var expectedJson = `a = bob cool`
	var actualProps = yamlToProps(sampleYaml)
	test.AssertResult(t, expectedJson, actualProps)
}

func TestPropertiesEncoderSimpleWithComments(t *testing.T) {
	var sampleYaml = `a: 'bob cool' # line`

	var expectedJson = `# line
a = bob cool`
	var actualProps = yamlToProps(sampleYaml)
	test.AssertResult(t, expectedJson, actualProps)
}

func TestPropertiesEncoderDeep(t *testing.T) {
	var sampleYaml = `a: 
  b: "bob cool"
`

	var expectedJson = `a.b = bob cool`
	var actualProps = yamlToProps(sampleYaml)
	test.AssertResult(t, expectedJson, actualProps)
}

func TestPropertiesEncoderDeepWithComments(t *testing.T) {
	var sampleYaml = `a:  # a thing
  b: "bob cool" # b thing
`

	var expectedJson = `# b thing
a.b = bob cool`
	var actualProps = yamlToProps(sampleYaml)
	test.AssertResult(t, expectedJson, actualProps)
}

func TestPropertiesEncoderArray(t *testing.T) {
	var sampleYaml = `a: 
  b: [{c: dog}, {c: cat}]
`

	var expectedJson = `a.b.0.c = dog
a.b.1.c = cat`
	var actualProps = yamlToProps(sampleYaml)
	test.AssertResult(t, expectedJson, actualProps)
}
