package main

import (
	"testing"
)

var getValueTests = []struct {
	argument        string
	expectedResult  interface{}
	testDescription string
}{
	{"true", true, "boolean"},
	{"\"true\"", "true", "boolean as string"},
	{"3.4", 3.4, "number"},
	{"\"3.4\"", "3.4", "number as string"},
}

func TestGetValue(t *testing.T) {
	for _, tt := range getValueTests {
		assertResultWithContext(t, tt.expectedResult, getValue(tt.argument), tt.testDescription)
	}
}
