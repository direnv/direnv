package main

import (
	"log"
	"os"
	"time"
)

func main() {
	env := GetEnv()

	done := make(chan bool, 1)
	go func() {
		select {
		case <-done:
			return
		case <-time.After(2 * time.Second):
			log.Printf("direnv(%v) is taking a while to execute. Use CTRL-C to give up.", os.Args)
		}
	}()

	err := CommandsDispatch(env, os.Args)
	done <- true
	if err != nil {
		os.Exit(1)
	}
}
