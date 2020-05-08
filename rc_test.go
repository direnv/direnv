package main

import (
	"testing"
)

func TestSomething(t *testing.T) {
	paths := eachDir("/foo/b//bar/")
	if len(paths) != 4 {
		t.Fail()
	}
	paths = eachDir("/")
	if len(paths) != 1 && paths[0] != "/" {
		t.Fail()
	}
}
