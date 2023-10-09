package cmd

import (
	"fmt"
	"regexp"
)

type pwsh struct{}

// Pwsh shell instance
var Pwsh Shell = pwsh{}

func (sh pwsh) Hook() (string, error) {
	const hook = `using namespace System;
using namespace System.Management.Automation;

$hook = [EventHandler[LocationChangedEventArgs]] {
  param([object] $source, [LocationChangedEventArgs] $eventArgs)
  end {
    $export = {{.SelfPath}} export pwsh;
    if ($export) {
      Invoke-Expression -Command $export;
    }
  }
};
$currentAction = $ExecutionContext.SessionState.InvokeCommand.LocationChangedAction;
if ($currentAction) {
  $ExecutionContext.SessionState.InvokeCommand.LocationChangedAction = [Delegate]::Combine($currentAction, $hook);
}
else {
  $ExecutionContext.SessionState.InvokeCommand.LocationChangedAction = $hook;
};

`
	return hook, nil
}

func (sh pwsh) Export(e ShellExport) (out string) {
	for key, value := range e {
		if value == nil {
			out += sh.unset(key)
		} else {
			out += sh.export(key, *value)
		}
	}
	return out
}

func (sh pwsh) Dump(env Env) (out string) {
	for key, value := range env {
		out += sh.export(key, value)
	}
	return
}

func (sh pwsh) export(key, value string) string {
	value = sh.escape(value)
	if !regexp.MustCompile(`'.*'`).MatchString(value) {
		value = fmt.Sprintf("'%s'", value)
	}
	return fmt.Sprintf("$env:%s=%s;", sh.escape(key), value)
}

func (sh pwsh) unset(key string) string {
	return fmt.Sprintf("Remove-Item -Path 'env:/%s';", sh.escape(key))
}

func (pwsh) escape(str string) string {
	return PowerShellEscape(str)
}

func PowerShellEscape(str string) string {
	if str == "" {
		return "''"
	}
	in := []byte(str)
	out := ""
	i := 0
	l := len(in)
	escape := false

	hex := func(char byte) {
		escape = true
		out += fmt.Sprintf("\\x%02x", char)
	}

	backslash := func(char byte) {
		escape = true
		out += string([]byte{BACKTICK, char})
	}

	escaped := func(str string) {
		escape = true
		out += str
	}

	quoted := func(char byte) {
		escape = true
		out += string([]byte{char})
	}

	literal := func(char byte) {
		out += string([]byte{char})
	}

	for i < l {
		char := in[i]
		switch {
		case char == ACK:
			hex(char)
		case char == TAB:
			escaped("`t")
		case char == LF:
			escaped("`n")
		case char == CR:
			escaped("`r")
		case char <= US:
			hex(char)
		// case char <= AMPERSTAND:
		// 	quoted(char)
		case char == SINGLE_QUOTE:
			backslash(char)
		case char <= PLUS:
			quoted(char)
		case char <= NINE:
			literal(char)
		// case char <= QUESTION:
		// 	quoted(char)
		case char <= UPPERCASE_Z:
			literal(char)
		// case char == OPEN_BRACKET:
		// 	quoted(char)
		// case char == BACKSLASH:
		// 	quoted(char)
		case char == UNDERSCORE:
			literal(char)
		// case char <= CLOSE_BRACKET:
		// 	quoted(char)
		// case char <= BACKTICK:
		// 	quoted(char)
		// case char <= TILDA:
		// 	quoted(char)
		case char == DEL:
			hex(char)
		default:
			quoted(char)
		}
		i++
	}

	if escape {
		out = "'" + out + "'"
	}

	return out
}
