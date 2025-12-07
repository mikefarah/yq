package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var nestedExample = `service "http" "web_proxy" {
  listen_addr = "127.0.0.1:8080"
}`

var nestedExampleYaml = "service:\n  http:\n    web_proxy:\n      listen_addr: \"127.0.0.1:8080\"\n"

var multipleBlockLabelKeys = `service "cat" {
  process "main" {
    command = ["/usr/local/bin/awesome-app", "server"]
  }

  process "management" {
    command = ["/usr/local/bin/awesome-app", "management"]
  }
}
`
var multipleBlockLabelKeysExpected = `service "cat" {
  process "main" {
    command = ["/usr/local/bin/awesome-app", "server"]
  }
  process "management" {
    command = ["/usr/local/bin/awesome-app", "management"]
  }
}
`

var multipleBlockLabelKeysExpectedYaml = `service:
  cat:
    process:
      main:
        command:
          - "/usr/local/bin/awesome-app"
          - "server"
      management:
        command:
          - "/usr/local/bin/awesome-app"
          - "management"
`

var roundtripSample = `# Arithmetic with literals and application-provided variables
sum = 1 + addend

# String interpolation and templates
message = "Hello, ${name}!"

# Application-provided functions
shouty_message = upper(message)`

var roundtripSampleExpected = `# Arithmetic with literals and application-provided variables
sum = 1 + addend
# String interpolation and templates
message = "Hello, ${name}!"
# Application-provided functions
shouty_message = upper(message)
`

var hclFormatScenarios = []formatScenario{
	{
		description:  "Simple decode",
		input:        `io_mode = "async"`,
		expected:     "io_mode: \"async\"\n",
		scenarioType: "decode",
	},
	{
		description:  "Simple decode, no quotes",
		input:        `io_mode = async`,
		expected:     "io_mode: async\n",
		scenarioType: "decode",
	},
	{
		description:  "Simple roundtrip, no quotes",
		input:        `io_mode = async`,
		expected:     "io_mode = async\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "Nested decode",
		input:        nestedExample,
		expected:     nestedExampleYaml,
		scenarioType: "decode",
	},
	{
		description:  "Template decode",
		input:        `message = "Hello, ${name}!"`,
		expected:     "message: \"Hello, ${name}!\"\n",
		scenarioType: "decode",
	},
	{
		description:  "Template roundtrip",
		input:        `message = "Hello, ${name}!"`,
		expected:     "message = \"Hello, ${name}!\"\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "Function roundtrip",
		input:        `shouty_message = upper(message)`,
		expected:     "shouty_message = upper(message)\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "Arithmetic roundtrip",
		input:        `sum = 1 + addend`,
		expected:     "sum = 1 + addend\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "number attribute",
		input:        `port = 8080`,
		expected:     "port: 8080\n",
		scenarioType: "decode",
	},
	{
		description:  "float attribute",
		input:        `pi = 3.14`,
		expected:     "pi: 3.14\n",
		scenarioType: "decode",
	},
	{
		description:  "boolean attribute",
		input:        `enabled = true`,
		expected:     "enabled: true\n",
		scenarioType: "decode",
	},
	{
		description:  "list of strings",
		input:        `tags = ["a", "b"]`,
		expected:     "tags:\n  - \"a\"\n  - \"b\"\n",
		scenarioType: "decode",
	},
	{
		description:  "object/map attribute",
		input:        `obj = { a = 1, b = "two" }`,
		expected:     "obj: {a: 1, b: \"two\"}\n",
		scenarioType: "decode",
	},
	{
		description:  "nested block",
		input:        `server { port = 8080 }`,
		expected:     "server:\n  port: 8080\n",
		scenarioType: "decode",
	},
	{
		description:  "multiple attributes",
		input:        "name = \"app\"\nversion = 1\nenabled = true",
		expected:     "name: \"app\"\nversion: 1\nenabled: true\n",
		scenarioType: "decode",
	},
	{
		description:  "binary expression",
		input:        `count = 0 - 42`,
		expected:     "count: -42\n",
		scenarioType: "decode",
	},
	{
		description:  "negative number",
		input:        `count = -42`,
		expected:     "count: -42\n",
		scenarioType: "decode",
	},
	{
		description:  "scientific notation",
		input:        `value = 1e-3`,
		expected:     "value: 0.001\n",
		scenarioType: "decode",
	},
	{
		description:  "nested object",
		input:        `config = { db = { host = "localhost", port = 5432 } }`,
		expected:     "config: {db: {host: \"localhost\", port: 5432}}\n",
		scenarioType: "decode",
	},
	{
		description:  "mixed list",
		input:        `values = [1, "two", true]`,
		expected:     "values:\n  - 1\n  - \"two\"\n  - true\n",
		scenarioType: "decode",
	},
	{
		description:  "multiple block label keys roundtrip",
		input:        multipleBlockLabelKeys,
		expected:     multipleBlockLabelKeysExpected,
		scenarioType: "roundtrip",
	},
	{
		description:  "multiple block label keys decode",
		input:        multipleBlockLabelKeys,
		expected:     multipleBlockLabelKeysExpectedYaml,
		scenarioType: "decode",
	},
	{
		description:  "block with labels",
		input:        `resource "aws_instance" "example" { ami = "ami-12345" }`,
		expected:     "resource:\n  aws_instance:\n    example:\n      ami: \"ami-12345\"\n",
		scenarioType: "decode",
	},
	{
		description:  "block with labels roundtrip",
		input:        `resource "aws_instance" "example" { ami = "ami-12345" }`,
		expected:     "resource \"aws_instance\" \"example\" {\n  ami = \"ami-12345\"\n}\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip simple attribute",
		input:        `io_mode = "async"`,
		expected:     `io_mode = "async"` + "\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip number attribute",
		input:        `port = 8080`,
		expected:     "port = 8080\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip float attribute",
		input:        `pi = 3.14`,
		expected:     "pi = 3.14\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip boolean attribute",
		input:        `enabled = true`,
		expected:     "enabled = true\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip list of strings",
		input:        `tags = ["a", "b"]`,
		expected:     "tags = [\"a\", \"b\"]\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip object/map attribute",
		input:        `obj = { a = 1, b = "two" }`,
		expected:     "obj = {\n  a = 1\n  b = \"two\"\n}\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip nested block",
		input:        `server { port = 8080 }`,
		expected:     "server {\n  port = 8080\n}\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip multiple attributes",
		input:        "name = \"app\"\nversion = 1\nenabled = true",
		expected:     "name = \"app\"\nversion = 1\nenabled = true\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "decode with comments",
		input:        "# Configuration\nport = 8080 # server port",
		expected:     "# Configuration\nport: 8080 # server port\n",
		scenarioType: "decode",
	},
	{
		description:  "roundtrip with comments",
		input:        "# Configuration\nport = 8080",
		expected:     "# Configuration\nport = 8080\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip example",
		input:        roundtripSample,
		expected:     roundtripSampleExpected,
		scenarioType: "roundtrip",
	},
}

func testHclScenario(t *testing.T, s formatScenario) {
	switch s.scenarioType {
	case "decode":
		result := mustProcessFormatScenario(s, NewHclDecoder(), NewYamlEncoder(ConfiguredYamlPreferences))
		test.AssertResultWithContext(t, s.expected, result, s.description)
	case "roundtrip":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewHclDecoder(), NewHclEncoder()), s.description)
	}
}

func TestHclFormatScenarios(t *testing.T) {
	for _, tt := range hclFormatScenarios {
		testHclScenario(t, tt)
	}
}
