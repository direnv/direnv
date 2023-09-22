package cmd

import "fmt"

type powershell struct{}

// PowerShell shell instance
var PowerShell Shell = powershell{}

func (sh powershell) Hook() (string, error) {
	const hook = `using namespace System;
using namespace System.Management.Automation;

$hook = [EventHandler[LocationChangedEventArgs]] {
  param([object] $source, [LocationChangedEventArgs] $eventArgs)
  end {
    $export = {{.SelfPath}} export powershell;
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

func (sh powershell) Export(e ShellExport) (out string) {
	for key, value := range e {
		if value == nil {
			out += sh.unset(key)
		} else {
			out += sh.export(key, *value)
		}
	}
	return out
}

func (sh powershell) Dump(env Env) (out string) {
	for key, value := range env {
		out += sh.export(key, value)
	}
	return
}

func (sh powershell) export(key, value string) string {
	return fmt.Sprintf("$env:%s = %s;", sh.escape(key), sh.escape(value))
}

func (sh powershell) unset(key string) string {
	return fmt.Sprintf("Remove-Item -Path Env:/%s;", sh.escape(key))
}

func (powershell) escape(str string) string {
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
		out += string([]byte{BACKSLASH, char})
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
		case char <= AMPERSTAND:
			quoted(char)
		case char == SINGLE_QUOTE:
			backslash(char)
		case char <= PLUS:
			quoted(char)
		case char <= NINE:
			literal(char)
		case char <= QUESTION:
			quoted(char)
		case char <= UPPERCASE_Z:
			literal(char)
		case char == OPEN_BRACKET:
			quoted(char)
		case char == BACKSLASH:
			backslash(char)
		case char == UNDERSCORE:
			literal(char)
		case char <= CLOSE_BRACKET:
			quoted(char)
		case char <= BACKTICK:
			quoted(char)
		case char <= TILDA:
			quoted(char)
		case char == DEL:
			hex(char)
		default:
			hex(char)
		}
		i++
	}

	if escape {
		out = "'" + out + "'"
	}

	return out
}
