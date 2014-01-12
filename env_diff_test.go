package main

import (
	"testing"
)

func TestEnvDiff(t *testing.T) {
	diff := &EnvDiff{map[string]string{"FOO": "bar"}, map[string]string{"BAR": "baz"}}

	out := diff.Serialize()

	diff2, err := LoadEnvDiff(out)
	if err != nil {
		t.Errorf("parse error", err)
	}

	if len(diff2.Prev) != 1 {
		t.Errorf("len(diff2.prev) != 1", len(diff2.Prev))
	}

	if len(diff2.Next) != 1 {
		t.Errorf("len(diff2.next) != 0", len(diff2.Next))
	}
}

func TestIgnoredEnv(t *testing.T) {
	if !IgnoredEnv(DIRENV_BASH) {
		t.Fail()
	}
	if IgnoredEnv(DIRENV_DIFF) {
		t.Fail()
	}
	if !IgnoredEnv("_") {
		t.Fail()
	}
	if !IgnoredEnv("__fish_foo") {
		t.Fail()
	}
	if !IgnoredEnv("__fishx") {
		t.Fail()
	}
}
