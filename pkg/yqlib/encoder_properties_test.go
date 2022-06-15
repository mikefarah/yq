package yqlib

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

func yamlToProps(sampleYaml string, unwrapScalar bool) string {
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)

	var propsEncoder = NewPropertiesEncoder(unwrapScalar)
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

func TestPropertiesEncoderSimple_Unwrapped(t *testing.T) {
	var sampleYaml = `a: 'bob cool'`

	var expectedProps = `a = bob cool`
	var actualProps = yamlToProps(sampleYaml, true)
	test.AssertResult(t, expectedProps, actualProps)
}

func TestPropertiesEncoderSimple_Wrapped(t *testing.T) {
	var sampleYaml = `a: 'bob cool'`

	var expectedProps = `a = "bob cool"`
	var actualProps = yamlToProps(sampleYaml, false)
	test.AssertResult(t, expectedProps, actualProps)
}

func TestPropertiesEncoderSimpleWithComments_Unwrapped(t *testing.T) {
	var sampleYaml = `a: 'bob cool' # line`

	var expectedProps = `# line
a = bob cool`
	var actualProps = yamlToProps(sampleYaml, true)
	test.AssertResult(t, expectedProps, actualProps)
}

func TestPropertiesEncoderSimpleWithComments_Wrapped(t *testing.T) {
	var sampleYaml = `a: 'bob cool' # line`

	var expectedProps = `# line
a = "bob cool"`
	var actualProps = yamlToProps(sampleYaml, false)
	test.AssertResult(t, expectedProps, actualProps)
}

func TestPropertiesEncoderDeep_Unwrapped(t *testing.T) {
	var sampleYaml = `a: 
  b: "bob cool"
`

	var expectedProps = `a.b = bob cool`
	var actualProps = yamlToProps(sampleYaml, true)
	test.AssertResult(t, expectedProps, actualProps)
}

func TestPropertiesEncoderDeep_Wrapped(t *testing.T) {
	var sampleYaml = `a: 
  b: "bob cool"
`

	var expectedProps = `a.b = "bob cool"`
	var actualProps = yamlToProps(sampleYaml, false)
	test.AssertResult(t, expectedProps, actualProps)
}

func TestPropertiesEncoderDeepWithComments_Unwrapped(t *testing.T) {
	var sampleYaml = `a:  # a thing
  b: "bob cool" # b thing
`

	var expectedProps = `# b thing
a.b = bob cool`
	var actualProps = yamlToProps(sampleYaml, true)
	test.AssertResult(t, expectedProps, actualProps)
}

func TestPropertiesEncoderDeepWithComments_Wrapped(t *testing.T) {
	var sampleYaml = `a:  # a thing
  b: "bob cool" # b thing
`

	var expectedProps = `# b thing
a.b = "bob cool"`
	var actualProps = yamlToProps(sampleYaml, false)
	test.AssertResult(t, expectedProps, actualProps)
}

func TestPropertiesEncoderArray_Unwrapped(t *testing.T) {
	var sampleYaml = `a: 
  b: [{c: dog}, {c: cat}]
`

	var expectedProps = `a.b.0.c = dog
a.b.1.c = cat`
	var actualProps = yamlToProps(sampleYaml, true)
	test.AssertResult(t, expectedProps, actualProps)
}

func TestPropertiesEncoderArray_Wrapped(t *testing.T) {
	var sampleYaml = `a: 
  b: [{c: dog named jim}, {c: cat named jam}]
`

	var expectedProps = `a.b.0.c = "dog named jim"
a.b.1.c = "cat named jam"`
	var actualProps = yamlToProps(sampleYaml, false)
	test.AssertResult(t, expectedProps, actualProps)
}
