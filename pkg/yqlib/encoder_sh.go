package yqlib

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

var unsafeChars = regexp.MustCompile(`[^\w@%+=:,./-]`)

type shEncoder struct {
	quoteAll bool
}

func NewShEncoder() Encoder {
	return &shEncoder{false}
}

func (e *shEncoder) CanHandleAliases() bool {
	return false
}

func (e *shEncoder) PrintDocumentSeparator(_ io.Writer) error {
	return nil
}

func (e *shEncoder) PrintLeadingContent(_ io.Writer, _ string) error {
	return nil
}

func (e *shEncoder) Encode(writer io.Writer, node *CandidateNode) error {
	if node.guessTagFromCustomType() != "!!str" {
		return fmt.Errorf("cannot encode %v as URI, can only operate on strings. Please first pipe through another encoding operator to convert the value to a string", node.Tag)
	}

	return writeString(writer, e.encode(node.Value))
}

// put any (shell-unsafe) characters into a single-quoted block, close the block lazily
func (e *shEncoder) encode(input string) string {
	const quote = '\''
	var inQuoteBlock = false
	var encoded strings.Builder
	encoded.Grow(len(input))

	for _, ir := range input {
		// open or close a single-quote block
		if ir == quote {
			if inQuoteBlock {
				// get out of a quote block for an input quote
				encoded.WriteRune(quote)
				inQuoteBlock = !inQuoteBlock
			}
			// escape the quote with a backslash
			encoded.WriteRune('\\')
		} else {
			if e.shouldQuote(ir) && !inQuoteBlock {
				// start a quote block for any (unsafe) characters
				encoded.WriteRune(quote)
				inQuoteBlock = !inQuoteBlock
			}
		}
		// pass on the input character
		encoded.WriteRune(ir)
	}
	// close any pending quote block
	if inQuoteBlock {
		encoded.WriteRune(quote)
	}
	return encoded.String()
}

func (e *shEncoder) shouldQuote(ir rune) bool {
	return e.quoteAll || unsafeChars.MatchString(string(ir))
}
