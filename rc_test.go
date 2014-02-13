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

func TestRootDir(t *testing.T) {
	var r string
	r = rootDir("/foo")
	if r != "/foo" {
		t.Error(r)
	}

	r = rootDir("/foo/bar")
	if r != "/foo" {
		t.Error(r)
	}
}
