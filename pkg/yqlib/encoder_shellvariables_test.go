package yqlib

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

func assertEncodesTo(t *testing.T, yaml string, shellvars string) {
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)

	var encoder = NewShellVariablesEncoder()
	inputs, err := readDocuments(strings.NewReader(yaml), "test.yml", 0, NewYamlDecoder(ConfiguredYamlPreferences))
	if err != nil {
		panic(err)
	}
	node := inputs.Front().Value.(*CandidateNode)
	err = encoder.Encode(writer, node)
	if err != nil {
		panic(err)
	}
	writer.Flush()

	test.AssertResult(t, shellvars, strings.TrimSuffix(output.String(), "\n"))
}

func TestShellVariablesEncoderNonquoting(t *testing.T) {
	assertEncodesTo(t, "a: alice", "a=alice")
}

func TestShellVariablesEncoderQuoting(t *testing.T) {
	assertEncodesTo(t, "a: Lewis Carroll", "a='Lewis Carroll'")
}

func TestShellVariablesEncoderQuotesQuoting(t *testing.T) {
	assertEncodesTo(t, "a: Lewis Carroll's Alice", "a='Lewis Carroll'\"'\"'s Alice'")
}

func TestShellVariablesEncoderStripComments(t *testing.T) {
	assertEncodesTo(t, "a: Alice # comment", "a=Alice")
}

func TestShellVariablesEncoderMap(t *testing.T) {
	assertEncodesTo(t, "a:\n b: Lewis\n c: Carroll", "a_b=Lewis\na_c=Carroll")
}

func TestShellVariablesEncoderArray_Unwrapped(t *testing.T) {
	assertEncodesTo(t, "a: [{n: Alice}, {n: Bob}]", "a_0_n=Alice\na_1_n=Bob")
}

func TestShellVariablesEncoderKeyNonPrintable(t *testing.T) {
	assertEncodesTo(t, `"be\all": ring!`, "bell='ring!'")
}

func TestShellVariablesEncoderKeyPrintableNonAlphaNumeric(t *testing.T) {
	assertEncodesTo(t, `"b-e l=l": ring!`, "b_e_l_l='ring!'")
}

func TestShellVariablesEncoderKeyPrintableNonAscii(t *testing.T) {
	assertEncodesTo(t, `"b\u00e9ll": ring!`, "bell='ring!'")
}

func TestShellVariablesEncoderRootKeyStartingWithDigit(t *testing.T) {
	assertEncodesTo(t, "1a: onea", "_1a=onea")
}

func TestShellVariablesEncoderRootKeyStartingWithUnderscore(t *testing.T) {
	assertEncodesTo(t, "_key: value", "_key=value")
}

func TestShellVariablesEncoderChildStartingWithUnderscore(t *testing.T) {
	assertEncodesTo(t, "root:\n _child: value", "root__child=value")
}

func TestShellVariablesEncoderEmptyValue(t *testing.T) {
	assertEncodesTo(t, "empty:", "empty=")
}

func TestShellVariablesEncoderEmptyArray(t *testing.T) {
	assertEncodesTo(t, "empty: []", "")
}

func TestShellVariablesEncoderEmptyMap(t *testing.T) {
	assertEncodesTo(t, "empty: {}", "")
}

func TestShellVariablesEncoderScalarNode(t *testing.T) {
	assertEncodesTo(t, "some string", "value='some string'")
}

func assertEncodesToWithSeparator(t *testing.T, yaml string, shellvars string, separator string) {
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)

	// Save the original separator
	originalSeparator := ConfiguredShellVariablesPreferences.KeySeparator
	defer func() {
		ConfiguredShellVariablesPreferences.KeySeparator = originalSeparator
	}()

	// Set the custom separator
	ConfiguredShellVariablesPreferences.KeySeparator = separator

	var encoder = NewShellVariablesEncoder()
	inputs, err := readDocuments(strings.NewReader(yaml), "test.yml", 0, NewYamlDecoder(ConfiguredYamlPreferences))
	if err != nil {
		panic(err)
	}
	node := inputs.Front().Value.(*CandidateNode)
	err = encoder.Encode(writer, node)
	if err != nil {
		panic(err)
	}
	writer.Flush()

	test.AssertResult(t, shellvars, strings.TrimSuffix(output.String(), "\n"))
}

func TestShellVariablesEncoderCustomSeparator(t *testing.T) {
	assertEncodesToWithSeparator(t, "a:\n b: Lewis\n c: Carroll", "a__b=Lewis\na__c=Carroll", "__")
}

func TestShellVariablesEncoderCustomSeparatorNested(t *testing.T) {
	assertEncodesToWithSeparator(t, "my_app:\n db_config:\n  host: localhost", "my_app__db_config__host=localhost", "__")
}

func TestShellVariablesEncoderCustomSeparatorArray(t *testing.T) {
	assertEncodesToWithSeparator(t, "a: [{n: Alice}, {n: Bob}]", "a__0__n=Alice\na__1__n=Bob", "__")
}

func TestShellVariablesEncoderCustomSeparatorSingleChar(t *testing.T) {
	assertEncodesToWithSeparator(t, "a:\n b: value", "aXb=value", "X")
}

