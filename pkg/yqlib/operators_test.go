package yqlib

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
	yaml "gopkg.in/yaml.v3"
)

type expressionScenario struct {
	description           string
	subdescription        string
	environmentVariable   string
	document              string
	document2             string
	expression            string
	expected              []string
	skipDoc               bool
	dontFormatInputForDoc bool // dont format input doc for documentation generation
}

func readDocumentWithLeadingContent(content string, fakefilename string, fakeFileIndex int) (*list.List, error) {
	reader, firstFileLeadingContent, err := processReadStream(bufio.NewReader(strings.NewReader(content)))
	if err != nil {
		return nil, err
	}

	inputs, err := readDocuments(reader, fakefilename, fakeFileIndex)
	if err != nil {
		return nil, err
	}
	inputs.Front().Value.(*CandidateNode).LeadingContent = firstFileLeadingContent
	return inputs, nil
}

func testScenario(t *testing.T, s *expressionScenario) {
	var err error

	node, err := NewExpressionParser().ParseExpression(s.expression)
	if err != nil {
		t.Error(fmt.Errorf("Error parsing expression %v of %v: %w", s.expression, s.description, err))
		return
	}
	inputs := list.New()

	if s.document != "" {
		inputs, err = readDocumentWithLeadingContent(s.document, "sample.yml", 0)

		if err != nil {
			t.Error(err, s.document, s.expression)
			return
		}

		if s.document2 != "" {
			moreInputs, err := readDocumentWithLeadingContent(s.document2, "another.yml", 1)
			if err != nil {
				t.Error(err, s.document2, s.expression)
				return
			}
			inputs.PushBackList(moreInputs)
		}
	} else {
		candidateNode := &CandidateNode{
			Document:  0,
			Filename:  "",
			Node:      &yaml.Node{Tag: "!!null", Kind: yaml.ScalarNode},
			FileIndex: 0,
		}
		inputs.PushBack(candidateNode)

	}

	if s.environmentVariable != "" {
		os.Setenv("myenv", s.environmentVariable)
	}

	context, err := NewDataTreeNavigator().GetMatchingNodes(Context{MatchingNodes: inputs}, node)

	if err != nil {
		t.Error(fmt.Errorf("%w: %v", err, s.expression))
		return
	}
	test.AssertResultComplexWithContext(t, s.expected, resultsToString(t, context.MatchingNodes), fmt.Sprintf("desc: %v\nexp: %v\ndoc: %v", s.description, s.expression, s.document))
}

func resultsToString(t *testing.T, results *list.List) []string {
	var pretty []string = make([]string, 0)

	for el := results.Front(); el != nil; el = el.Next() {
		n := el.Value.(*CandidateNode)
		var valueBuffer bytes.Buffer
		printer := NewPrinterWithSingleWriter(bufio.NewWriter(&valueBuffer), YamlOutputFormat, true, false, 4, true)

		err := printer.PrintResults(n.AsList())
		if err != nil {
			t.Error(err)
			return nil
		}

		tag := n.Node.Tag
		if n.Node.Kind == yaml.DocumentNode {
			tag = "doc"
		} else if n.Node.Kind == yaml.AliasNode {
			tag = "alias"
		}
		output := fmt.Sprintf(`D%v, P%v, (%v)::%v`, n.Document, n.Path, tag, valueBuffer.String())
		pretty = append(pretty, output)
	}
	return pretty
}

func writeOrPanic(w *bufio.Writer, text string) {
	_, err := w.WriteString(text)
	if err != nil {
		panic(err)
	}
}

func copyFromHeader(title string, out *os.File) error {
	source := fmt.Sprintf("doc/headers/%v.md", title)
	_, err := os.Stat(source)
	if os.IsNotExist(err) {
		return nil
	}
	in, err := os.Open(source) // nolint gosec
	if err != nil {
		return err
	}
	defer safelyCloseFile(in)
	_, err = io.Copy(out, in)
	return err
}

func formatYaml(yaml string, filename string) string {
	var output bytes.Buffer
	printer := NewPrinterWithSingleWriter(bufio.NewWriter(&output), YamlOutputFormat, true, false, 2, true)

	node, err := NewExpressionParser().ParseExpression(".. style= \"\"")
	if err != nil {
		panic(err)
	}
	streamEvaluator := NewStreamEvaluator()
	_, err = streamEvaluator.Evaluate(filename, strings.NewReader(yaml), node, printer, "")
	if err != nil {
		panic(err)
	}
	return output.String()
}

