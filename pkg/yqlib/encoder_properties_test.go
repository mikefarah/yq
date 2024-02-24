package yqlib

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

type keyValuePair struct {
	key     string
	value   string
	comment string
}

func (kv *keyValuePair) String(unwrap bool, sep string) string {
	builder := strings.Builder{}

	if kv.comment != "" {
		builder.WriteString(kv.comment)
		builder.WriteString("\n")
	}

	builder.WriteString(kv.key)
	builder.WriteString(sep)

	if unwrap {
		builder.WriteString(kv.value)
	} else {
		builder.WriteString("\"")
		builder.WriteString(kv.value)
		builder.WriteString("\"")
	}

	return builder.String()
}

type testProperties struct {
	pairs []keyValuePair
}

func (tp *testProperties) String(unwrap bool, sep string) string {
	kvs := []string{}

	for _, kv := range tp.pairs {
		kvs = append(kvs, kv.String(unwrap, sep))
	}

	return strings.Join(kvs, "\n")
}

func yamlToProps(sampleYaml string, unwrapScalar bool, separator string) string {
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)

	var propsEncoder = NewPropertiesEncoder(PropertiesPreferences{KeyValueSeparator: separator, UnwrapScalar: unwrapScalar})
	inputs, err := readDocuments(strings.NewReader(sampleYaml), "sample.yml", 0, NewYamlDecoder(ConfiguredYamlPreferences))
	if err != nil {
		panic(err)
	}
	node := inputs.Front().Value.(*CandidateNode)
	err = propsEncoder.Encode(writer, node)
	if err != nil {
		panic(err)
	}
	writer.Flush()

	return strings.TrimSuffix(output.String(), "\n")
}

func doTest(t *testing.T, sampleYaml string, props testProperties, testUnwrapped, testWrapped bool) {
	wraps := []bool{}
	if testUnwrapped {
		wraps = append(wraps, true)
	}
	if testWrapped {
		wraps = append(wraps, false)
	}

	for _, unwrap := range wraps {
		for _, sep := range []string{" = ", ";", "=", " "} {
			var actualProps = yamlToProps(sampleYaml, unwrap, sep)
			test.AssertResult(t, props.String(unwrap, sep), actualProps)
		}
	}
}

func TestPropertiesEncoderSimple(t *testing.T) {
	var sampleYaml = `a: 'bob cool'`

	doTest(
		t, sampleYaml,
		testProperties{
			pairs: []keyValuePair{
				{
					key:   "a",
					value: "bob cool",
				},
			},
		},
		true, true,
	)
}

func TestPropertiesEncoderSimpleWithComments(t *testing.T) {
	var sampleYaml = `a: 'bob cool' # line`

	doTest(
		t, sampleYaml,
		testProperties{
			pairs: []keyValuePair{
				{
					key:     "a",
					value:   "bob cool",
					comment: "# line",
				},
			},
		},
		true, true,
	)
}

func TestPropertiesEncoderDeep(t *testing.T) {
	var sampleYaml = `a: 
  b: "bob cool"
`

	doTest(
		t, sampleYaml,
		testProperties{
			pairs: []keyValuePair{
				{
					key:   "a.b",
					value: "bob cool",
				},
			},
		},
		true, true,
	)
}

func TestPropertiesEncoderDeepWithComments(t *testing.T) {
	var sampleYaml = `a:  # a thing
  b: "bob cool" # b thing
`

	doTest(
		t, sampleYaml,
		testProperties{
			pairs: []keyValuePair{
				{
					key:     "a.b",
					value:   "bob cool",
					comment: "# b thing",
				},
			},
		},
		true, true,
	)
}

func TestPropertiesEncoderArray_Unwrapped(t *testing.T) {
	var sampleYaml = `a: 
  b: [{c: dog}, {c: cat}]
`

	doTest(
		t, sampleYaml,
		testProperties{
			pairs: []keyValuePair{
				{
					key:   "a.b.0.c",
					value: "dog",
				},
				{
					key:   "a.b.1.c",
					value: "cat",
				},
			},
		},
		true, false,
	)
}

func TestPropertiesEncoderArray_Wrapped(t *testing.T) {
	var sampleYaml = `a: 
  b: [{c: dog named jim}, {c: cat named jim}]
`

	doTest(
		t, sampleYaml,
		testProperties{
			pairs: []keyValuePair{
				{
					key:   "a.b.0.c",
					value: "dog named jim",
				},
				{
					key:   "a.b.1.c",
					value: "cat named jim",
				},
			},
		},
		false, true,
	)
}
