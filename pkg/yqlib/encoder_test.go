//go:build !yq_nojson

package yqlib

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

func yamlToJSON(sampleYaml string, indent int) string {
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)

	var jsonEncoder = NewJSONEncoder(indent, false, false)
	inputs, err := readDocuments(strings.NewReader(sampleYaml), "sample.yml", 0, NewYamlDecoder(ConfiguredYamlPreferences))
	if err != nil {
		panic(err)
	}
	node := inputs.Front().Value.(*CandidateNode).Node
	err = jsonEncoder.Encode(writer, node)
	if err != nil {
		panic(err)
	}
	writer.Flush()

	return strings.TrimSuffix(output.String(), "\n")
}

func TestJSONEncoderPreservesObjectOrder(t *testing.T) {
	var sampleYaml = `zabbix: winner
apple: great
banana:
- {cobra: kai, angus: bob}
`
	var expectedJSON = `{
  "zabbix": "winner",
  "apple": "great",
  "banana": [
    {
      "cobra": "kai",
      "angus": "bob"
    }
  ]
}`
	var actualJSON = yamlToJSON(sampleYaml, 2)
	test.AssertResult(t, expectedJSON, actualJSON)
}

func TestJsonNullInArray(t *testing.T) {
	var sampleYaml = `[null]`
	var actualJSON = yamlToJSON(sampleYaml, 0)
	test.AssertResult(t, sampleYaml, actualJSON)
}

func TestJsonNull(t *testing.T) {
	var sampleYaml = `null`
	var actualJSON = yamlToJSON(sampleYaml, 0)
	test.AssertResult(t, sampleYaml, actualJSON)
}

func TestJsonNullInObject(t *testing.T) {
	var sampleYaml = `{x: null}`
	var actualJSON = yamlToJSON(sampleYaml, 0)
	test.AssertResult(t, `{"x":null}`, actualJSON)
}

func TestJsonEncoderDoesNotEscapeHTMLChars(t *testing.T) {
	var sampleYaml = `build: "( ./lint && ./format && ./compile ) < src.code"`
	var expectedJSON = `{"build":"( ./lint && ./format && ./compile ) < src.code"}`
	var actualJSON = yamlToJSON(sampleYaml, 0)
	test.AssertResult(t, expectedJSON, actualJSON)
}
