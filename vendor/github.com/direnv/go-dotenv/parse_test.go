package dotenv_test

import (
	"testing"

	"github.com/direnv/go-dotenv"
)

func shouldNotHaveEmptyKey(t *testing.T, env map[string]string) {
	if _, ok := env[""]; ok {
		t.Error("should not have empty key")
	}
}

// See the reference implementation:
//   https://github.com/bkeepers/dotenv/blob/master/lib/dotenv/environment.rb
// TODO: support shell variable expansions

const TEST_EXPORTED = `export OPTION_A=2
export OPTION_B='\n' # foo
#export OPTION_C=3
export OPTION_D=
export OPTION_E="foo"
`

func TestDotEnvExported(t *testing.T) {
	env := dotenv.MustParse(TEST_EXPORTED)
	shouldNotHaveEmptyKey(t, env)

	if env["OPTION_A"] != "2" {
		t.Error("OPTION_A")
	}
	if env["OPTION_B"] != "\\n" {
		t.Error("OPTION_B")
	}
	if env["OPTION_C"] != "" {
		t.Error("OPTION_C", env["OPTION_C"])
	}
	if v, ok := env["OPTION_D"]; !(v == "" && ok) {
		t.Error("OPTION_D")
	}
	if env["OPTION_E"] != "foo" {
		t.Error("OPTION_E")
	}
}

const TEST_PLAIN = `OPTION_A=1
OPTION_B=2
OPTION_C= 3
OPTION_D =4
OPTION_E = 5
OPTION_F=
OPTION_G =
`

func TestDotEnvPlain(t *testing.T) {
	env := dotenv.MustParse(TEST_PLAIN)
	shouldNotHaveEmptyKey(t, env)

	if env["OPTION_A"] != "1" {
		t.Error("OPTION_A")
	}
	if env["OPTION_B"] != "2" {
		t.Error("OPTION_B")
	}
	if env["OPTION_C"] != "3" {
		t.Error("OPTION_C")
	}
	if env["OPTION_D"] != "4" {
		t.Error("OPTION_D")
	}
	if env["OPTION_E"] != "5" {
		t.Error("OPTION_E")
	}
	if v, ok := env["OPTION_F"]; !(v == "" && ok) {
		t.Error("OPTION_F")
	}
	if v, ok := env["OPTION_G"]; !(v == "" && ok) {
		t.Error("OPTION_G")
	}
}

const TEST_SOLO_EMPTY = "SOME_VAR="

func TestSoloEmpty(t *testing.T) {
	env := dotenv.MustParse(TEST_SOLO_EMPTY)
	shouldNotHaveEmptyKey(t, env)

	v, ok := env["SOME_VAR"]
	if !ok {
		t.Error("SOME_VAR missing")
	}
	if v != "" {
		t.Error("SOME_VAR should be empty")
	}
}

const TEST_QUOTED = `OPTION_A='1'
OPTION_B='2'
OPTION_C=''
OPTION_D='\n'
OPTION_E="1"
OPTION_F="2"
OPTION_G=""
OPTION_H="\n"
#OPTION_I="3"
`

func TestDotEnvQuoted(t *testing.T) {
	env := dotenv.MustParse(TEST_QUOTED)
	shouldNotHaveEmptyKey(t, env)

	if env["OPTION_A"] != "1" {
		t.Error("OPTION_A")
	}
	if env["OPTION_B"] != "2" {
		t.Error("OPTION_B")
	}
	if env["OPTION_C"] != "" {
		t.Error("OPTION_C")
	}
	if env["OPTION_D"] != "\\n" {
		t.Error("OPTION_D")
	}
	if env["OPTION_E"] != "1" {
		t.Error("OPTION_E")
	}
	if env["OPTION_F"] != "2" {
		t.Error("OPTION_F")
	}
	if env["OPTION_G"] != "" {
		t.Error("OPTION_G")
	}
	if env["OPTION_H"] != "\n" {
		t.Error("OPTION_H")
	}
	if env["OPTION_I"] != "" {
		t.Error("OPTION_I")
	}
}

const TEST_YAML = `OPTION_A: 1
OPTION_B: '2'
OPTION_C: ''
OPTION_D: '\n'
#OPTION_E: '333'
OPTION_F: 
`

func TestDotEnvYAML(t *testing.T) {
	env := dotenv.MustParse(TEST_YAML)
	shouldNotHaveEmptyKey(t, env)

	if env["OPTION_A"] != "1" {
		t.Error("OPTION_A")
	}
	if env["OPTION_B"] != "2" {
		t.Error("OPTION_B")
	}
	if env["OPTION_C"] != "" {
		t.Error("OPTION_C")
	}
	if env["OPTION_D"] != "\\n" {
		t.Error("OPTION_D")
	}
	if env["OPTION_E"] != "" {
		t.Error("OPTION_E")
	}
	if v, ok := env["OPTION_F"]; !(v == "" && ok) {
		t.Error("OPTION_F")
	}
}

func TestFailingMustParse(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("should panic")
		}
	}()
	dotenv.MustParse("...")
}

const TEST_COMMENT_OVERRIDE = `
VARIABLE=value
#VARIABLE=disabled_value
`

func TestCommentOverride(t *testing.T) {
	env := dotenv.MustParse(TEST_COMMENT_OVERRIDE)
	shouldNotHaveEmptyKey(t, env)

	if env["VARIABLE"] != "value" {
		t.Error("VARIABLE should == value, not", env["VARIABLE"])
	}
}