func documentScenarios(t *testing.T, title string, scenarios []expressionScenario) {
	f, err := os.Create(fmt.Sprintf("doc/%v.md", title))

	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()

	err = copyFromHeader(title, f)
	if err != nil {
		t.Error(err)
		return
	}

	w := bufio.NewWriter(f)
	writeOrPanic(w, "\n")

	for _, s := range scenarios {
		if !s.skipDoc {
			documentScenario(t, w, s)
		}
	}
	w.Flush()
}

func documentScenario(t *testing.T, w *bufio.Writer, s expressionScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	formattedDoc, formattedDoc2 := documentInput(w, s)

	writeOrPanic(w, "will output\n")

	documentOutput(t, w, s, formattedDoc, formattedDoc2)
}

func documentInput(w *bufio.Writer, s expressionScenario) (string, string) {
	formattedDoc := ""
	formattedDoc2 := ""
	command := "eval"

	envCommand := ""

	if s.environmentVariable != "" {
		envCommand = fmt.Sprintf("myenv=\"%v\" ", s.environmentVariable)
		os.Setenv("myenv", s.environmentVariable)
	}

	if s.document != "" {
		if s.dontFormatInputForDoc {
			formattedDoc = s.document + "\n"
		} else {
			formattedDoc = formatYaml(s.document, "sample.yml")
		}

		writeOrPanic(w, "Given a sample.yml file of:\n")
		writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n", formattedDoc))

		files := "sample.yml"

		if s.document2 != "" {
			if s.dontFormatInputForDoc {
				formattedDoc2 = s.document2 + "\n"
			} else {
				formattedDoc2 = formatYaml(s.document2, "another.yml")
			}

			writeOrPanic(w, "And another sample another.yml file of:\n")
			writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n", formattedDoc2))
			files = "sample.yml another.yml"
			command = "eval-all"
		}

		writeOrPanic(w, "then\n")

		if s.expression != "" {
			writeOrPanic(w, fmt.Sprintf("```bash\n%vyq %v '%v' %v\n```\n", envCommand, command, s.expression, files))
		} else {
			writeOrPanic(w, fmt.Sprintf("```bash\n%vyq %v %v\n```\n", envCommand, command, files))
		}
	} else {
		writeOrPanic(w, "Running\n")
		writeOrPanic(w, fmt.Sprintf("```bash\n%vyq %v --null-input '%v'\n```\n", envCommand, command, s.expression))
	}
	return formattedDoc, formattedDoc2
}

func documentOutput(t *testing.T, w *bufio.Writer, s expressionScenario, formattedDoc string, formattedDoc2 string) {
	var output bytes.Buffer
	var err error
	printer := NewPrinterWithSingleWriter(bufio.NewWriter(&output), YamlOutputFormat, true, false, 2, true)

	node, err := NewExpressionParser().ParseExpression(s.expression)
	if err != nil {
		t.Error(fmt.Errorf("Error parsing expression %v of %v: %w", s.expression, s.description, err))
		return
	}

	inputs := list.New()

	if s.document != "" {

		inputs, err = readDocumentWithLeadingContent(formattedDoc, "sample.yml", 0)
		if err != nil {
			t.Error(err, s.document, s.expression)
			return
		}
		if s.document2 != "" {
			moreInputs, err := readDocumentWithLeadingContent(formattedDoc2, "another.yml", 1)
			if err != nil {
				t.Error(err, s.document, s.expression)
				return
			}
			inputs.PushBackList(moreInputs)
		}
	} else {
		candidateNode := &CandidateNode{
			Document:  0,
			Filename:  "",
			Node:      &yaml.Node{Tag: "!!null", Kind: yaml.ScalarNode},
			FileIndex: 0,
		}
		inputs.PushBack(candidateNode)

	}

	context, err := NewDataTreeNavigator().GetMatchingNodes(Context{MatchingNodes: inputs}, node)
	if err != nil {
		t.Error(err, s.expression)
	}

	err = printer.PrintResults(context.MatchingNodes)
	if err != nil {
		t.Error(err, s.expression)
	}

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", output.String()))
}
