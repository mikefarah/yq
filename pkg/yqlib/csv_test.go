package yqlib

import (
	"bufio"
	"fmt"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

const csvSimple = `name,numberOfCats,likesApples,height
Gary,1,true,168.8
Samantha's Rabbit,2,false,-188.8
`

const csvSimpleShort = `Name,Number of Cats
Gary,1
Samantha's Rabbit,2
`

const tsvSimple = `name	numberOfCats	likesApples	height
Gary	1	true	168.8
Samantha's Rabbit	2	false	-188.8
`

const expectedYamlFromCSV = `- name: Gary
  numberOfCats: 1
  likesApples: true
  height: 168.8
- name: Samantha's Rabbit
  numberOfCats: 2
  likesApples: false
  height: -188.8
`

const csvTestSimpleYaml = `- [i, like, csv]
- [because, excel, is, cool]`

const csvTestExpectedSimpleCsv = `i,like,csv
because,excel,is,cool
`

const tsvTestExpectedSimpleCsv = `i	like	csv
because	excel	is	cool
`

var csvScenarios = []formatScenario{
	{
		description:  "Encode CSV simple",
		input:        csvTestSimpleYaml,
		expected:     csvTestExpectedSimpleCsv,
		scenarioType: "encode-csv",
	},
	{
		description:  "Encode TSV simple",
		input:        csvTestSimpleYaml,
		expected:     tsvTestExpectedSimpleCsv,
		scenarioType: "encode-tsv",
	},
	{
		description:    "Encode array of objects to csv",
		subdescription: "Add the header row manually, then the we convert each object into an array of values - resulting in an array of arrays. Nice thing about this method is you can pick the columns and call the header whatever you like.",
		input:          expectedYamlFromCSV,
		expected:       csvSimpleShort,
		expression:     `[["Name", "Number of Cats"]] +  [.[] | [.name, .numberOfCats ]]`,
		scenarioType:   "encode-csv",
	},
	{
		description:    "Encode array of objects to csv - generic",
		subdescription: "This is a little trickier than the previous example - we dynamically work out the $header, and use that to automatically create the value arrays.",
		input:          expectedYamlFromCSV,
		expected:       csvSimple,
		expression:     `(.[0] | keys | .[] ) as $header |  [[$header]] +  [.[] | [ .[$header] ]]`,
		scenarioType:   "encode-csv",
	},
	{
		description:    "Parse CSV into an array of objects",
		subdescription: "First row is assumed to define the fields",
		input:          csvSimple,
		expected:       expectedYamlFromCSV,
		scenarioType:   "decode-csv-object",
	},
	{
		description:    "Parse TSV into an array of objects",
		subdescription: "First row is assumed to define the fields",
		input:          tsvSimple,
		expected:       expectedYamlFromCSV,
		scenarioType:   "decode-tsv-object",
	},
}

func testCSVScenario(t *testing.T, s formatScenario) {
	switch s.scenarioType {
	case "encode-csv":
		test.AssertResultWithContext(t, s.expected, processFormatScenario(s, NewYamlDecoder(), NewCsvEncoder(',')), s.description)
	case "encode-tsv":
		test.AssertResultWithContext(t, s.expected, processFormatScenario(s, NewYamlDecoder(), NewCsvEncoder('\t')), s.description)
	case "decode-csv-object":
		test.AssertResultWithContext(t, s.expected, processFormatScenario(s, NewCSVObjectDecoder(','), NewYamlEncoder(2, false, true, true)), s.description)
	case "decode-tsv-object":
		test.AssertResultWithContext(t, s.expected, processFormatScenario(s, NewCSVObjectDecoder('\t'), NewYamlEncoder(2, false, true, true)), s.description)
	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentCSVDecodeObjectScenario(t *testing.T, w *bufio.Writer, s formatScenario, formatType string) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, fmt.Sprintf("Given a sample.%v file of:\n", formatType))
	writeOrPanic(w, fmt.Sprintf("```%v\n%v\n```\n", formatType, s.input))

	writeOrPanic(w, "then\n")
	writeOrPanic(w, fmt.Sprintf("```bash\nyq -p=%v sample.%v\n```\n", formatType, formatType))
	writeOrPanic(w, "will output\n")

	separator := ','
	if formatType == "tsv" {
		separator = '\t'
	}

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n",
		processFormatScenario(s, NewCSVObjectDecoder(separator), NewYamlEncoder(s.indent, false, true, true))),
	)
}

func documentCSVEncodeScenario(w *bufio.Writer, s formatScenario, formatType string) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.yml file of:\n")
	writeOrPanic(w, fmt.Sprintf("```yaml\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")

	expression := s.expression

	if expression != "" {
		writeOrPanic(w, fmt.Sprintf("```bash\nyq -o=%v '%v' sample.yml\n```\n", formatType, expression))
	} else {
		writeOrPanic(w, fmt.Sprintf("```bash\nyq -o=%v sample.yml\n```\n", formatType))
	}
	writeOrPanic(w, "will output\n")

	separator := ','
	if formatType == "tsv" {
		separator = '\t'
	}

	writeOrPanic(w, fmt.Sprintf("```%v\n%v```\n\n", formatType,
		processFormatScenario(s, NewYamlDecoder(), NewCsvEncoder(separator))),
	)
}

func documentCSVScenario(t *testing.T, w *bufio.Writer, i interface{}) {
	s := i.(formatScenario)
	if s.skipDoc {
		return
	}
	switch s.scenarioType {
	case "encode-csv":
		documentCSVEncodeScenario(w, s, "csv")
	case "encode-tsv":
		documentCSVEncodeScenario(w, s, "tsv")
	case "decode-csv-object":
		documentCSVDecodeObjectScenario(t, w, s, "csv")
	case "decode-tsv-object":
		documentCSVDecodeObjectScenario(t, w, s, "tsv")

	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func TestCSVScenarios(t *testing.T) {
	for _, tt := range csvScenarios {
		testCSVScenario(t, tt)
	}
	genericScenarios := make([]interface{}, len(csvScenarios))
	for i, s := range csvScenarios {
		genericScenarios[i] = s
	}
	documentScenarios(t, "usage", "csv-tsv", genericScenarios, documentCSVScenario)
}
