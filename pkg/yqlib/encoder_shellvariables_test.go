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
	node := inputs.Front().Value.(*CandidateNode).Node
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
	assertEncodesTo(t, `"bèll": ring!`, "bell='ring!'")
}

func TestShellVariablesEncoderRootKeyStartingWithDigit(t *testing.T) {
	assertEncodesTo(t, "1a: onea", "_1a=onea")
}

func TestShellVariablesEncoderEmptyValue(t *testing.T) {
	assertEncodesTo(t, "empty:", "empty=")
}

func TestShellVariablesEncoderScalarNode(t *testing.T) {
	assertEncodesTo(t, "some string", "value='some string'")
}
