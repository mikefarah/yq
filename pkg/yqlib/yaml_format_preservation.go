package yqlib

import (
	"regexp"
	"strings"

	yaml "go.yaml.in/yaml/v4"
)

// This file implements preservation of the original document layout when
// re-encoding a yaml node tree that was decoded from an existing document
// (e.g. `yq -i 'expr' file.yml`, see https://github.com/mikefarah/yq/issues/465).
//
// The yaml parser keeps comments (HeadComment/LineComment/FootComment) and the
// line/column of every node, but drops blank lines between nodes and the
// whitespace before inline comments. Both can be reconstructed from the
// original source lines:
//
//   - blank lines: the gap between the end line of the previous sibling and
//     the (head comment adjusted) start line of the next sibling tells us how
//     many blank lines preceded it. Prepending that many "\n" to the node's
//     HeadComment makes the emitter re-insert them.
//   - spacing before an inline comment: read the original line the node was
//     on and count the whitespace before the '#'.
//   - indentation: the column of a block collection relative to its key gives
//     the indent used by the original document.

// splitSourceLines splits raw file content into lines, matching the 1-based
// line numbers reported by the yaml parser. Carriage returns are kept so the
// original line ending style can be detected.
func splitSourceLines(content string) []string {
	return strings.Split(content, "\n")
}

// commentLineCount returns how many source lines a comment occupies.
func commentLineCount(comment string) int {
	if comment == "" {
		return 0
	}
	return strings.Count(comment, "\n") + 1
}

// scalarEndLine returns the last source line occupied by a scalar node.
func scalarEndLine(n *yaml.Node) int {
	if n.Kind != yaml.ScalarNode {
		return n.Line
	}
	value := n.Value
	if value == "" {
		return n.Line
	}
	breaks := strings.Count(value, "\n")
	if n.Style&(yaml.LiteralStyle|yaml.FoldedStyle) != 0 {
		// block scalars: content starts on the line after the header line,
		// a trailing newline does not add an extra line.
		if strings.HasSuffix(value, "\n") {
			return n.Line + breaks
		}
		return n.Line + breaks + 1
	}
	// flow scalars live on the node's own line (plus any embedded breaks)
	return n.Line + breaks
}

// nodeEndLine returns the last source line occupied by the node, including
// its subtree and any foot comment.
func nodeEndLine(n *yaml.Node) int {
	end := scalarEndLine(n)
	for _, child := range n.Content {
		if child.Line > 0 {
			end = max(end, nodeEndLine(child))
		}
	}
	if n.FootComment != "" {
		// a foot comment sits on the lines directly after the node
		end += commentLineCount(n.FootComment)
	}
	return end
}

// countBlankSourceLines counts the blank lines in sourceLines (1-based,
// inclusive) between fromLine and toLine.
func countBlankSourceLines(sourceLines []string, fromLine, toLine int) int {
	count := 0
	for line := fromLine; line <= toLine; line++ {
		if line < 1 || line > len(sourceLines) {
			continue
		}
		if strings.TrimSpace(sourceLines[line-1]) == "" {
			count++
		}
	}
	return count
}

