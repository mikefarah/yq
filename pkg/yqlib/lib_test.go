package yqlib

import "testing"

func TestGetLogger(t *testing.T) {
	l := GetLogger()
	if l != log {
		t.Fatal("GetLogger should return the yq logger instance, not a copy")
	}
}
