package dotenv_test

import (
	"os"
	"testing"

	dotenv "github.com/direnv/go-dotenv"
)

func shouldNotHaveEmptyKey(t *testing.T, env map[string]string) {
	if _, ok := env[""]; ok {
		t.Error("should not have empty key")
	}
}

func envShouldContain(t *testing.T, env map[string]string, key string, value string) {
	if env[key] != value {
		t.Errorf("%s: %s, expected %s", key, env[key], value)
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

const TEST_VARIABLE_EXPANSION = `
OPTION_A=$FOO
OPTION_B="$FOO"
OPTION_C=${FOO}
OPTION_D="${FOO}"
OPTION_E='$FOO'
OPTION_F=$FOO/bar
OPTION_G="$FOO/bar"
OPTION_H=${FOO}/bar
OPTION_I="${FOO}/bar"
OPTION_J='$FOO/bar'
OPTION_K=$BAR
OPTION_L="$BAR"
OPTION_M=${BAR}
OPTION_N="${BAR}"
OPTION_O='$BAR'
OPTION_P=$BAR/baz
OPTION_Q="$BAR/baz"
OPTION_R=${BAR}/baz
OPTION_S="${BAR}/baz"
OPTION_T='$BAR/baz'
OPTION_U="$OPTION_A/bar"
OPTION_V=$OPTION_A/bar
OPTION_W="$OPTION_A/bar"
OPTION_X=${OPTION_A}/bar
OPTION_Y="${OPTION_A}/bar"
OPTION_Z='$OPTION_A/bar'
OPTION_A1="$OPTION_A/bar/${OPTION_H}/$FOO"
`

func TestVariableExpansion(t *testing.T) {
	err := os.Setenv("FOO", "foo")
	if err != nil {
		t.Fatalf("unable to set environment variable for testing: %s", err)
	}

	env := dotenv.MustParse(TEST_VARIABLE_EXPANSION)
	shouldNotHaveEmptyKey(t, env)

	envShouldContain(t, env, "OPTION_A", "foo")
	envShouldContain(t, env, "OPTION_B", "foo")
	envShouldContain(t, env, "OPTION_C", "foo")
	envShouldContain(t, env, "OPTION_D", "foo")
	envShouldContain(t, env, "OPTION_E", "$FOO")
	envShouldContain(t, env, "OPTION_F", "foo/bar")
	envShouldContain(t, env, "OPTION_G", "foo/bar")
	envShouldContain(t, env, "OPTION_H", "foo/bar")
	envShouldContain(t, env, "OPTION_I", "foo/bar")
	envShouldContain(t, env, "OPTION_J", "$FOO/bar")
	envShouldContain(t, env, "OPTION_K", "")
	envShouldContain(t, env, "OPTION_L", "")
	envShouldContain(t, env, "OPTION_M", "")
	envShouldContain(t, env, "OPTION_N", "")
	envShouldContain(t, env, "OPTION_O", "$BAR")
	envShouldContain(t, env, "OPTION_P", "/baz")
	envShouldContain(t, env, "OPTION_Q", "/baz")
	envShouldContain(t, env, "OPTION_R", "/baz")
	envShouldContain(t, env, "OPTION_S", "/baz")
	envShouldContain(t, env, "OPTION_T", "$BAR/baz")
	envShouldContain(t, env, "OPTION_U", "foo/bar")
	envShouldContain(t, env, "OPTION_V", "foo/bar")
	envShouldContain(t, env, "OPTION_W", "foo/bar")
	envShouldContain(t, env, "OPTION_X", "foo/bar")
	envShouldContain(t, env, "OPTION_Y", "foo/bar")
	envShouldContain(t, env, "OPTION_Z", "$OPTION_A/bar")
	envShouldContain(t, env, "OPTION_A1", "foo/bar/foo/bar/foo")
}
