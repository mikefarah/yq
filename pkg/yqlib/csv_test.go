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
const csvMissing = `name,numberOfCats,likesApples,height
,null,,168.8
`
const expectedUpdatedSimpleCsv = `name,numberOfCats,likesApples,height
Gary,3,true,168.8
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

const expectedYamlFromCSVMissingData = `- name: Gary
  numberOfCats: 1
  height: 168.8
- name: Samantha's Rabbit
  height: -188.8
  likesApples: false
`

const csvSimpleMissingData = `name,numberOfCats,height
Gary,1,168.8
Samantha's Rabbit,,-188.8
`

const csvTestSimpleYaml = `- [i, like, csv]
- [because, excel, is, cool]`

const expectedSimpleCsv = `i,like,csv
because,excel,is,cool
`

const tsvTestExpectedSimpleCsv = `i	like	csv
because	excel	is	cool
`

var csvScenarios = []formatScenario{
	{
		description:  "Encode CSV simple",
		input:        csvTestSimpleYaml,
		expected:     expectedSimpleCsv,
		scenarioType: "encode-csv",
	},
	{
		description:  "Encode TSV simple",
		input:        csvTestSimpleYaml,
		expected:     tsvTestExpectedSimpleCsv,
		scenarioType: "encode-tsv",
	},
	{
		description:  "Encode Empty",
		skipDoc:      true,
		input:        `[]`,
		expected:     "",
		scenarioType: "encode-csv",
	},
	{
		description:  "Comma in value",
		skipDoc:      true,
		input:        `["comma, in, value", things]`,
		expected:     "\"comma, in, value\",things\n",
		scenarioType: "encode-csv",
	},
	{
		description:  "Encode array of objects to csv",
		input:        expectedYamlFromCSV,
		expected:     csvSimple,
		scenarioType: "encode-csv",
	},
	{
		description:    "Encode array of objects to custom csv format",
		subdescription: "Add the header row manually, then the we convert each object into an array of values - resulting in an array of arrays. Pick the columns and call the header whatever you like.",
		input:          expectedYamlFromCSV,
		expected:       csvSimpleShort,
		expression:     `[["Name", "Number of Cats"]] +  [.[] | [.name, .numberOfCats ]]`,
		scenarioType:   "encode-csv",
	},
	{
		description:    "Encode array of objects to csv - missing fields behaviour",
		subdescription: "First entry is used to determine the headers, and it is missing 'likesApples', so it is not included in the csv. Second entry does not have 'numberOfCats' so that is blank",
		input:          expectedYamlFromCSVMissingData,
		expected:       csvSimpleMissingData,
		scenarioType:   "encode-csv",
	},
	{
		description:  "decode csv missing",
		skipDoc:      true,
		input:        csvMissing,
		expected:     csvMissing,
		scenarioType: "roundtrip-csv",
	},
	{
		description:    "Parse CSV into an array of objects",
		subdescription: "First row is assumed to be the header row.",
		input:          csvSimple,
		expected:       expectedYamlFromCSV,
		scenarioType:   "decode-csv-object",
	},
	{
		description:  "Scalar roundtrip",
		skipDoc:      true,
		input:        "mike\ncat",
		expression:   ".[0].mike",
		expected:     "cat\n",
		scenarioType: "roundtrip-csv",
	},
	{
		description:    "Parse TSV into an array of objects",
		subdescription: "First row is assumed to be the header row.",
		input:          tsvSimple,
		expected:       expectedYamlFromCSV,
		scenarioType:   "decode-tsv-object",
	},
	{
		description:  "Round trip",
		input:        csvSimple,
		expected:     expectedUpdatedSimpleCsv,
		expression:   `(.[] | select(.name == "Gary") | .numberOfCats) = 3`,
		scenarioType: "roundtrip-csv",
	},
}

func testCSVScenario(t *testing.T, s formatScenario) {
	switch s.scenarioType {
	case "encode-csv":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewCsvEncoder(',')), s.description)
	case "encode-tsv":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewCsvEncoder('\t')), s.description)
	case "decode-csv-object":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewCSVObjectDecoder(','), NewYamlEncoder(2, false, ConfiguredYamlPreferences)), s.description)
	case "decode-tsv-object":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewCSVObjectDecoder('\t'), NewYamlEncoder(2, false, ConfiguredYamlPreferences)), s.description)
	case "roundtrip-csv":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewCSVObjectDecoder(','), NewCsvEncoder(',')), s.description)
	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentCSVDecodeObjectScenario(w *bufio.Writer, s formatScenario, formatType string) {
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
		mustProcessFormatScenario(s, NewCSVObjectDecoder(separator), NewYamlEncoder(s.indent, false, ConfiguredYamlPreferences))),
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
		mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewCsvEncoder(separator))),
	)
}

func documentCSVRoundTripScenario(w *bufio.Writer, s formatScenario, formatType string) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, fmt.Sprintf("Given a sample.%v file of:\n", formatType))
	writeOrPanic(w, fmt.Sprintf("```%v\n%v\n```\n", formatType, s.input))

	writeOrPanic(w, "then\n")

	expression := s.expression

	if expression != "" {
		writeOrPanic(w, fmt.Sprintf("```bash\nyq -p=%v -o=%v '%v' sample.%v\n```\n", formatType, formatType, expression, formatType))
	} else {
		writeOrPanic(w, fmt.Sprintf("```bash\nyq -p=%v -o=%v sample.%v\n```\n", formatType, formatType, formatType))
	}
	writeOrPanic(w, "will output\n")

	separator := ','
	if formatType == "tsv" {
		separator = '\t'
	}

	writeOrPanic(w, fmt.Sprintf("```%v\n%v```\n\n", formatType,
		mustProcessFormatScenario(s, NewCSVObjectDecoder(separator), NewCsvEncoder(separator))),
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
		documentCSVDecodeObjectScenario(w, s, "csv")
	case "decode-tsv-object":
		documentCSVDecodeObjectScenario(w, s, "tsv")
	case "roundtrip-csv":
		documentCSVRoundTripScenario(w, s, "csv")

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
