package yqlib

import (
	"io/ioutil"
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
	bytes, err := ioutil.ReadFile(filename)
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

	contentBytes, err := ioutil.ReadAll(fmHandler.GetContentReader())
	if err != nil {
		panic(err)
	}
	test.AssertResult(t, expectedContent, string(contentBytes))

	tryRemoveFile(file)
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

	contentBytes, err := ioutil.ReadAll(fmHandler.GetContentReader())
	if err != nil {
		panic(err)
	}
	test.AssertResult(t, expectedContent, string(contentBytes))

	tryRemoveFile(file)
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

	contentBytes, err := ioutil.ReadAll(fmHandler.GetContentReader())
	if err != nil {
		panic(err)
	}
	test.AssertResult(t, expectedContent, string(contentBytes))

	tryRemoveFile(file)
	fmHandler.CleanUp()
}