func assertEncodesToUnwrapped(t *testing.T, yaml string, shellvars string) {
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)

	originalUnwrapScalar := ConfiguredShellVariablesPreferences.UnwrapScalar
	defer func() {
		ConfiguredShellVariablesPreferences.UnwrapScalar = originalUnwrapScalar
	}()

	ConfiguredShellVariablesPreferences.UnwrapScalar = true

	var encoder = NewShellVariablesEncoder()
	inputs, err := readDocuments(strings.NewReader(yaml), "test.yml", 0, NewYamlDecoder(ConfiguredYamlPreferences))
	if err != nil {
		panic(err)
	}
	node := inputs.Front().Value.(*CandidateNode)
	err = encoder.Encode(writer, node)
	if err != nil {
		panic(err)
	}
	writer.Flush()

	test.AssertResult(t, shellvars, strings.TrimSuffix(output.String(), "\n"))
}

func TestShellVariablesEncoderUnwrapScalar(t *testing.T) {
	assertEncodesToUnwrapped(t, "a: Lewis Carroll", "a=Lewis Carroll")
	assertEncodesToUnwrapped(t, "b: 123", "b=123")
	assertEncodesToUnwrapped(t, "c: true", "c=true")
	assertEncodesToUnwrapped(t, "d: value with spaces", "d=value with spaces")
}

func TestShellVariablesEncoderNumericSeparators(t *testing.T) {
	cases := []struct {
		name      string
		yaml      string
		quoted    string
		unwrapped string
	}{
		{
			name:      "nested integer",
			yaml:      "app:\n  server:\n    host: localhost\n    port: 8080\n    buffer_size: 20_000",
			quoted:    "app_server_host=localhost\napp_server_port=8080\napp_server_buffer_size=20000",
			unwrapped: "app_server_host=localhost\napp_server_port=8080\napp_server_buffer_size=20000",
		},
		{
			name:      "signed integers",
			yaml:      "negative: -20_000\npositive: +20_000",
			quoted:    "negative='-20000'\npositive='+20000'",
			unwrapped: "negative=-20000\npositive=+20000",
		},
		{
			name:      "numeric bases",
			yaml:      "hex: 0xA_B\noctal: 0o1_7\nbinary: 0b1_0",
			quoted:    "hex=0xAB\noctal=0o17\nbinary=0b10",
			unwrapped: "hex=0xAB\noctal=0o17\nbinary=0b10",
		},
		{
			name:      "float and exponent",
			yaml:      "float: 1_000.5_0\nexponent: 1_0.5e+1_0",
			quoted:    "float='1000.50'\nexponent='10.5e+10'",
			unwrapped: "float=1000.50\nexponent=10.5e+10",
		},
		{
			name:      "large integer",
			yaml:      "value: 18_446_744_073_709_551_615",
			quoted:    "value=18446744073709551615",
			unwrapped: "value=18446744073709551615",
		},
		{
			name:      "quoted strings",
			yaml:      "single: '20_000'\ndouble: \"20_000\"\nplain: text_with_underscores",
			quoted:    "single=20_000\ndouble=20_000\nplain=text_with_underscores",
			unwrapped: "single=20_000\ndouble=20_000\nplain=text_with_underscores",
		},
		{
			name:      "custom tag",
			yaml:      "custom: !amount 20_000",
			quoted:    "custom=20_000",
			unwrapped: "custom=20_000",
		},
		{
			name:      "key underscores",
			yaml:      "key_1_000: 2_000\n1_000: 3_000",
			quoted:    "key_1_000=2000\n_1_000=3000",
			unwrapped: "key_1_000=2000\n_1_000=3000",
		},
		{
			name:      "sequence",
			yaml:      "values: [1_000, '2_000', 3_000.0]",
			quoted:    "values_0=1000\nvalues_1=2_000\nvalues_2='3000.0'",
			unwrapped: "values_0=1000\nvalues_1=2_000\nvalues_2=3000.0",
		},
		{
			name:      "root scalar",
			yaml:      "20_000",
			quoted:    "value=20000",
			unwrapped: "value=20000",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("quoted", func(t *testing.T) {
				assertEncodesTo(t, tc.yaml, tc.quoted)
			})
			t.Run("unwrapped", func(t *testing.T) {
				assertEncodesToUnwrapped(t, tc.yaml, tc.unwrapped)
			})
		})
	}
}

func TestShellVariablesEncoderRetainsNumericSourceValue(t *testing.T) {
	inputs, err := readDocuments(strings.NewReader("number: 20_000"), "test.yml", 0, NewYamlDecoder(ConfiguredYamlPreferences))
	if err != nil {
		t.Fatal(err)
	}
	node := inputs.Front().Value.(*CandidateNode)
	value := node.Content[1]
	var output bytes.Buffer

	if err := NewShellVariablesEncoder().Encode(&output, node); err != nil {
		t.Fatal(err)
	}

	test.AssertResult(t, "number=20000\n", output.String())
	test.AssertResult(t, "20_000", value.Value)
	test.AssertResult(t, "!!int", value.Tag)
}
