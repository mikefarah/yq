package yqlib

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

func yamlToCsv(sampleYaml string, separator rune) string {
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)

	var jsonEncoder = NewCsvEncoder(writer, separator)
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

var sampleYaml = `["apple", apple2, "comma, in, value", "new
line", 3, 3.40, true, "tab	here"]`

func TestCsvEncoder(t *testing.T) {
	var expectedCsv = `apple,apple2,"comma, in, value",new line,3,3.40,true,tab	here`

	var actualCsv = yamlToCsv(sampleYaml, ',')
	test.AssertResult(t, expectedCsv, actualCsv)
}

func TestTsvEncoder(t *testing.T) {

	var expectedCsv = `apple	apple2	comma, in, value	new line	3	3.40	true	"tab	here"`

	var actualCsv = yamlToCsv(sampleYaml, '\t')
	test.AssertResult(t, expectedCsv, actualCsv)
}
