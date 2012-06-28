package main

import (
	"direnv"
	"fmt"
	"log"
)

func main() {
	env := direnv.FilteredEnv()
	str, err := env.Serialize()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Print(str)
}
