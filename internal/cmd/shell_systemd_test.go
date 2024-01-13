package cmd

import (
	"fmt"
	"strings"
	"testing"
)

func StringPtr(value string) *string {
	tmp := value
	return &tmp
}

func TestCutEncapsulated_ok(t *testing.T) {
	input := "\"encapsulated string with \n line return\""
	cutInput, isEncapsulated := cutEncapsulated(input, "\"")
	if !isEncapsulated {
		t.Error("Test TestCutEncapsulated_ok failing.")
	}
	assertEqual(t, "encapsulated string with \n line return", cutInput)
}

func TestExport_ok(t *testing.T) {

	env := Env{
		"Key":  " just a Value",
		"Ex1":  `'single quotes ' works like that'`,
		"Ex2":  `however, you can't use quotes inline`,
		"Ex3":  `double quotes " are doing the "same" way`,
		"Ex4":  `quotes allows escapes: \n \x`,
		"Ex5":  `quotes are doing trick for chars: ", \, $`,
		"Ex6":  `naked values allows escapes: \a, \b, \c`,
		"Ex7":  `and even \$, and even at the beginning`,
		"Ex8":  `\x all the rest`,
		"Ex9":  `in naked values backslash \allows splitting values`,
		"Ex10": `quotes\nallow multi lines`,
		"Ex11": `'with single quotes around it and '' single quotes in it '''`,
		"Ex12": `"with quotes around it and quotes "" in "" it"`,
	}

	systemdExporter := Systemd
	actualOutput := env.ToShell(systemdExporter)

	expectedOutputMap := map[string]string{
		"Key":  " just a Value",
		"Ex1":  "'single quotes \\' works like that'",
		"Ex2":  "\"however, you can't use quotes inline\"",
		"Ex3":  "\"double quotes \\\" are doing the \\\"same\\\" way\"",
		"Ex4":  `"quotes allows escapes: \n \x"`,
		"Ex5":  "\"quotes are doing trick for chars: \\\", \\, $\"",
		"Ex6":  "\"naked values allows escapes: \\a, \\b, \\c\"",
		"Ex7":  "\"and even \\$, and even at the beginning\"",
		"Ex8":  "\"\\x all the rest\"",
		"Ex9":  "\"in naked values backslash \\allows splitting values\"",
		"Ex10": `"quotes\nallow multi lines"`,
		"Ex11": `'with single quotes around it and \'\' single quotes in it \'\''`,
		"Ex12": `"with quotes around it and quotes \"\" in \"\" it"`,
	}
	for key, expectedValue := range expectedOutputMap {
		_, afterPrefix, PrefixFound := strings.Cut(actualOutput, key+"=")
		actualValue, _, SuffixFound := strings.Cut(afterPrefix, "\n")
		if !PrefixFound || !SuffixFound {
			t.Error("Test for systemd shell failed, couldn't not found expected key=value")
		}
		assertEqual(t, fmt.Sprint(expectedValue), fmt.Sprint(actualValue))
	}

}
