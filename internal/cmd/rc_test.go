package cmd

import (
	"runtime"
	"testing"
)

func TestSomething(t *testing.T) {
	paths := eachDir("/foo/b//bar/")
	if len(paths) != 4 {
		t.Fail()
	}
	// TODO: fix me for windows
	if runtime.GOOS != "windows" {
		paths = eachDir("/")
		if len(paths) != 1 && paths[0] != "/" {
			t.Fail()
		}
	}
}