// restoreBlankLines walks the node tree and prepends "\n" to head comments so
// that the emitter re-inserts the blank lines that existed between nodes in
// the original document. When the original source lines are available, the
// blank lines are counted directly from the source (which also copes with
// nodes being deleted in the meantime). Nodes without line information
// (newly created) are left alone.
func restoreBlankLines(node *yaml.Node, sourceLines []string) {
	switch node.Kind {
	case yaml.MappingNode, yaml.SequenceNode:
		if node.Style&yaml.FlowStyle != 0 {
			// flow collections are rendered on a single line
			return
		}
		prevEnd := 0
		// the emitter always writes one blank line after a foot comment, so
		// that one must not be counted again when restoring blanks
		prevHadFoot := false
		for i, child := range node.Content {
			if child.Line > 0 && prevEnd > 0 {
				// trailing newlines in a head comment (produced by the parser
				// for CRLF documents) stand for blank lines before the
				// comment; normalise them so they are re-added explicitly.
				headComment := strings.TrimRight(child.HeadComment, "\n")
				trailingBreaks := len(child.HeadComment) - len(headComment)
				child.HeadComment = headComment
				start := child.Line - commentLineCount(headComment)
				blanks := 0
				if start > prevEnd {
					if sourceLines != nil {
						blanks = countBlankSourceLines(sourceLines, prevEnd+1, start-1)
					} else {
						blanks = start - prevEnd - 1 + trailingBreaks
					}
					if prevHadFoot {
						blanks = max(0, blanks-1)
					}
				}
				if blanks > 0 {
					child.HeadComment = strings.Repeat("\n", blanks) + child.HeadComment
				}
			}
			if child.Line > 0 {
				prevEnd = max(prevEnd, nodeEndLine(child))
			}
			// track foot comments of the preceding siblings: a new key starts
			// a new key/value pair, so reset the tracking there
			if node.Kind == yaml.MappingNode && i%2 == 0 {
				prevHadFoot = child.FootComment != ""
			} else {
				prevHadFoot = prevHadFoot || child.FootComment != ""
			}
			restoreBlankLines(child, sourceLines)
		}
	default:
		for _, child := range node.Content {
			restoreBlankLines(child, sourceLines)
		}
	}
}

// detectOriginalIndent inspects the columns of block collections relative to
// their mapping keys to work out the indentation used by the original
// document. Returns false when it cannot be determined (e.g. flat documents).
func detectOriginalIndent(node *yaml.Node) (int, bool) {
	best := 0
	var walk func(n *yaml.Node)
	walk = func(n *yaml.Node) {
		if n.Kind == yaml.MappingNode && n.Style&yaml.FlowStyle == 0 {
			for i := 0; i+1 < len(n.Content); i += 2 {
				key, value := n.Content[i], n.Content[i+1]
				if key.Line <= 0 || value.Column <= key.Column {
					continue
				}
				if value.Kind != yaml.MappingNode && value.Kind != yaml.SequenceNode {
					continue
				}
				if value.Style&yaml.FlowStyle != 0 {
					continue
				}
				delta := value.Column - key.Column
				if delta > 0 && (best == 0 || delta < best) {
					best = delta
				}
			}
		}
		for _, child := range n.Content {
			walk(child)
		}
	}
	walk(node)
	if best >= 2 && best <= 9 {
		return best, true
	}
	return 0, false
}

// collectLineCommentNodes returns, in document order, the nodes that carry an
// inline comment and were decoded from a known source line.
func collectLineCommentNodes(node *yaml.Node) []*yaml.Node {
	var result []*yaml.Node
	var walk func(n *yaml.Node)
	walk = func(n *yaml.Node) {
		if n.LineComment != "" && n.Line > 0 {
			result = append(result, n)
		}
		for _, child := range n.Content {
			walk(child)
		}
	}
	walk(node)
	return result
}

// blockScalarHeaderRe matches lines that end with a block scalar indicator
// (e.g. `key: |`, `- >-`, `key: |2+`). Lines following such a header that are
// more indented belong to the scalar content and must not be touched.
var blockScalarHeaderRe = regexp.MustCompile(`(?:^|\s)[|>][0-9]?[+-]?[0-9]?\s*$`)

func leadingSpaces(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}

// blockScalarHeaderHasPlus reports whether the block scalar header on the
// line uses keep (+) chomping, meaning trailing blank lines are content.
func blockScalarHeaderHasPlus(line string) bool {
	idx := strings.LastIndexAny(line, "|>")
	if idx < 0 {
		return false
	}
	return strings.Contains(line[idx:], "+")
}

// sourceUsesCRLF reports whether the original source used CRLF line endings.
func sourceUsesCRLF(sourceLines []string) bool {
	for _, line := range sourceLines {
		if strings.HasSuffix(line, "\r") {
			return true
		}
	}
	return false
}

