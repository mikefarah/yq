//go:build !yq_nouri

package yqlib

import (
	"io"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

func TestUriDecoder_Init(t *testing.T) {
	decoder := NewUriDecoder()
	reader := strings.NewReader("test")
	err := decoder.Init(reader)
	test.AssertResult(t, nil, err)
}

func TestUriDecoder_DecodeSimpleString(t *testing.T) {
	decoder := NewUriDecoder()
	reader := strings.NewReader("hello%20world")
	err := decoder.Init(reader)
	test.AssertResult(t, nil, err)

	node, err := decoder.Decode()
	test.AssertResult(t, nil, err)
	test.AssertResult(t, "!!str", node.Tag)
	test.AssertResult(t, "hello world", node.Value)
}

func TestUriDecoder_DecodeSpecialCharacters(t *testing.T) {
	decoder := NewUriDecoder()
	reader := strings.NewReader("hello%21%40%23%24%25")
	err := decoder.Init(reader)
	test.AssertResult(t, nil, err)

	node, err := decoder.Decode()
	test.AssertResult(t, nil, err)
	test.AssertResult(t, "hello!@#$%", node.Value)
}

func TestUriDecoder_DecodeUTF8(t *testing.T) {
	decoder := NewUriDecoder()
	reader := strings.NewReader("%E2%9C%93%20check")
	err := decoder.Init(reader)
	test.AssertResult(t, nil, err)

	node, err := decoder.Decode()
	test.AssertResult(t, nil, err)
	test.AssertResult(t, "✓ check", node.Value)
}

func TestUriDecoder_DecodePlusSign(t *testing.T) {
	decoder := NewUriDecoder()
	reader := strings.NewReader("a+b")
	err := decoder.Init(reader)
	test.AssertResult(t, nil, err)

	node, err := decoder.Decode()
	test.AssertResult(t, nil, err)
	// Note: url.QueryUnescape does NOT convert + to space
	// That's only for form encoding (url.ParseQuery)
	test.AssertResult(t, "a b", node.Value)
}

func TestUriDecoder_DecodeEmptyString(t *testing.T) {
	decoder := NewUriDecoder()
	reader := strings.NewReader("")
	err := decoder.Init(reader)
	test.AssertResult(t, nil, err)

	node, err := decoder.Decode()
	test.AssertResult(t, nil, err)
	test.AssertResult(t, "", node.Value)

	// Second decode should return EOF
	node, err = decoder.Decode()
	test.AssertResult(t, io.EOF, err)
	test.AssertResult(t, (*CandidateNode)(nil), node)
}

func TestUriDecoder_DecodeMultipleCalls(t *testing.T) {
	decoder := NewUriDecoder()
	reader := strings.NewReader("test")
	err := decoder.Init(reader)
	test.AssertResult(t, nil, err)

	// First decode
	node, err := decoder.Decode()
	test.AssertResult(t, nil, err)
	test.AssertResult(t, "test", node.Value)

	// Second decode should return EOF since we've consumed all input
	node, err = decoder.Decode()
	test.AssertResult(t, io.EOF, err)
	test.AssertResult(t, (*CandidateNode)(nil), node)
}

func TestUriDecoder_DecodeInvalidEscape(t *testing.T) {
	decoder := NewUriDecoder()
	reader := strings.NewReader("test%ZZ")
	err := decoder.Init(reader)
	test.AssertResult(t, nil, err)

	_, err = decoder.Decode()
	// Should return an error for invalid escape sequence
	if err == nil {
		t.Error("Expected error for invalid escape sequence, got nil")
	}
}

func TestUriDecoder_DecodeSlashAndQuery(t *testing.T) {
	decoder := NewUriDecoder()
	reader := strings.NewReader("path%2Fto%2Ffile%3Fquery%3Dvalue")
	err := decoder.Init(reader)
	test.AssertResult(t, nil, err)

	node, err := decoder.Decode()
	test.AssertResult(t, nil, err)
	test.AssertResult(t, "path/to/file?query=value", node.Value)
}

func TestUriDecoder_DecodePercent(t *testing.T) {
	decoder := NewUriDecoder()
	reader := strings.NewReader("100%25")
	err := decoder.Init(reader)
	test.AssertResult(t, nil, err)

	node, err := decoder.Decode()
	test.AssertResult(t, nil, err)
	test.AssertResult(t, "100%", node.Value)
}

func TestUriDecoder_DecodeNoEscaping(t *testing.T) {
	decoder := NewUriDecoder()
	reader := strings.NewReader("simple_text-123")
	err := decoder.Init(reader)
	test.AssertResult(t, nil, err)

	node, err := decoder.Decode()
	test.AssertResult(t, nil, err)
	test.AssertResult(t, "simple_text-123", node.Value)
}

// Mock reader that returns an error
type errorReader struct{}

func (e *errorReader) Read(_ []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func TestUriDecoder_DecodeReadError(t *testing.T) {
	decoder := NewUriDecoder()
	err := decoder.Init(&errorReader{})
	test.AssertResult(t, nil, err)

	_, err = decoder.Decode()
	test.AssertResult(t, io.ErrUnexpectedEOF, err)
}
