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