// postProcessPreservedYAML cleans up the encoded output so it matches the
// original layout:
//   - blank lines re-inserted via head comments come out indented; strip the
//     trailing whitespace (except inside block scalar content, where
//     whitespace-only lines are meaningful).
//   - inline comments are always emitted with a single space before the '#';
//     restore the original spacing using the source lines.
//   - line endings are converted back to CRLF when the source used CRLF.
func postProcessPreservedYAML(out []byte, sourceLines []string, root *yaml.Node) []byte {
	if len(out) == 0 {
		return out
	}
	crlf := len(sourceLines) > 0 && sourceUsesCRLF(sourceLines)
	text := string(out)
	// split into lines, remembering each line's original ending
	type line struct {
		content string
		ending  string
	}
	var lines []line
	for len(text) > 0 {
		idx := strings.IndexByte(text, '\n')
		if idx < 0 {
			lines = append(lines, line{content: text, ending: ""})
			break
		}
		content := text[:idx]
		ending := "\n"
		if strings.HasSuffix(content, "\r") {
			content = content[:len(content)-1]
			ending = "\r\n"
		}
		lines = append(lines, line{content: content, ending: ending})
		text = text[idx+1:]
	}

	// protect block scalar content lines from whitespace stripping
	protected := make([]bool, len(lines))
	for i := 0; i < len(lines); i++ {
		if !blockScalarHeaderRe.MatchString(lines[i].content) {
			continue
		}
		keyIndent := leadingSpaces(lines[i].content)
		// with keep (+) chomping, trailing blank lines are scalar content
		keepTrailingBlanks := blockScalarHeaderHasPlus(lines[i].content)
		j := i + 1
		for j < len(lines) {
			content := lines[j].content
			if strings.TrimSpace(content) == "" {
				j++
				continue
			}
			if leadingSpaces(content) > keyIndent {
				j++
				continue
			}
			break
		}
		// blank lines at the end of the block are not content (unless '+'):
		// they are usually blank lines re-inserted after the block, which we
		// do want to normalise.
		end := j
		if !keepTrailingBlanks {
			for end > i+1 && strings.TrimSpace(lines[end-1].content) == "" {
				end--
			}
		}
		for k := i + 1; k < end; k++ {
			protected[k] = true
		}
		i = j - 1
	}

	// strip trailing whitespace from blank lines (outside block scalars)
	for i := range lines {
		if protected[i] {
			continue
		}
		if lines[i].content != "" && strings.TrimRight(lines[i].content, " \t") == "" {
			lines[i].content = ""
		}
	}

	// restore the original spacing before inline comments
	if len(sourceLines) > 0 {
		pos := 0
		for _, node := range collectLineCommentNodes(root) {
			comment := node.LineComment
			if !strings.HasPrefix(comment, "#") {
				comment = "#" + comment
			}
			suffix := " " + comment
			for k := pos; k < len(lines); k++ {
				if !strings.HasSuffix(lines[k].content, suffix) {
					continue
				}
				if node.Line-1 < len(sourceLines) {
					if gap := inlineCommentGap(sourceLines[node.Line-1]); gap > 1 {
						lines[k].content = lines[k].content[:len(lines[k].content)-len(comment)-1] +
							strings.Repeat(" ", gap) + comment
					}
				}
				pos = k + 1
				break
			}
		}
	}

	var sb strings.Builder
	for _, l := range lines {
		sb.WriteString(l.content)
		if crlf && l.ending != "" {
			sb.WriteString("\r\n")
		} else {
			sb.WriteString(l.ending)
		}
	}
	return []byte(sb.String())
}

// inlineCommentGap returns the number of whitespace characters between the
// end of the content and the '#' of the first comment on the line.
func inlineCommentGap(sourceLine string) int {
	idx := strings.Index(sourceLine, "#")
	if idx <= 0 {
		return 1
	}
	pre := sourceLine[:idx]
	gap := len(pre) - len(strings.TrimRight(pre, " \t"))
	if gap < 1 {
		return 1
	}
	return gap
}
