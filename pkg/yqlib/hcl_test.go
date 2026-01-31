//go:build !yq_nohcl

package yqlib

import (
	"bufio"
	"bytes"
	"fmt"
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

var multipleBlockLabelKeysExpectedUpdate = `service "cat" {
  process "main" {
    command = ["/usr/local/bin/awesome-app", "server", "meow"]
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

var simpleSample = `# Arithmetic with literals and application-provided variables
sum = 1 + addend

# String interpolation and templates
message = "Hello, ${name}!"

# Application-provided functions
shouty_message = upper(message)`

var simpleSampleExpected = `# Arithmetic with literals and application-provided variables
sum = 1 + addend
# String interpolation and templates
message = "Hello, ${name}!"
# Application-provided functions
shouty_message = upper(message)
`

var simpleSampleExpectedYaml = `# Arithmetic with literals and application-provided variables
sum: 1 + addend
# String interpolation and templates
message: "Hello, ${name}!"
# Application-provided functions
shouty_message: upper(message)
`

var hclFormatScenarios = []formatScenario{
	{
		description:  "Parse HCL",
		input:        `io_mode = "async"`,
		expected:     "io_mode: \"async\"\n",
		scenarioType: "decode",
	},
	{
		description:  "Simple decode, no quotes",
		skipDoc:      true,
		input:        `io_mode = async`,
		expected:     "io_mode: async\n",
		scenarioType: "decode",
	},
	{
		description:  "Simple roundtrip, no quotes",
		skipDoc:      true,
		input:        `io_mode = async`,
		expected:     "io_mode = async\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "Nested decode",
		skipDoc:      true,
		input:        nestedExample,
		expected:     nestedExampleYaml,
		scenarioType: "decode",
	},
	{
		description:  "Template decode",
		skipDoc:      true,
		input:        `message = "Hello, ${name}!"`,
		expected:     "message: \"Hello, ${name}!\"\n",
		scenarioType: "decode",
	},
	{
		description:  "Roundtrip: with template",
		skipDoc:      true,
		input:        `message = "Hello, ${name}!"`,
		expected:     "message = \"Hello, ${name}!\"\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip: with function",
		skipDoc:      true,
		input:        `shouty_message = upper(message)`,
		expected:     "shouty_message = upper(message)\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip: with arithmetic",
		skipDoc:      true,
		input:        `sum = 1 + addend`,
		expected:     "sum = 1 + addend\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "Arithmetic decode",
		skipDoc:      true,
		input:        `sum = 1 + addend`,
		expected:     "sum: 1 + addend\n",
		scenarioType: "decode",
	},
	{
		description:  "number attribute",
		skipDoc:      true,
		input:        `port = 8080`,
		expected:     "port: 8080\n",
		scenarioType: "decode",
	},
	{
		description:  "float attribute",
		skipDoc:      true,
		input:        `pi = 3.14`,
		expected:     "pi: 3.14\n",
		scenarioType: "decode",
	},
	{
		description:  "boolean attribute",
		skipDoc:      true,
		input:        `enabled = true`,
		expected:     "enabled: true\n",
		scenarioType: "decode",
	},
	{
		description:  "object/map attribute",
		skipDoc:      true,
		input:        `obj = { a = 1, b = "two" }`,
		expected:     "obj: {a: 1, b: \"two\"}\n",
		scenarioType: "decode",
	},
	{
		description:  "nested block",
		skipDoc:      true,
		input:        `server { port = 8080 }`,
		expected:     "server:\n  port: 8080\n",
		scenarioType: "decode",
	},
	{
		description:  "multiple attributes",
		skipDoc:      true,
		input:        "name = \"app\"\nversion = 1\nenabled = true",
		expected:     "name: \"app\"\nversion: 1\nenabled: true\n",
		scenarioType: "decode",
	},
	{
		description:  "binary expression",
		skipDoc:      true,
		input:        `count = 0 - 42`,
		expected:     "count: -42\n",
		scenarioType: "decode",
	},
	{
		description:  "negative number",
		skipDoc:      true,
		input:        `count = -42`,
		expected:     "count: -42\n",
		scenarioType: "decode",
	},
	{
		description:  "scientific notation",
		skipDoc:      true,
		input:        `value = 1e-3`,
		expected:     "value: 0.001\n",
		scenarioType: "decode",
	},
	{
		description:  "nested object",
		skipDoc:      true,
		input:        `config = { db = { host = "localhost", port = 5432 } }`,
		expected:     "config: {db: {host: \"localhost\", port: 5432}}\n",
		scenarioType: "decode",
	},
	{
		description:  "mixed list",
		skipDoc:      true,
		input:        `values = [1, "two", true]`,
		expected:     "values:\n  - 1\n  - \"two\"\n  - true\n",
		scenarioType: "decode",
	},
	{
		description:  "Roundtrip: Sample Doc",
		input:        multipleBlockLabelKeys,
		expected:     multipleBlockLabelKeysExpected,
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip: With an update",
		input:        multipleBlockLabelKeys,
		expression:   `.service.cat.process.main.command += "meow"`,
		expected:     multipleBlockLabelKeysExpectedUpdate,
		scenarioType: "roundtrip",
	},
	{
		description:  "Parse HCL: Sample Doc",
		input:        multipleBlockLabelKeys,
		expected:     multipleBlockLabelKeysExpectedYaml,
		scenarioType: "decode",
	},
	{
		description:  "block with labels",
		skipDoc:      true,
		input:        `resource "aws_instance" "example" { ami = "ami-12345" }`,
		expected:     "resource:\n  aws_instance:\n    example:\n      ami: \"ami-12345\"\n",
		scenarioType: "decode",
	},
	{
		description:  "block with labels roundtrip",
		skipDoc:      true,
		input:        `resource "aws_instance" "example" { ami = "ami-12345" }`,
		expected:     "resource \"aws_instance\" \"example\" {\n  ami = \"ami-12345\"\n}\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip simple attribute",
		skipDoc:      true,
		input:        `io_mode = "async"`,
		expected:     `io_mode = "async"` + "\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip number attribute",
		skipDoc:      true,
		input:        `port = 8080`,
		expected:     "port = 8080\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip float attribute",
		skipDoc:      true,
		input:        `pi = 3.14`,
		expected:     "pi = 3.14\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip boolean attribute",
		skipDoc:      true,
		input:        `enabled = true`,
		expected:     "enabled = true\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip list of strings",
		skipDoc:      true,
		input:        `tags = ["a", "b"]`,
		expected:     "tags = [\"a\", \"b\"]\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip object/map attribute",
		skipDoc:      true,
		input:        `obj = { a = 1, b = "two" }`,
		expected:     "obj = {\n  a = 1\n  b = \"two\"\n}\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip nested block",
		skipDoc:      true,
		input:        `server { port = 8080 }`,
		expected:     "server {\n  port = 8080\n}\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip multiple attributes",
		skipDoc:      true,
		input:        "name = \"app\"\nversion = 1\nenabled = true",
		expected:     "name = \"app\"\nversion = 1\nenabled = true\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "Parse HCL: with comments",
		input:        "# Configuration\nport = 8080 # server port",
		expected:     "# Configuration\nport: 8080 # server port\n",
		scenarioType: "decode",
	},
	{
		description:  "Roundtrip: with comments",
		input:        "# Configuration\nport = 8080",
		expected:     "# Configuration\nport = 8080\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip:  extraction",
		skipDoc:      true,
		input:        simpleSample,
		expression:   ".shouty_message",
		expected:     "upper(message)\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip: With templates, functions and arithmetic",
		input:        simpleSample,
		expected:     simpleSampleExpected,
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip example",
		skipDoc:      true,
		input:        simpleSample,
		expected:     simpleSampleExpectedYaml,
		scenarioType: "decode",
	},
	{
		description:  "Parse HCL: List of strings",
		skipDoc:      true,
		input:        `tags = ["a", "b"]`,
		expected:     "tags:\n  - \"a\"\n  - \"b\"\n",
		scenarioType: "decode",
	},
	{
		description:  "roundtrip list of objects",
		skipDoc:      true,
		input:        `items = [{ name = "a", value = 1 }, { name = "b", value = 2 }]`,
		expected:     "items = [{\n  name = \"a\"\n  value = 1\n  }, {\n  name = \"b\"\n  value = 2\n}]\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip nested blocks with same name",
		skipDoc:      true,
		input:        "database \"primary\" {\n  host = \"localhost\"\n  port = 5432\n}\ndatabase \"replica\" {\n  host = \"replica.local\"\n  port = 5433\n}",
		expected:     "database \"primary\" {\n  host = \"localhost\"\n  port = 5432\n}\ndatabase \"replica\" {\n  host = \"replica.local\"\n  port = 5433\n}\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip mixed nested structure",
		skipDoc:      true,
		input:        "servers \"web\" {\n  addresses = [\"10.0.1.1\", \"10.0.1.2\"]\n  port = 8080\n}",
		expected:     "servers \"web\" {\n  addresses = [\"10.0.1.1\", \"10.0.1.2\"]\n  port = 8080\n}\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip null value",
		skipDoc:      true,
		input:        `value = null`,
		expected:     "value = null\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip empty list",
		skipDoc:      true,
		input:        `items = []`,
		expected:     "items = []\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip empty object",
		skipDoc:      true,
		input:        `config = {}`,
		expected:     "config = {}\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip: Separate blocks with same name.",
		input:        "resource \"aws_instance\" \"web\" {\n  ami = \"ami-12345\"\n}\nresource \"aws_instance\" \"db\" {\n  ami = \"ami-67890\"\n}",
		expected:     "resource \"aws_instance\" \"web\" {\n  ami = \"ami-12345\"\n}\nresource \"aws_instance\" \"db\" {\n  ami = \"ami-67890\"\n}\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip deeply nested structure",
		skipDoc:      true,
		input:        "app \"database\" \"primary\" \"connection\" {\n  host = \"db.local\"\n  port = 5432\n}",
		expected:     "app \"database\" \"primary\" \"connection\" {\n  host = \"db.local\"\n  port = 5432\n}\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip with leading comments",
		skipDoc:      true,
		input:        "# Main config\nenabled = true\nport = 8080",
		expected:     "# Main config\nenabled = true\nport = 8080\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "Multiple attributes with comments (comment safety with safe path separator)",
		skipDoc:      true,
		input:        "# Database config\ndb_host = \"localhost\"\n# Connection pool\ndb_pool = 10",
		expected:     "# Database config\ndb_host = \"localhost\"\n# Connection pool\ndb_pool = 10\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "Nested blocks with head comments",
		skipDoc:      true,
		input:        "service \"api\" {\n  # Listen address\n  listen = \"0.0.0.0:8080\"\n  # TLS enabled\n  tls = true\n}",
		expected:     "service \"api\" {\n  # Listen address\n  listen = \"0.0.0.0:8080\"\n  # TLS enabled\n  tls = true\n}\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "Multiple blocks with EncodeSeparate preservation",
		skipDoc:      true,
		input:        "resource \"aws_s3_bucket\" \"bucket1\" {\n  bucket = \"my-bucket-1\"\n}\nresource \"aws_s3_bucket\" \"bucket2\" {\n  bucket = \"my-bucket-2\"\n}",
		expected:     "resource \"aws_s3_bucket\" \"bucket1\" {\n  bucket = \"my-bucket-1\"\n}\nresource \"aws_s3_bucket\" \"bucket2\" {\n  bucket = \"my-bucket-2\"\n}\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "Blocks with same name handled separately",
		skipDoc:      true,
		input:        "server \"primary\" { port = 8080 }\nserver \"backup\" { port = 8081 }",
		expected:     "server \"primary\" {\n  port = 8080\n}\nserver \"backup\" {\n  port = 8081\n}\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "Block label with dot roundtrip (commentPathSep)",
		skipDoc:      true,
		input:        "service \"api.service\" {\n  port = 8080\n}",
		expected:     "service \"api.service\" {\n  port = 8080\n}\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "Nested template expression",
		skipDoc:      true,
		input:        `message = "User: ${username}, Role: ${user_role}"`,
		expected:     "message = \"User: ${username}, Role: ${user_role}\"\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "Empty object roundtrip",
		skipDoc:      true,
		input:        `obj = {}`,
		expected:     "obj = {}\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "Null value in block",
		skipDoc:      true,
		input:        `service { optional_field = null }`,
		expected:     "service {\n  optional_field = null\n}\n",
		scenarioType: "roundtrip",
	},
}

func testHclScenario(t *testing.T, s formatScenario) {
	switch s.scenarioType {
	case "decode":
		result := mustProcessFormatScenario(s, NewHclDecoder(), NewYamlEncoder(ConfiguredYamlPreferences))
		test.AssertResultWithContext(t, s.expected, result, s.description)
	case "roundtrip":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewHclDecoder(), NewHclEncoder(ConfiguredHclPreferences)), s.description)
	}
}

func documentHclScenario(_ *testing.T, w *bufio.Writer, i interface{}) {
	s := i.(formatScenario)

	if s.skipDoc {
		return
	}
	switch s.scenarioType {
	case "", "decode":
		documentHclDecodeScenario(w, s)
	case "roundtrip":
		documentHclRoundTripScenario(w, s)
	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentHclDecodeScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.hcl file of:\n")
	writeOrPanic(w, fmt.Sprintf("```hcl\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")
	expression := s.expression
	if s.expression != "" {
		expression = fmt.Sprintf(" '%v'", s.expression)
	}
	writeOrPanic(w, fmt.Sprintf("```bash\nyq -oy%v sample.hcl\n```\n", expression))
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", mustProcessFormatScenario(s, NewHclDecoder(), NewYamlEncoder(ConfiguredYamlPreferences))))
}

func documentHclRoundTripScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.hcl file of:\n")
	writeOrPanic(w, fmt.Sprintf("```hcl\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")
	expression := s.expression
	if s.expression != "" {
		expression = fmt.Sprintf(" '%v'", s.expression)
	}
	writeOrPanic(w, fmt.Sprintf("```bash\nyq%v sample.hcl\n```\n", expression))
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```hcl\n%v```\n\n", mustProcessFormatScenario(s, NewHclDecoder(), NewHclEncoder(ConfiguredHclPreferences))))
}

func TestHclEncoderPrintDocumentSeparator(t *testing.T) {
	encoder := NewHclEncoder(ConfiguredHclPreferences)
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	err := encoder.PrintDocumentSeparator(writer)
	writer.Flush()

	test.AssertResult(t, nil, err)
	test.AssertResult(t, "", buf.String())
}

func TestHclEncoderPrintLeadingContent(t *testing.T) {
	encoder := NewHclEncoder(ConfiguredHclPreferences)
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	err := encoder.PrintLeadingContent(writer, "some content")
	writer.Flush()

	test.AssertResult(t, nil, err)
	test.AssertResult(t, "", buf.String())
}

func TestHclEncoderCanHandleAliases(t *testing.T) {
	encoder := NewHclEncoder(ConfiguredHclPreferences)
	test.AssertResult(t, false, encoder.CanHandleAliases())
}

func TestHclFormatScenarios(t *testing.T) {
	for _, tt := range hclFormatScenarios {
		testHclScenario(t, tt)
	}
	genericScenarios := make([]interface{}, len(hclFormatScenarios))
	for i, s := range hclFormatScenarios {
		genericScenarios[i] = s
	}
	documentScenarios(t, "usage", "hcl", genericScenarios, documentHclScenario)
}
