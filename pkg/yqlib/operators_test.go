package yqlib

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"io"
	"io/fs"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/mikefarah/yq/v4/test"
	logging "gopkg.in/op/go-logging.v1"
)

type expressionScenario struct {
	description           string
	subdescription        string
	explanation           []string
	environmentVariables  map[string]string
	document              string
	document2             string
	expression            string
	expected              []string
	skipDoc               bool
	expectedError         string
	dontFormatInputForDoc bool // dont format input doc for documentation generation
	requiresFormat        string
	skipForGoccy          bool
}

var goccyTesting = false

var testingDecoder = NewYamlDecoder(ConfiguredYamlPreferences)

func TestMain(m *testing.M) {
	logging.SetLevel(logging.WARNING, "")
	if os.Getenv("DEBUG") == "true" {
		logging.SetLevel(logging.DEBUG, "")
	}
	ConfiguredYamlPreferences.ColorsEnabled = false
	ConfiguredJSONPreferences.ColorsEnabled = false

	goccyTesting = os.Getenv("GOCCY") == "true"

	if goccyTesting {
		testingDecoder = NewGoccyYAMLDecoder()
	}

	Now = func() time.Time {
		return time.Date(2021, time.May, 19, 1, 2, 3, 4, time.UTC)
	}
	code := m.Run()
	os.Exit(code)
}

func NewSimpleYamlPrinter(writer io.Writer, unwrapScalar bool, indent int, printDocSeparators bool) Printer {
	prefs := ConfiguredYamlPreferences.Copy()
	prefs.PrintDocSeparators = printDocSeparators
	prefs.UnwrapScalar = unwrapScalar
	prefs.Indent = indent
	return NewPrinter(NewYamlEncoder(prefs), NewSinglePrinterWriter(writer))
}

func readDocument(content string, fakefilename string, fakeFileIndex int) (*list.List, error) {
	reader := bufio.NewReader(strings.NewReader(content))

	return readDocuments(reader, fakefilename, fakeFileIndex, testingDecoder)
}

func testScenario(t *testing.T, s *expressionScenario) {
	if s.skipForGoccy {
		return
	}
	log.Debugf("\n\ntesting scenario %v", s.description)
	var err error
	node, err := getExpressionParser().ParseExpression(s.expression)
	if err != nil {
		t.Error(fmt.Errorf("Error parsing expression %v of %v: %w", s.expression, s.description, err))
		return
	}
	inputs := list.New()

	if s.document != "" {
		inputs, err = readDocument(s.document, "sample.yml", 0)

		if err != nil {
			t.Error(err, s.document, s.expression)
			return
		}

		if s.document2 != "" {
			moreInputs, err := readDocument(s.document2, "another.yml", 1)
			if err != nil {
				t.Error(err, s.document2, s.expression)
				return
			}
			inputs.PushBackList(moreInputs)
		}
	} else {
		candidateNode := &CandidateNode{
			document:  0,
			filename:  "",
			Tag:       "!!null",
			Kind:      ScalarNode,
			fileIndex: 0,
		}
		inputs.PushBack(candidateNode)

	}

	for name, value := range s.environmentVariables {
		os.Setenv(name, value)
	}

	context, err := NewDataTreeNavigator().GetMatchingNodes(Context{MatchingNodes: inputs}, node)

	if s.expectedError != "" {
		if err == nil {
			t.Errorf("Expected error '%v' but it worked!", s.expectedError)
		} else {
			test.AssertResultComplexWithContext(t, s.expectedError, err.Error(), fmt.Sprintf("desc: %v\nexp: %v\ndoc: %v", s.description, s.expression, s.document))
		}
		return
	}

	if s.requiresFormat != "" {
		format := s.requiresFormat
		inputFormat, err := FormatFromString(format)
		if err != nil {
			t.Error(err)
		}
		if decoder := inputFormat.DecoderFactory(); decoder == nil {
			t.Skipf("no support for %s input format", format)
		}
		outputFormat, err := FormatFromString(format)
		if err != nil {
			t.Error(err)
		}
		if encoder := configureEncoder(outputFormat, 4); encoder == nil {
			t.Skipf("no support for %s output format", format)
		}
	}
	if err != nil {
		t.Error(fmt.Errorf("%w: %v: %v", err, s.description, s.expression))
		return
	}
	test.AssertResultComplexWithContext(t, s.expected, resultsToString(t, context.MatchingNodes), fmt.Sprintf("desc: %v\nexp: %v\ndoc: %v", s.description, s.expression, s.document))
}

func resultToString(t *testing.T, n *CandidateNode) string {
	var valueBuffer bytes.Buffer
	log.Debugf("printing result %v", NodeToString(n))
	printer := NewSimpleYamlPrinter(bufio.NewWriter(&valueBuffer), true, 4, true)

	err := printer.PrintResults(n.AsList())
	if err != nil {
		t.Error(err)
		return ""
	}

	tag := n.Tag
	if n.Kind == AliasNode {
		tag = "alias"
	}
	return fmt.Sprintf(`D%v, P%v, (%v)::%v`, n.GetDocument(), n.GetPath(), tag, valueBuffer.String())
}

