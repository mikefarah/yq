package cmd

import (
	"testing"

	"github.com/mikefarah/yq/v3/test"
)

func TestCompareSameCmd(t *testing.T) {
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "compare ../examples/data1.yaml ../examples/data1.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := ``
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestCompareDifferentCmd(t *testing.T) {
	forceOsExit = false
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "compare ../examples/data1.yaml ../examples/data3.yaml")

	expectedOutput := `-a: simple # just the best
-b: [1, 2]
+a: "simple" # just the best
+b: [1, 3]
 c:
   test: 1
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestComparePrettyCmd(t *testing.T) {
	forceOsExit = false
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "compare -P ../examples/data1.yaml ../examples/data3.yaml")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := ` a: simple # just the best
 b:
   - 1
-  - 2
+  - 3
 c:
   test: 1
`
	test.AssertResult(t, expectedOutput, result.Output)
}

func TestComparePathsCmd(t *testing.T) {
	forceOsExit = false
	cmd := getRootCommand()
	result := test.RunCmd(cmd, "compare -P -ppv ../examples/data1.yaml ../examples/data3.yaml **")
	if result.Error != nil {
		t.Error(result.Error)
	}
	expectedOutput := ` a: simple # just the best
 b.[0]: 1
-b.[1]: 2
+b.[1]: 3
 c.test: 1
`
	test.AssertResult(t, expectedOutput, result.Output)
}
