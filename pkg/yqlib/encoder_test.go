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

	var jsonEncoder = NewJsonEncoder(writer, indent)
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
	var sampleYaml = `zabbix: winner
apple: great
banana:
- {cobra: kai, angus: bob}
`
	var expectedJson = `{
  "zabbix": "winner",
  "apple": "great",
  "banana": [
    {
      "cobra": "kai",
      "angus": "bob"
    }
  ]
}`
	var actualJson = yamlToJson(sampleYaml, 2)
	test.AssertResult(t, expectedJson, actualJson)
}

func TestJsonEncoderDoesNotEscapeHTMLChars(t *testing.T) {
	var sampleYaml = `build: "( ./lint && ./format && ./compile ) < src.code"`
	var expectedJson = `{"build":"( ./lint && ./format && ./compile ) < src.code"}`
	var actualJson = yamlToJson(sampleYaml, 0)
	test.AssertResult(t, expectedJson, actualJson)
}
