package yqlib

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

func yamlToJson(sampleYaml string, indent int) string {
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)

	jsonEncoder := NewJsonEncoder(writer, indent)
	inputs, err := readDocuments(strings.NewReader(sampleYaml), "sample.yml", 0)
	if err != nil {
		panic(err)
	}
	node := inputs.Front().Value.(*CandidateNode).Node
	err = jsonEncoder.Encode(node)
	if err != nil {
		panic(err)
	}
	writer.Flush()

	return strings.TrimSuffix(output.String(), "\n")
}

func TestJsonEncoderPreservesObjectOrder(t *testing.T) {
	sampleYaml := `zabbix: winner
apple: great
banana:
- {cobra: kai, angus: bob}
`
	expectedJson := `{
  "zabbix": "winner",
  "apple": "great",
  "banana": [
    {
      "cobra": "kai",
      "angus": "bob"
    }
  ]
}`
	actualJson := yamlToJson(sampleYaml, 2)
	test.AssertResult(t, expectedJson, actualJson)
}

func TestJsonNullInArray(t *testing.T) {
	sampleYaml := `[null]`
	actualJson := yamlToJson(sampleYaml, 0)
	test.AssertResult(t, sampleYaml, actualJson)
}

func TestJsonNull(t *testing.T) {
	sampleYaml := `null`
	actualJson := yamlToJson(sampleYaml, 0)
	test.AssertResult(t, sampleYaml, actualJson)
}

func TestJsonNullInObject(t *testing.T) {
	sampleYaml := `{x: null}`
	actualJson := yamlToJson(sampleYaml, 0)
	test.AssertResult(t, `{"x":null}`, actualJson)
}

func TestJsonEncoderDoesNotEscapeHTMLChars(t *testing.T) {
	sampleYaml := `build: "( ./lint && ./format && ./compile ) < src.code"`
	expectedJson := `{"build":"( ./lint && ./format && ./compile ) < src.code"}`
	actualJson := yamlToJson(sampleYaml, 0)
	test.AssertResult(t, expectedJson, actualJson)
}
