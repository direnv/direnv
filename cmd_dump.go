package main

import (
	"fmt"
)

// Used to export the inner bash state at the end of execution.
func Dump(env Env, args []string) (err error) {
	fmt.Println(env.Filtered().Serialize())
	return
}
