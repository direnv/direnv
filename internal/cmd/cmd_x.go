package cmd

// CmdX is `direnv x COMMAND [...ARGS]`.
var CmdX = &Cmd{
	Name: "x",
	Desc: "Executes a command in the current directory's environment (shorthand for exec .)",
	Args: []string{"COMMAND", "[...ARGS]"},
	Action: actionWithConfig(func(env Env, args []string, config *Config) error {
		return cmdExecAction(env, append([]string{args[0], "."}, args[1:]...), config)
	}),
}
