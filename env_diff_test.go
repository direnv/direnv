package main

import (
	"reflect"
	"testing"
)

func TestEnvDiff(t *testing.T) {
	diff := &EnvDiff{map[string]string{"FOO": "bar"}, map[string]string{"BAR": "baz"}}

	out := diff.Serialize()

	diff2, err := LoadEnvDiff(out)
	if err != nil {
		t.Error("parse error", err)
	}

	if len(diff2.Prev) != 1 {
		t.Error("len(diff2.prev) != 1", len(diff2.Prev))
	}

	if len(diff2.Next) != 1 {
		t.Error("len(diff2.next) != 0", len(diff2.Next))
	}
}

// Issue #114
// Check that empty environment variables correctly appear in the diff
func TestEnvDiffEmptyValue(t *testing.T) {
	before := Env{}
	after := Env{"FOO": ""}

	diff := BuildEnvDiff(before, after)

	if !reflect.DeepEqual(diff.Next, map[string]string(after)) {
		t.Errorf("diff.Next != after (%#+v != %#+v)", diff.Next, after)
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
