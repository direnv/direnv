package main

import (
	"fmt"
	"github.com/direnv/direnv/gzenv"
)

// `direnv gzenv encode <str>` and `direnv gzenv decode <str>`
var CmdGzEnv = &Cmd{
	Name:    "gzenv",
	Desc:    "Encode and decode a string using the gzenv format",
	Args:    []string{"<encode|decode>", "STRING", "[...STRING]"},
	Private: true,
	Fn: func(env Env, args []string) (err error) {
		if len(args) < 2 || (args[1] != "encode" && args[1] != "decode") {
			return fmt.Errorf("missing <encode|decode> action")
		}

		encode := args[1] == "encode"

		for _, s := range args[2:] {
			var r interface{}
			if encode {
				r = gzenv.Marshal(s)
			} else {
				err := gzenv.Unmarshal(s, &r)
				if err != nil {
					fmt.Println("Bad encoding: ", s)
					continue
				}
			}
			fmt.Println(r)
		}

		return
	},
}
