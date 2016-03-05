package main

var CommandWatch = &Cmd{
	Name:    "watch",
	Desc:    "Adds a path to the list that direnv watches for changes",
	Args:    []string{"PATH"},
	Private: false,
	Fn:      watchCommand,
}

func watchCommand(env Env, args []string) (err error) {
	return
}
