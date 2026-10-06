package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"text/template"
)

func TestHookPreservesSIGINT(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("requires Unix signal handling")
	}
	for _, shellName := range []string{"bash", "zsh"} {
		t.Run(shellName, func(t *testing.T) {
			shellPath, err := exec.LookPath(shellName)
			if err != nil {
				t.Skipf("%s is not installed", shellName)
			}
			hook, err := DetectShell(shellName).Hook()
			if err != nil {
				t.Fatal(err)
			}
			tmpl, err := template.New("hook").Parse(hook)
			if err != nil {
				t.Fatal(err)
			}
			var rendered bytes.Buffer
			if err := tmpl.Execute(&rendered, HookContext{SelfPath: "__direnv_export"}); err != nil {
				t.Fatal(err)
			}
			cases := []struct {
				name   string
				setup  string
				signal string
			}{
				{"default", "trap - INT", ":"},
				{"ignored", "trap '' INT", `kill -INT "$$"; test "$signal_received" = no`},
				{"custom", `trap 'signal_received="handled with spaces"' INT`, `kill -INT "$$"; test "$signal_received" = "handled with spaces"`},
			}
			if shellName == "zsh" {
				cases = append(cases, struct {
					name   string
					setup  string
					signal string
				}{"function", `TRAPINT() { signal_received="handled with spaces"; }`, `kill -INT "$$"; test "$signal_received" = "handled with spaces"`})
			}
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					// Bash preserves the previous command's status; Zsh's hook returns zero.
					expectedStatus := 0
					args := []string{"-f", "-c"}
					if shellName == "bash" {
						expectedStatus = 23
						args = []string{"--noprofile", "--norc", "-c"}
					}
					script := fmt.Sprintf(`
__direnv_export() {
  if test "$mode" = failed_export; then
    return 1
  fi
  printf 'kill -INT $$; export hook_ran=yes;\n'
  if test "$mode" = failed_eval; then
    printf 'false\n'
  fi
}
%s
%s
expected_traps=$(trap)
if test -n "${ZSH_VERSION:-}"; then
  expected_options=$(setopt)
fi
return_status() { return 23; }
for mode in success failed_export failed_eval; do
  for repeat in 1 2; do
    signal_received=no
    hook_ran=no
    return_status
    _direnv_hook
    actual_status=$?
    test "$actual_status" -eq %d || { echo "changed exit status: $actual_status"; exit 1; }
    test "$(trap)" = "$expected_traps" || { echo "changed SIGINT trap ($mode, invocation $repeat)"; exit 1; }
    test "$signal_received" = no || { echo "SIGINT was not ignored during eval"; exit 1; }
    if test "$mode" != failed_export; then
      test "$hook_ran" = yes || exit 1
    fi
    if test -n "${ZSH_VERSION:-}"; then
      test "$(setopt)" = "$expected_options" || { echo "changed shell options"; exit 1; }
    fi
    { %s; } || { echo "original SIGINT handler was not preserved"; exit 1; }
  done
done
`, rendered.String(), tc.setup, expectedStatus, tc.signal)
					command := exec.Command(shellPath, append(args, script)...)
					command.Env = append(os.Environ(), "BASH_ENV=", "ENV=", "ZDOTDIR="+t.TempDir())
					if output, err := command.CombinedOutput(); err != nil {
						t.Fatalf("hook failed: %v\n%s", err, output)
					}
				})
			}
		})
	}
}
