package cmd

import (
	"fmt"
)

type pwsh struct{}

// Pwsh shell instance
var Pwsh Shell = pwsh{}

func (sh pwsh) Hook() (string, error) {
	const hook = `using namespace System;
using namespace System.Management.Automation;

if ($PSVersionTable.PSVersion.Major -lt 7 -or ($PSVersionTable.PSVersion.Major -eq 7 -and $PSVersionTable.PSVersion.Minor -lt 2)) {
    throw "direnv: PowerShell version $($PSVersionTable.PSVersion) does not meet the minimum required version 7.2!"
}

$hook = [EventHandler[LocationChangedEventArgs]] {
  param([object] $source, [LocationChangedEventArgs] $eventArgs)
  end {
    $export = ({{.SelfPath}} export pwsh) -join [Environment]::NewLine;
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

func (sh pwsh) Export(e ShellExport) (string, error) {
	var out string
	for key, value := range e {
		if key != "" {
			if value == nil {
				out += sh.unset(key)
			} else {
				out += sh.export(key, *value)
			}
		}
	}
	return out, nil
}

func (sh pwsh) Dump(env Env) (string, error) {
	var out string
	for key, value := range env {
		out += sh.export(key, value)
	}
	return out, nil
}

func (sh pwsh) export(key, value string) string {
	return fmt.Sprintf("${env:%s}='%s';", sh.escapeEnvKey(key), sh.escapeVerbatimString(value))
}

func (sh pwsh) unset(key string) string {
	return fmt.Sprintf("Remove-Item -LiteralPath 'env:/%s';", sh.escapeVerbatimEnvKey(key))
}

func (pwsh) escapeEnvKey(str string) string {
	return PowerShellEscapeEnvKey(str)
}

// PowerShellEscapeEnvKey escapes environment variable keys for PowerShell.
func PowerShellEscapeEnvKey(str string) string {
	if str == "" {
		return "__DiReNv_UnReAcHaBlE__"
	}
	in := []byte(str)
	out := ""
	i := 0
	l := len(in)

	escaped := func(str string) {
		out += str
	}

	hex := func(char byte) {
		out += fmt.Sprintf("\\x%02x", char)
	}

	literal := func(char byte) {
		out += string([]byte{char})
	}

	for i < l {
		char := in[i]
		switch char {
		case STAR:
			hex(char)
		case COLON:
			hex(char)
		case EQUALS:
			hex(char)
		case QUESTION:
			hex(char)
		case OPEN_BRACKET:
			hex(char)
		case CLOSE_BRACKET:
			hex(char)
		case OPEN_CURLY_BRACE:
			escaped("`{")
		case CLOSE_CURLY_BRACE:
			escaped("`}")
		default:
			literal(char)
		}
		i++
	}

	return out
}

func (pwsh) escapeVerbatimEnvKey(str string) string {
	return PowerShellEscapeVerbatimEnvKey(str)
}

// PowerShellEscapeVerbatimEnvKey escapes environment variable keys using verbatim strings for PowerShell.
func PowerShellEscapeVerbatimEnvKey(str string) string {
	if str == "" {
		return "__DiReNv_UnReAcHaBlE__"
	}
	in := []byte(str)
	out := ""
	i := 0
	l := len(in)

	escaped := func(str string) {
		out += str
	}

	literal := func(char byte) {
		out += string([]byte{char})
	}

	for i < l {
		char := in[i]
		switch char {
		case SINGLE_QUOTE:
			escaped("''")
		default:
			literal(char)
		}
		i++
	}

	return out
}
func (pwsh) escapeVerbatimString(str string) string {
	return PowerShellEscapeVerbatimString(str)
}

// PowerShellEscapeVerbatimString escapes strings using verbatim string literals for PowerShell.
func PowerShellEscapeVerbatimString(str string) string {
	if str == "" {
		return ""
	}
	in := []byte(str)
	out := ""
	i := 0
	l := len(in)

	escaped := func(str string) {
		out += str
	}

	literal := func(char byte) {
		out += string([]byte{char})
	}

	for i < l {
		char := in[i]
		switch char {
		case SINGLE_QUOTE:
			escaped("''")
		default:
			literal(char)
		}
		i++
	}

	return out
}

/*
   1. Minimal handling required for verbatim strings:
   Characters in a verbatim string (e.g.: 'a single quoted string') don't require escaping
   except for the single quote character itself which is escaped by doubling it
   (i.e.: '''' -eq "'" -and ''''.Length -eq 1).

   2. Handling any exported newline or carriage return characters from the PowerShell hook:
   Newline or carriage return characters in any part of the output of `direnv export pwsh`
   will produce an array of strings when imported into PowerShell. To join all parts of the
   array into a single string, the following is done in the PowerShell hook:

   `(direnv export pwsh) -join [Environment]::NewLine`

   3. Allowing PowerShell variable names with special characters:
   PowerShell environment variable names may contain "special characters" when enclosed in
   curly braces like this:

   ${env:name-with-special-chars-like-dashes} = 'value'

   The following special characters may NOT be used in such names: *, ?, :, =, [, ]
   These invalid special characters are mapped to hex codes (e.g.: "*" -> "\x2A").

   Curly braces may be used, if escaped with a backtick: `{, `}

   For more info on Pwsh variable names that include special characters see:
   https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_variables#variable-names-that-include-special-characters

   4. Paranoid handling of paths when removing environment variables:
   Use `Remove-Item -LiteralPath <PATH>` instead of `Remove-Item -Path <PATH>` to avoid any
   potential wildcard interpretations.

   5. Paranoid handling of potentially empty key names:
   I'm not sure if `Remove-Item -LiteralPath 'env:'` or `${env:} = 'value'` could be abused so two
   overlapping steps are taken to avoid this issue:
     a. Empty key names are skipped.
     b. Empty key names are replaced with "__DiReNv_UnReAcHaBlE__".
*/