func resultsToString(t *testing.T, results *list.List) []string {
	var pretty = make([]string, 0)

	for el := results.Front(); el != nil; el = el.Next() {
		n := el.Value.(*CandidateNode)

		output := resultToString(t, n)
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

func copySnippet(source string, out *os.File) error {
	_, err := os.Stat(source)
	if os.IsNotExist(err) {
		return nil
	}
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer safelyCloseFile(in)
	_, err = io.Copy(out, in)
	return err
}

func formatYaml(yaml string, filename string) string {
	var output bytes.Buffer
	printer := NewSimpleYamlPrinter(bufio.NewWriter(&output), true, 2, true)

	node, err := getExpressionParser().ParseExpression(".. style= \"\"")
	if err != nil {
		panic(err)
	}
	streamEvaluator := NewStreamEvaluator()
	_, err = streamEvaluator.Evaluate(filename, strings.NewReader(yaml), node, printer, testingDecoder)
	if err != nil {
		panic(err)
	}
	return output.String()
}

type documentScenarioFunc func(t *testing.T, writer *bufio.Writer, scenario interface{})

func documentScenarios(t *testing.T, folder string, title string, scenarios []interface{}, documentScenario documentScenarioFunc) {
	filename := fmt.Sprintf("doc/%v/%v.md", folder, title)
	f, err := os.Create(filename)

	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()

	source := fmt.Sprintf("doc/%v/headers/%v.md", folder, title)
	err = copySnippet(source, f)
	if err != nil {
		t.Error(err)
		return
	}

	err = copySnippet("doc/notification-snippet.md", f)
	if err != nil {
		t.Error(err)
		return
	}

	w := bufio.NewWriter(f)
	writeOrPanic(w, "\n")

	for _, s := range scenarios {
		documentScenario(t, w, s)
	}
	w.Flush()
}

func appendOperatorDocumentScenario(t *testing.T, title string, scenarios []expressionScenario) {
	filename := fmt.Sprintf("doc/%v/%v.md", "operators", title)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, fs.ModeAppend)
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, s := range scenarios {
		documentOperatorScenario(t, w, s)
	}
	w.Flush()

}

func documentOperatorScenarios(t *testing.T, title string, scenarios []expressionScenario) {
	genericScenarios := make([]interface{}, len(scenarios))
	for i, s := range scenarios {
		genericScenarios[i] = s
	}

	documentScenarios(t, "operators", title, genericScenarios, documentOperatorScenario)
}

func documentOperatorScenario(t *testing.T, w *bufio.Writer, i interface{}) {
	s := i.(expressionScenario)

	if s.skipDoc {
		return
	}
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	formattedDoc, formattedDoc2 := documentInput(w, s)

	writeOrPanic(w, "will output\n")

	documentOutput(t, w, s, formattedDoc, formattedDoc2)

	if len(s.explanation) > 0 {
		writeOrPanic(w, "### Explanation:\n")
		for _, text := range s.explanation {
			writeOrPanic(w, fmt.Sprintf("- %v\n", text))
		}
		writeOrPanic(w, "\n")
	}
}

func documentInput(w *bufio.Writer, s expressionScenario) (string, string) {
	formattedDoc := ""
	formattedDoc2 := ""
	command := ""

	envCommand := ""

	envKeys := make([]string, 0, len(s.environmentVariables))
	for k := range s.environmentVariables {
		envKeys = append(envKeys, k)
	}
	sort.Strings(envKeys)

	for _, name := range envKeys {
		value := s.environmentVariables[name]
		if envCommand == "" {
			envCommand = fmt.Sprintf("%v=\"%v\" ", name, value)
		} else {
			envCommand = fmt.Sprintf("%v %v=\"%v\" ", envCommand, name, value)
		}
		os.Setenv(name, value)
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
			command = "eval-all "
		}

		writeOrPanic(w, "then\n")

		if s.expression != "" {
			writeOrPanic(w, fmt.Sprintf("```bash\n%vyq %v'%v' %v\n```\n", envCommand, command, strings.ReplaceAll(s.expression, "'", `'\''`), files))
		} else {
			writeOrPanic(w, fmt.Sprintf("```bash\n%vyq %v%v\n```\n", envCommand, command, files))
		}
	} else {
		writeOrPanic(w, "Running\n")
		writeOrPanic(w, fmt.Sprintf("```bash\n%vyq %v--null-input '%v'\n```\n", envCommand, command, s.expression))
	}
	return formattedDoc, formattedDoc2
}

func documentOutput(t *testing.T, w *bufio.Writer, s expressionScenario, formattedDoc string, formattedDoc2 string) {
	var output bytes.Buffer
	var err error
	printer := NewSimpleYamlPrinter(bufio.NewWriter(&output), true, 2, true)

	node, err := getExpressionParser().ParseExpression(s.expression)
	if err != nil {
		t.Error(fmt.Errorf("Error parsing expression %v of %v: %w", s.expression, s.description, err))
		return
	}

	inputs := list.New()

	if s.document != "" {

		inputs, err = readDocument(formattedDoc, "sample.yml", 0)
		if err != nil {
			t.Error(err, formattedDoc, "exp: "+s.expression)
			return
		}
		if s.document2 != "" {
			moreInputs, err := readDocument(formattedDoc2, "another.yml", 1)
			if err != nil {
				t.Error(err, formattedDoc2, "exp: "+s.expression)
				return
			}
			inputs.PushBackList(moreInputs)
		}
	} else {
		candidateNode := &CandidateNode{
			document:  0,
			filename:  "",
			Tag:       "!!null",
			Kind:      ScalarNode,
			fileIndex: 0,
		}
		inputs.PushBack(candidateNode)

	}

	context, err := NewDataTreeNavigator().GetMatchingNodes(Context{MatchingNodes: inputs}, node)

	if s.expectedError != "" && err != nil {
		writeOrPanic(w, fmt.Sprintf("```bash\nError: %v\n```\n\n", err.Error()))
		return
	} else if err != nil {
		t.Error(err, s.expression)
		return
	}

	err = printer.PrintResults(context.MatchingNodes)
	if err != nil {
		t.Error(err, s.expression)
	}

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", output.String()))
}
