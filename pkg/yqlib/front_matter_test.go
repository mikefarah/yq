package yqlib

import (
	"io"
	"os"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

func createTestFile(content string) string {
	tempFile, err := createTempFile()
	if err != nil {
		panic(err)
	}

	_, err = tempFile.Write([]byte(content))
	if err != nil {
		panic(err)
	}

	safelyCloseFile(tempFile)

	return tempFile.Name()
}

func readFile(filename string) string {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func TestFrontMatterSplitWithLeadingSep(t *testing.T) {
	file := createTestFile(`---
a: apple
b: banana
---
not a 
yaml: doc
`)

	expectedYamlFm := `---
a: apple
b: banana
`

	expectedContent := `---
not a 
yaml: doc
`

	fmHandler := NewFrontMatterHandler(file)
	err := fmHandler.Split()
	if err != nil {
		panic(err)
	}

	yamlFm := readFile(fmHandler.GetYamlFrontMatterFilename())

	test.AssertResult(t, expectedYamlFm, yamlFm)

	contentBytes, err := io.ReadAll(fmHandler.GetContentReader())
	if err != nil {
		panic(err)
	}
	test.AssertResult(t, expectedContent, string(contentBytes))

	tryRemoveTempFile(file)
	fmHandler.CleanUp()
}

func TestFrontMatterSplitWithNoLeadingSep(t *testing.T) {
	file := createTestFile(`a: apple
b: banana
---
not a 
yaml: doc
`)

	expectedYamlFm := `a: apple
b: banana
`

	expectedContent := `---
not a 
yaml: doc
`

	fmHandler := NewFrontMatterHandler(file)
	err := fmHandler.Split()
	if err != nil {
		panic(err)
	}

	yamlFm := readFile(fmHandler.GetYamlFrontMatterFilename())

	test.AssertResult(t, expectedYamlFm, yamlFm)

	contentBytes, err := io.ReadAll(fmHandler.GetContentReader())
	if err != nil {
		panic(err)
	}
	test.AssertResult(t, expectedContent, string(contentBytes))

	tryRemoveTempFile(file)
	fmHandler.CleanUp()
}

func TestFrontMatterFilenamePreserved(t *testing.T) {
	// Regression test for https://github.com/mikefarah/yq/issues/2538
	// When using --front-matter, the filename operator should return
	// the original filename, not the path to the temporary file.
	file := createTestFile(`---
name: john
---
Some content
`)
	originalFilename := "/path/to/original/file.md"

	fmHandler := NewFrontMatterHandler(file)
	err := fmHandler.Split()
	if err != nil {
		panic(err)
	}

	tempFilename := fmHandler.GetYamlFrontMatterFilename()

	// Register the alias (as the command code does)
	SetFilenameAlias(tempFilename, originalFilename)
	defer ClearFilenameAliases()

	// Verify resolveFilename returns the original name
	resolved := resolveFilename(tempFilename)
	test.AssertResult(t, originalFilename, resolved)

	// Read documents using the temp file, verify they get the original filename
	reader, err := readStream(tempFilename)
	if err != nil {
		panic(err)
	}
	decoder := NewYamlDecoder(ConfiguredYamlPreferences)
	docs, err := readDocuments(reader, tempFilename, 0, decoder)
	if err != nil {
		panic(err)
	}

	if docs.Len() == 0 {
		t.Fatal("expected at least one document")
	}

	firstDoc := docs.Front().Value.(*CandidateNode)
	test.AssertResult(t, originalFilename, firstDoc.filename)

	tryRemoveTempFile(file)
	fmHandler.CleanUp()
}

func TestFrontMatterSplitWithArray(t *testing.T) {
	file := createTestFile(`[1,2,3]
---
not a 
yaml: doc
`)

	expectedYamlFm := "[1,2,3]\n"

	expectedContent := `---
not a 
yaml: doc
`

	fmHandler := NewFrontMatterHandler(file)
	err := fmHandler.Split()
	if err != nil {
		panic(err)
	}

	yamlFm := readFile(fmHandler.GetYamlFrontMatterFilename())

	test.AssertResult(t, expectedYamlFm, yamlFm)

	contentBytes, err := io.ReadAll(fmHandler.GetContentReader())
	if err != nil {
		panic(err)
	}
	test.AssertResult(t, expectedContent, string(contentBytes))

	tryRemoveTempFile(file)
	fmHandler.CleanUp()
}
