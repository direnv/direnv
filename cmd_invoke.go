package main

import (
	"flag"
	"fmt"
)

var CmdInvoke = &Cmd{
	Name: "invoke",
	Desc: "run EXECUTABLE with appropriate environment for WORKDIR",
	Args: []string{"WORKDIR EXECUTABLE [ARGS...]"},
	Fn: func(env Env, args []string) (err error) {

		flagset := flag.NewFlagSet(args[0], flag.ExitOnError)
		flagset.Parse(args[1:])

		workdir := flagset.Arg(0)
		if workdir == "" {
			return fmt.Errorf("WORKDIR missing")
		}
		program := flagset.Args()[1:]
		if len(program) < 1 {
			return fmt.Errorf("EXECUTABLE missing")
		}

		var config *Config
		if config, err = LoadConfig(env); err != nil {
			return
		}

		script := `
	        DIRENV_PATH="%s"
    	    eval "$(${DIRENV_PATH} export bash)"    
    	    exec "$@"
	    `
		script = fmt.Sprintf(script, config.SelfPath)

		// Invoke `bash -c "eval \"$($DIRENV export bash)\"; exec $@" -- program [args]`
		// in the destination directory so as to have the correct environment
		bash_args := append([]string{"-c", script, "--"}, program...)

		err = Invoke(workdir, bash_args)
		return err
	},
}
