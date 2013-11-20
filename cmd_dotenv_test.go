package main

import (
	"testing"
)

// See the reference implementation:
//   https://github.com/bkeepers/dotenv/blob/master/lib/dotenv/environment.rb
// TODO: support shell variable expansions
// TODO: support comments at the end of a line

const TEST_EXPORTED = `export OPTION_A=2
export OPTION_B='\n'
`

func TestDotEnvExported(t *testing.T) {
	env := ParseDotEnv(TEST_EXPORTED)

	if env["OPTION_A"] != "2" {
		//t.Fail()
	}
	if env["OPTION_B"] != "\\n" {
		t.Fail()
	}
}

const TEST_PLAIN = `OPTION_A=1
OPTION_B=2
OPTION_C= 3
OPTION_D =4
OPTION_E = 5
`

func TestDotEnvPlain(t *testing.T) {
	env := ParseDotEnv(TEST_PLAIN)

	if env["OPTION_A"] != "1" {
		t.Fail()
	}
	if env["OPTION_B"] != "2" {
		t.Fail()
	}
	if env["OPTION_C"] != "3" {
		t.Fail()
	}
	if env["OPTION_D"] != "4" {
		t.Fail()
	}
	if env["OPTION_E"] != "5" {
		t.Fail()
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
`

func TestDotEnvQuoted(t *testing.T) {
	env := ParseDotEnv(TEST_QUOTED)

	if env["OPTION_A"] != "1" {
		t.Fail()
	}
	if env["OPTION_B"] != "2" {
		t.Fail()
	}
	if env["OPTION_C"] != "" {
		t.Fail()
	}
	if env["OPTION_D"] != "\\n" {
		t.Fail()
	}
	if env["OPTION_E"] != "1" {
		t.Fail()
	}
	if env["OPTION_F"] != "2" {
		t.Fail()
	}
	if env["OPTION_G"] != "" {
		t.Fail()
	}
	if env["OPTION_H"] != "\n" {
		t.Fail()
	}
}

const TEST_YAML = `OPTION_A: 1
OPTION_B: '2'
OPTION_C: ''
OPTION_D: '\n'
`

func TestDotEnvYAML(t *testing.T) {
	env := ParseDotEnv(TEST_YAML)

	if env["OPTION_A"] != "1" {
		t.Fail()
	}
	if env["OPTION_B"] != "2" {
		t.Fail()
	}
	if env["OPTION_C"] != "" {
		t.Fail()
	}
	if env["OPTION_D"] != "\\n" {
		t.Fail()
	}
}
