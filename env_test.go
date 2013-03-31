package main

import (
	"testing"
)

func TestEnv(t *testing.T) {
	env := Env{"FOO": "bar"}

	out := env.Serialize()

	env2, err := ParseEnv(out)
	if err != nil {
		t.Fail()
	}

	if env2["FOO"] != "bar" {
		t.Fail()
	}

	if len(env2) != 1 {
		t.Fail()
	}
}

func TestIgnoredKeys(t *testing.T) {
	if ignoredKey("DIRENV_FOOBAR") {
		t.Fail()
	}
	if ignoredKey("DIRENV_") {
		t.Fail()
	}
	if !ignoredKey("_") {
		t.Fail()
	}
}
