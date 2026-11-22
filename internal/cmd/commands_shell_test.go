package cmd

import "testing"

func TestCommandsRejectUnknownShell(t *testing.T) {
	t.Setenv(DIRENV_DUMP_FILE_PATH, "")
	commands := []struct {
		command *Cmd
		args    []string
	}{
		{CmdExport, nil},
		{CmdHook, nil},
		{CmdDump, nil},
		{CmdWatch, []string{"file"}},
		{CmdWatchDir, []string{"directory"}},
		{CmdWatchList, nil},
		{CmdCheckRequired, []string{".envrc", "file"}},
	}
	for _, tc := range commands {
		t.Run(tc.command.Name, func(t *testing.T) {
			for _, shellName := range []string{"unknown", ""} {
				t.Run(shellName, func(t *testing.T) {
					args := append([]string{tc.command.Name, shellName}, tc.args...)
					err := tc.command.Action.Call(Env{}, args, &Config{})
					want := "unknown target shell '" + shellName + "'"
					if err == nil || err.Error() != want {
						t.Fatalf("expected %q, got %v", want, err)
					}
				})
			}
		})
	}
}
